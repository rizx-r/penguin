package group

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"penguin/apps/social/social_models"
	"penguin/pkg/constants"
	"penguin/pkg/xerr"

	"penguin/apps/social/rpc/internal/svc"
	"penguin/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupPutInHandleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupPutInHandleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupPutInHandleLogic {
	return &GroupPutInHandleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GroupPutInHandle 加群申请处理
func (l *GroupPutInHandleLogic) GroupPutInHandle(in *social.GroupPutInHandleReq) (*social.GroupPutInHandleResp, error) {
	groupReq, err := l.svcCtx.GroupRequestsModel.FindOne(l.ctx, uint64(in.GroupReqId))
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "finding group request err:%v,  id %d", err, in.GroupReqId)
	}

	// 如果已有处理过
	switch constants.HandlerResult(groupReq.HandleResult.Int64) {
	case constants.PassHandleResult:
		return nil, errors.WithStack(xerr.ErrGroupRequestAlreadyPassed)
	case constants.RefuseHandleResult:
		return nil, errors.WithStack(xerr.ErrGroupRequestAlreadyRefused)
	default:
	}

	groupReq.HandleResult = sql.NullInt64{
		Int64: int64(in.HandleResult),
		Valid: true,
	}

	err = l.svcCtx.GroupRequestsModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		if err := l.svcCtx.GroupRequestsModel.TransUpdate(l.ctx, session, groupReq); err != nil {
			return errors.Wrapf(xerr.NewDBErr(), "updating group request err:%v,  id %d", err, in.GroupReqId)
		}
		if constants.HandlerResult(groupReq.HandleResult.Int64) != constants.PassHandleResult {
			return nil
		}

		groupMember := &social_models.GroupMembers{
			GroupId:     groupReq.GroupId,
			UserId:      groupReq.ReqId,
			RoleLevel:   int(constants.OrdinaryGroupRoleLevel),
			OperatorUid: in.HandleUid,
		}
		_, err = l.svcCtx.GroupMembersModel.TransInsert(l.ctx, session, groupMember)
		if err != nil {
			return errors.Wrapf(xerr.NewDBErr(), "inserting group members err:%v, req: %v", err, groupMember)
		}
		return nil
	})

	return &social.GroupPutInHandleResp{}, nil
}
