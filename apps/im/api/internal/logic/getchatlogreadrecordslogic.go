package logic

import (
	"context"
	"penguin/apps/im/rpc/im"
	"penguin/apps/social/rpc/socialclient"
	"penguin/pkg/bitmap"
	"penguin/pkg/constants"

	"penguin/apps/im/api/internal/svc"
	"penguin/apps/im/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetChatLogReadRecordsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetChatLogReadRecordsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetChatLogReadRecordsLogic {
	return &GetChatLogReadRecordsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetChatLogReadRecordsLogic) GetChatLogReadRecords(req *types.GetChatLogReadRecordsReq) (resp *types.GetChatLogReadRecordsResp, err error) {
	// todo: add your logic here and delete this line
	chatLogs, err := l.svcCtx.Im.GetChatLog(l.ctx, &im.GetChatLogReq{
		MsgId: req.MsgId,
	})
	if err != nil || len(chatLogs.List) == 0 {
		return nil, err
	}

	var (
		chatlog = chatLogs.List[0]
		reads   = []string{chatlog.SendId}
		unreads []string
		ids     []string
	)

	// 分别设置 已读/未读
	switch constants.ChatType(chatlog.ChatType) {
	case constants.SingleChatType:
		if len(chatlog.ReadRecords) == 0 || chatlog.ReadRecords[0] == 0 {
			unreads = []string{chatlog.RecvId}
		} else {
			reads = append(reads, chatlog.SendId)
		}
		ids = []string{chatlog.RecvId, chatlog.SendId}
	case constants.GroupChatType:
		groupUsers, err := l.svcCtx.Social.GroupUsers(l.ctx, &socialclient.GroupUsersReq{
			GroupId: chatlog.RecvId,
		})
		if err != nil {
			return nil, err
		}
		bitmaps := bitmap.LoadBitmap(chatlog.ReadRecords)
		for _, members := range groupUsers.List {
			ids = append(ids, members.UserId)

			if members.UserId == chatlog.SendId {
				// 自己的消息 则已读
				continue
			}

			if bitmaps.IsSet(members.UserId) {
				reads = append(reads, members.UserId)
			} else {
				unreads = append(unreads, members.UserId)
			}
		}
	}
	return
}
