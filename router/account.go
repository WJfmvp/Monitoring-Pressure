package router

import (
	accountController "Monitoring-Pressure/controller/account"
	adminController "Monitoring-Pressure/controller/admin"
	studentController "Monitoring-Pressure/controller/student"
	teacherController "Monitoring-Pressure/controller/teacher"
	accountMiddleware "Monitoring-Pressure/middleware/account"
	"Monitoring-Pressure/models/users"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 账号相关路由
	accountGroup := r.Group("/account")
	{
		accountGroup.GET("/sendVerifyCode", accountController.SendVerifyCodeHandle)
		accountGroup.POST("/login", accountController.LoginHandle)
		accountGroup.POST("/login/code", accountController.Login2Handle)
		accountGroup.POST("/register", accountController.RegisterHandle)
		accountGroup.POST("/resetPassword", accountController.ResetPasswordHandle)
		accountGroup.POST("/completeInformation", accountController.CompleteInformation)
	}

	// 需要登录才能访问
	authGroup := r.Group("/")
	authGroup.Use(accountMiddleware.JWTAuthMiddleware())
	{
		authGroup.GET("/account/profile", accountController.GetProfile)
	}

	// 学生路由
	studentGroup := r.Group("/student")
	studentGroup.Use(accountMiddleware.JWTAuthMiddleware())
	studentGroup.Use(accountMiddleware.RoleMiddleware(string(users.RoleStudent), string(users.RoleAdmin)))
	{
		studentGroup.GET("/home", studentController.GetHomeHandle)
	}

	// 老师路由
	teacherGroup := r.Group("/teacher")
	teacherGroup.Use(accountMiddleware.JWTAuthMiddleware())
	teacherGroup.Use(accountMiddleware.RoleMiddleware(string(users.RoleTeacher), string(users.RoleAdmin)))
	{
		teacherGroup.GET("/home", teacherController.GetHomeHandle)
	}

	// 管理员路由
	adminGroup := r.Group("/admin")
	adminGroup.Use(accountMiddleware.JWTAuthMiddleware())
	adminGroup.Use(accountMiddleware.RoleMiddleware(string(users.RoleAdmin)))
	{
		adminGroup.GET("/home", adminController.GetHomeHandle)
		adminGroup.POST("/user/updateRole", adminController.UpdateUserRoleHandle)
		adminGroup.GET("/user/list", adminController.GetUserListHandle)
	}

	return r
}
