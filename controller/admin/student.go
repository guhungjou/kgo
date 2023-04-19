package admin

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/ykgk/kgo/controller/context"
	kindergartendb "gitlab.com/ykgk/kgo/db/kindergarten"
	"gitlab.com/ykgk/kgo/x"
)

type FindKindergartenStudentsRequest struct {
	Query          string `json:"query" form:"query" query:"query"`
	Gender         string `json:"gender" form:"gender" query:"gender"`
	KindergartenID int64  `json:"kindergarten_id" form:"kindergarten_id" query:"kindergarten_id"`
	ClassID        int64  `json:"class_id" form:"class_id" query:"class_id"`
	Page           int    `json:"page" form:"page" query:"page"`
	PageSize       int    `json:"page_size" form:"page_size" query:"page_size"`
}

func FindKindergartenStudents(c echo.Context) error {
	ctx := c.(*context.Context)
	req := FindKindergartenStudentsRequest{}
	if err := ctx.Bind(&req); err != nil {
		return ctx.BadRequest()
	}
	req.Page, req.PageSize = x.Pagination(req.Page, req.PageSize)

	students, total, err := kindergartendb.FindKindergartenStudents(req.Query, req.Gender, req.KindergartenID, req.ClassID, req.Page, req.PageSize)
	if err != nil {
		return ctx.InternalServerError()
	}
	return ctx.List(students, req.Page, req.PageSize, total)
}
