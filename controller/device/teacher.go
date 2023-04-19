package device

import (
	"time"

	"github.com/labstack/echo/v4"
	"gitlab.com/ykgk/kgo/controller/context"
	"gitlab.com/ykgk/kgo/controller/errors"
	kindergartendb "gitlab.com/ykgk/kgo/db/kindergarten"
)

type KindergartenTeacherLoginRequest struct {
	Username string `json:"username" form:"username" query:"username" validate:"gt=0"`
	Password string `json:"password" form:"password" query:"password" validate:"gt=0"`
}

func KindergartenTeacherLogin(c echo.Context) error {
	ctx := c.(*context.Context)
	req := KindergartenTeacherLoginRequest{}
	if err := ctx.BindAndValidate(&req); err != nil {
		return ctx.BadRequest()
	}

	teacher, err := kindergartendb.GetKindergartenTeacherByUsername(req.Username)
	if err != nil {
		return ctx.InternalServerError()
	} else if teacher == nil || teacher.Deleted || teacher.Role != kindergartendb.KindergartenTeacherRoleManager { /* 用户不存在 */
		/* 只有园长可以登录 */
		return ctx.Fail(errors.ErrKindergartenTeacherNotFound, nil)
	}

	if !teacher.Auth(req.Password) { /* 密码错误 */
		return ctx.Fail(errors.ErrKindergartenTeacherPasswordIncorrect, nil)
	}

	token, err := kindergartendb.CreateKindergartenTeacherToken(teacher.ID, time.Time{}, ctx.Request().UserAgent(), ctx.RealIP())
	if err != nil {
		return ctx.InternalServerError()
	}

	type Response struct {
		*kindergartendb.KindergartenTeacher
		Token string `json:"token" form:"token"`
	}

	r := &Response{
		KindergartenTeacher: teacher,
		Token:               token.ID,
	}

	return ctx.Success(r)
}

func FindKindergartenClasses(c echo.Context) error {
	ctx := c.(*context.Context)

	classes, _, err := kindergartendb.FindKindergartenClasses("", ctx.Teacher.KindergartenID, 0, 0)
	if err != nil {
		return ctx.InternalServerError()
	}
	return ctx.Success(classes)
}

func FindKindergartenClassStudents(c echo.Context) error {
	ctx := c.(*context.Context)

	id := ctx.IntParam(`id`)

	students, _, err := kindergartendb.FindKindergartenStudents("", "", ctx.Teacher.KindergartenID, id, 0, 0)
	if err != nil {
		return ctx.InternalServerError()
	}
	return ctx.Success(students)
}
