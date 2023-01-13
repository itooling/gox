package oth

import (
	"fmt"
	"reflect"
	"unsafe"
)

type User struct {
	Name string `say:"hello world"`
	Age  int
}

func Ref() {
	u := User{"lulu", 18}
	t := reflect.TypeOf(u)
	f, _ := t.FieldByName("Name")
	tag := f.Tag
	fmt.Println(tag)
	fmt.Println(*(*string)(unsafe.Pointer(uintptr(unsafe.Pointer(&u)) + f.Offset)))

	ff, _ := t.FieldByName("Age")
	fmt.Println(*(*int)(unsafe.Pointer(uintptr(unsafe.Pointer(&u)) + ff.Offset)))
}
