package main

import (
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"log"
	"penguin/apps/task/mq/internal/config"
	"penguin/apps/task/mq/internal/handler"
	"penguin/apps/task/mq/internal/svc"
)

var (
	configFile = flag.String("f", "etc/dev/task.yaml", "path to config file")
)

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	if err := c.SetUp(); err != nil {
		log.Fatal(err)
	}

	ctx := svc.NewServiceContext(c)
	listen := handler.NewListen(ctx)

	serviceGroup := service.NewServiceGroup()
	for _, s := range listen.Services() {
		serviceGroup.Add(s)
	}

	fmt.Println("Task Service starting ...")
	serviceGroup.Start()
}

/*
# 测试发送消息
docker exec -it 4c039136c177 kafka-console-producer.sh --broker-list 127.0.0.1:9092 --topic msgChatTransfer
hello world
*/
