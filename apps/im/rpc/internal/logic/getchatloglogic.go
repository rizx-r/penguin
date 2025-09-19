package logic

import (
	"context"
	"github.com/pkg/errors"
	//"go/types"
	"penguin/pkg/xerr"

	"penguin/apps/im/rpc/im"
	"penguin/apps/im/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetChatLogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetChatLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetChatLogLogic {
	return &GetChatLogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取会话记录
func (l *GetChatLogLogic) GetChatLog(in *im.GetChatLogReq) (*im.GetChatLogResp, error) {
	// todo: add your logic here and delete this line

	// 根据id

	if in.MsgId != "" {
		chatLog, err := l.svcCtx.ChatLogModel.FindOne(l.ctx, in.MsgId)
		if err != nil {
			return nil, errors.Wrapf(xerr.NewDBErr(), "find chat log by msgId error: %v, msgId: %s", err, in.MsgId)
		}
		return &im.GetChatLogResp{List: []*im.ChatLog{
			{
				Id:             chatLog.ID.Hex(), // 转成string
				ConversationId: chatLog.ConversationID,
				SendId:         chatLog.SendID,
				RecvId:         chatLog.RecvID,
				MsgType:        int32(chatLog.MsgType),
				MsgContent:     chatLog.MsgContent,
				ChatType:       int32(chatLog.ChatType),
				SendTime:       chatLog.SendTime,
				ReadRecords:    chatLog.ReadRecords,
			}}}, nil
	}

	// 时间分段查询
	data, err := l.svcCtx.ChatLogModel.ListBySendTime(l.ctx, in.ConversationId, in.StartSendTime, in.EndSendTime, in.Count)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "list chat log by sendTime error: %v, req: %v", err, in)
	}
	res := make([]*im.ChatLog, 0, len(data))
	for _, d := range data {
		res = append(res, &im.ChatLog{
			Id:             d.ID.Hex(),
			ConversationId: d.ConversationID,
			SendId:         d.SendID,
			RecvId:         d.RecvID,
			MsgType:        int32(d.MsgType),
			MsgContent:     d.MsgContent,
			ChatType:       int32(d.ChatType),
			SendTime:       d.SendTime,
			ReadRecords:    d.ReadRecords,
		})
	}
	return &im.GetChatLogResp{List: res}, nil
}
