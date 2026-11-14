package conditionally

func Use[T any](value T) use[T] {
	return use[T]{value: value}
}

type use[T any] struct {
	value T
}

func (u use[T]) If(condition bool) elseUse[T] {
	return elseUse[T]{
		condition: condition,
		value:     u.value,
	}
}

type elseUse[T any] struct {
	condition bool
	value     T
}

func (u elseUse[T]) ElseUse(elseValue T) T {
	if u.condition {
		return u.value
	}
	return elseValue
}
