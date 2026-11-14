# goods

Small, generic Go utility packages: a fluent slice `View`, `result`/`optional` types, and concurrency helpers.

Requires **Go 1.27+** (uses generic methods).

[![CI](https://github.com/eitanha/goods/actions/workflows/ci.yml/badge.svg)](https://github.com/eitanha/goods/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/eitanha/goods.svg)](https://pkg.go.dev/github.com/eitanha/goods)

## Install

```
go get github.com/eitanha/goods
```

## slices

A fluent `View` over any slice, with chainable operations:

```go
names := slices.From(users).
    Where(func(u User) bool { return u.Active }).
    Select(func(u User) string { return u.Name }).
    Distinct().
    Take(10).
    Slice()

first := slices.From(users).First(func(u User) bool { return u.Admin }) // result.Result[User]
total := slices.From(numbers).AggregateWith(func(n, acc int) int { return n + acc }, 0)
```

Operations come in a few flavors:

- Transforming (chainable, returns a new `View`): `Where`, `WhereNot`, `Without`, `Select`, `SelectWhere`, `SelectOk`, `SelectSpecified`, `SelectNotNil`, `Distinct`, `DistinctBy`, `Take`, `Sort`, `Reverse`, `Flatten`
- Mutating in place: `Map`, `MapWhere`, `SortFunc`
- Terminal: `Slice`, `OrNil`, `First`, `Exists`, `Contains`, `Aggregate`, `AggregateWith`, `ForEach`

Iterators work too: `slices.FromSeq(seq)` and `slices.SelectSeq(seq, selector)`.

## result

`result.Result[T]` carries either a value or an error, with fluent combinators:

```go
res := result.AccessSlice(items, index).Wrapf("accessing item")

res.Do(func(item Item) { fmt.Println(item) })

name := res.Match(func(item Item) string {
    return item.Name
}, func(err error) string {
    return "unknown"
})
```

Extract with `OrZero`, `Or`, `OrElse`, `OrNil`, `OrPanic`, or `OrError`. Transform with `Map`, `FlatMap`, `Filter`, `Validate`, `Recover`.

## optional

`optional.Ref[T]` is a JSON/YAML-friendly optional value:

```go
port := optional.Make(8080).OrValue(3000)
```

## async

Goroutines that return typed futures instead of using channels by hand:

```go
value := async.Go(func() int { return compute() }).Await()

// panics inside the goroutine become result errors
res := await.Safely(async.Go(func() string { return risky() })) // result.Result[string]

// run many futures concurrently, or sequentially
results := await.All(
    async.Go(func() int { return fetchA() }),
    async.Go(func() int { return fetchB() }),
)
```

`async.Task` gives you a settable future you can complete yourself, and `async/workerpool` provides bounded worker pools.

## sugar

Tiny syntax helpers:

```go
p := pointer.From(42)                          // *int
v := conditionally.Use("fast").If(cond).ElseUse("slow")
```

## License

[MIT](LICENSE)
