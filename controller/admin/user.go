package admin

import (
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"gitlab.com/ykgk/kgo/controller/context"
	"gitlab.com/ykgk/kgo/controller/errors"
	admindb "gitlab.com/ykgk/kgo/db/admin"
	"gitlab.com/ykgk/kgo/x"
)

type LoginRequest struct {
	Username string `json:"username" form:"username" query:"username" validate:"gt=0"`
	Password string `json:"password" form:"password" query:"password" validate:"gt=0"`
}

/* 管理员帐号登录 */
func Login(c echo.Context) error {
	ctx := c.(*context.Context)
	req := LoginRequest{}
	if err := ctx.BindAndValidate(&req); err != nil {
		return ctx.BadRequest()
	}

	user, err := admindb.GetAdminUserByUsername(req.Username)
	if err != nil {
		return ctx.InternalServerError()
	} else if user == nil { /* 用户不存在 */
		return ctx.Fail(errors.ErrAdminUserNotFound, nil)
	} else if !user.IsSuperuser && user.Status != admindb.AdminUserStatusOK { /* 用户状态不合法 */
		return ctx.Fail(errors.ErrAdminUserStatusInvalid, nil)
	}

	if !user.Auth(req.Password) { /* 密码错误 */
		return ctx.Fail(errors.ErrAdminUserPasswordIncorrect, nil)
	}

	expiresAt := time.Now().AddDate(0, 0, 1)
	token, err := admindb.CreateAdminToken(user.ID, expiresAt, ctx.Request().UserAgent(), ctx.RealIP())
	if err != nil {
		return ctx.InternalServerError()
	}

	sess, _ := session.Get("adminuser", c)
	sess.Options = &sessions.Options{
		Path:   "/",
		MaxAge: 3600 * 24,
	}
	sess.Values["token"] = token.ID
	sess.Save(ctx.Request(), ctx.Response())
	return ctx.Success(nil)
}

/* 注销登录 */
func Logout(c echo.Context) error {
	ctx := c.(*context.Context)

	sess, _ := session.Get("adminuser", c)
	sess.Options = &sessions.Options{
		Path:   "/",
		MaxAge: 0,
	}
	sess.Values["token"] = ""
	sess.Save(ctx.Request(), ctx.Response())
	return ctx.Success(nil)
}

func GetSelf(c echo.Context) error {
	ctx := c.(*context.Context)

	return ctx.Success(ctx.AdminUser)
}

func GetUser(c echo.Context) error {
	ctx := c.(*context.Context)

	id := ctx.IntParam(`id`)

	user, err := admindb.GetAdminUser(id)
	if err != nil {
		return ctx.InternalServerError()
	} else if user == nil {
		return ctx.NotFound()
	}
	return ctx.Success(user)
}

type FindAdminUsersRequest struct {
	Query    string `json:"query" form:"query" query:"query"`
	Status   string `json:"status" form:"status" query:"status"`
	Page     int    `json:"page" form:"page" query:"page"`
	PageSize int    `json:"page_size" form:"page_size" query:"page_size"`
}

func FindAdminUsers(c echo.Context) error {
	ctx := c.(*context.Context)

	req := FindAdminUsersRequest{}
	if err := ctx.BindAndValidate(&req); err != nil {
		return ctx.BadRequest()
	}

	req.Page, req.PageSize = x.Pagination(req.Page, req.PageSize)

	users, total, err := admindb.FindAdminUsers(req.Query, req.Status, req.Page, req.PageSize)
	if err != nil {
		return ctx.InternalServerError()
	}
	return ctx.List(users, req.Page, req.PageSize, total)
}

type CreateAdminUserRequest struct {
	Name     string `json:"name" form:"name" query:"name" validate:"gt=0"`
	Username string `json:"username" form:"username" query:"username" validate:"gt=0"`
	Phone    string `json:"phone" form:"phone" query:"phone"`
	Password string `json:"password" form:"password" query:"password"`
}

func CreateAdminUser(c echo.Context) error {
	ctx := c.(*context.Context)
	req := CreateAdminUserRequest{}
	if err := ctx.BindAndValidate(&req); err != nil {
		return ctx.BadRequest()
	}

	if !x.ValidateUsername(req.Username) {
		return ctx.BadRequest()
	}

	user, err := admindb.CreateAdminUser(req.Username, req.Name, req.Phone, req.Password, ctx.AdminUser.ID)
	if err != nil {
		return ctx.InternalServerError()
	}
	return ctx.Success(user)
}

type UpdateAdminUserRequest struct {
	Name   string `json:"name" form:"name" query:"name" validate:"gt=0"`
	Status string `json:"status" form:"status" query:"status" validate:"gt=0"`
	Phone  string `json:"phone" form:"phone" query:"phone"`
}

func UpdateAdminUser(c echo.Context) error {
	ctx := c.(*context.Context)

	id := ctx.IntParam(`id`)
	req := UpdateAdminUserRequest{}
	if err := ctx.BindAndValidate(&req); err != nil {
		return ctx.BadRequest()
	}

	user, err := admindb.UpdateAdminUser(id, req.Name, req.Status, req.Phone)
	if err != nil {
		return ctx.InternalServerError()
	}
	return ctx.Success(user)
}
