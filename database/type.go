package database

import "reflect"

func CheckPointerType(model any) (ok bool) {
	typeOfModel := reflect.TypeOf(model)
	return typeOfModel.Kind() == reflect.Ptr
}
