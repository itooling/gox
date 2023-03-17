package oth

import "reflect"

func CopyStruct(dst, src interface{}) {
	dstValue := reflect.ValueOf(dst).Elem()
	srcValue := reflect.ValueOf(src).Elem()

	for i := 0; i < srcValue.NumField(); i++ {
		srcField := srcValue.Field(i)
		srcName := srcValue.Type().Field(i).Name
		dstFieldByName := dstValue.FieldByName(srcName)

		if dstFieldByName.IsValid() {
			switch dstFieldByName.Kind() {
			case reflect.Ptr:
				switch srcField.Kind() {
				case reflect.Ptr:
					if srcField.IsNil() {
						dstFieldByName.Set(reflect.New(dstFieldByName.Type().Elem()))
					} else {
						dstFieldByName.Set(srcField)
					}
				default:
					dstFieldByName.Set(srcField.Addr())
				}
			default:
				switch srcField.Kind() {
				case reflect.Ptr:
					if srcField.IsNil() {
						dstFieldByName.Set(reflect.Zero(dstFieldByName.Type()))
					} else {
						dstFieldByName.Set(srcField.Elem())
					}
				default:
					dstFieldByName.Set(srcField)
				}
			}
		}
	}
}
