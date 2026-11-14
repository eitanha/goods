package result

type Result[T any] struct {
	value *T
	err   error
}
