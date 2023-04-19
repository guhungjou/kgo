package teacher

import (
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
	"gitlab.com/ykgk/kgo/controller/context"
	healthdb "gitlab.com/ykgk/kgo/db/health"
	kindergartendb "gitlab.com/ykgk/kgo/db/kindergarten"
	"gitlab.com/ykgk/kgo/x"
)

type FindKindergartenStudentFitnessTestsRequest struct {
	Query     string    `json:"query" form:"query" query:"query"`
	StudentID int64     `json:"student_id" form:"student_id" query:"student_id"`
	ClassID   int64     `json:"class_id" form:"class_id" query:"class_id"`
	StartTime time.Time `json:"start_time" form:"start_time" query:"start_time"`
	EndTime   time.Time `json:"end_time" form:"end_time" query:"end_time"`
	Page      int       `json:"page" form:"page" query:"page"`
	PageSize  int       `json:"page_size" form:"page_size" query:"page_size"`

	HeightWeightFilters     []int    `json:"height_weight_filters" form:"height_weight_filters" query:"height_weight_filters"`
	ShuttleRun10Filters     []int    `json:"shuttle_run_10_filters" form:"shuttle_run_10_filters" query:"shuttle_run_10_filters"`
	StandingLongJumpFilters []int    `json:"standing_long_jump_filters" form:"standing_long_jump_filters" query:"standing_long_jump_filters"`
	BaseballThrowFilters    []int    `json:"baseball_throw_filters" form:"baseball_throw_filters" query:"baseball_throw_filters"`
	BunnyHoppingFilters     []int    `json:"bunny_hopping_filters" form:"bunny_hopping_filters" query:"bunny_hopping_filters"`
	SitAndReachFilters      []int    `json:"sit_and_reach_filters" form:"sit_and_reach_filters" query:"sit_and_reach_filters"`
	BalanceBeamFilters      []int    `json:"balance_beam_filters" form:"balance_beam_filters" query:"balance_beam_filters"`
	TotalStatusFilters      []string `json:"total_status_filters" form:"total_status_filters" query:"total_status_filters"`
}

func FindKindergartenStudentFitnessTests(c echo.Context) error {
	ctx := c.(*context.Context)
	req := FindKindergartenStudentFitnessTestsRequest{}
	if err := ctx.Bind(&req); err != nil {
		return ctx.BadRequest()
	}

	req.Page, req.PageSize = x.Pagination(req.Page, req.PageSize)

	exams, total, err := healthdb.FindKindergartenStudentFitnessTests(
		req.Query, ctx.Teacher.KindergartenID, ctx.Teacher.TeacherClassID(req.ClassID), req.StudentID,
		req.StartTime, req.EndTime, req.HeightWeightFilters, req.ShuttleRun10Filters, req.StandingLongJumpFilters,
		req.BaseballThrowFilters, req.BunnyHoppingFilters, req.SitAndReachFilters, req.BalanceBeamFilters,
		req.TotalStatusFilters, req.Page, req.PageSize)
	if err != nil {
		return ctx.InternalServerError()
	}
	return ctx.List(exams, req.Page, req.PageSize, total)
}

func ExportKindergartenStudentFitnessTests(c echo.Context) error {
	ctx := c.(*context.Context)
	req := FindKindergartenStudentFitnessTestsRequest{}
	if err := ctx.Bind(&req); err != nil {
		return ctx.BadRequest()
	}

	tests, _, err := healthdb.FindKindergartenStudentFitnessTests(
		req.Query, ctx.Teacher.KindergartenID, ctx.Teacher.TeacherClassID(req.ClassID), req.StudentID,
		req.StartTime, req.EndTime, req.HeightWeightFilters, req.ShuttleRun10Filters, req.StandingLongJumpFilters,
		req.BaseballThrowFilters, req.BunnyHoppingFilters, req.SitAndReachFilters, req.BalanceBeamFilters,
		req.TotalStatusFilters, 0, 0)
	if err != nil {
		return ctx.InternalServerError()
	}

	headers := []string{"日期", "班级", "学生", "性别", "10米折返跑(秒)", "立定跳远(厘米)", "网球掷远(米)", "双脚连续跳(秒)", "坐位体前屈(厘米)", "走平衡木(秒)", "总分", "创建时间"}
	rows := make([][]interface{}, 0)
	for _, test := range tests {
		row := make([]interface{}, 0)
		row = append(row, fmt.Sprintf("%04d-%02d-%02d", test.Date.Year(), test.Date.Month(), test.Date.Day()))
		row = append(row, test.Student.Class.Name)
		row = append(row, test.Student.Name)
		row = append(row, test.Student.GenderName())
		if !test.ShuttleRun10UpdatedAt.IsZero() {
			row = append(row, fmt.Sprintf("%.1f (%.0f)", test.ShuttleRun10, test.ShuttleRun10Score))
		} else {
			row = append(row, "###")
		}

		if !test.StandingLongJumpUpdatedAt.IsZero() {
			row = append(row, fmt.Sprintf("%.1f (%.0f)", test.StandingLongJump, test.StandingLongJumpScore))
		} else {
			row = append(row, "###")
		}

		if !test.BaseballThrowUpdatedAt.IsZero() {
			row = append(row, fmt.Sprintf("%.1f (%.0f)", test.BaseballThrow, test.BaseballThrowScore))
		} else {
			row = append(row, "###")
		}

		if !test.BunnyHoppingUpdatedAt.IsZero() {
			row = append(row, fmt.Sprintf("%.1f (%.0f)", test.BunnyHopping, test.BunnyHoppingScore))
		} else {
			row = append(row, "###")
		}

		if !test.SitAndReachUpdatedAt.IsZero() {
			row = append(row, fmt.Sprintf("%.1f (%.0f)", test.SitAndReach, test.SitAndReachScore))
		} else {
			row = append(row, "###")
		}

		if !test.BalanceBeamUpdatedAt.IsZero() {
			row = append(row, fmt.Sprintf("%.1f (%.0f)", test.BalanceBeam, test.BalanceBeamScore))
		} else {
			row = append(row, "###")
		}

		row = append(row, test.TotalScore)

		row = append(row, test.CreatedAt)
		rows = append(rows, row)
	}
	return ctx.XLSX("体测记录", headers, rows)
}

