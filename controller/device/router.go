package device

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/ykgk/kgo/controller/middleware"
)

func Register(deviceAPI *echo.Group) {
	deviceAPI.GET("/student/:device", GetKindergartenStudent)
	/* 上传晨检信息 */
	deviceAPI.POST("/student/morning/check", CreateKindergartenStudentMorningCheck)
	/* 上传体检信息 */
	deviceAPI.POST("/student/medical/exam/height", CreateKindergartenStudentMedicalExaminationHeight)         // 身高
	deviceAPI.POST("/student/medical/exam/weight", CreateKindergartenStudentMedicalExaminationWeight)         // 体重
	deviceAPI.POST("/student/medical/exam/hemoglobin", CreateKindergartenStudentMedicalExaminationHemoglobin) // 血红蛋白
	deviceAPI.POST("/student/medical/exam/sight", CreateKindergartenStudentMedicalExaminationSight)           // 视力
	deviceAPI.POST("/student/medical/exam/tooth", CreateKindergartenStudentMedicalExaminationTooth)           // 牙齿
	deviceAPI.POST("/student/medical/exam/alt", CreateKindergartenStudentMedicalExaminationALT)               // 谷丙转氨酶
	deviceAPI.GET("/student/:id/medical/exam/status", GetKindergartenStudentMedicalExaminationTodayStatus)    // 获取学生体检状态
	// Create new left sight examination
	deviceAPI.GET("/student/medical/exam/newLSight", CreateKindergartenStudentMedicalExaminationNewLSight)
	// Create new right sight examination
	deviceAPI.GET("/student/medical/exam/newRSight", CreateKindergartenStudentMedicalExaminationNewRSight)
	/* 上传体侧信息 */
	deviceAPI.POST("/student/fitness/test/height", CreateKindergartenStudentFitnessTestHeight)                       // 身高
	deviceAPI.POST("/student/fitness/test/weight", CreateKindergartenStudentFitnessTestWeight)                       // 体重
	deviceAPI.POST("/student/fitness/test/shuttle_run_10", CreateKindergartenStudentFitnessTestShuttleRun10)         // 十米折返跑
	deviceAPI.POST("/student/fitness/test/standing_long_jump", CreateKindergartenStudentFitnessTestStandingLongJump) // 立定跳远
	deviceAPI.POST("/student/fitness/test/baseball_throw", CreateKindergartenStudentFitnessTestBaseballThrow)        // 网球掷远
	deviceAPI.POST("/student/fitness/test/bunny_hopping", CreateKindergartenStudentFitnessTestBunnyHopping)          // 双脚连续跳
	deviceAPI.POST("/student/fitness/test/sit_and_reach", CreateKindergartenStudentFitnessTestSitAndReach)           // 坐位体前屈
	deviceAPI.POST("/student/fitness/test/balance_beam", CreateKindergartenStudentFitnessTestBalanceBeam)            // 走平衡木
	deviceAPI.GET("/student/:id/fitness/test/status", GetKindergartenStudentFitnessTestTodayStatus)                  // 获取学生体测状态

	/* 老师登录 */
	deviceAPI.PUT("/teacher/login", KindergartenTeacherLogin)
	teacherAPI := deviceAPI.Group("/teacher", middleware.KindergartenTeacherXAuthMiddleware)
	teacherAPI.GET("/classes", FindKindergartenClasses)
	teacherAPI.GET("/class/:id/students", FindKindergartenClassStudents)
}
