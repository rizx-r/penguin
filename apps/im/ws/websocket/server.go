package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
	"sync"
	"time"
)

type (
	Server struct {
		sync.RWMutex
		opt            *serverOption
		authentication Authentication
		routes         map[string]HandlerFunc
		addr           string
		patten         string
		connToUser     map[*Conn]string
		userToConn     map[string]*Conn
		upgrader       *websocket.Upgrader // ?? *
		logx.Logger
	}

	AckType int
)

const (
	NoAck    AckType = iota // 不进行Ack验证
	OnlyAck                 // 只有一次Ack
	RigorAck                // 严格Ack
)

func (t AckType) ToString() string {
	switch t {
	case NoAck:
		return "NoAck"
	case OnlyAck:
		return "OnlyAck"
	case RigorAck:
		return "RigorAck"
	}
	return ""
}

func NewServer(addr string, opts ...ServerOptions) *Server {
	opt := newServerOptions(opts...)
	return &Server{
		authentication: opt.Authentication,
		opt:            &opt,
		routes:         make(map[string]HandlerFunc),
		addr:           addr,
		patten:         opt.patten,
		connToUser:     make(map[*Conn]string),
		userToConn:     make(map[string]*Conn),
		upgrader:       &websocket.Upgrader{},
		Logger:         logx.WithContext(context.Background()),
	}
}

func (srv *Server) ServerWs(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			srv.Errorf("server handler ws error: %v", err)
		}
	}()

	conn := NewConn(srv, w, r)
	if conn == nil {
		return
	}
	//conn, err := srv.upgrader.Upgrade(w, r, nil)
	//if err != nil {
	//	srv.Errorf("server upgrader ws error: %v", err)
	//	return
	//}

	if !srv.authentication.Authenticate(w, r) {
		srv.Send(&Message{FrameType: FrameData, Data: fmt.Sprintf("无访问权限")}, conn)
		//conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprint("无访问权限")))
		conn.Close()
		return
	}

	// record connection
	srv.addConn(conn, r)
	// handle connection
	go srv.handleConnection(conn)
}

func (srv *Server) addConn(conn *Conn, req *http.Request) {
	uid := srv.authentication.UserId(req)
	srv.RWMutex.Lock()
	defer srv.RWMutex.Unlock() // defer 会在函数返回前自动执行，无论是正常返回还是发生 panic：

	// 验证用户之前是否登入过
	if c := srv.userToConn[uid]; c != nil {
		//关闭之前的连接
		c.Close()
	}
	srv.connToUser[conn] = uid
	srv.userToConn[uid] = conn
}

func (srv *Server) GetConnection(uid string) *Conn {
	srv.RWMutex.RLock()
	defer srv.RWMutex.RUnlock()
	return srv.userToConn[uid]
}

func (srv *Server) GetConnections(uids ...string) []*Conn {
	if len(uids) == 0 {
		return nil
	}

	srv.RWMutex.RLock()
	defer srv.RWMutex.RUnlock()

	res := make([]*Conn, 0, len(uids))
	for _, uid := range uids {
		res = append(res, srv.GetConnection(uid))
	}
	return res
}

func (srv *Server) GetUsers(conns ...*Conn) []string {
	srv.RWMutex.RLock()
	defer srv.RWMutex.RUnlock()

	var res []string
	if len(conns) == 0 {
		// get all
		res = make([]string, 0, len(srv.connToUser))
		for _, uid := range srv.connToUser {
			res = append(res, uid)
		}
	} else {
		// get a part of
		res = make([]string, 0, len(conns))
		for _, conn := range conns {
			res = append(res, srv.connToUser[conn])
		}
	}
	return res
}

func (srv *Server) CloseConnection(conn *Conn) {

	srv.RWMutex.Lock()
	defer srv.RWMutex.Unlock()

	uid := srv.connToUser[conn]
	if uid == "" {
		// 已经被关闭了
		return
	}
	delete(srv.connToUser, conn)
	delete(srv.userToConn, uid)
	err := conn.Close()
	if err != nil {
		fmt.Println("close connection error:", err)
		return
	}
}

func (srv *Server) SendByUserId(msg interface{}, sendIds ...string) error {
	if len(sendIds) == 0 {
		return nil
	}
	return srv.Send(msg, srv.GetConnections(sendIds...)...)
}

func (srv *Server) Send(msg interface{}, conns ...*Conn) error {
	if len(conns) == 0 {
		return nil
	}
	data, err := json.Marshal(msg) // 把 msg 转成 JSON。
	if err != nil {
		return err
	}

	for _, conn := range conns {
		/*
			遍历所有目标连接（conns），逐个调用 conn.WriteMessage 写入 WebSocket。
			 单发：传一个 rconn
			 群发：传多个连接
		*/
		if err = conn.WriteMessage(websocket.TextMessage, data); err != nil {
			return err
		}
	}
	return nil
}

