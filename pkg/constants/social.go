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
	NoHandleResult     HandlerResult = iota + 1 // 未处理
	PassHandleResult                            // 已通过
	RefuseHandleResult                          // 已拒绝
	CancelHandleResult                          // 已取消
)

// GroupRoleLevel
const (
	MasterGroupRoleLevel   GroupRoleLevel = iota + 1 // 群主
	ManagerGroupRoleLevel                            // 管理员
	OrdinaryGroupRoleLevel                           // 普通群成员
)

// GroupJoinSource
const (
	InviteGroupJoinSource GroupJoinSource = iota + 1 // 被邀请
	PutInGroupJoinSource                             // 自己申请
)
