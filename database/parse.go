package database

import (
	crudErrors "github.com/amrebada/go-crud/errors"

	"github.com/gofiber/fiber/v2"
)

func ParseFromBody(model interface{}, ctx *fiber.Ctx) (ok bool, err error) {
	if !CheckPointerType(model) {
		err = crudErrors.SendError(ctx, crudErrors.PARSE_BODY_ERROR, nil)
		crudErrors.LogError(err)
		return false, err
	}
	if err := ctx.BodyParser(model); err != nil {
		crudErrors.LogError(err)
		err = crudErrors.SendError(ctx, crudErrors.PARSE_BODY_ERROR, nil)
		return false, err
	}
	return true, nil
}

func ParseIdFromParam(id *string, ctx *fiber.Ctx) (ok bool, err error) {
	if id == nil {
		err = crudErrors.SendError(ctx, crudErrors.PARSE_ID_ERROR, nil)
		crudErrors.LogError(err)
		return false, err
	}
	*id = ctx.Params("id")
	return true, nil
}
