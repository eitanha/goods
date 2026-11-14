package void

type Type struct{}

func Func(f func()) func() Type {
	return func() (v Type) {
		f()
		return
	}
}
