package admin

import (
	"math"

	"gitlab.com/ykgk/kgo/controller/context"
	systemdb "gitlab.com/ykgk/kgo/db/system"

	"github.com/labstack/echo/v4"
)

func FindStandardScaleScoresByName(c echo.Context) error {
	ctx := c.(*context.Context)

	name := ctx.StringParam(`name`)

	scores, err := systemdb.FindStandardScaleScoresByName(name)
	if err != nil {
		return ctx.InternalServerError()
	}
	for _, s := range scores {
		if math.IsInf(s.Max, 1) {
			s.Max = -1
		}
	}
	return ctx.Success(scores)
}

type FindDistrictsRequest struct {
	Query    string `json:"query" form:"query" query:"query"`
	ParentID string `json:"parent_id" form:"parent_id" query:"parent_id"`
}

func FindDistricts(c echo.Context) error {
	ctx := c.(*context.Context)
	req := FindDistrictsRequest{}
	if err := ctx.Bind(&req); err != nil {
		return ctx.BadRequest()
	}

	districts, err := systemdb.FindDistricts(req.Query, req.ParentID)
	if err != nil {
		return ctx.InternalServerError()
	}
	return ctx.Success(districts)
}

func GetDistrict(c echo.Context) error {
	ctx := c.(*context.Context)

	id := ctx.StringParam(`id`)

	district, err := systemdb.GetDistrict(id)
	if err != nil {
		return ctx.InternalServerError()
	}
	return ctx.Success(district)
}
