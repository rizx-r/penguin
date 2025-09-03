package websocket

import (
	"fmt"
	"net/http"
	"time"
)

type (
	Authentication interface {
		Authenticate(w http.ResponseWriter, r *http.Request) (ok bool)
		UserId(r *http.Request) (uid string)
	}
	authentication struct{}
)

func (a *authentication) Authenticate(w http.ResponseWriter, r *http.Request) (ok bool) {
	return true
}

func (a *authentication) UserId(r *http.Request) (uid string) {
	query := r.URL.Query()
	/*	if query != nil && query["userId"] != nil { // ???
		return fmt.Sprintf("%v", query["userId"])
	}*/
	if uid := query.Get("userId"); uid != "" {
		return uid
	}
	return fmt.Sprintf("%v", time.Now().UnixMilli())
}
