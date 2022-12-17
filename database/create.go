package database

import (
	crudErrors "sandbox/go-crud/errors"

	"github.com/gofiber/fiber/v2"
)

func CreateRecord(model any, ctx *fiber.Ctx) (ok bool, err error) {
	if !CheckPointerType(model) {
		err = crudErrors.SendError(ctx, crudErrors.CREATE_RECORD_ERROR, nil)
		crudErrors.LogError(err)
		return false, err
	}
	if err := DB.instance.Create(model).Error; err != nil {
		crudErrors.LogError(err)
		err = crudErrors.SendError(ctx, crudErrors.CREATE_RECORD_ERROR, map[string]string{"db": err.Error()})
		return false, err
	}
	return true, nil
}
