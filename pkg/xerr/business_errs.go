package xerr

import "errors"

var (
	ErrPhoneNotRegistered          = errors.New("phone not registered")
	ErrPhoneRegistered             = errors.New("phone registered")
	ErrUserPwdNotMatched           = errors.New("user password not matched")
	ErrFriendRequestAlreadyPassed  = NewMsg("friend request already passed")
	ErrFriendRequestAlreadyRefused = NewMsg("friend request already refused")
)
