package controller

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"gitlab.com/ykgk/kgo/controller/admin"
	"gitlab.com/ykgk/kgo/controller/device"
	"gitlab.com/ykgk/kgo/controller/middleware"
	"gitlab.com/ykgk/kgo/controller/teacher"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func Register(e *echo.Echo) {
	e.Validator = &CustomValidator{validator: validator.New()}
	api := e.Group("/api", middleware.APIMiddleware)

	device.Register(api.Group("/device"))
	/************************ 教师后台接口 *********************/
	teacher.Register(api.Group("/teacher"))
	/************************ 内部管理系统接口 *********************/
	admin.Register(api.Group("/admin"))
}
