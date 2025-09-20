package msgTransfer

import (
	"penguin/apps/im/ws/ws"
	"sync"
	"time"
)

type (
	GroupMsgRead struct {
		mu       sync.Mutex
		push     *ws.Push
		pushChan chan *ws.Push
		count    int
		pushTime time.Time
		done     chan struct{}
	}
)

func NewGroupMsgRead(push *ws.Push, pushCh chan *ws.Push) *GroupMsgRead {
	return &GroupMsgRead{
		push:     push,
		pushChan: pushCh,
		count:    0,
		pushTime: time.Now(),
		done:     make(chan struct{}),
	}
}
