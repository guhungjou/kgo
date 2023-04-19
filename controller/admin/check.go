package admin

import (
	"time"

	"github.com/labstack/echo/v4"
	"gitlab.com/ykgk/kgo/controller/context"
	healthdb "gitlab.com/ykgk/kgo/db/health"
	"gitlab.com/ykgk/kgo/x"
)

type FindKindergartenStudentMorningChecksRequest struct {
	Query              string    `json:"query" form:"query" query:"query"`
	KindergartenID     int64     `json:"kindergarten_id" form:"kindergarten_id" query:"kindergarten_id"`
	ClassID            int64     `json:"class_id" form:"class_id" query:"class_id"`
	StudentID          int64     `json:"student_id" form:"student_id" query:"student_id"`
	StartTime          time.Time `json:"start_time" form:"start_time" query:"start_time"`
	EndTime            time.Time `json:"end_time" form:"end_time" query:"end_time"`
	Page               int       `json:"page" form:"page" query:"page"`
	PageSize           int       `json:"page_size" form:"page_size" query:"page_size"`
	TemperatureFilters []string  `json:"temperature_filters" form:"temperature_filters" query:"temperature_filters"`
}

func FindKindergartenStudentMorningChecks(c echo.Context) error {
	ctx := c.(*context.Context)

	req := FindKindergartenStudentMorningChecksRequest{}
	if err := ctx.BindAndValidate(&req); err != nil {
		return ctx.BadRequest()
	}
	req.Page, req.PageSize = x.Pagination(req.Page, req.PageSize)

	checks, total, err := healthdb.FindKindergartenStudentMorningChecks(req.Query, req.KindergartenID,
		req.ClassID, req.StudentID, req.TemperatureFilters, req.StartTime, req.EndTime, req.Page, req.PageSize)
	if err != nil {
		return ctx.InternalServerError()
	}
	return ctx.List(checks, req.Page, req.PageSize, total)
}
