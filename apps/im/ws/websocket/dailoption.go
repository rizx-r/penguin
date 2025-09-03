package websocket

import "net/http"

type (
	DailOption struct {
		pattern string
		header  http.Header
	}
	DailOptions func(option *DailOption)
)

func NewDailOptions(opts ...DailOptions) DailOption {
	o := DailOption{
		pattern: "ws",
		header:  nil,
	}
	for _, opt := range opts {
		opt(&o)
	}
	return o
}

func WithClientPatten(pattern string) DailOptions {
	return func(o *DailOption) {
		o.pattern = pattern
	}
}

func WithClientHeader(header http.Header) DailOptions {
	return func(o *DailOption) {
		o.header = header
	}
}
