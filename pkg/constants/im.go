package constants

type (
	MType    int
	ChatType int
)

const (
	TextMtype              = iota // iota 表示 从 0 开始，每行递增 1 的常量生成器。
	GroupChatType ChatType = iota
	SingleChatType
)