// 根据连接对象执行任务处理
func (srv *Server) handleConnection(conn *Conn) {
	uids := srv.GetUsers(conn)
	conn.Uid = uids[0]

	// 处理任务
	go srv.handlerWrite(conn)

	if /*srv.opt.ack != NoAck*/ srv.needAck(nil) {
		go srv.readAck(conn)
	}

	for {
		// 获取请求消息
		_, msg, err := conn.ReadMessage()
		if err != nil {
			srv.Errorf("server read message error: %v", err)
			srv.CloseConnection(conn)
			return
		}

		// 解析消息
		var message Message
		if err = json.Unmarshal(msg, &message); err != nil {
			srv.Errorf("server unmarshal message error: %v", err)
			srv.CloseConnection(conn)
			return
		}

		// todo: 给客户端回复一个ack

		//根据消息处理执行
		if srv.needAck(&message) {
			srv.Infof("conn message read ack msg: %v", message)
			conn.appendMsgMq(&message)
		} else {
			conn.message <- &message
		}
	}
}

// 是否要进行ack验证 ? isAck
func (srv *Server) needAck(msg *Message) bool {
	if msg == nil {
		return srv.opt.ack != NoAck
	}
	return srv.opt.ack != NoAck && msg.FrameType != FrameNoAck
}

// 读取消息的ack
func (srv *Server) readAck(conn *Conn) {
	for {
		select {
		case <-conn.done:
			srv.Infof("colse message ack. uid: %v", conn.Uid)
			return
		default:

		}

		// 从队列中读取新的消息
		conn.messageMu.Lock()
		if len(conn.readMessage) == 0 {
			conn.messageMu.Unlock()
			// sleep
			time.Sleep(time.Second * 3)
			continue
		}

		// 读取第一条
		message := conn.readMessage[0]

		// 判断ack方式
		switch srv.opt.ack {
		case OnlyAck:
			// 直接给客户端回复
			err := srv.Send(&Message{
				Id:        message.Id,
				FrameType: FrameAck,
				AckSeq:    message.AckSeq + 1,
				AckTime:   time.Time{},
				ErrCount:  0,
				Method:    "",
				FormId:    "",
				Data:      nil,
			}, conn)
			if err != nil {
				fmt.Println("send ack error on OnlyAck mod:", err)
				return
			}
			// 进行业务处理，把消息从队列中移除
			conn.readMessage = conn.readMessage[1:]
			conn.messageMu.Unlock()
			conn.message <- message
		case RigorAck:
			// 先回复
			if message.AckSeq == 0 {
				conn.readMessage[0].AckSeq++
				conn.readMessage[0].AckTime = time.Now()
				err := srv.Send(&Message{
					FrameType: FrameAck,
					Id:        message.Id,
					AckSeq:    message.AckSeq + 1,
				})
				if err != nil {
					fmt.Println("send ack error on RigorAck mod:", err)
					return
				}
				srv.Infof("message ack RigorAck mod send mid(message id): %v, seq: %v, time: %v", message.Id, message.AckSeq, message.AckTime.Unix())
				conn.messageMu.Unlock()
				continue
			}

			// 再验证

			// 1. 客户端返回结果，再一次确认
			msgSeq := conn.readMessageSeq[message.Id]
			if msgSeq.AckSeq > message.AckSeq {
				// 确认
				conn.readMessage = conn.readMessage[1:]
				conn.messageMu.Unlock()
				conn.message <- message
				srv.Infof("message ack RigorAck mod success. mid: %v, seq: %v, time: %v", message.Id, message.AckSeq, message.AckTime)
				continue
			}
			// 2. 客户端没有确认，是否超过了确认时间
			val := srv.opt.ackTimeout - time.Since(message.AckTime)
			if !message.AckTime.IsZero() && val <= 0 {
				//    2.1 超过，结束确认
				delete(conn.readMessageSeq, message.Id)
				conn.readMessage = conn.readMessage[1:]
				conn.messageMu.Unlock()
				continue
			}
			//    2.2 不超过 重新发送
			conn.messageMu.Unlock()
			srv.Send(&Message{
				FrameType: FrameAck,
				Id:        message.Id,
				AckSeq:    message.AckSeq,
			}, conn)
			// 睡眠一定时间
			time.Sleep(time.Second * 3)
		default:
			panic("unhandled default case")
		}
	}
}

// 任务的处理
func (srv *Server) handlerWrite(conn *Conn) {
	for {
		select {
		case <-conn.done: // 当前连接已关闭
			return
		case message := <-conn.message: // 存在最新的消息
			switch message.FrameType {
			case FramePing:
				srv.Send(&Message{FrameType: FramePing}, conn)
			case FrameData:
				// 根据请求的Method分发路由并执行
				if handler, ok := srv.routes[message.Method]; ok {
					handler(srv, conn, message)
				} else {
					srv.Send(&Message{FrameType: FrameData, Data: fmt.Sprintf("不存在执行的方法[im/ws/websocket/server.go] %v", message.Method)}, conn)
				}
			}

			// 避免发送重复的数据
			if srv.needAck(message) {
				conn.messageMu.Lock()
				delete(conn.readMessageSeq, message.Id)
				conn.messageMu.Unlock()
			}
		}
	}
}

func (srv *Server) AddRoutes(rs []Route) {
	for _, r := range rs {
		srv.routes[r.Method] = r.Handler
	}
}

func (srv *Server) Start() {
	http.HandleFunc(srv.patten, srv.ServerWs)
	srv.Info(http.ListenAndServe(srv.addr, nil))
}

func (srv *Server) Stop() {
	fmt.Println("server stop")

}
