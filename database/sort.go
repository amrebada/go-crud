package database

import (
	"errors"
	"fmt"
	"strings"

	"github.com/amrebada/go-crud/crud"
	crudErrors "github.com/amrebada/go-crud/errors"
	"github.com/amrebada/go-crud/slices"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetSortQuery[T any](model *T, entity crud.IEngineType, ctx *fiber.Ctx, query *gorm.DB) (ok bool, err error) {
	sortQuery := ctx.Query("sort")
	if sortQuery != "" {

		if len(sortQuery) > 100 {
			return false, crudErrors.SendError(ctx, crudErrors.INVALID_SORT_ERROR, nil)
		}
		fields := entity.GetSortableFields()
		if len(fields) == 0 {
			fields = crud.DefaultSortableFields
		}
		sortParts := strings.Split(sortQuery, ":")
		if len(sortParts) != 2 {
			return false, crudErrors.SendError(ctx, crudErrors.INVALID_SORT_ERROR, map[string]string{
				"server": "sort query is invalid",
			})
		}
		field := sortParts[0]
		sort := strings.ToUpper(sortParts[1])
		if !slices.Contains(fields, field) {
			return false, crudErrors.SendError(ctx, crudErrors.INVALID_SORT_ERROR, map[string]string{
				"server": "sort field is not supported",
			})
		}

		if !slices.Contains(crud.SortDirection, crud.SortOrder(sort)) {
			crudErrors.LogError(errors.New("sort direction is not supported"))
			return false, crudErrors.SendError(ctx, crudErrors.INVALID_SORT_ERROR, map[string]string{
				"server": "sort direction is not supported",
			})
		}

		query.Order(fmt.Sprintf("%s %s", field, sort))
	}
	return true, nil
}
