package database

import (
	"reflect"
)

func ConvertJsonTagToColumnName[T any](model T, jsonName string) string {
	modelType := reflect.TypeOf(model)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}
	if modelType.Kind() != reflect.Struct {
		return ""
	}
	for i := 0; i < modelType.NumField(); i++ {
		if modelType.Field(i).Tag.Get("json") == jsonName {
			return modelType.Field(i).Name
		}
	}
	return ""
}

func ConvertColumnNameToJsonTag[T any](model T, columnName string) string {
	modelType := reflect.TypeOf(model)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}
	if modelType.Kind() != reflect.Struct {
		return ""
	}
	for i := 0; i < modelType.NumField(); i++ {
		if modelType.Field(i).Name == columnName {
			return modelType.Field(i).Tag.Get("json")
		}
	}
	return ""
}

func MapJsonRelations[T any](relations []string, model T) []string {
	result := make([]string, len(relations))
	for _, relation := range relations {
		result = append(result, ConvertJsonTagToColumnName(model, relation))
	}
	return result
}
