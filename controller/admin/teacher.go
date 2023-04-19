package admin

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/ykgk/kgo/controller/context"
	kindergartendb "gitlab.com/ykgk/kgo/db/kindergarten"
	"gitlab.com/ykgk/kgo/x"
)

type FindKindergartenTeachersRequest struct {
	Query          string `json:"query" form:"query" query:"query"`
	KindergartenID int64  `json:"kindergarten_id" form:"kindergarten_id" query:"kindergarten_id"`
	ClassID        int64  `json:"class_id" form:"class_id" query:"class_id"`
	Role           string `json:"role" form:"role" query:"role"`
	Page           int    `json:"page" form:"page" query:"page"`
	PageSize       int    `json:"page_size" form:"page_size" query:"page_size"`
}

func FindKindergartenTeachers(c echo.Context) error {
	ctx := c.(*context.Context)

	req := FindKindergartenTeachersRequest{}
	if err := ctx.BindAndValidate(&req); err != nil {
		return ctx.BadRequest()
	}
	req.Page, req.PageSize = x.Pagination(req.Page, req.PageSize)

	teachers, total, err := kindergartendb.FindKindergartenTeachers(req.Query, req.KindergartenID, req.ClassID, req.Role, req.Page, req.PageSize)
	if err != nil {
		return ctx.InternalServerError()
	}
	return ctx.List(teachers, req.Page, req.PageSize, total)
}
