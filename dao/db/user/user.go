package user

import (
	"errors"
	"fmt"

	"Monitoring-Pressure/dao/db"
	"Monitoring-Pressure/models/users"
	"Monitoring-Pressure/util"

	"gorm.io/gorm"
)

const (
	PasswordSalt = "i87czXc8IjEZRBRBSqpYudI1xEGHYWjx"
)

func Register(user *users.UserInfo) (err error) {

	// 唯一字段是 telephone，按 telephone 查重
	var count int64
	err = db.DB.Model(&users.User{}).
		Where("telephone = ?", user.Telephone).
		Count(&count).Error
	if err != nil {
		return
	}

	if count > 0 {
		err = ErrUserExists
		return
	}

	passwd := user.Password + PasswordSalt
	dbPassword := util.Md5([]byte(passwd))

	newUser := users.User{
		UserID:    user.UserID,
		Username:  user.Username,
		Password:  dbPassword,
		Telephone: user.Telephone,
		Role:      user.Role,
	}
	result := db.DB.Create(&newUser)
	if result.Error != nil {
		err = fmt.Errorf("insert user error:%v", result.Error)
		return
	}

	return
}

func Login(user *users.UserInfo) (err error) {

	originPassword := user.Password

	// 不要用 Select() 限定列，否则 Role/Username/Email/Sex 全是零值
	err = db.DB.Where("telephone = ?", user.Telephone).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = ErrUserNotExists
		return
	}
	if err != nil {
		return
	}

	passwd := originPassword + PasswordSalt
	originPasswordSalt := util.Md5([]byte(passwd))

	if originPasswordSalt != user.Password {
		err = ErrUserPasswordWrong
		return
	}

	return
}

func ResetPassword(user *users.UserInfo) (err error) {

	passwd := user.Password + PasswordSalt
	dbPassword := util.Md5([]byte(passwd))

	result := db.DB.Model(&users.User{}).
		Where("telephone = ?", user.Telephone).
		Update("password", dbPassword)
	if result.Error != nil {
		err = fmt.Errorf("update user error:%v", result.Error)
		return
	}

	if result.RowsAffected == 0 {
		err = ErrUserNotExists
		return
	}

	return
}
