package user

import (
	"Monitoring-Pressure/dao/db"
	"Monitoring-Pressure/models/users"
	"Monitoring-Pressure/util"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
)

const (
	PasswordSalt = "i87czXc8IjEZRBRBSqpYudI1xEGHYWjx"
)

func Register(user *users.UserInfo) (err error) {

	var count int64
	err = db.DB.Model(&users.UserInfo{}).Where("username = ?", user.Username).Count(&count).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
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
	err = db.DB.Select("telephone, password, user_id").Where("telephone = ?", user.Telephone).First(&user).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		fmt.Printf("error: %v\n", err)
		return
	} else {
		fmt.Printf("db:%p user:%#v\n", db.DB, user)
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = ErrUserNotExists
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

	result := db.DB.Model(&users.UserInfo{}).Where("telephone=?", user.Telephone).Update("password", dbPassword)
	if result.Error != nil {
		err = fmt.Errorf("update user error:%v", result.Error)
		return
	}

	return
}
