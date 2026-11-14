package async

type Panicable interface {
	Panic(panic any)
}

func RecoverInto(panicables ...Panicable) {
	if r := recover(); r != nil {
		for _, panicable := range panicables {
			panicable.Panic(r)
		}
	}
}
