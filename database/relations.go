package database

import (
	"strings"

	"github.com/amrebada/go-crud/crud"
	crudErrors "github.com/amrebada/go-crud/errors"
	"github.com/amrebada/go-crud/slices"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetRelationsQuery[T any](model T, entity crud.IEngineType, ctx *fiber.Ctx, query *gorm.DB) (relationsInGo []string, ok bool, err error) {
	relationsQuery := ctx.Query("with")
	if relationsQuery != "" {
		if len(relationsQuery) > 100 {
			return relationsInGo, false, crudErrors.SendError(ctx, crudErrors.INVALID_RELATIONS_ERROR, map[string]string{
				"server": "relations query is invalid",
			})
		}
		relations := strings.Split(relationsQuery, ",")
		for _, relation := range relations {
			if !slices.Contains(entity.GetRelations(), relation) {
				return relationsInGo, false, crudErrors.SendError(ctx, crudErrors.INVALID_RELATIONS_ERROR, map[string]string{
					"server": "relation is not supported",
				})
			}
			query.Preload(ConvertJsonTagToColumnName(model, relation))
			relationsInGo = append(relationsInGo, relation)
		}
	}
	return relationsInGo, true, nil
}
