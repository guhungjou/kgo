package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
	"gitlab.com/ykgk/kgo/controller/context"
	admindb "gitlab.com/ykgk/kgo/db/admin"
)

func AdminAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.(*context.Context)
		/* AdminUser 认证 */
		tokenID := GetSessionToken(ctx, "adminuser")
		if tokenID == "" {
			return ctx.Unauthorized()
		}
		token, err := admindb.GetAdminToken(tokenID)
		if err != nil {
			return ctx.InternalServerError()
		} else if token == nil || token.Status != admindb.AdminTokenStatusOK {
			return ctx.Unauthorized()
		} else if !token.ExpiresAt.IsZero() && token.ExpiresAt.Before(time.Now()) { /* TOKEN已过期 */
			return ctx.Unauthorized()
		}
		user := token.User
		if user == nil || (!user.IsSuperuser && user.Status == admindb.AdminUserStatusBanned) {
			return ctx.Unauthorized()
		}
		ctx.AdminUser = user
		return next(ctx)
	}
}
