package database

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

const (
	DB_FILTERS_KEY        = "filters"
	DB_FILTERS_PARAMS_KEY = "filtersParams"
)

func addQueryFilter(query *gorm.DB, ctx *fiber.Ctx) (ok bool, err error) {
	filters, okFilter := ctx.Locals(DB_FILTERS_KEY).(string)
	filterParams, okValues := ctx.Locals(DB_FILTERS_PARAMS_KEY).([]interface{})
	if okFilter && filters != "" {
		if !okValues {
			filterParams = []interface{}{}
		}
		query.Where(filters, filterParams...)
	}
	return true, nil

}

func AssignFilters(ctx *fiber.Ctx, filters string, filterParams ...[]interface{}) {
	ctx.Locals(DB_FILTERS_KEY, filters)
	ctx.Locals(DB_FILTERS_PARAMS_KEY, filterParams)
}
