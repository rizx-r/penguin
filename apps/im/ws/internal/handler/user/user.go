package user

import (
	"penguin/apps/im/ws/internal/svc"
	websocketx "penguin/apps/im/ws/websocket"
)

func Online(svc *svc.ServiceContext) websocketx.HandlerFunc {
	return func(srv *websocketx.Server, conn *websocketx.Conn, msg *websocketx.Message) {
		uids := srv.GetUsers()
		u := srv.GetUsers(conn)
		err := srv.Send(websocketx.NewMessage(u[0], uids), conn)
		srv.Info("Error on sending online users: %v", err)
	}
}
