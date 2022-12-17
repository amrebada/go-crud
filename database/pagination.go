package database

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetPaginationQuery[T any](model *T, ctx *fiber.Ctx, query *gorm.DB) (page int, limit int) {

	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil {
		page = 1
	}

	if page < 1 {
		page = 1
	}
	limit, err = strconv.Atoi(ctx.Query("limit"))
	if err != nil {
		limit = 10
	}
	if limit < 1 {
		limit = 10
	}

	query.Offset((page - 1) * limit).Limit(limit)
	return page, limit
}
