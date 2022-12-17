package database

import (
	crudErrors "sandbox/go-crud/errors"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

func Validate(object any, ctx *fiber.Ctx) (ok bool, err error) {
	validate := validator.New()
	err = validate.Struct(object)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		if len(errors) > 0 {
			validationParams := make(map[string]string)
			for _, err := range errors {
				validationParams[err.Field()] = err.Tag()
			}
			return false, crudErrors.SendError(ctx, crudErrors.VALIDATION_ERROR, validationParams)
		}
	}
	return true, nil
}
