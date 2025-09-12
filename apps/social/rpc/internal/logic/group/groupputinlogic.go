package group

// 加群申请业务逻辑

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
	"penguin/apps/social/social_models"
	"penguin/pkg/constants"
	"penguin/pkg/xerr"
	"time"

	"penguin/apps/social/rpc/internal/svc"
	"penguin/apps/social/rpc/social"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupPutinLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGroupPutinLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupPutinLogic {
	return &GroupPutinLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GroupPutin 加群申请处理
func (l *GroupPutinLogic) GroupPutin(in *social.GroupPutinReq) (*social.GroupPutinResp, error) {
	/*
		1. 普通用户申请：
			群无验证：直接入群
			群有验证：发送申请
		2. 普通群成员邀请：
			群无验证：直接入群
			群有验证：发送申请
		3. 群主/管理员 邀请：直接入群
	*/

	var (
		inviteGroupMember *social_models.GroupMembers
		userGroupMember   *social_models.GroupMembers
		groupInfo         *social_models.Groups
		err               error
	)

	userGroupMember, err = l.svcCtx.GroupMembersModel.FindByGroudIdAndUserId(l.ctx, in.ReqId, in.GroupId)
	if err != nil && !errors.Is(err, social_models.ErrNotFound) {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find group members by groups ids err: %v, req: %v", err, in)
	}
	if userGroupMember != nil {
		// 如果已经在群里了
		return &social.GroupPutinResp{}, nil
	}

	groupReq, err := l.svcCtx.GroupRequestsModel.FindByGroupIdAndReqId(l.ctx, in.GroupId, in.ReqId)
	if err != nil && !errors.Is(err, social_models.ErrNotFound) {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find group requests by groups ids err: %v, req: %v", err, in)
	}
	if groupReq != nil {
		// 如果已经有申请了
		return &social.GroupPutinResp{}, nil
	}

	// 构建申请入群的请求
	groupReq = &social_models.GroupRequests{
		ReqId:   in.ReqId,
		GroupId: in.GroupId,
		ReqMsg: sql.NullString{
			String: in.ReqMsg,
			Valid:  true,
		},
		ReqTime: sql.NullTime{
			Time:  time.Unix(in.ReqTime, 0),
			Valid: true,
		},
		JoinSource: sql.NullInt64{
			Int64: int64(in.JoinSource),
			Valid: true,
		},
		InviterUserId: sql.NullString{
			String: in.InviterUid,
			Valid:  true,
		},
		HandleResult: sql.NullInt64{
			Int64: int64(constants.NoHandleResult),
			Valid: true,
		},
	}

	// 群信息
	groupInfo, err = l.svcCtx.GroupsModel.FindOne(l.ctx, in.GroupId)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find group by groups ids err: %v, req: %v", err, in)
	}

	addGroupMember := func() {
		if err != nil {
			return
		}
		err = l.addGroupMember(in)
	}

	// 1. 如果这个群进群不需要验证，则直接入群
	if !groupInfo.IsVerify {
		defer addGroupMember()
		groupReq.HandleResult = sql.NullInt64{
			Int64: int64(constants.PassHandleResult),
			Valid: true,
		}

		return l.createGroupReq(groupReq, true)
	}

	// 2. 如果这个群进群需要验验证
	// 2.1 如果是申请者自己申请入群
	if constants.GroupJoinSource(in.JoinSource) == constants.PutInGroupJoinSource {
		return l.createGroupReq(groupReq, false)
	}

	// 2.2  如果是被邀请入群的
	// 邀请者
	inviteGroupMember, err = l.svcCtx.GroupMembersModel.FindByGroudIdAndUserId(l.ctx, in.InviterUid, in.GroupId)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find group members by groups ids err: %v, req: %v", err, in)
	}

	// 如果是管理员或群主邀请进群
	if constants.GroupRoleLevel(inviteGroupMember.RoleLevel) == constants.MasterGroupRoleLevel || constants.GroupRoleLevel(inviteGroupMember.RoleLevel) == constants.ManagerGroupRoleLevel {
		defer addGroupMember()
		// 填写请求结果和处理者id
		groupReq.HandleResult = sql.NullInt64{
			Int64: int64(constants.PassHandleResult),
			Valid: true,
		}
		groupReq.HandleUserId = sql.NullString{
			String: in.InviterUid,
			Valid:  true,
		}
		return l.createGroupReq(groupReq, true)
	}
	// 如果是普通群成员邀请入群
	return l.createGroupReq(groupReq, false)
}

// addGroupMember 添加群成员
func (l *GroupPutinLogic) addGroupMember(in *social.GroupPutinReq) error {
	groupMember := &social_models.GroupMembers{
		GroupId:     in.GroupId,
		UserId:      in.ReqId,
		RoleLevel:   int(constants.OrdinaryGroupRoleLevel),
		OperatorUid: in.InviterUid,
	}
	//_, err := l.svcCtx.GroupMembersModel.TransInsert(l.ctx, nil, groupMember)
	_, err := l.svcCtx.GroupMembersModel.Insert(l.ctx, groupMember)
	if err != nil {
		return errors.Wrapf(xerr.NewDBErr(), "insert group members err: %v, req: %v", err, in)
	}
	return nil
}

// createGroupReq 添加群申请信息到数据库
// isPass 是否通过
func (l *GroupPutinLogic) createGroupReq(groupReq *social_models.GroupRequests, isPass bool) (*social.GroupPutinResp, error) {
	_, err := l.svcCtx.GroupRequestsModel.Insert(l.ctx, groupReq)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "insert group requests err: %v, req: %v", err, groupReq)
	}
	if isPass {
		return &social.GroupPutinResp{GroupId: groupReq.GroupId}, nil
	}
	return &social.GroupPutinResp{}, nil
}
