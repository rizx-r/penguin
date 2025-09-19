package websocket

import (
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
	"time"
)

type (
	Conn struct {
		idleMu sync.Mutex

		Uid string

		*websocket.Conn
		s *Server

		idle              time.Time     // 记录这个 WebSocket 连接最近一次活跃的时间戳
		maxConnectionIdle time.Duration // 允许的最大空闲时长

		messageMu      sync.Mutex
		readMessage    []*Message
		readMessageSeq map[string]*Message

		message chan *Message

		done chan struct{}
	}
)

func NewConn(s *Server, w http.ResponseWriter, r *http.Request) *Conn {
	c, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.Errorf("websocket upgrader error: %s", err)
		return nil
	}

	conn := &Conn{
		Conn:              c,
		s:                 s,
		idle:              time.Now(), // 新建连接时，设置为当前时间
		maxConnectionIdle: s.opt.maxConnectionIdle,
		readMessage:       make([]*Message, 0, 2),
		readMessageSeq:    make(map[string]*Message, 2),
		message:           make(chan *Message, 2),
		done:              make(chan struct{}),
	}
	go conn.keepalive()
	return conn
}

func (c *Conn) appendMsgMq(msg *Message) {
	c.messageMu.Lock()
	defer c.messageMu.Unlock()
	// 读队列中
	if m, ok := c.readMessageSeq[msg.Id]; ok {
		// 记录已有消息记录，改消息已有Ack确认
		if len(c.readMessage) == 0 {
			// 队列中没有该消息
			return
		}

		// msg.AckSeq > m.Ack
		if m.AckSeq >= msg.AckSeq {
			// 没有进行Ack确认，重复
			return
		}
		c.readMessageSeq[msg.Id] = msg
		return
	}

	// 没有进行ack确，避免客户端重复发送多余的ack消息
	if msg.FrameType == FrameAck {
		return
	}

	c.readMessage = append(c.readMessage, msg)
	c.readMessageSeq[msg.Id] = msg
}

func (c *Conn) ReadMessage() (messageType int, p []byte, err error) {
	messageType, p, err = c.Conn.ReadMessage()
	c.idleMu.Lock()
	defer c.idleMu.Unlock()
	c.idle = time.Time{} // 清零，代表刚刚有消息活跃
	return
}

func (c *Conn) WriteMessage(messageType int, p []byte) error {
	c.idleMu.Lock()
	defer c.idleMu.Unlock()
	err := c.Conn.WriteMessage(messageType, p)
	c.idle = time.Time{} // 同样清零
	return err
}

func (c *Conn) Close() error {
	select {
	case <-c.done:
	default:
		close(c.done)
	}
	return c.Conn.Close()
}

func (c *Conn) keepalive() {
	idleTimer := time.NewTimer(c.maxConnectionIdle)
	defer idleTimer.Stop()

	for {
		select {
		case <-idleTimer.C:
			c.idleMu.Lock()
			idle := c.idle
			if idle.IsZero() {
				c.idleMu.Unlock()
				idleTimer.Reset(c.maxConnectionIdle)
				continue
			}
			val := c.maxConnectionIdle - time.Since(idle)
			c.idleMu.Unlock()
			if val <= 0 {
				// 空闲超时，关闭连接
				c.s.CloseConnection(c)
				return
			}
			idleTimer.Reset(val)
		case <-c.done:
			return
		}
	}
}
