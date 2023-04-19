package admin

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/ykgk/kgo/controller/context"
	"gitlab.com/ykgk/kgo/controller/errors"
	kindergartendb "gitlab.com/ykgk/kgo/db/kindergarten"
	"gitlab.com/ykgk/kgo/x"
)

type CreateKindergartenRequest struct {
	DistrictID string `json:"district_id" form:"district_id" query:"district_id" validate:"gt=0"`
	Name       string `json:"name" form:"name" query:"name" validate:"gt=0"`
	Remark     string `json:"remark" form:"remark" query:"remark"`

	ManagerUsername string `json:"manager_username" form:"manager_username" query:"manager_username" validate:"gt=0"`
	ManagerName     string `json:"manager_name" form:"manager_name" query:"manager_name" validate:"gt=0"`
	ManagerGender   string `json:"manager_gender" form:"manager_gender" query:"manager_gender"`
	ManagerPassword string `json:"manager_password" form:"manager_password" query:"manager_password" validate:"gt=0"`
	ManagerPhone    string `json:"manager_phone" form:"manager_phone" query:"manager_phone"`
}

/* 创建幼儿园 */
func CreateKindergarten(c echo.Context) error {
	ctx := c.(*context.Context)
	req := CreateKindergartenRequest{}
	if err := ctx.BindAndValidate(&req); err != nil {
		return ctx.BadRequest()
	}

	if !x.ValidateUsername(req.ManagerUsername) {
		return ctx.BadRequest()
	}

	teacher, err := kindergartendb.GetKindergartenTeacherByUsername(req.ManagerUsername)
	if err != nil {
		return ctx.InternalServerError()
	} else if teacher != nil {
		return ctx.Fail(errors.ErrKindergartenTeacherUsernameDuplicated, nil)
	}

	kg, _, err := kindergartendb.CreateKindergarten(req.DistrictID, req.Name, req.Remark, req.ManagerUsername, req.ManagerPassword,
		req.ManagerName, req.ManagerGender, req.ManagerPhone, ctx.AdminUser.ID)
	if err != nil {
		return ctx.InternalServerError()
	}
	return ctx.Success(kg)
}

type UpdateKindergartenRequest struct {
	DistrictID string `json:"district_id" form:"district_id" query:"district_id" validate:"gt=0"`
	Name       string `json:"name" form:"name" query:"name" validate:"gt=0"`
	Remark     string `json:"remark" form:"remark" query:"remark"`
}

func UpdateKindergarten(c echo.Context) error {
	ctx := c.(*context.Context)

	id := ctx.IntParam(`id`)

	req := UpdateKindergartenRequest{}
	if err := ctx.BindAndValidate(&req); err != nil {
		return ctx.BadRequest()
	}

	kg, err := kindergartendb.UpdateKindergarten(id, req.DistrictID, req.Name, req.Remark, ctx.AdminUser.ID)
	if err != nil {
		return ctx.InternalServerError()
	}
	return ctx.Success(kg)
}

type FindKindergartensRequest struct {
	Query      string `json:"query" form:"query" query:"query"`
	DistrictID string `json:"district_id" form:"district_id" query:"district_id"`
	Page       int    `json:"page" form:"page" query:"page"`
	PageSize   int    `json:"page_size" form:"page_size" query:"page_size"`
}

func FindKindergartens(c echo.Context) error {
	ctx := c.(*context.Context)
	req := FindKindergartensRequest{}
	if err := ctx.BindAndValidate(&req); err != nil {
		return ctx.BadRequest()
	}
	req.Page, req.PageSize = x.Pagination(req.Page, req.PageSize)

	kindergartens, total, err := kindergartendb.FindKindergartens(req.Query, req.DistrictID, req.Page, req.PageSize)
	if err != nil {
		return ctx.InternalServerError()
	}
	return ctx.List(kindergartens, req.Page, req.PageSize, total)
}

func GetKindergarten(c echo.Context) error {
	ctx := c.(*context.Context)
	id := ctx.IntParam(`id`)

	kindergarten, err := kindergartendb.GetKindergarten(id)
	if err != nil {
		return ctx.InternalServerError()
	} else if kindergarten == nil {
		return ctx.NotFound()
	}
	return ctx.Success(kindergarten)
}
