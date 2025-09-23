package msgTransfer

import (
	"github.com/zeromicro/go-zero/core/logx"
	"penguin/apps/im/ws/ws"
	"penguin/pkg/constants"
	"sync"
	"time"
)

type (
	GroupMsgRead struct {
		mu             sync.Mutex
		push           *ws.Push
		ConversationId string
		pushChan       chan *ws.Push
		count          int
		pushTime       time.Time
		done           chan struct{}
	}
)

func NewGroupMsgRead(push *ws.Push, pushCh chan *ws.Push) *GroupMsgRead {
	m := &GroupMsgRead{
		push:           push,
		ConversationId: push.ConversationId,
		pushChan:       pushCh,
		count:          0,
		pushTime:       time.Now(),
		done:           make(chan struct{}),
	}
	go m.transfer()
	return m
}

// mergePush 合并消息推送
func (g *GroupMsgRead) mergePush(push *ws.Push) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.count++
	for msgId, read := range push.ReadRecords {
		g.push.ReadRecords[msgId] = read
	}
}

func (g *GroupMsgRead) transfer() {
	// 1. 超时发送
	// 2. 超量发送

	timer := time.NewTimer(constants.GroupMsgReadHandlerDelayTransfer / 2)
	defer timer.Stop()

	for {
		select {
		case <-g.done:
			return
		case <-timer.C: // 定时器被触发的时候
			g.mu.Lock()
			pushTime := g.pushTime
			val := GroupMsgReadRecordDelayTime - time.Since(pushTime)
			push := g.push
			if val > 0 && int64(g.count) < GroupMsgReadRecordDelayCount || push == nil {
				// 未达标
				if val > 0 {
					timer.Reset(val)
				}
				g.mu.Unlock()
				continue
			}
			g.pushTime = time.Now()
			g.push = nil
			g.count = 0
			timer.Reset(GroupMsgReadRecordDelayTime / 2)
			g.mu.Unlock()
		default:
			g.mu.Lock()
			if int64(g.count) >= GroupMsgReadRecordDelayCount {
				push := g.push
				g.push = nil
				g.count = 0
				g.mu.Unlock()

				// 推送
				logx.Infof("default 超过 合并的条件推送 %v", push)
				g.pushChan <- push
				continue
			}
			if g.isIdle() {
				g.mu.Unlock()
				// 使 msgReadTransfer 释放
				g.pushChan <- &ws.Push{
					ChatType:       constants.GroupChatType,
					ConversationId: g.ConversationId,
				}
				continue
			}
			g.mu.Unlock()
			tempDelay := GroupMsgReadRecordDelayTime / 4
			if tempDelay > time.Second {
				tempDelay = time.Second
			}
			time.Sleep(tempDelay)
		}
	}
}

func (g *GroupMsgRead) IsIdle() bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.isIdle()
}

// isIdle 是否为活跃状态
func (g *GroupMsgRead) isIdle() bool {
	pushTime := g.pushTime
	val := GroupMsgReadRecordDelayTime*2 - time.Since(pushTime)
	return val <= 0 && g.push == nil && g.count == 0
}

func (g *GroupMsgRead) clear() {
	select {
	case <-g.done:
	default:
		close(g.done)
	}
	g.push = nil
}
