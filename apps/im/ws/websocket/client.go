package websocket

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"net/url"
)

type (
	Client interface {
		Close() error
		Send(v any) error
		Read(v any) error
	}
	client struct {
		*websocket.Conn
		host string
		opt  DailOption
	}
)

func NewClient(host string, opts ...DailOptions) *client {
	opt := NewDailOptions(opts...)
	c := client{
		Conn: nil,
		host: host,
		opt:  opt,
	}
	conn, err := c.dail()
	if err != nil {
		panic(err)
	}

	c.Conn = conn
	return &c
}

func (c *client) dail() (*websocket.Conn, error) {
	u := url.URL{Scheme: "ws", Host: c.host, Path: c.opt.pattern}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), c.opt.header)
	return conn, err
}

func (c *client) Send(v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	err = c.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		return err
	}

	// todo: 再增加一个重连发送
	conn, err := c.dail()
	if err != nil {
		return err
	}
	c.Conn = conn
	return c.WriteMessage(websocket.TextMessage, data)
}

func (c *client) Read(v any) error {
	_, msg, err := c.ReadMessage()
	if err != nil {
		return err
	}

	return json.Unmarshal(msg, v)
}

func (c *client) Close() error {
	return c.Conn.Close()
}
