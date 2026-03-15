package users

type UserRole string

const (
	RoleStudent UserRole = "student"
	RoleTeacher UserRole = "teacher"
	RoleAdmin   UserRole = "admin"
)

type User struct {
	UserID    int64    `json:"user_id" gorm:"primaryKey;autoIncrement"`
	Username  string   `json:"username" gorm:"type:varchar(50);not null"`
	Password  string   `json:"password" gorm:"type:varchar(255)"`
	Telephone string   `json:"telephone" gorm:"type:varchar(20);uniqueIndex;not null"`
	Sex       int      `json:"sex" gorm:"not null;default:0"`
	Email     string   `json:"email" gorm:"type:varchar(100)"`
	Role      UserRole `json:"role" gorm:"type:varchar(20);not null;default:'student'"`
}

type UserInfo struct {
	User
	Token string `json:"token" gorm:"-"`
}

func (User) TableName() string {
	return "user"
}
