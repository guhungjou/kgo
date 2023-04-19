package device

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/ykgk/kgo/controller/context"
	kindergartendb "gitlab.com/ykgk/kgo/db/kindergarten"
)

type GetKindergartenStudentRequest struct {
	Type string `json:"type" form:"type" query:"type"`
}

func GetKindergartenStudent(c echo.Context) error {
	ctx := c.(*context.Context)

	var req GetKindergartenStudentRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.BadRequest()
	}

	var student *kindergartendb.KindergartenStudent
	var err error

	if req.Type == "id" {
		id := ctx.IntParam(`device`)
		student, err = kindergartendb.GetKindergartenStudent(id)
	} else {
		device := ctx.StringParam(`device`)
		student, err = kindergartendb.GetKindergartenStudentByDevice(device)
	}
	if err != nil {
		return ctx.InternalServerError()
	}
	return ctx.Success(student)
}
