package context

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	admindb "gitlab.com/ykgk/kgo/db/admin"
	kindergartendb "gitlab.com/ykgk/kgo/db/kindergarten"
	"gitlab.com/ykgk/kgo/x"
)

type Context struct {
	echo.Context

	AdminUser *admindb.AdminUser
	Teacher   *kindergartendb.KindergartenTeacher
}

func (c *Context) IntParam(key string) int64 {
	v := c.Param(key)

	i, _ := strconv.ParseInt(v, 10, 64)
	return i
}

func (c *Context) StringParam(key string) string {
	s, _ := url.PathUnescape(c.Param(key))
	return s
}

func (c *Context) IntFormParam(key string) int64 {
	v := c.FormValue(key)
	i, _ := strconv.ParseInt(v, 10, 64)
	return i
}

func (c *Context) Success(data interface{}) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": 0,
		"data":   data,
	})
}

func (c *Context) Fail(status int, data interface{}) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": status,
		"data":   data,
	})
}

func (c *Context) List(data interface{}, page, pageSize, total int, v ...interface{}) error {
	d := map[string]interface{}{
		"page":      page,
		"page_size": pageSize,
		"total":     total,
		"list":      data,
	}
	if len(v) > 0 {
		d["extra"] = v[0]
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": 0,
		"data":   d,
	})
}

func (c *Context) Bind(data interface{}) error {
	if err := c.Context.Bind(data); err != nil {
		return err
	}

	/* 过滤字符串的首尾空格 */
	v := reflect.ValueOf(data).Elem()
	for i := 0; i < v.NumField(); i++ {
		vv := v.Field(i)
		if vv.Kind() == reflect.String {
			vv.SetString(strings.TrimSpace(vv.String()))
		}
	}
	return nil
}

func (c *Context) BindAndValidate(data interface{}) error {
	if err := c.Bind(data); err != nil {
		return err
	}

	return c.Validate(data)
}

func (c *Context) BadRequest() error {
	return c.NoContent(http.StatusBadRequest)
}

func (c *Context) NotFound() error {
	return c.NoContent(http.StatusNotFound)
}

func (c *Context) InternalServerError() error {
	return c.NoContent(http.StatusInternalServerError)
}

func (c *Context) Unauthorized() error {
	return c.NoContent(http.StatusUnauthorized)
}

func (c *Context) Forbidden() error {
	return c.NoContent(http.StatusForbidden)
}

func (c *Context) XLSX(filename string, headers []string, rows [][]interface{}) error {

	wb, err := x.CreateXLSXBook(filename, headers, rows)
	if err != nil {
		return c.InternalServerError()
	}

	buf := bytes.NewBuffer(nil)
	if err := wb.Write(buf); err != nil {
		return err
	}

	c.Response().Header().Add("Content-Disposition", fmt.Sprintf("attachment;filename=%s.xlsx", filename))
	return c.Stream(200, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buf)
}
