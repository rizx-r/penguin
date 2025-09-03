package xerr

import "github.com/zeromicro/x/errors"

func New(code int, msg string) error {
	return errors.New(code, msg)
}

func NewDBErr() error {
	return errors.New(DB_ERROR, "DB Error")
}

func NewInternalErr() error {
	return errors.New(SERVER_COMMON_ERROR, "Server Internal Error")
}

func NewMsg(msg string) error {
	return errors.New(SERVER_COMMON_ERROR, msg)
}
