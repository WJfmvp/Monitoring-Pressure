package util

const (
	ErrCodeSuccess           = 0
	ErrCodeParameter         = 1001
	ErrCodeUserExist         = 1002
	ErrCodeServerBusy        = 1003
	ErrCodeUserNotExist      = 1004
	ErrCodeUserPasswordWrong = 1005
	ErrCodeInvalidVerifyCode = 1006
)

func GetMessage(code int) (message string) {
	switch code {
	case ErrCodeSuccess:
		message = "success"
	case ErrCodeParameter:
		message = "参数错误"
	case ErrCodeUserExist:
		message = "用户存在"
	case ErrCodeServerBusy:
		message = "服务器繁忙"
	case ErrCodeUserNotExist:
		message = "用户不存在"
	case ErrCodeUserPasswordWrong:
		message = "用户或密码错误"
	case ErrCodeInvalidVerifyCode:
		message = "验证码不正确"
	default:
		message = "未知错误"
	}
	return
}
