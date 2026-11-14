package async

import (
	"github.com/eitanha/goods/result"
)

type Promise[TResult any] = Task[result.Result[TResult]]

func MakePromise[TResult any]() Promise[TResult] {
	return MakeTask[result.Result[TResult]]()
}
