package middleware

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"gitlab.com/ykgk/kgo/controller/context"
	kindergartendb "gitlab.com/ykgk/kgo/db/kindergarten"
)

func KindergartenTeacherAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.(*context.Context)
		tokenID := GetSessionToken(ctx, "teacher")
		if tokenID == "" {
			return c.NoContent(http.StatusUnauthorized)
		}
		token, err := kindergartendb.GetKindergartenTeacherToken(tokenID)
		if err != nil {
			return ctx.InternalServerError()
		} else if token == nil || token.Status != kindergartendb.KindergartenTeacherTokenStatusOK {
			return ctx.Unauthorized()
		} else if !token.ExpiresAt.IsZero() && token.ExpiresAt.Before(time.Now()) { /* TOKEN已过期 */
			return ctx.Unauthorized()
		}
		teacher := token.Teacher
		if teacher == nil || teacher.Deleted {
			return ctx.Unauthorized()
		}
		ctx.Teacher = teacher
		return next(ctx)
	}
}

func KindergartenTeacherXAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.(*context.Context)
		tokenID := ctx.Request().Header.Get("X-Auth-Token")
		if tokenID == "" {
			return ctx.Unauthorized()
		}
		token, err := kindergartendb.GetKindergartenTeacherToken(tokenID)
		if err != nil {
			return ctx.InternalServerError()
		} else if token == nil || token.Status != kindergartendb.KindergartenTeacherTokenStatusOK {
			return ctx.Unauthorized()
		} else if !token.ExpiresAt.IsZero() && token.ExpiresAt.Before(time.Now()) { /* TOKEN已过期 */
			return ctx.Unauthorized()
		}
		teacher := token.Teacher
		if teacher == nil || teacher.Deleted {
			return ctx.Unauthorized()
		}
		ctx.Teacher = teacher
		return next(ctx)
	}
}
