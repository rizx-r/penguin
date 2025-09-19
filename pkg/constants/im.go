package constants

type (
	MType       int
	ChatType    int
	ContentType int // 内容类型
)

const (
	TextMtype              = iota // iota 表示 从 0 开始，每行递增 1 的常量生成器。
	GroupChatType ChatType = iota
	SingleChatType
)

// ContentType
const (
	// ContentChatMsg 消息聊天
	ContentChatMsg ContentType = iota
	// ContentMarkRead 记录已读
	ContentMarkRead
)
