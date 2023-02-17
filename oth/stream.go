package oth

import (
	"errors"
)

type S[V any] []V

type Optional[T any] struct {
	Val S[T]
}

func (o *Optional[any]) Of(v []any) (*Optional[any], error) {
	if v != nil && len(v) > 0 {
		o.Val = v
		return o, nil
	}
	return nil, errors.New("data is empty")
}

func (o *Optional[any]) OfNilable(v []any) *Optional[any] {
	if v != nil && len(v) > 0 {
		o.Val = v
		return o
	}
	return nil
}

func (o *Optional[any]) IfPresent(f func(any)) {
	if o.Val != nil && len(o.Val) > 0 {
		for _, v := range o.Val {
			f(v)
		}
	}
}

func (o *Optional[any]) IsPresent() bool {
	if o.Val != nil && len(o.Val) > 0 {
		return true
	}
	return false
}

func (o *Optional[any]) ForEach(f func(any)) {
	for _, v := range o.Val {
		f(v)
	}
}
