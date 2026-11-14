# goods
Go but good
## What?

>`good` doesn't just not suck — its good. No bad, no yucks, just pure goodness, goodies and good.
<br>Today, it is common to see "vibe based" software. And honestly? — thats OK. But `good` doesn't just compromise — it delivers.

`good` is a small, type safe utility library for Go that makes working with Go *good*

## Install

```
go get github.com/eitanha/goods
```

Requires **Go 1.27+** (uses generic methods).

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

Iterators work too: `slices.FromSeq(seq)` and `slices.SelectSeq(seq, selector)`.

## maps

The same fluent style, for maps:

```go
ages := maps.From(users).
	Where(func(name string, u User) bool { return u.Active }).
	SelectValues(func(name string, u User) int { return u.Age }).
	ToMap()

first := maps.From(config).First(func(k, v string) bool { return strings.HasPrefix(k, "env_") })
total := maps.From(counts).Aggregate(func(a, b int) int { return a + b })
```

`Get`/`GetOr` for lookups, `Keys`/`Values`/`Pairs` for slice views, `Merge`, `Clone`, `Invert`, and `FromPairs`/`FromSeq` for building maps.

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

// panics inside the goroutine become result errors, right off the future
res := async.Go(func() string { return risky() }).Safely().Await() // result.Result[string]

err := async.VoidGo(func() { cleanup() }).Safely().Await() // error

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
p := pointer.From(42) // *int
v := conditionally.Use("fast").If(cond).ElseUse("slow")
```

## License

[MIT](LICENSE)
