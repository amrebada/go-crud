package errors

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

type ApplicationErrorType int

const (
	PARSE_BODY_ERROR ApplicationErrorType = iota
	VALIDATION_ERROR
	CREATE_RECORD_ERROR
	UPDATE_RECORD_ERROR
	DELETE_RECORD_ERROR
	GET_RECORD_ERROR
	GET_RECORDS_ERROR
	PARSE_ID_ERROR
	INVALID_PAGE_ERROR
	INVALID_SORT_ERROR
	INVALID_SEARCH_QUERY_ERROR
	INVALID_RELATIONS_ERROR
)

var (
	ErrorMessages = map[ApplicationErrorType]string{
		PARSE_BODY_ERROR:           "Error parsing body",
		CREATE_RECORD_ERROR:        "Error creating record",
		VALIDATION_ERROR:           "Validation error",
		UPDATE_RECORD_ERROR:        "Error updating record",
		DELETE_RECORD_ERROR:        "Error deleting record",
		GET_RECORD_ERROR:           "Error getting record",
		GET_RECORDS_ERROR:          "Error getting all records",
		PARSE_ID_ERROR:             "Error parsing id",
		INVALID_PAGE_ERROR:         "Invalid page",
		INVALID_SORT_ERROR:         "Invalid sort",
		INVALID_SEARCH_QUERY_ERROR: "Invalid search query",
		INVALID_RELATIONS_ERROR:    "Invalid relations",
	}
)

func SendError(ctx *fiber.Ctx, errType ApplicationErrorType, params map[string]string) error {
	return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"error":   errType,
		"message": ErrorMessages[errType],
		"params":  params,
	})
}

func LogError(err error) {
	if err != nil {
		log.Println(err)
	}
}
