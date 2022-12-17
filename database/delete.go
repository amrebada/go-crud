package database

import (
	crudErrors "github.com/amrebada/go-crud/errors"

	"github.com/gofiber/fiber/v2"
)

func DeleteRecord(model any, ctx *fiber.Ctx) (ok bool, err error) {
	if !CheckPointerType(model) {
		err = crudErrors.SendError(ctx, crudErrors.DELETE_RECORD_ERROR, nil)
		crudErrors.LogError(err)
		return false, err
	}

	query := DB.Instance

	addQueryFilter(query, ctx)

	if err := query.Delete(model).Error; err != nil {
		crudErrors.LogError(err)
		err = crudErrors.SendError(ctx, crudErrors.DELETE_RECORD_ERROR, map[string]string{"db": err.Error()})
		return false, err
	}
	return true, nil
}
