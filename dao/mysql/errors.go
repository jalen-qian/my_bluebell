package mysql

import "github.com/pkg/errors"

var (
	ErrUserExist       = errors.New("用户已存在")
	ErrUserNotExist    = errors.New("用户不存在")
	ErrServerBusy      = errors.New("服务繁忙")
	ErrInvalidPassword = errors.New("密码错误")
)
