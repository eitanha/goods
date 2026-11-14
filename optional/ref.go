package optional

import (
	"encoding/json"

	"github.com/eitanha/goods/result"
)

type Ref[T any] struct {
	value *T
}

func Empty[T any]() Ref[T] {
	return Ref[T]{value: nil}
}

func Make[T any](value T) Ref[T] {
	return Ref[T]{value: &value}
}

func From[T any](value *T) Ref[T] {
	return Ref[T]{value: value}
}

func FromResult[T any](result result.Result[T]) Ref[T] {
	return Ref[T]{value: result.OrNil()}
}

func (o Ref[T]) IsZero() bool {
	return o.value == nil
}

func (o Ref[T]) IsSpecified() bool {
	return o.value != nil
}

func (o Ref[T]) IsntSpecified() bool {
	return !o.IsSpecified()
}

func (o Ref[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.value)
}

func (o *Ref[T]) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &o.value)
}

func (o Ref[T]) MarshalYAML() (any, error) {
	return o.value, nil
}

func (o *Ref[T]) UnmarshalYAML(unmarshal func(any) error) error {
	return unmarshal(&o.value)
}
