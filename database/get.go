package database

import (
	"reflect"

	"github.com/amrebada/go-crud/crud"
	crudErrors "github.com/amrebada/go-crud/errors"

	"github.com/gofiber/fiber/v2"
)

func GetRecord[T any](model T, entity crud.IEngineType, id string, ctx *fiber.Ctx) (relations []string, record interface{}, ok bool, err error) {
	record = reflect.New(reflect.TypeOf(model)).Interface()

	query := DB.Instance.Where("id = ?", id)

	addQueryFilter(query, ctx)
	relations, ok, err = GetRelationsQuery(model, entity, ctx, query)
	if !ok {
		crudErrors.LogError(err)
		return relations, nil, false, err
	}
	if err := query.First(record).Error; err != nil {
		crudErrors.LogError(err)
		return relations, nil, false, err
	}
	return relations, record, true, nil
}
