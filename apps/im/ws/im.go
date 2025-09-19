package main

import (
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
	"log"
	"penguin/apps/im/ws/internal/config"
	"penguin/apps/im/ws/internal/handler"
	"penguin/apps/im/ws/internal/svc"
	"penguin/apps/im/ws/websocket"
)

var configFile = flag.String("f", "etc/dev/im.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	if err := c.SetUp(); err != nil {
		log.Fatal(err)
	}
	ctx := svc.NewServiceContext(c)

	srv := websocket.NewServer(c.ListenOn,
		websocket.WithServerAuthentication(handler.NewJwtAuth(ctx)),
		//websocket.WithServerAck(websocket.NoAck),
		//websocket.WithServerAck(websocket.OnlyAck),
		//websocket.WithServerAck(websocket.RigorAck),
		/*websocket.WithServerMaxConnectionIdle(10*time.Second)*/
	)
	defer srv.Stop()

	handler.RegisterHandlers(srv, ctx)

	fmt.Println("start websocket server on", c.ListenOn)
	srv.Start()
}
