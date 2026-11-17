# async — futures, tasks, and pools

[← goods index](README.md)

A small futures kit. Vocabulary:

- `Waitable` — `Done() <-chan struct{}`; the minimal awaitable.
- `Future[T]` — Waitable + `Await() T`. **`Await` re-panics** anything
  the task recorded; every `async.Go*` starter already recovers the
  goroutine's panic into the task, so nothing crashes the process —
  the panic surfaces at the await point.
- `Task[T]` — a Future you can complete yourself:
  - `Set(value)` — complete successfully.
  - `Panic(p)` — complete with a panic; aborts waiters.
  - Both are once-only — first call wins.
- Type aliases: `Result[T] = Future[result.Result[T]]`,
  `Void = Future[void.Type]`, `Error = Future[error]`,
  `Promise[T] = Task[result.Result[T]]` (make one with
  `MakePromise[T]()`).
- `IsDone(w) bool` — non-blocking check.

## Starters — return the future immediately

```go
f := async.Go(func() int { … })             // routineView: f.Safely()
g := async.VoidGo(func() { … })             // f.Safely() → Error
h := async.GoResult(func() result.Result[int] { … })
e := async.GoError(func() error { … })
t := async.MakeTask[Token](); async.GoFor(t, fetch)  // complete t yourself
async.GoAndRecover(fn, panicables…)         // bare goroutine, panics routed
```

## Guards — panic becomes an error at the await point

- `Guard(f)` — lazily await as `result.Result[T]`.
- `GuardResult(f)` — same for `Result` futures, flattened.
- `GuardError(f)` — same for `Error` futures.
- `GuardVoid(f)` — same for `Void` futures.
- `await.Safely(f)` — the eager one-shot: await + guard.
- `await.ResultSafely(f)` — eager, flattened.
- `await.ErrorSafely(f)` — eager for `Error` futures.

## Cancellation and hand-completion

- `task.Panic(p)` — records the panic and closes `Done()`; waiters see
  it re-panicked at `Await`.
- `RecoverInto(panicables… )` — deferred helper: routes a recovered
  panic into every listed task (the starters use it).
- `Send(ch, v) Task[void.Type]` — send as a task; selects on its own
  completion.
- `Receive(ch) Task[T]` — receive as a task; same shape.
- Completing a `Send`/`Receive` task yourself (e.g. `task.Panic(err)`)
  cancels a stuck operation.

## await — groups and safe waiting

- `Add(futures… )` — extend the group.
- `All()` — collect every value; a panic in any future aborts the
  collection early.
- `AllSequentially()` — await in order.
- `First()` — first to finish wins.

Free functions:

- `All(futures… ) []T`.
- `ForFirstOf(futures… ) T`.
- `AllSequentially(futures… )`.
- `UntilAllDone(waitables… )`.
- `DoneGroup{…}.All()` — wait on bare `Waitable`s.

## after — continuations on a completed future

```go
after.Finishing(load()).
        Map(validate).    // runs in a new goroutine after load
        Do(publish).
        Catch(func(p any) { log })
```

- `Do(fn)` — run after completion; the value passes through.
- `Map(fn)` — transform the value.
- `Recover(fn)` — replace a panic with a value (callback sees the raw
  panic).
- `RecoverError(fn)` — same, with the panic already converted to an
  error (`errorutils.FromPanic`).
- `Catch(fn)` — observe a panic and drop it (zero value through).

## workerpool — bounded concurrent execution

`workerpool.Make(ctx, options… )` — bounded pool; `Run(job) async.Error`
queues and returns once queued (error when ctx is done).

Options:

- `Workers(n)` — pool size, default 16.
- `BufferSize(n)` — queue capacity, default 64.

```go
pool := workerpool.Make(ctx, workerpool.Workers(8))
res := pool.Go(func() int { … })            // async.Result[int]
err := pool.GoError(func() error { … }).Await()
```

`Go`/`GoResult`/`GoError` run the job on the pool with panic recovery;
a panicking job surfaces at the await point — wrap with
`await.Safely(...)` to get it as an eris error.

### worker — keepalive

`worker.Go(ctx, fn, options… )` reruns `fn` every retry interval until
it finishes without panicking, stops when ctx is done.

- `WaitBeforeRetry(d)` — retry interval, default 1s.

This is what powers the workerpool's own workers.
