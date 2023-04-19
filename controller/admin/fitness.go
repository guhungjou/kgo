package admin

import (
	"time"

	"github.com/labstack/echo/v4"
	"gitlab.com/ykgk/kgo/controller/context"
	healthdb "gitlab.com/ykgk/kgo/db/health"
	"gitlab.com/ykgk/kgo/x"
)

type FindKindergartenStudentFitnessTestsRequest struct {
	Query          string    `json:"query" form:"query" query:"query"`
	KindergartenID int64     `json:"kindergarten_id" form:"kindergarten_id" query:"kindergarten_id"`
	ClassID        int64     `json:"class_id" form:"class_id" query:"class_id"`
	StudentID      int64     `json:"student_id" form:"student_id" query:"student_id"`
	StartTime      time.Time `json:"start_time" form:"start_time" query:"start_time"`
	EndTime        time.Time `json:"end_time" form:"end_time" query:"end_time"`
	Page           int       `json:"page" form:"page" query:"page"`
	PageSize       int       `json:"page_size" form:"page_size" query:"page_size"`

	HeightWeightFilters     []int `json:"height_weight_filters" form:"height_weight_filters" query:"height_weight_filters"`
	ShuttleRun10Filters     []int `json:"shuttle_run_10_filters" form:"shuttle_run_10_filters" query:"shuttle_run_10_filters"`
	StandingLongJumpFilters []int `json:"standing_long_jump_filters" form:"standing_long_jump_filters" query:"standing_long_jump_filters"`
	BaseballThrowFilters    []int `json:"baseball_throw_filters" form:"baseball_throw_filters" query:"baseball_throw_filters"`
	BunnyHoppingFilters     []int `json:"bunny_hopping_filters" form:"bunny_hopping_filters" query:"bunny_hopping_filters"`
	SitAndReachFilters      []int `json:"sit_and_reach_filters" form:"sit_and_reach_filters" query:"sit_and_reach_filters"`
	BalanceBeamFilters      []int `json:"balance_beam_filters" form:"balance_beam_filters" query:"balance_beam_filters"`

	TotalStatusFilters []string `json:"total_status_filters" form:"total_status_filters" query:"total_status_filters"`
}

func FindKindergartenStudentFitnessTests(c echo.Context) error {
	ctx := c.(*context.Context)

	req := FindKindergartenStudentFitnessTestsRequest{}
	if err := ctx.Bind(&req); err != nil {
		return ctx.BadRequest()
	}

	req.Page, req.PageSize = x.Pagination(req.Page, req.PageSize)

	tests, total, err := healthdb.FindKindergartenStudentFitnessTests(
		req.Query, req.KindergartenID, req.ClassID, req.StudentID,
		req.StartTime, req.EndTime, req.HeightWeightFilters, req.ShuttleRun10Filters, req.StandingLongJumpFilters,
		req.BaseballThrowFilters, req.BunnyHoppingFilters, req.SitAndReachFilters, req.BalanceBeamFilters,
		req.TotalStatusFilters, req.Page, req.PageSize,
	)
	if err != nil {
		return ctx.InternalServerError()
	}
	return ctx.List(tests, req.Page, req.PageSize, total)
}
