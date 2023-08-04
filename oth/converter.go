package oth

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/petersunbag/coven"
)

var (
	mtx   sync.Mutex
	c_map = make(map[string]*coven.Converter)
)

func Convert(dst, src any) (err error) {
	key := fmt.Sprintf("%v_%v", reflect.TypeOf(src).String(), reflect.TypeOf(dst).String())
	if _, ok := c_map[key]; !ok {
		mtx.Lock()
		defer mtx.Unlock()
		if c_map[key], err = coven.NewConverter(dst, src); err != nil {
			return
		}
	}
	if err = c_map[key].Convert(dst, src); err != nil {
		return
	}
	return
}
