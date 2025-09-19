package push

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"penguin/apps/im/ws/internal/svc"
	"penguin/apps/im/ws/websocket"
	"penguin/apps/im/ws/ws"
	"penguin/pkg/constants"
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
		switch data.ChatType {
		case constants.SingleChatType:
			err := single(srv, &data, data.RecvId)
			if err != nil {
				fmt.Println("err: ", err)
			}
		case constants.GroupChatType:
			err := group(srv, &data)
			if err != nil {
				fmt.Println("err: ", err)
			}
		}
	}
}

func single(srv *websocket.Server, data *ws.Push, recvId string) error {
	rconn := srv.GetConnection(recvId)
	if rconn == nil {
		// todo: 目标离线
		srv.Infof("目标离线")
		return nil
	}

	return srv.Send(websocket.NewMessage(data.SendId, &ws.Chat{
		ConversationId: data.ConversationId,
		ChatType:       data.ChatType,
		SendTime:       data.SendTime,
		Msg: ws.Msg{
			MType:   data.MType,
			Content: data.Content,
		},
	}), rconn)
}

func group(srv *websocket.Server, data *ws.Push) error {
	for _, id := range data.RecvIds {
		func(id string) {
			srv.Schedule(func() {
				single(srv, data, id)
			})
		}(id)
	}
	return nil
}
