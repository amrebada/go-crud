package database

import (
	"reflect"
	"sandbox/go-crud/crud"
	crudErrors "sandbox/go-crud/errors"

	"github.com/gofiber/fiber/v2"
)

func GetRecord[T any](model T, entity crud.IEngineType, id string, ctx *fiber.Ctx) (relations []string, record interface{}, ok bool, err error) {
	record = reflect.New(reflect.TypeOf(model)).Interface()

	query := DB.instance.Where("id = ?", id)
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
