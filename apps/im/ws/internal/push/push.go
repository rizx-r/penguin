package push

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"penguin/apps/im/ws/internal/svc"
	"penguin/apps/im/ws/websocket"
	"penguin/apps/im/ws/ws"
)

func Push(svc *svc.ServiceContext) websocket.HandlerFunc {
	return func(srv *websocket.Server, conn *websocket.Conn, msg *websocket.Message) {
		var data ws.Push
		if err := mapstructure.Decode(msg.Data, &data); err != nil {
			err := srv.Send(websocket.NewErrMessage(err))
			if err != nil {
				fmt.Println("Send Error Message Error: at fun-im\\fun-chat\\apps\\im\\ws\\internal\\push\\push.go: ", err)
				return
			}
			return
		}
		// 发送的目标
		rconn := srv.GetConnection(data.RecvId)
		if rconn == nil {
			// todo: 目标离线
			return
		}

		srv.Infof("push msg: %v", data)
		srv.Send(websocket.NewMessage(data.SendId, &ws.Chat{
			ConversationId: data.ConversationId,
			ChatType:       data.ChatType,
			SendTime:       data.SendTime,
			Msg: ws.Msg{
				MType:   data.MType,
				Content: data.Content,
			},
		}), rconn)
	}
}
