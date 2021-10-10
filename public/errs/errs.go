package errs

import (
	"errors"
)

var (
	CountEmpty = errors.New("用户 ID 不能空")
	PwdEmpty   = errors.New("用户密码不能空")

	INVAILD_MSG_INSUFFICIENT_LENGTH = errors.New("无效的信息长度")
	INVAILD_MSG_PARSE_FAIL          = errors.New("信息解析失败")
	COUNT_NOT_EXIST                 = errors.New("用户不存在")
	COUNT_WRONG_PASSWORD            = errors.New("用户密码错误")
	SERVER_ERROR                    = errors.New("服务器错误")

	COUNT_EXIST           = errors.New("用户已存在")
	COUNT_ID_LENGTH_WRONG = errors.New("用户 ID 长度错误")

	SERVER_TIMEOUT = errors.New("服务器超时")
)
