package account

import (
	"Monitoring-Pressure/dao/db"
	"Monitoring-Pressure/dao/db/user"
	"Monitoring-Pressure/id_gen"
	myjwt "Monitoring-Pressure/jwt"
	sessionAccount "Monitoring-Pressure/middleware/account"
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
	})
}

// LoginHandle 账号密码登录
func LoginHandle(c *gin.Context) {
	sessionAccount.ProcessRequest(c)

	var userInfo users.UserInfo
	err := c.BindJSON(&userInfo)
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

	// 登录成功后签发 JWT
	token, err := myjwt.GenerateToken(userInfo.UserID, userInfo.Telephone, userInfo.Role)
	if err != nil {
		util.ResponseError(c, util.ErrCodeServerBusy)
		return
	}

	// 如果你原来的 session 逻辑还要保留，可以继续保留
	sessionAccount.SetUserId(userInfo.UserID, c)
	sessionAccount.ProcessResponse(c)

	util.ResponseSuccess(c, gin.H{
		"token": token,
		"user_info": gin.H{
			"user_id":   userInfo.UserID,
			"username":  userInfo.Username,
			"telephone": userInfo.Telephone,
			"email":     userInfo.Email,
			"sex":       userInfo.Sex,
			"role":      userInfo.Role,
		},
	})
}

// Login2Handle 手机验证码登录
func Login2Handle(c *gin.Context) {
	var req users.UserInfo
	err := c.BindJSON(&req)
	if err != nil {
		util.ResponseError(c, util.ErrCodeParameter)
		return
	}

	if req.Telephone == "" || req.Token == "" {
		util.ResponseError(c, util.ErrCodeParameter)
		return
	}

	// 校验验证码
	if !verifycode.VerifyCode(req.Telephone, req.Token) {
		util.ResponseError(c, util.ErrCodeInvalidVerifyCode)
		return
	}

	// 查询完整用户信息
	var userInfo users.UserInfo
	err = db.DB.Where("telephone = ?", req.Telephone).First(&userInfo).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		util.ResponseError(c, util.ErrCodeUserNotExist)
		return
	}
	if err != nil {
		util.ResponseError(c, util.ErrCodeServerBusy)
		return
	}

	token, err := myjwt.GenerateToken(userInfo.UserID, userInfo.Telephone, userInfo.Role)
	if err != nil {
		util.ResponseError(c, util.ErrCodeServerBusy)
		return
	}

	util.ResponseSuccess(c, gin.H{
		"token": token,
		"user_info": gin.H{
			"user_id":   userInfo.UserID,
			"username":  userInfo.Username,
			"telephone": userInfo.Telephone,
			"email":     userInfo.Email,
			"sex":       userInfo.Sex,
			"role":      userInfo.Role,
		},
	})
}

// RegisterHandle 注册接口
func RegisterHandle(c *gin.Context) {
	var userInfo users.UserInfo
	err := c.BindJSON(&userInfo)
	if err != nil {
		util.ResponseError(c, util.ErrCodeParameter)
		return
	}

	if len(userInfo.Username) == 0 || len(userInfo.Password) == 0 || len(userInfo.Telephone) == 0 || len(userInfo.Token) == 0 {
		util.ResponseError(c, util.ErrCodeParameter)
		return
	}

	// 验证码校验
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

	// 默认注册为 student
	userInfo.Role = users.RoleStudent

	err = user.Register(&userInfo)
	if errors.Is(err, user.ErrUserExists) {
		util.ResponseError(c, util.ErrCodeUserExist)
		return
	}
	if err != nil {
		util.ResponseError(c, util.ErrCodeServerBusy)
		return
	}

	util.ResponseSuccess(c, gin.H{
		"user_info": gin.H{
			"user_id":   userInfo.UserID,
			"username":  userInfo.Username,
			"telephone": userInfo.Telephone,
			"role":      userInfo.Role,
		},
	})
}

// ResetPasswordHandle 修改密码接口
func ResetPasswordHandle(c *gin.Context) {
	var userInfo users.UserInfo
	err := c.BindJSON(&userInfo)
	if err != nil {
		util.ResponseError(c, util.ErrCodeParameter)
		return
	}

	if len(userInfo.Telephone) == 0 || len(userInfo.Password) == 0 || len(userInfo.Token) == 0 {
		util.ResponseError(c, util.ErrCodeParameter)
		return
	}

	if !verifycode.VerifyCode(userInfo.Telephone, userInfo.Token) {
		util.ResponseError(c, util.ErrCodeInvalidVerifyCode)
		return
	}

	err = user.ResetPassword(&userInfo)
	if errors.Is(err, user.ErrUserNotExists) {
		util.ResponseError(c, util.ErrCodeUserNotExist)
		return
	}
	if err != nil {
		util.ResponseError(c, util.ErrCodeServerBusy)
		return
	}

	util.ResponseSuccess(c, gin.H{
		"message": "密码修改成功",
	})
}

// CompleteInformation 完善用户信息
func CompleteInformation(c *gin.Context) {
	var userInfo users.UserInfo
	err := c.BindJSON(&userInfo)
	if err != nil {
		util.ResponseError(c, util.ErrCodeParameter)
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

// GetProfile 获取当前登录用户信息
func GetProfile(c *gin.Context) {
	userIDValue, exists := c.Get("user_id")
	if !exists {
		util.ResponseError(c, util.ErrCodeNeedLogin)
		return
	}

	userID, ok := userIDValue.(int64)
	if !ok {
		util.ResponseError(c, util.ErrCodeNeedLogin)
		return
	}

	var userInfo users.UserInfo
	err := db.DB.Where("user_id = ?", userID).First(&userInfo).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		util.ResponseError(c, util.ErrCodeUserNotExist)
		return
	}
	if err != nil {
		util.ResponseError(c, util.ErrCodeServerBusy)
		return
	}

	util.ResponseSuccess(c, gin.H{
		"user_info": gin.H{
			"user_id":   userInfo.UserID,
			"username":  userInfo.Username,
			"telephone": userInfo.Telephone,
			"email":     userInfo.Email,
			"sex":       userInfo.Sex,
			"role":      userInfo.Role,
		},
	})
}
