package device

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/ykgk/kgo/controller/context"
	healthdb "gitlab.com/ykgk/kgo/db/health"
)

type CreateKindergartenStudentMorningCheckRequest struct {
	StudentID   int64   `json:"student_id" form:"student_id" query:"student_id" validate:"gt=0"`
	Temperature float64 `json:"temperature" form:"temperature" query:"temperature" validate:"gt=0"`
	Hand        string  `json:"hand" form:"hand" query:"hand"`
	Mouth       string  `json:"mouth" form:"mouth" query:"mouth"`
	Eye         string  `json:"eye" form:"eye" query:"eye"`
}

func CreateKindergartenStudentMorningCheck(c echo.Context) error {
	ctx := c.(*context.Context)
	req := CreateKindergartenStudentMorningCheckRequest{}
	if err := ctx.BindAndValidate(&req); err != nil {
		return ctx.BadRequest()
	}

	check, err := healthdb.CreateKindergartenStudentMorningCheck(req.StudentID, req.Temperature, req.Hand, req.Mouth, req.Eye)
	if err != nil {
		return ctx.InternalServerError()
	}
	return ctx.Success(check)
}
