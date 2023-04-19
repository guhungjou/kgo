package admin

import (
	"time"

	"github.com/labstack/echo/v4"
	"gitlab.com/ykgk/kgo/controller/context"
	healthdb "gitlab.com/ykgk/kgo/db/health"
	"gitlab.com/ykgk/kgo/x"
)

type FindKindergartenStudentMedicalExaminationsRequest struct {
	Query          string    `json:"query" form:"query" query:"query"`
	StudentID      int64     `json:"student_id" form:"student_id" query:"student_id"`
	KindergartenID int64     `json:"kindergarten_id" form:"kindergarten_id" query:"kindergarten_id"`
	ClassID        int64     `json:"class_id" form:"class_id" query:"class_id"`
	StartTime      time.Time `json:"start_time" form:"start_time" query:"start_time"`
	EndTime        time.Time `json:"end_time" form:"end_time" query:"end_time"`
	Page           int       `json:"page" form:"page" query:"page"`
	PageSize       int       `json:"page_size" form:"page_size" query:"page_size"`

	HeightFilters     []string `json:"height_filters" form:"height_filters" query:"height_filters"`
	WeightFilters     []string `json:"weight_filters" form:"weight_filters" query:"weight_filters"`
	HemoglobinFilters []string `json:"hemoglobin_filters" form:"hemoglobin_filters" query:"hemoglobin_filters"`
	SightFilters      []string `json:"sight_filters" form:"sight_filters" query:"sight_filters"`
	ALTFilters        []string `json:"alt_filters" form:"alt_filters" query:"alt_filters"`
	BMIFilters        []string `json:"bmi_filters" form:"bmi_filters" query:"bmi_filters"`
}

func FindKindergartenStudentMedicalExaminations(c echo.Context) error {
	ctx := c.(*context.Context)
	req := FindKindergartenStudentMedicalExaminationsRequest{}
	if err := ctx.Bind(&req); err != nil {
		return ctx.BadRequest()
	}

	req.Page, req.PageSize = x.Pagination(req.Page, req.PageSize)

	exams, total, err := healthdb.FindKindergartenStudentMedicalExaminations(
		req.Query, req.KindergartenID, req.ClassID, req.StudentID,
		req.HeightFilters, req.WeightFilters, req.HemoglobinFilters, req.SightFilters, req.ALTFilters, req.BMIFilters,
		req.StartTime, req.EndTime, req.Page, req.PageSize)
	if err != nil {
		return ctx.InternalServerError()
	}
	return ctx.List(exams, req.Page, req.PageSize, total)
}
