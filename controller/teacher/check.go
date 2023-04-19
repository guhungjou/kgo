package teacher

import (
	"fmt"
	"time"

	healthdb "gitlab.com/ykgk/kgo/db/health"
	kindergartendb "gitlab.com/ykgk/kgo/db/kindergarten"

	"github.com/labstack/echo/v4"
	"gitlab.com/ykgk/kgo/controller/context"
	"gitlab.com/ykgk/kgo/x"
)

// type FindKindergartenStudentMorningCheckStatsRequest struct {
// 	StartDate time.Time `json:"start_date" form:"start_date" query:"start_date" validate:"required"`
// 	EndDate   time.Time `json:"end_date" form:"end_date" query:"end_date" validate:"required"`
// 	ClassID   int64     `json:"class_id" form:"class_id" query:"class_id"`
// }

// func FindKindergartenStudentMorningCheckStats(c echo.Context) error {
// 	ctx := c.(*context.Context)
// 	req := FindKindergartenStudentMorningCheckStatsRequest{}
// 	if err := ctx.BindAndValidate(&req); err != nil {
// 		return ctx.BadRequest()
// 	}

// 	stats, err := healthdb.FindKindergartenStudentMorningCheckStats(ctx.Teacher.KindergartenID, req.ClassID, req.StartDate, req.EndDate)
// 	if err != nil {
// 		return ctx.InternalServerError()
// 	}
// 	return ctx.Success(stats)
// }

type GetKindergartenStudentMorningCheckStatRequest struct {
	Date    time.Time `json:"date" form:"date" query:"date" validate:"required"`
	ClassID int64     `json:"class_id" form:"class_id" query:"class_id"`
}

/* 获取晨检统计 */
func GetKindergartenStudentMorningCheckStat(c echo.Context) error {
	ctx := c.(*context.Context)

	req := GetKindergartenStudentMorningCheckStatRequest{}
	if err := ctx.BindAndValidate(&req); err != nil {
		return ctx.BadRequest()
	}

	stat, err := healthdb.GetKindergartenStudentMorningCheckStat(ctx.Teacher.KindergartenID, req.ClassID, req.Date)
	if err != nil {
		return ctx.InternalServerError()
	}
	return ctx.Success(stat)
}

type FindKindergartenStudentMorningChecksRequest struct {
	Query              string    `json:"query" form:"query" query:"query"`
	StudentID          int64     `json:"student_id" form:"student_id" query:"student_id"`
	ClassID            int64     `json:"class_id" form:"class_id" query:"class_id"`
	StartTime          time.Time `json:"start_time" form:"start_time" query:"start_time"`
	EndTime            time.Time `json:"end_time" form:"end_time" query:"end_time"`
	Page               int       `json:"page" form:"page" query:"page"`
	PageSize           int       `json:"page_size" form:"page_size" query:"page_size"`
	TemperatureFilters []string  `json:"temperature_filters" form:"temperature_filters" query:"temperature_filters"`
}

func FindKindergartenStudentMorningChecks(c echo.Context) error {
	ctx := c.(*context.Context)

	req := FindKindergartenStudentMorningChecksRequest{}
	if err := ctx.Bind(&req); err != nil {
		return ctx.BadRequest()
	}

	req.Page, req.PageSize = x.Pagination(req.Page, req.PageSize)

	checks, total, err := healthdb.FindKindergartenStudentMorningChecks(req.Query,
		ctx.Teacher.KindergartenID, ctx.Teacher.TeacherClassID(req.ClassID), req.StudentID,
		req.TemperatureFilters,
		req.StartTime, req.EndTime, req.Page, req.PageSize)
	if err != nil {
		return ctx.InternalServerError()
	}
	return ctx.List(checks, req.Page, req.PageSize, total)
}

func ExportKindergartenStudentMorningChecks(c echo.Context) error {
	ctx := c.(*context.Context)

	req := FindKindergartenStudentMorningChecksRequest{}
	if err := ctx.Bind(&req); err != nil {
		return ctx.BadRequest()
	}

	checks, _, err := healthdb.FindKindergartenStudentMorningChecks(req.Query, ctx.Teacher.KindergartenID,
		ctx.Teacher.TeacherClassID(req.ClassID), req.StudentID, req.TemperatureFilters, req.StartTime, req.EndTime, 0, 0)
	if err != nil {
		return ctx.InternalServerError()
	}

	headers := []string{"日期", "班级", "学生", "性别", "体温", "手", "口", "眼", "创建时间"}
	rows := make([][]interface{}, 0)
	for _, check := range checks {
		row := make([]interface{}, 0)
		row = append(row, fmt.Sprintf("%04d-%02d-%02d", check.Date.Year(), check.Date.Month(), check.Date.Day()))
		row = append(row, check.Student.Class.Name)
		row = append(row, check.Student.Name)
		row = append(row, check.Student.GenderName())
		row = append(row, check.Temperature)
		row = append(row, check.Hand)
		row = append(row, check.Mouth)
		row = append(row, check.Eye)
		row = append(row, check.CreatedAt)
		rows = append(rows, row)
	}
	return ctx.XLSX("晨检记录", headers, rows)
}

type FindKindergartenStudentMorningCheckTemperatureVisionRequest struct {
	Date    time.Time `json:"date" form:"date" query:"date" validate:"required"`
	ClassID int64     `json:"class_id" form:"class_id" query:"class_id"`
}

func FindKindergartenStudentMorningCheckTemperatureVision(c echo.Context) error {
	ctx := c.(*context.Context)

	req := FindKindergartenStudentMorningCheckTemperatureVisionRequest{}
	if err := ctx.BindAndValidate(&req); err != nil {
		return ctx.BadRequest()
	}

	if ctx.Teacher.Role != kindergartendb.KindergartenTeacherRoleManager {
		if ctx.Teacher.ClassID > 0 {
			req.ClassID = ctx.Teacher.ClassID
		} else {
			return ctx.Success([]interface{}{})
		}
	}

	datas, err := healthdb.FindKindergartenStudentMorningCheckTemperatureVision(ctx.Teacher.KindergartenID, req.ClassID, req.Date)
	if err != nil {
		return ctx.InternalServerError()
	}
	return ctx.Success(datas)
}
