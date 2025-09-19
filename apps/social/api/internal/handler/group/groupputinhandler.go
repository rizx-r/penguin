package group

import (
	"fmt"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"penguin/apps/social/api/internal/logic/group"
	"penguin/apps/social/api/internal/svc"
	"penguin/apps/social/api/internal/types"
)

// 申请进群
func GroupPutInHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GroupPutInReq
		if err := httpx.Parse(r, &req); err != nil {
			fmt.Println("0xadwawdefaf", req)

			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := group.NewGroupPutInLogic(r.Context(), svcCtx)
		resp, err := l.GroupPutIn(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
