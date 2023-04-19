package teacher

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/ykgk/kgo/controller/middleware"
)

func Register(teacherAPI *echo.Group) {
	teacherAPI.PUT("/login", Login)
	teacherAuthAPI := teacherAPI.Group("", middleware.KindergartenTeacherAuthMiddleware)
	teacherAuthAPI.PUT("/logout", Logout)
	teacherAuthAPI.GET("/self", GetSelf)
	/* 修改个人信息 */
	teacherAuthAPI.PUT("/self", UpdateSelf)
	/* 修改密码 */
	teacherAuthAPI.PUT("/password", UpdateSelfPassword)
	/* 获取当前账号的班级 */
	teacherAuthAPI.GET("/self/classes", FindSelfClasses)

	/* 获取班级列表 */
	teacherAuthAPI.GET("/classes", FindKindergartenClasses)
	/* 创建班级 */
	teacherAuthAPI.POST("/class", CreateKindergartenClass)
	/* 获取班级详情 */
	teacherAuthAPI.GET("/class/:id", GetKindergartenClass)
	/* 编辑班级 */
	teacherAuthAPI.PUT("/class/:id", UpdateKindergartenClass)
	/* 上传班级的模板 */
	teacherAPI.GET("/class/load/template", DownloadKindergartenClassTemplate)
	/* 解析XLSX文件，获取班级信息 */
	teacherAuthAPI.POST("/class/load", LoadKindergartenClass)
	/* 批量创建班级 */
	teacherAuthAPI.POST("/classes", CreateKindergartenClassLoad)
	/* 上传老师的模板 */
	teacherAPI.GET("/teacher/load/template", DownloadKindergartenTeacherTemplate)
	/* 解析XLSX文件，获取老师信息 */
	teacherAuthAPI.POST("/teacher/load", LoadKindergartenTeacher)
	/* 批量创建老师 */
	teacherAuthAPI.POST("/teachers", CreateKindergartenTeacherLoad)
	/* 获取老师列表 */
	teacherAuthAPI.GET("/teachers", FindKindergartenTeachers)
	/* 创建老师 */
	teacherAuthAPI.POST("/teacher", CreateKindergartenTeacher)
	/* 编辑老师 */
	teacherAuthAPI.PUT("/teacher/:id", UpdateKindergartenTeacher)
	/* 上传学生信息的模板 */
	teacherAPI.GET("/student/load/template", DownloadKindergartenStudentTemplate)
	teacherAuthAPI.POST("/student/load", LoadKindergartenStudent)
	/* 批量创建学生 */
	teacherAuthAPI.POST("/students", CreateKindergartenStudentLoad)
	/* 获取学生列表 */
	teacherAuthAPI.GET("/students", FindKindergartenStudents)
	teacherAuthAPI.GET("/students/export", ExportKindergartenStudents)
	/* 创建学生 */
	teacherAuthAPI.POST("/student", CreateKindergartenStudent)
	/* 编辑学生 */
	teacherAuthAPI.PUT("/student/:id", UpdateKindergartenStudent)
	/* 获取学生晨检记录 */
	teacherAuthAPI.GET("/student/morning/checks", FindKindergartenStudentMorningChecks)
	teacherAuthAPI.GET("/student/morning/checks/export", ExportKindergartenStudentMorningChecks) // 导出晨检记录
	/* 获取学生体检记录 */
	teacherAuthAPI.GET("/student/medical/examinations", FindKindergartenStudentMedicalExaminations)
	teacherAuthAPI.GET("/student/medical/examinations/export", ExportKindergartenStudentMedicalExaminations) // 导出体检记录
	teacherAuthAPI.GET("/student/medical/examination/:id", GetKindergartenStudentMedicalExamination)
	teacherAuthAPI.POST("/student/medical/examination/ai", GetKindergartenAnswer)
	/* 获取学生体测记录 */
	teacherAuthAPI.GET("/student/fitness/tests", FindKindergartenStudentFitnessTests)
	teacherAuthAPI.GET("/student/fitness/tests/export", ExportKindergartenStudentFitnessTests)

	/* 删除学生 */
	teacherAuthAPI.DELETE("/student/:id", DeleteKindergartenStudent)
	/* 删除老师 */
	teacherAuthAPI.DELETE("/teacher/:id", DeleteKindergartenTeacher)
	/* 删除班级 */
	teacherAuthAPI.DELETE("/class/:id", DeleteKindergartenClass)

	/* 晨检统计 */
	teacherAuthAPI.GET("/student/morning/check/stat", GetKindergartenStudentMorningCheckStat)
	// teacherAuthAPI.GET("/student/morning/check/stats", FindKindergartenStudentMorningCheckStats)
	teacherAuthAPI.GET("/student/morning/check/temperature/vision", FindKindergartenStudentMorningCheckTemperatureVision)
	/* 体检统计 */
	teacherAuthAPI.GET("/student/medical/examination/dates", FindKindergartenStudentMedicalExaminationDates)
	teacherAuthAPI.GET("/student/medical/examination/height/vision", FindKindergartenStudentMedicalExaminationHeightVision)
	teacherAuthAPI.GET("/student/medical/examination/weight/vision", FindKindergartenStudentMedicalExaminationWeightVision)
	teacherAuthAPI.GET("/student/medical/examination/bmi/vision", FindKindergartenStudentMedicalExaminationBMIVision)
	teacherAuthAPI.GET("/student/medical/examination/hemoglobin/vision", FindKindergartenStudentMedicalExaminationHemoglobinVision)
	teacherAuthAPI.GET("/student/medical/examination/alt/vision", FindKindergartenStudentMedicalExaminationALTVision)
	teacherAuthAPI.GET("/student/medical/examination/sight/vision", FindKindergartenStudentMedicalExaminationSightVision)
	teacherAuthAPI.GET("/student/medical/examination/eye/vision", FindKindergartenStudentMedicalExaminationEyeVision)
	teacherAuthAPI.POST("/student/medical/examination/alt/batch", BatchCreateKindergartenStudentMedicalExaminationALT)
	/* 体测统计 */
	teacherAuthAPI.GET("/student/fitness/test/dates", FindKindergartenStudentFitnessTestDates)
	teacherAuthAPI.GET("/student/fitness/test/height/vision", FindKindergartenStudentFitnessTestHeightVision)
	teacherAuthAPI.GET("/student/fitness/test/weight/vision", FindKindergartenStudentFitnessTestWeightVision)
	teacherAuthAPI.GET("/student/fitness/test/shuttle_run_10/vision", FindKindergartenStudentFitnessTestScoreShuttleRun10Vision)
	teacherAuthAPI.GET("/student/fitness/test/standing_long_jump/vision", FindKindergartenStudentFitnessTestScoreStandingLongJumpVision)
	teacherAuthAPI.GET("/student/fitness/test/baseball_throw/vision", FindKindergartenStudentFitnessTestScoreBaseballThrowVision)
	teacherAuthAPI.GET("/student/fitness/test/bunny_hopping/vision", FindKindergartenStudentFitnessTestScoreBunnyHoppingVision)
	teacherAuthAPI.GET("/student/fitness/test/sit_and_reach/vision", FindKindergartenStudentFitnessTestScoreSitAndReachVision)
	teacherAuthAPI.GET("/student/fitness/test/balance_beam/vision", FindKindergartenStudentFitnessTestScoreBalanceBeamVision)
	teacherAuthAPI.GET("/student/fitness/test/status/vision", FindKindergartenStudentFitnessTestStatusVision)

	/* 获取省市区数据 */
	teacherAuthAPI.GET("/districts", FindDistricts)
	teacherAuthAPI.GET("/district/:id", GetDistrict)
}
