package resultx

import (
	"context"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	zerr "github.com/zeromicro/x/errors"
	"google.golang.org/grpc/status"
	"net/http"
	"penguin/pkg/xerr"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func Success(data interface{}) *Response {
	return &Response{
		Code: 200,
		Msg:  "success",
		Data: data,
	}
}

func Fail(code int, msg string) *Response {
	return &Response{
		Code: 200,
		Msg:  msg,
		Data: nil,
	}
}

func OkHandler(_ context.Context, data interface{}) any {
	return Success(data)
}

func ErrHandler(name string) func(ctx context.Context, err error) (int, any) {
	return func(ctx context.Context, err error) (int, any) {
		err_code := xerr.SERVER_COMMON_ERROR
		err_msg := xerr.ErrMsg(err_code)

		causeErr := errors.Cause(err)
		if e, ok := causeErr.(*zerr.CodeMsg); ok {
			err_code = e.Code
			err_msg = e.Msg
		} else {
			if gstatus, ok := status.FromError(causeErr); ok {
				err_code = int(gstatus.Code())
				err_msg = gstatus.Message()
			}
		}

		logx.WithContext(ctx).Errorf("[%s] error: %s", name, err_msg)

		return http.StatusBadRequest, Fail(err_code, err_msg)
	}
}
