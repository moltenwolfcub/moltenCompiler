package main

import "errors"

type Optional[T any] struct {
	value T
	has   bool
}

func NewOptional[T any]() Optional[T] {
	return Optional[T]{}
}

func ToOptional[T any](start T) Optional[T] {
	return Optional[T]{
		value: start,
		has:   true,
	}
}

func (o Optional[T]) HasValue() bool {
	return o.has
}

func (o *Optional[T]) SetValue(newVal T) {
	o.value = newVal
	o.has = true
}

func (o *Optional[T]) Clear() {
	o.value = *new(T)
	o.has = false
}

func (o Optional[T]) GetValue() (T, error) {
	if !o.has {
		return *new(T), errors.New("optional has no value:\n\t- check with HasValue() before trying to get a value")
	}
	return o.value, nil
}

func (o Optional[T]) MustGetValue() T {
	if !o.has {
		panic("optional has no value:\n\t- check with HasValue() before trying to get a value\n\t- use GetValue() to have the error returned for processing (E.G. catching)")
	}
	return o.value
}
