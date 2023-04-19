package admin

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/ykgk/kgo/controller/middleware"
)

func Register(adminAPI *echo.Group) {
	adminAPI.GET("/systeminfo", GetSystemInfo)
	adminAPI.POST("/superuser", CreateSuperAdminUser)
	adminAPI.PUT("/login", Login)

	adminAuthAPI := adminAPI.Group("", middleware.AdminAuthMiddleware)
	adminAuthAPI.PUT("/logout", Logout)
	adminAuthAPI.GET("/self", GetSelf)
	adminAuthAPI.PUT("/self/password", UpdateSelfPassword)
	/* 获取用户信息 */
	adminAuthAPI.GET("/user/:id", GetUser)
	/* 获取管理员帐号列表 */
	adminAuthAPI.GET("/admin/users", FindAdminUsers)
	/* 创建管理员帐号 */
	adminAuthAPI.POST("/admin/user", CreateAdminUser)
	/* 更新管理员帐号 */
	adminAuthAPI.PUT("/admin/user/:id", UpdateAdminUser)

	/* 获取幼儿园详情 */
	adminAuthAPI.GET("/kindergarten/:id", GetKindergarten)
	/* 创建幼儿园 */
	adminAuthAPI.POST("/kindergarten", CreateKindergarten)
	/* 更新幼儿园 */
	adminAuthAPI.PUT("/kindergarten/:id", UpdateKindergarten)
	/* 获取幼儿园列表 */
	adminAuthAPI.GET("/kindergartens", FindKindergartens)
	/* 获取老师列表 */
	adminAuthAPI.GET("/kindergarten/teachers", FindKindergartenTeachers)
	/* 获取班级列表 */
	adminAuthAPI.GET("/kindergarten/classes", FindKindergartenClasses)
	/* 获取班级详情 */
	adminAuthAPI.GET("/kindergarten/class/:id", GetKindergartenClass)
	/* 获取学生列表 */
	adminAuthAPI.GET("/kindergarten/students", FindKindergartenStudents)
	/* 获取学生晨检记录 */
	adminAuthAPI.GET("/kindergarten/student/morning/checks", FindKindergartenStudentMorningChecks)
	/* 获取学生体检记录 */
	adminAuthAPI.GET("/kindergarten/student/medical/examinations", FindKindergartenStudentMedicalExaminations)

	/* 获取体侧标准数据 */
	adminAuthAPI.GET("/standard/scale/score/:name", FindStandardScaleScoresByName)

	/* 获取学生体测记录 */
	adminAuthAPI.GET("/kindergarten/student/fitness/tests", FindKindergartenStudentFitnessTests)

	/* 获取省市区数据 */
	adminAuthAPI.GET("/districts", FindDistricts)
	adminAuthAPI.GET("/district/:id", GetDistrict)
}
