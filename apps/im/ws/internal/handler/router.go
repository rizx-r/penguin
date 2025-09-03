package handler

import (
	"penguin/apps/im/ws/internal/handler/conversion"
	"penguin/apps/im/ws/internal/handler/user"
	"penguin/apps/im/ws/internal/push"
	"penguin/apps/im/ws/internal/svc"
	"penguin/apps/im/ws/websocket"
)

func RegisterHandlers(srv *websocket.Server, svc *svc.ServiceContext) {
	srv.AddRoutes([]websocket.Route{
		{
			Method:  "user.online",
			Handler: user.Online(svc),
		},
		{
			Method:  "conversation.chat",
			Handler: conversion.Chat(svc),
		},
		{
			Method:  "push",
			Handler: push.Push(svc),
		},
	})
}
