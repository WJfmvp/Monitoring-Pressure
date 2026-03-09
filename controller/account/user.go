package account

import (
	"Monitoring-Pressure/dao/db"
	"Monitoring-Pressure/dao/db/user"
	"Monitoring-Pressure/id_gen"
	"Monitoring-Pressure/middleware/account"
	"Monitoring-Pressure/models/users"
	"Monitoring-Pressure/util"
	"Monitoring-Pressure/verifycode"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SendVerifyCodeHandle 验证码发送接口
func SendVerifyCodeHandle(c *gin.Context) {
	telephone := c.Query("telephone")
	if telephone == "" {
		util.ResponseError(c, util.ErrCodeParameter)
		return
	}

	_, err := verifycode.GenerateCode(telephone)
	if err != nil {
		util.ResponseError(c, util.ErrCodeServerBusy)
		return
	}

	util.ResponseSuccess(c, gin.H{
		"message": "验证码已发送",
		//"token":   token,
	})
}

// LoginHandle 登陆接口1 通过账号密码进行登陆
func LoginHandle(c *gin.Context) {

	account.ProcessRequest(c)

	var err error
	var userInfo users.UserInfo
	defer func() {
		if err != nil {
			return
		}
		// 用户登陆成功之后，而要把user_id设置到session中
		account.SetUserId(userInfo.UserID, c)
		// 当调用responseSuccess的时候，gin框架已经把数据发送给浏览器了
		// 所以在responseSuccess之后，SetCookie就不会生效
		account.ProcessResponse(c)
		util.ResponseSuccess(c, nil)
	}()
	err = c.BindJSON(&userInfo)
	if err != nil {
		util.ResponseError(c, util.ErrCodeParameter)
		return
	}

	if len(userInfo.Password) == 0 || len(userInfo.Telephone) == 0 {
		util.ResponseError(c, util.ErrCodeParameter)
		return
	}

	err = user.Login(&userInfo)
	if errors.Is(err, user.ErrUserNotExists) {
		util.ResponseError(c, util.ErrCodeUserNotExist)
		return
	}

	if errors.Is(err, user.ErrUserPasswordWrong) {
		util.ResponseError(c, util.ErrCodeUserPasswordWrong)
		return
	}
	if err != nil {
		util.ResponseError(c, util.ErrCodeServerBusy)
		return
	}

	// 用户登陆成功之后，而要把user_id设置到session中
	account.SetUserId(userInfo.UserID, c)
	// 当调用responseSuccess的时候，gin框架已经把数据发送给浏览器了
	// 所以在responseSuccess之后，SetCookie就不会生效
	account.ProcessResponse(c)
	util.ResponseSuccess(c, nil)
}

// Login2Handle 登陆接口2 通过手机验证码进行登陆
func Login2Handle(c *gin.Context) {

	var userInfo users.UserInfo
	err := c.BindJSON(&userInfo)
	if err != nil {
		util.ResponseError(c, util.ErrCodeParameter)
		return
	}

	if userInfo.Telephone == "" {
		util.ResponseError(c, util.ErrCodeParameter)
		return
	}

	err = db.DB.Select("user_id").Where("telephone = ?", userInfo.Telephone).First(&userInfo).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		util.ResponseError(c, util.ErrCodeUserNotExist)
		return
	}

	// 验证用户输入的验证码是否正确
	if !verifycode.VerifyCode(userInfo.Telephone, userInfo.Token) {
		util.ResponseError(c, util.ErrCodeInvalidVerifyCode)
		return
	}

	util.ResponseSuccess(c, nil)
}

// RegisterHandle 注册接口
func RegisterHandle(c *gin.Context) {

	var userInfo users.UserInfo
	err := c.BindJSON(&userInfo)
	if err != nil {
		util.ResponseError(c, util.ErrCodeParameter)
		return
	}

	if len(userInfo.Username) == 0 || len(userInfo.Password) == 0 || len(userInfo.Telephone) == 0 ||
		len(userInfo.Token) == 0 {
		util.ResponseError(c, util.ErrCodeParameter)
		return
	}

	// 验证用户输入的验证码是否正确
	if !verifycode.VerifyCode(userInfo.Telephone, userInfo.Token) {
		util.ResponseError(c, util.ErrCodeInvalidVerifyCode)
		return
	}

	userId, err := id_gen.GetId()
	if err != nil {
		util.ResponseError(c, util.ErrCodeServerBusy)
		return
	}
	userInfo.UserID = int64(userId)

	err = user.Register(&userInfo)
	if errors.Is(err, user.ErrUserExists) {
		util.ResponseError(c, util.ErrCodeUserExist)
		return
	}
	if err != nil {
		util.ResponseError(c, util.ErrCodeServerBusy)
		return
	}

	util.ResponseSuccess(c, nil)
}

// ResetPasswordHandle 修改密码接口
func ResetPasswordHandle(c *gin.Context) {
	var userInfo users.UserInfo
	err := c.BindJSON(&userInfo)
	if err != nil {
		util.ResponseError(c, util.ErrCodeParameter)
		return
	}

	if len(userInfo.Telephone) == 0 || len(userInfo.Password) == 0 ||
		len(userInfo.Token) == 0 {
		util.ResponseError(c, util.ErrCodeParameter)
		return
	}

	// 验证用户输入的验证码是否正确
	/*	if !verifycode.VerifyCode(userInfo.Telephone, userInfo.Token) {
		util.ResponseError(c, util.ErrCodeInvalidVerifyCode)
		return
	}*/

	err = user.ResetPassword(&userInfo)
	if errors.Is(err, user.ErrUserNotExists) {
		util.ResponseError(c, util.ErrCodeUserNotExist)
		return
	}

	util.ResponseSuccess(c, nil)
}

// CompleteInformation 完善用户信息
func CompleteInformation(c *gin.Context) {

	var userInfo users.UserInfo
	err := c.BindJSON(&userInfo)
	if err != nil {
		util.ResponseError(c, util.ErrCodeParameter)
		panic(err)
		return
	}

	result := db.DB.Model(&users.UserInfo{}).
		Where("telephone = ?", userInfo.Telephone).
		Updates(map[string]interface{}{
			"username": userInfo.Username,
			"sex":      userInfo.Sex,
			"email":    userInfo.Email,
		})
	err = result.Error
	if err != nil {
		util.ResponseError(c, util.ErrCodeServerBusy)
		return
	}

	util.ResponseSuccess(c, gin.H{
		"result":   "信息修改成功",
		"username": userInfo.Username,
		"sex":      userInfo.Sex,
		"email":    userInfo.Email,
	})
}
