package admin

import (
	"github.com/labstack/echo/v4"
	"gitlab.com/ykgk/kgo/controller/context"
	"gitlab.com/ykgk/kgo/controller/errors"
	admindb "gitlab.com/ykgk/kgo/db/admin"
)

type UpdateSelfPasswordRequest struct {
	Old string `json:"old" form:"old" query:"old" validate:"gt=0"`
	New string `json:"new" form:"new" query:"new" validate:"gt=0"`
}

/* 更新帐号密码 */
func UpdateSelfPassword(c echo.Context) error {
	ctx := c.(*context.Context)

	req := UpdateSelfPasswordRequest{}
	if err := ctx.BindAndValidate(&req); err != nil {
		return ctx.BadRequest()
	}

	if !ctx.AdminUser.Auth(req.Old) { /* 密码错误 */
		return ctx.Fail(errors.ErrAdminUserPasswordIncorrect, nil)
	}

	if err := admindb.UpdateAdminUserPassword(ctx.AdminUser.ID, req.New); err != nil {
		return ctx.InternalServerError()
	}
	return ctx.Success(nil)
}
