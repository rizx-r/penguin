package logic

import (
	"context"
	"penguin/apps/user/rpc/user"
	"testing"
)

func TestRegisterLogic_Register(t *testing.T) {

	type args struct {
		in *user.RegisterReq
	}
	tests := []struct {
		name      string
		args      args
		want      *user.RegisterResp
		wantErr   bool
		wantPrint bool
	}{
		// TODO: Add test cases.
		{
			name: "ok",
			args: args{
				in: &user.RegisterReq{
					Phone:    "13700001111",
					Nickname: "test-user",
					Password: "123456",
					Avatar:   "https://avatars.githubusercontent.com/u/123456",
					Sex:      0,
				},
			},
			want: &user.RegisterResp{
				Token:  "",
				Expire: 0,
			},
			wantErr:   false,
			wantPrint: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewRegisterLogic(context.Background(), svcCtx)
			got, err := l.Register(tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantPrint {
				t.Log(tt.name, got)
			}
		})
	}
}
