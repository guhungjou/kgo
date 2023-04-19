package admin

import (
	"time"

	"github.com/labstack/echo/v4"
	"gitlab.com/ykgk/kgo/controller/context"
	"gitlab.com/ykgk/kgo/controller/middleware"
	admindb "gitlab.com/ykgk/kgo/db/admin"
	"gitlab.com/ykgk/kgo/x"
)

/* 获取系统信息 */
func GetSystemInfo(c echo.Context) error {
	ctx := c.(*context.Context)

	superuser, err := admindb.GetSuperAdminUser()
	if err != nil {
		return ctx.InternalServerError()
	}

	var user *admindb.AdminUser = nil

	tokenID := middleware.GetSessionToken(ctx, "adminuser")
	if tokenID != "" {
		token, err := admindb.GetAdminToken(tokenID)
		if err != nil {
			return ctx.InternalServerError()
		} else if token != nil && token.Status == admindb.AdminTokenStatusOK {
			if token.ExpiresAt.IsZero() || token.ExpiresAt.After(time.Now()) {
				if token.User != nil && token.User.Status == admindb.AdminUserStatusOK {
					user = token.User
				}
			}
		}
	}

	return ctx.Success(map[string]interface{}{
		"superuser": superuser != nil,
		"user":      user,
	})
}

type CreateSuperAdminUserRequest struct {
	Username string `json:"username" form:"username" query:"username" validate:"gte=6"`
	Name     string `json:"name" form:"name" query:"name" validate:"required"`
	Phone    string `json:"phone" form:"phone" query:"phone"`
	Password string `json:"password" form:"password" query:"password" validate:"gte=8"`
}

/* 创建超级用户 */
func CreateSuperAdminUser(c echo.Context) error {
	ctx := c.(*context.Context)
	req := CreateSuperAdminUserRequest{}
	if err := ctx.BindAndValidate(&req); err != nil {
		return ctx.BadRequest()
	}

	/* 验证用户名有效性 */
	if !x.ValidateUsername(req.Username) {
		return ctx.BadRequest()
	}

	user, err := admindb.CreateSuperAdminUser(req.Username, req.Name, req.Phone, req.Password)
	if err != nil {
		return ctx.InternalServerError()
	}
	return ctx.Success(user)
}
