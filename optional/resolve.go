package optional

import (
	"github.com/eitanha/goods/result"
)

func (o Ref[T]) OrZero() (v T) {
	if o.IsSpecified() {
		v = *o.value
	}
	return
}

func (o Ref[T]) OrNot() (v T, ok bool) {
	if o.IsSpecified() {
		v = *o.value
		ok = true
	}
	return
}

func (o Ref[T]) OrValue(value T) T {
	if o.IsSpecified() {
		return *o.value
	}
	return value
}

func (o Ref[T]) OrNil() *T {
	return o.value
}

func (o Ref[T]) Result() result.Result[T] {
	return result.OkIfNotNil(o.OrNil())
}

func (o Ref[T]) IfSpecified(action func(T)) {
	if o.IsSpecified() {
		action(*o.value)
	}
}
