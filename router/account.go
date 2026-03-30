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

		// 查询成绩
		studentGroup.GET("/home", studentController.GetMyAcademicRecordListHandle)

		// 心理问卷
		studentGroup.POST("/psychological/submit", studentController.SubmitPsychologicalAssessmentHandle)
		studentGroup.GET("/psychological/list", studentController.GetPsychologicalAssessmentListHandle)

		// 压力评估结果
		studentGroup.GET("/stress/result/latest", studentController.GetLatestStressAssessmentResultHandle)

		// 干预建议记录
		studentGroup.GET("/intervention/list", studentController.GetMyInterventionRecordListHandle)
	}

	// 老师路由
	teacherGroup := r.Group("/teacher")
	teacherGroup.Use(accountMiddleware.JWTAuthMiddleware())
	teacherGroup.Use(accountMiddleware.RoleMiddleware(string(users.RoleTeacher), string(users.RoleAdmin)))
	{
		teacherGroup.GET("/home", teacherController.GetHomeHandle)

		// 查看学生相关数据
		teacherGroup.GET("/student/academic/list", teacherController.GetStudentAcademicRecordListHandle)
		teacherGroup.GET("/student/stress/list", teacherController.GetStudentStressAssessmentListHandle)
		teacherGroup.GET("/student/intervention/list", teacherController.GetStudentInterventionRecordListHandle)
	}

	// 管理员路由
	adminGroup := r.Group("/admin")
	adminGroup.Use(accountMiddleware.JWTAuthMiddleware())
	adminGroup.Use(accountMiddleware.RoleMiddleware(string(users.RoleAdmin)))
	{
		adminGroup.GET("/home", adminController.GetHomeHandle)
		adminGroup.POST("/user/updateRole", adminController.UpdateUserRoleHandle)
		adminGroup.GET("/user/list", adminController.GetUserListHandle)

		// 成绩 Excel 导入
		adminGroup.POST("/academic/import", adminController.ImportAcademicExcelHandle)
		adminGroup.GET("/academic/import/list", adminController.GetAcademicImportRecordListHandle)

		// 成绩数据
		adminGroup.GET("/academic/list", adminController.GetAcademicRecordListHandle)

		// 压力评估结果
		adminGroup.GET("/stress/result/list", adminController.GetStressAssessmentResultListHandle)

		// 干预建议模板
		adminGroup.POST("/intervention/create", adminController.CreateInterventionSuggestionHandle)
		adminGroup.GET("/intervention/list", adminController.GetInterventionSuggestionListHandle)
	}

	return r
}
