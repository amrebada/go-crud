package crud

import (
	"fmt"

	"github.com/amrebada/go-crud/errors"

	"github.com/gofiber/fiber/v2"
)

type IEngineType interface {
	GetEntityName() string
	TableName() string
	GetInstance() IEngineType
	GetSliceOfInstances() []IEngineType
	GetSearchableFields() []string
	GetSortableFields() []string
	GetRelations() []string
	GetOneRecord(*fiber.Ctx) ([]string, interface{}, bool, error)
	GetAllRecords(*fiber.Ctx) ([]string, interface{}, PaginationDetails, bool, error)
	ParseIdFromParam(*fiber.Ctx) (bool, error)
	ParseFromBody(*fiber.Ctx) (bool, error)
	Validate(*fiber.Ctx) (bool, error)
	CreateRecord(*fiber.Ctx) (bool, error)
	UpdateRecord(*fiber.Ctx) (bool, error)
	DeleteRecord(*fiber.Ctx) (bool, error)
}

type Engine struct {
	entities []IEngineType
}

func NewEngine(entities ...IEngineType) *Engine {
	return &Engine{
		entities: entities,
	}
}

func (e *Engine) Generate(app *fiber.App) error {
	for _, entity := range e.entities {
		if err := GenerateRoutes(app, entity); err != nil {
			return err
		}
	}
	return nil
}

func GenerateRoutes(app *fiber.App, entity IEngineType) error {
	routes, err := GetRoutes(entity)
	if err != nil {
		return err
	}

	for _, route := range routes {
		app.Add(route.Method, route.Path, route.Handler)
	}
	return nil
}

func GetRoutes(entity IEngineType) ([]Route, error) {
	routes := GetDefaultRoutes(entity)
	routes = append(routes, routes...)
	return routes, nil
}

type Route struct {
	Method  string
	Path    string
	Handler fiber.Handler
}

func GetDefaultRoutes(entity IEngineType) []Route {

	entityName := entity.GetEntityName()
	return []Route{
		{
			Method:  "GET",
			Path:    fmt.Sprintf("/%s/", entityName),
			Handler: generateGetAll(entity),
		},
		{
			Method:  "GET",
			Path:    fmt.Sprintf("/%s/:id", entityName),
			Handler: generateGetOne(entity),
		},
		{
			Method:  "POST",
			Path:    fmt.Sprintf("/%s/", entityName),
			Handler: generateCreate(entity),
		},
		{
			Method:  "PUT",
			Path:    fmt.Sprintf("/%s/:id", entityName),
			Handler: generateUpdate(entity),
		},
		{
			Method:  "DELETE",
			Path:    fmt.Sprintf("/%s/:id", entityName),
			Handler: generateDelete(entity),
		},
	}
}

func generateGetAll[T IEngineType](entity T) fiber.Handler {
	return func(c *fiber.Ctx) error {
		relations, records, paginationDetails, ok, err := (entity).GetAllRecords(c)
		if !ok {
			return err
		}
		records, err = MapRecords(records, relations)
		if err != nil {
			errors.LogError(err)
			return fiber.NewError(fiber.StatusInternalServerError, "Error mapping records")
		}
		return c.JSON(fiber.Map{
			"type":       "GET All",
			"pagination": paginationDetails,
			"data":       records,
		})
	}
}

func generateGetOne[T IEngineType](entity T) fiber.Handler {
	return func(c *fiber.Ctx) error {
		instance := entity.GetInstance()
		ok, err := instance.ParseIdFromParam(c)
		if !ok {
			return err
		}

		relations, record, ok, err := instance.GetOneRecord(c)
		if !ok {
			return err
		}
		record, err = MapRelations(record, relations)
		if err != nil {
			errors.LogError(err)
			return fiber.NewError(fiber.StatusInternalServerError, "Error mapping relations")
		}
		return c.JSON(fiber.Map{
			"type": "GET One",
			"data": record,
		})
	}
}

func generateCreate[T IEngineType](entity T) fiber.Handler {
	return func(c *fiber.Ctx) error {
		instance := entity.GetInstance()
		ok, err := instance.ParseFromBody(c)
		if !ok {
			return err
		}

		ok, err = instance.Validate(c)
		if !ok {
			return err
		}

		ok, err = instance.CreateRecord(c)
		if !ok {
			return err
		}

		return c.JSON(fiber.Map{
			"type": "Create",
			"data": instance,
		})
	}
}

func generateUpdate[T IEngineType](entity T) fiber.Handler {
	return func(c *fiber.Ctx) error {
		instance := entity.GetInstance()

		ok, err := instance.ParseIdFromParam(c)
		if !ok {
			return err
		}

		ok, err = instance.ParseFromBody(c)
		if !ok {
			return err
		}

		ok, err = instance.Validate(c)
		if !ok {
			return err
		}

		ok, err = instance.UpdateRecord(c)
		if !ok {
			return err
		}

		return c.JSON(fiber.Map{
			"type": "Update",
			"data": instance,
		})
	}
}

func generateDelete[T IEngineType](entity T) fiber.Handler {
	return func(c *fiber.Ctx) error {
		instance := entity.GetInstance()

		ok, err := instance.ParseIdFromParam(c)
		if !ok {
			return err
		}

		ok, err = instance.DeleteRecord(c)
		if !ok {
			return err
		}

		return c.JSON(fiber.Map{
			"type": "Delete",
		})
	}
}
