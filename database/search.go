package database

import (
	"fmt"
	"math"
	"reflect"
	"strings"

	"github.com/amrebada/go-crud/crud"
	crudErrors "github.com/amrebada/go-crud/errors"
	"github.com/amrebada/go-crud/slices"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetRecords[T any](model *T, entity crud.IEngineType, ctx *fiber.Ctx) (relations []string, records interface{}, paginationDetails crud.PaginationDetails, ok bool, err error) {
	models := reflect.New(reflect.SliceOf(reflect.TypeOf(model))).Interface()
	query := DB.Instance.Model(model)
	addQueryFilter(query, ctx)
	ok, err = GetSearchQuery(model, entity, ctx, query)
	if !ok {
		return relations, nil, crud.PaginationDetails{}, false, err
	}
	ok, err = GetSortQuery(model, entity, ctx, query)
	if !ok {
		return relations, nil, crud.PaginationDetails{}, false, err
	}
	total := new(int64)
	err = query.Count(total).Error
	if err != nil {
		crudErrors.LogError(err)
		err = crudErrors.SendError(ctx, crudErrors.GET_RECORDS_ERROR, map[string]string{"db": err.Error()})
		return relations, nil, crud.PaginationDetails{}, false, err
	}
	relations, ok, err = GetRelationsQuery(*model, entity, ctx, query)
	if !ok {
		crudErrors.LogError(err)
		return relations, nil, crud.PaginationDetails{}, false, err
	}
	page, limit := GetPaginationQuery(model, ctx, query)
	if err != nil {
		crudErrors.LogError(err)
		err = crudErrors.SendError(ctx, crudErrors.GET_RECORDS_ERROR, map[string]string{"db": err.Error()})
		return relations, nil, crud.PaginationDetails{}, false, err
	}
	if err := query.Find(models).Error; err != nil {
		crudErrors.LogError(err)
		err = crudErrors.SendError(ctx, crudErrors.GET_RECORDS_ERROR, map[string]string{"db": err.Error()})
		return relations, nil, crud.PaginationDetails{}, false, err
	}
	return relations, models, crud.PaginationDetails{
		Total: *total,
		Page:  page,
		Limit: limit,
		Pages: int64(math.Ceil(float64(*total) / float64(limit))),
	}, true, nil
}

func GetSearchQuery[T any](model *T, entity crud.IEngineType, ctx *fiber.Ctx, query *gorm.DB) (ok bool, err error) {
	searchQuery := ctx.Query("search")
	if searchQuery != "" {
		if len(searchQuery) > 100 {
			return false, crudErrors.SendError(ctx, crudErrors.INVALID_SEARCH_QUERY_ERROR, map[string]string{
				"server": "Search query is too long",
			})
		}
		supportedSearchFields := entity.GetSearchableFields()
		if len(supportedSearchFields) == 0 {
			supportedSearchFields = crud.DefaultSearchableFields
		}

		searchParts := strings.Split(searchQuery, "|:|")
		if len(searchParts) != 3 {
			return false, crudErrors.SendError(ctx, crudErrors.INVALID_SEARCH_QUERY_ERROR, map[string]string{
				"server": "search query is invalid",
			})
		}

		searchFieldQuery := strings.ToLower(searchParts[0])
		searchOperator := strings.ToLower(searchParts[1])
		searchValue := searchParts[2]
		if searchFieldQuery == "" || searchValue == "" || searchOperator == "" {
			return false, crudErrors.SendError(ctx, crudErrors.INVALID_SEARCH_QUERY_ERROR, map[string]string{
				"server": "search query is invalid",
			})
		}
		searchFields := strings.Split(searchFieldQuery, ",")

		if searchFieldQuery != "all" && !slices.ContainsSlice(supportedSearchFields, searchFields) {
			return false, crudErrors.SendError(ctx, crudErrors.INVALID_SEARCH_QUERY_ERROR, map[string]string{
				"server": "search field is not supported",
			})
		}
		if !slices.Contains(crud.SearchOperators, searchOperator) {
			return false, crudErrors.SendError(ctx, crudErrors.INVALID_SEARCH_QUERY_ERROR, map[string]string{
				"server": "search operator is not supported",
			})
		}

		if slices.Contains(searchFields, "all") {
			searchFields = supportedSearchFields
		}

		condition := ""

		for _, field := range searchFields {
			if condition != "" {
				condition += " OR "
			}
			condition += fmt.Sprintf("%s %s ?", field, getSearchOperator(searchOperator))
		}

		query.Where(condition, getSearchValueWith(len(searchFields), searchOperator, searchValue)...)
	}
	return true, nil
}

func getSearchValueWith(repetitiveNumber int, operator string, value string) (searchValue []interface{}) {
	sqlValue := value
	switch operator {
	case "contains":
		sqlValue = fmt.Sprintf("%%%s%%", value)
	case "startswith":
		sqlValue = fmt.Sprintf("%s%%", value)
	case "endswith":
		sqlValue = fmt.Sprintf("%%%s", value)
	}

	for i := 0; i < repetitiveNumber; i++ {
		searchValue = append(searchValue, sqlValue)
	}
	return searchValue
}

func getSearchOperator(operator string) (searchOperator string) {
	switch operator {
	case "exact":
		return "="
	case "contains", "startswith", "endswith":
		return "ILIKE"
	default:
		return "="

	}
}
