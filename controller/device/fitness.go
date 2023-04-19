package device

import (
	"gitlab.com/ykgk/kgo/controller/context"
	healthdb "gitlab.com/ykgk/kgo/db/health"

	"github.com/labstack/echo/v4"
)

type CreateKindergartenStudentFitnessTestHeightRequest struct {
	StudentID int64   `json:"student_id" form:"student_id" query:"student_id" validate:"gt=0"`
	Height    float64 `json:"height" form:"height" query:"height" validate:"gt=0"`
}

func CreateKindergartenStudentFitnessTestHeight(c echo.Context) error {
	ctx := c.(*context.Context)

	req := CreateKindergartenStudentFitnessTestHeightRequest{}
	if err := ctx.BindAndValidate(&req); err != nil {
		return ctx.BadRequest()
	}

	test, err := healthdb.CreateKindergartenStudentFitnessTestHeight(req.StudentID, req.Height)
	if err != nil {
		return ctx.InternalServerError()
	}
	return ctx.Success(test)
}

type CreateKindergartenStudentFitnessTestWeightRequest struct {
	StudentID int64   `json:"student_id" form:"student_id" query:"student_id" validate:"gt=0"`
	Weight    float64 `json:"weight" form:"weight" query:"weight" validate:"gt=0"`
}

func CreateKindergartenStudentFitnessTestWeight(c echo.Context) error {
	ctx := c.(*context.Context)

	req := CreateKindergartenStudentFitnessTestWeightRequest{}
	if err := ctx.BindAndValidate(&req); err != nil {
		return ctx.BadRequest()
	}

	test, err := healthdb.CreateKindergartenStudentFitnessTestWeight(req.StudentID, req.Weight)
	if err != nil {
		return ctx.InternalServerError()
	}
	return ctx.Success(test)
}

type CreateKindergartenStudentFitnessTestShuttleRun10Request struct {
	StudentID    int64   `json:"student_id" form:"student_id" query:"student_id" validate:"gt=0"`
	ShuttleRun10 float64 `json:"shuttle_run_10" form:"shuttle_run_10" query:"shuttle_run_10" validate:"gt=0"`
}

func CreateKindergartenStudentFitnessTestShuttleRun10(c echo.Context) error {
	ctx := c.(*context.Context)

	req := CreateKindergartenStudentFitnessTestShuttleRun10Request{}
	if err := ctx.BindAndValidate(&req); err != nil {
		return ctx.BadRequest()
	}

	test, err := healthdb.CreateKindergartenStudentFitnessTestShuttleRun10(req.StudentID, req.ShuttleRun10)
	if err != nil {
		return ctx.InternalServerError()
	}
	return ctx.Success(test)
}

type CreateKindergartenStudentFitnessTestStandingLongJumpRequest struct {
	StudentID        int64   `json:"student_id" form:"student_id" query:"student_id" validate:"gt=0"`
	StandingLongJump float64 `json:"standing_long_jump" form:"standing_long_jump" query:"standing_long_jump" validate:"gt=0"`
}

func CreateKindergartenStudentFitnessTestStandingLongJump(c echo.Context) error {
	ctx := c.(*context.Context)

	req := CreateKindergartenStudentFitnessTestStandingLongJumpRequest{}
	if err := ctx.BindAndValidate(&req); err != nil {
		return ctx.BadRequest()
	}
	test, err := healthdb.CreateKindergartenStudentFitnessTestStandingLongJump(req.StudentID, req.StandingLongJump)
	if err != nil {
		return ctx.InternalServerError()
	}
	return ctx.Success(test)
}

type CreateKindergartenStudentFitnessTestBaseballThrowRequest struct {
	StudentID     int64   `json:"student_id" form:"student_id" query:"student_id" validate:"gt=0"`
	BaseballThrow float64 `json:"baseball_throw" form:"baseball_throw" query:"baseball_throw" validate:"gt=0"`
}

func CreateKindergartenStudentFitnessTestBaseballThrow(c echo.Context) error {
	ctx := c.(*context.Context)

	req := CreateKindergartenStudentFitnessTestBaseballThrowRequest{}
	if err := ctx.BindAndValidate(&req); err != nil {
		return ctx.BadRequest()
	}
	test, err := healthdb.CreateKindergartenStudentFitnessTestBaseballThrow(req.StudentID, req.BaseballThrow)
	if err != nil {
		return ctx.InternalServerError()
	}
	return ctx.Success(test)
}

type CreateKindergartenStudentFitnessTestBunnyHoppingRequest struct {
	StudentID    int64   `json:"student_id" form:"student_id" query:"student_id" validate:"gt=0"`
	BunnyHopping float64 `json:"bunny_hopping" form:"bunny_hopping" query:"bunny_hopping" validate:"gt=0"`
}

func CreateKindergartenStudentFitnessTestBunnyHopping(c echo.Context) error {
	ctx := c.(*context.Context)

	req := CreateKindergartenStudentFitnessTestBunnyHoppingRequest{}
	if err := ctx.BindAndValidate(&req); err != nil {
		return ctx.BadRequest()
	}
	test, err := healthdb.CreateKindergartenStudentFitnessTestBunnyHopping(req.StudentID, req.BunnyHopping)
	if err != nil {
		return ctx.InternalServerError()
	}
	return ctx.Success(test)
}

type CreateKindergartenStudentFitnessTestSitAndReachRequest struct {
	StudentID   int64   `json:"student_id" form:"student_id" query:"student_id" validate:"gt=0"`
	SitAndReach float64 `json:"sit_and_reach" form:"sit_and_reach" query:"sit_and_reach" validate:"gt=0"`
}

func CreateKindergartenStudentFitnessTestSitAndReach(c echo.Context) error {
	ctx := c.(*context.Context)

	req := CreateKindergartenStudentFitnessTestSitAndReachRequest{}
	if err := ctx.BindAndValidate(&req); err != nil {
		return ctx.BadRequest()
	}
	test, err := healthdb.CreateKindergartenStudentFitnessTestSitAndReach(req.StudentID, req.SitAndReach)
	if err != nil {
		return ctx.InternalServerError()
	}
	return ctx.Success(test)
}

type CreateKindergartenStudentFitnessTestBalanceBeamRequest struct {
	StudentID   int64   `json:"student_id" form:"student_id" query:"student_id" validate:"gt=0"`
	BalanceBeam float64 `json:"balance_beam" form:"balance_beam" query:"balance_beam" validate:"gt=0"`
}

func CreateKindergartenStudentFitnessTestBalanceBeam(c echo.Context) error {
	ctx := c.(*context.Context)

	req := CreateKindergartenStudentFitnessTestBalanceBeamRequest{}
	if err := ctx.BindAndValidate(&req); err != nil {
		return ctx.BadRequest()
	}
	test, err := healthdb.CreateKindergartenStudentFitnessTestBalanceBeam(req.StudentID, req.BalanceBeam)
	if err != nil {
		return ctx.InternalServerError()
	}
	return ctx.Success(test)
}

func GetKindergartenStudentFitnessTestTodayStatus(c echo.Context) error {
	ctx := c.(*context.Context)

	id := ctx.IntParam(`id`)

	test, err := healthdb.GetKindergartenStudentFitnessTestToday(id)
	if err != nil {
		return ctx.InternalServerError()
	}

	if test == nil {
		test = &healthdb.KindergartenStudentFitnessTest{
			StudentID: id,
		}
	}
	return ctx.Success(test)
}
