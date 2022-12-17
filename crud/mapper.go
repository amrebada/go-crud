package crud

import (
	"errors"
	"reflect"

	"github.com/amrebada/go-crud/slices"
)

func MapRelations(entity any, relations []string) (interface{}, error) {
	typeOfEntity := reflect.TypeOf(entity)
	if typeOfEntity.Kind() == reflect.Ptr {
		typeOfEntity = typeOfEntity.Elem()
	}
	if typeOfEntity.Kind() != reflect.Struct {
		return nil, errors.New("entity must be a pointer to a struct")
	}
	valueOfEntity := reflect.ValueOf(entity).Elem()
	result := make(map[string]interface{})
	for i := 0; i < valueOfEntity.NumField(); i++ {
		jsonTag := typeOfEntity.Field(i).Tag.Get("json")

		if CheckRelation(typeOfEntity.Field(i), relations) {
			result[jsonTag] = AssignEmptyObjectIfZero(valueOfEntity.Field(i))
		} else if ok := CheckSliceIfContainsStruct(typeOfEntity.Field(i)); ok && slices.Contains(relations, jsonTag) {
			result[jsonTag] = valueOfEntity.Field(i).Interface()
		} else if valueOfEntity.Field(i).Kind() != reflect.Slice && valueOfEntity.Field(i).Kind() != reflect.Struct {
			result[jsonTag] = valueOfEntity.Field(i).Interface()
		} else if CheckIfTypeModel(valueOfEntity.Field(i)) {
			model := valueOfEntity.Field(i)
			id := model.FieldByName("ID")
			createdAt := model.FieldByName("CreatedAt")
			updatedAt := model.FieldByName("UpdatedAt")
			result["id"] = id.Interface()
			result["createdAt"] = createdAt.Interface()
			result["updatedAt"] = updatedAt.Interface()

		}
	}
	return result, nil
}

func AssignEmptyObjectIfZero(field reflect.Value) interface{} {
	if field.IsZero() {
		return map[string]interface{}{}
	}
	return field.Interface()
}

func CheckIfTypeModel(field reflect.Value) (ok bool) {
	return field.Kind() == reflect.Struct && field.Type().Name() == "Model"
}

func CheckRelation(field reflect.StructField, relations []string) (ok bool) {
	return field.Type.Kind() == reflect.Struct && slices.Contains(relations, field.Tag.Get("json"))
}

func CheckSliceIfContainsStruct(field reflect.StructField) (ok bool) {
	return field.Type.Kind() == reflect.Slice && field.Type.Elem().Kind() == reflect.Struct
}

func MapRecords(records interface{}, relations []string) (interface{}, error) {
	typeOfRecords := reflect.TypeOf(records)
	valueOfRecords := reflect.ValueOf(records)
	if typeOfRecords.Kind() == reflect.Ptr {
		typeOfRecords = typeOfRecords.Elem()
		valueOfRecords = valueOfRecords.Elem()
	}
	if typeOfRecords.Kind() == reflect.Slice {
		result := make([]interface{}, valueOfRecords.Len())
		for i := 0; i < valueOfRecords.Len(); i++ {
			record := valueOfRecords.Index(i).Interface()
			if record == nil {
				continue
			}
			mappedRecord, err := MapRelations(record, relations)
			if err != nil {
				return nil, err
			}
			result[i] = mappedRecord
		}
		return result, nil
	}
	return nil, errors.New("records must be a slice")
}
