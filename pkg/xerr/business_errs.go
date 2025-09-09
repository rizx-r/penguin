package xerr

import "errors"

var (
	// 用户注册
	ErrPhoneNotRegistered = errors.New("phone not registered")
	ErrPhoneRegistered    = errors.New("phone registered")
	ErrUserPwdNotMatched  = errors.New("user password not matched")

	// 好友申请
	ErrFriendRequestAlreadyPassed  = NewMsg("friend request already passed")
	ErrFriendRequestAlreadyRefused = NewMsg("friend request already refused")

	// 群申请
	ErrGroupRequestAlreadyExists  = NewMsg("group request already exists")
	ErrGroupRequestAlreadyPassed  = NewMsg("group request already passed")
	ErrGroupRequestAlreadyRefused = NewMsg("group request already refused")
)
