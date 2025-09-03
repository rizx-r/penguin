package constants

// 处理结果  1未处理，2.处理，3.拒绝
type HandlerResult int

const (
	NoHandleResult HandlerResult = iota + 1
	PassHandleResult
	RefuseHandleResult
	CancelHandleResult
)
