package users

type UserInfo struct {
	User
	Token string `json:"token"` // 验证码
}

type User struct {
	UserID    int64  `json:"user_id"`
	Username  string `json:"username"`  // 用户名
	Password  string `json:"password"`  // 密码
	Telephone string `json:"telephone"` // 电话
	Sex       int    `json:"sex"`       // 性别
	Email     string `json:"email"`     // 邮箱
}

func (User) TableName() string {
	return "user"
}
