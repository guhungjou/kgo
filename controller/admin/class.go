package admin

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/ykgk/kgo/controller/context"
	kindergartendb "gitlab.com/ykgk/kgo/db/kindergarten"
	"gitlab.com/ykgk/kgo/x"
)

type FindKindergartenClassesRequest struct {
	Query          string `json:"query" form:"query" query:"query"`
	KindergartenID int64  `json:"kindergarten_id" form:"kindergarten_id" query:"kindergarten_id"`
	Page           int    `json:"page" form:"page" query:"page"`
	PageSize       int    `json:"page_size" form:"page_size" query:"page_size"`
}

/* 查询班级 */
func FindKindergartenClasses(c echo.Context) error {
	ctx := c.(*context.Context)
	req := FindKindergartenClassesRequest{}
	if err := ctx.BindAndValidate(&req); err != nil {
		return ctx.BadRequest()
	}
	req.Page, req.PageSize = x.Pagination(req.Page, req.PageSize)

	classes, total, err := kindergartendb.FindKindergartenClasses(req.Query, req.KindergartenID, req.Page, req.PageSize)
	if err != nil {
		return ctx.InternalServerError()
	}
	return ctx.List(classes, req.Page, req.PageSize, total)
}

func GetKindergartenClass(c echo.Context) error {
	ctx := c.(*context.Context)

	id := ctx.IntParam(`id`)

	class, err := kindergartendb.GetKindergartenClass(id)
	if err != nil {
		return ctx.InternalServerError()
	} else if class == nil {
		return ctx.NotFound()
	}
	return ctx.Success(class)
}
