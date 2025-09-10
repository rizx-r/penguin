package constants

type (
	// 处理结果  1未处理，2.处理，3.拒绝
	HandlerResult int

	// 群成员等级 1. 群主 2. 管理员，3. 普通群成员
	GroupRoleLevel int

	// 进群申请的方式： 1. 邀请， 2. 申请
	GroupJoinSource int
)

// HandlerResult
const (
	NoHandleResult HandlerResult = iota + 1
	PassHandleResult
	RefuseHandleResult
	CancelHandleResult
)

const (
	MasterGroupRoleLevel GroupRoleLevel = iota + 1
	ManagerGroupRoleLevel
	OrdinaryGroupRoleLevel
)

const (
	InviteGroupJoinSource GroupJoinSource = iota + 1
	PutInGroupJoinSource
)
