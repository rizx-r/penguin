package websocket

type (
	Route struct {
		Method  string `json:"method"`
		Handler HandlerFunc
	}

	HandlerFunc func(srv *Server, conn *Conn, msg *Message)
)
