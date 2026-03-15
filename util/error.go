package util

const (
	ErrCodeSuccess    = 0
	ErrCodeParameter  = 1001
	ErrCodeServerBusy = 1002

	ErrCodeUserExist         = 1003
	ErrCodeUserNotExist      = 1004
	ErrCodeUserPasswordWrong = 1005

	ErrCodeInvalidVerifyCode = 1006

	ErrCodeNeedLogin    = 1007
	ErrCodeNoPermission = 1008
	ErrCodeTokenInvalid = 1009
	ErrCodeRoleInvalid  = 1010
)

var codeMsgMap = map[int]string{
	ErrCodeSuccess:           "success",
	ErrCodeParameter:         "参数错误",
	ErrCodeServerBusy:        "服务器繁忙",
	ErrCodeUserExist:         "用户已存在",
	ErrCodeUserNotExist:      "用户不存在",
	ErrCodeUserPasswordWrong: "手机号或密码错误",
	ErrCodeInvalidVerifyCode: "验证码错误",
	ErrCodeNeedLogin:         "请先登录",
	ErrCodeNoPermission:      "权限不足",
	ErrCodeTokenInvalid:      "token无效",
	ErrCodeRoleInvalid:       "角色不合法",
}

func GetMessage(code int) string {
	msg, ok := codeMsgMap[code]
	if ok {
		return msg
	}
	return "未知错误"
}
