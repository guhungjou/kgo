package middleware

import (
	"github.com/labstack/echo-contrib/session"
	"gitlab.com/ykgk/kgo/controller/context"
)

func GetSessionToken(ctx *context.Context, name string) string {
	sess, _ := session.Get(name, ctx)
	if sess == nil {
		return ""
	}
	v, ok := sess.Values["token"]
	if !ok || v == nil {
		return ""
	}
	return v.(string)
}