/* 获取最近一百个有体检信息的日期 */
func FindKindergartenStudentFitnessTestDates(c echo.Context) error {
	ctx := c.(*context.Context)

	var classID int64
	if ctx.Teacher.Role != kindergartendb.KindergartenTeacherRoleManager {
		if ctx.Teacher.ClassID > 0 {
			classID = ctx.Teacher.ClassID
		} else {
			return ctx.Success([]interface{}{})
		}
	}
	dates, err := healthdb.FindKindergartenStudentFitnessTestDates(ctx.Teacher.KindergartenID, classID)
	if err != nil {
		return ctx.InternalServerError()
	}
	return ctx.Success(dates)
}

type FindKindergartenStudentFitnessTestScoreVisionRequest struct {
	Date    time.Time `json:"date" form:"date" query:"date" validate:"required"`
	ClassID int64     `json:"class_id" form:"class_id" query:"class_id"`
}

func FindKindergartenStudentFitnessTestScoreVision(c echo.Context, field string) error {
	ctx := c.(*context.Context)

	req := FindKindergartenStudentFitnessTestScoreVisionRequest{}
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

	datas, err := healthdb.FindKindergartenStudentFitnessTestScoreVision(field, ctx.Teacher.KindergartenID, req.ClassID, req.Date)
	if err != nil {
		return ctx.InternalServerError()
	}
	return ctx.Success(datas)
}

func FindKindergartenStudentFitnessTestScoreShuttleRun10Vision(c echo.Context) error {
	return FindKindergartenStudentFitnessTestScoreVision(c, "shuttle_run_10")
}

func FindKindergartenStudentFitnessTestScoreStandingLongJumpVision(c echo.Context) error {
	return FindKindergartenStudentFitnessTestScoreVision(c, "standing_long_jump")
}

func FindKindergartenStudentFitnessTestScoreBaseballThrowVision(c echo.Context) error {
	return FindKindergartenStudentFitnessTestScoreVision(c, "baseball_throw")
}

func FindKindergartenStudentFitnessTestScoreBunnyHoppingVision(c echo.Context) error {
	return FindKindergartenStudentFitnessTestScoreVision(c, "bunny_hopping")
}

func FindKindergartenStudentFitnessTestScoreSitAndReachVision(c echo.Context) error {
	return FindKindergartenStudentFitnessTestScoreVision(c, "sit_and_reach")
}

func FindKindergartenStudentFitnessTestScoreBalanceBeamVision(c echo.Context) error {
	return FindKindergartenStudentFitnessTestScoreVision(c, "balance_beam")
}

type FindKindergartenStudentFitnessTestHeightVisionRequest struct {
	Date    time.Time `json:"date" form:"date" query:"date" validate:"required"`
	ClassID int64     `json:"class_id" form:"class_id" query:"class_id"`
}

func FindKindergartenStudentFitnessTestHeightVision(c echo.Context) error {
	ctx := c.(*context.Context)

	req := FindKindergartenStudentFitnessTestHeightVisionRequest{}
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

	datas, err := healthdb.FindKindergartenStudentFitnessTestHeightVision(ctx.Teacher.KindergartenID, req.ClassID, req.Date)
	if err != nil {
		return ctx.InternalServerError()
	}
	return ctx.Success(datas)
}

type FindKindergartenStudentFitnessTestWeightVisionRequest struct {
	Date    time.Time `json:"date" form:"date" query:"date" validate:"required"`
	ClassID int64     `json:"class_id" form:"class_id" query:"class_id"`
}

func FindKindergartenStudentFitnessTestWeightVision(c echo.Context) error {
	ctx := c.(*context.Context)

	req := FindKindergartenStudentFitnessTestWeightVisionRequest{}
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

	datas, err := healthdb.FindKindergartenStudentFitnessTestWeightVision(ctx.Teacher.KindergartenID, req.ClassID, req.Date)
	if err != nil {
		return ctx.InternalServerError()
	}
	return ctx.Success(datas)
}

type FindKindergartenStudentFitnessTestStatusVisionRequest struct {
	Date    time.Time `json:"date" form:"date" query:"date" validate:"required"`
	ClassID int64     `json:"class_id" form:"class_id" query:"class_id"`
}

func FindKindergartenStudentFitnessTestStatusVision(c echo.Context) error {
	ctx := c.(*context.Context)

	req := FindKindergartenStudentFitnessTestStatusVisionRequest{}
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

	datas, err := healthdb.FindKindergartenStudentFitnessTestStatusVision(ctx.Teacher.KindergartenID, req.ClassID, req.Date)
	if err != nil {
		return ctx.InternalServerError()
	}
	return ctx.Success(datas)
}
