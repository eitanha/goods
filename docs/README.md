# goods — package reference

Reference for [`github.com/eitanha/goods`](https://github.com/eitanha/goods)
("Go but good") — a small, type-safe utility library for Go 1.27: a
`Result` monad, an optional type, fluent slice/map views, a few sugar
helpers, and a futures/workerpool kit.

- Added to this module in v0.1.2 (2026-11-14); currently on
  **v0.1.5** (2026-11-17). The module requires **go ≥ 1.27.4** (this is
  what bumped `backend/go.mod` to 1.27.4).
- Only transitive dependency: `github.com/rotisserie/eris` (error
  wrapping) — `result` errors are eris errors, so `Wrap`/`Wrapf` chains
  and `%+v` stacks work out of the box.
- Import paths are the subpackages (`result`, `optional`, …), never the
  root — there is no root API.
- Docs written from the v0.1.5 source; signatures are paraphrased, not
  quoted. When in doubt, read the module cache:
  `$(go env GOMODCACHE)/github.com/eitanha/goods@v0.1.5`.

## Package map

The small components are documented on this page; the larger ones have
their own page.

| Package | Purpose | Docs |
|---|---|---|
| `result` | `Result[T]` — value-or-error monad with chaining | [result.md](result.md) |
| `optional` | `Ref[T]` — nullable value with JSON marshaling | [below](#optional) |
| `void` | `Type` struct{} + `Func` adapter for generic contexts | [below](#void) |
| `slices` | `View[T]` — fluent slice operations | [slices.md](slices.md) |
| `maps` | `View[K,V]` — fluent map operations | [maps.md](maps.md) |
| `sugar/conditionally` | ternary expression: `Use(a).If(c).ElseUse(b)` | [below](#sugar) |
| `sugar/pointer` | `From(v) *T` — pointer literal | [below](#sugar) |
| `errorutils` | panic value → error conversion (`FromPanic`) | [below](#errorutils) |
| `async` | futures/tasks: `Go`, guards, cancellable channel ops | [async.md](async.md) |
| `async/await` | awaiting groups of futures, safely | [async.md](async.md) |
| `async/after` | continuation chains on a future | [async.md](async.md) |
| `async/workerpool` | bounded worker pool with panic recovery | [async.md](async.md) |
| `async/workerpool/worker` | self-healing (keepalive) workers | [async.md](async.md) |

## void

`void.Type` is `struct{}`; `void.Func(f)` lifts `func()` to
`func() void.Type` so no-value jobs fit the generic `Go`/task APIs:

```go
async.VoidGo(func() { cleanup() })   // VoidGo uses void.Func internally
```

That is the whole package — it exists so a `Task[TValue]`-shaped API
never needs a special no-value variant.

## sugar

Two one-purpose helpers:

- `conditionally.Use(a).If(cond).ElseUse(b)` — an expression-shaped
  ternary. Both branches are values, so they evaluate eagerly — fine for
  cheap operands, wrong for side effects:

  ```go
  level := conditionally.Use(LevelWarn).If(strict).ElseUse(LevelDebug)
  ```

- `pointer.From(v) *T` — "pointer to literal" for composite literals
  and call sites that take a pointer:

  ```go
  cfg := Config{Workers: pointer.From(8), Name: pointer.From("api")}
  ```

## optional

`Ref[T]` is a nullable `T` that marshals as its value (or `null`) for
JSON — useful for config structs that must distinguish "unset" from
zero.

```go
import "github.com/eitanha/goods/optional"

type Config struct {
    Workers optional.Ref[int] // omitted/empty → unspecified
}

w := cfg.Workers.OrValue(4)        // default when unspecified
cfg.Workers.IfSpecified(func(n int) { log.Printf("workers=%d", n) })
```

- **Construct**:
  - `Empty[T]()` — an unspecified ref.
  - `Make(v)` — a specified ref holding `v`.
  - `From(ptr)` — adopts an existing pointer, nil-aware.
  - `FromResult(r)` — from a `result.Result[T]`; error → unspecified.
- **Check**:
  - `IsSpecified()` — holds a value.
  - `IsntSpecified()` — empty.
  - `IsZero()` — alias of `IsntSpecified`.
- **Resolve**:
  - `OrValue(def)` — the value, or `def` when unspecified.
  - `OrZero()` — the value, or the zero `T`.
  - `OrNot() (T, bool)` — comma-ok form.
  - `OrNil() *T` — the raw pointer.
  - `Result()` — as a `result.Result[T]`.
  - `IfSpecified(action)` — run with the value when set.
- Implements `json.Marshaler/Unmarshaler` and the **goccy**
  `MarshalYAML`/`UnmarshalYAML` interfaces. Caution: this backend's
  config uses `gopkg.in/yaml.v3`, which knows neither method — under
  yaml.v3 a `Ref` field marshals as `{}`. Use JSON for wire formats;
  for yaml.v3 structs prefer plain pointers or resolve manually.

## errorutils

One function, new in v0.1.3:

- `FromPanic(v any) error` — turns a recovered panic value into an
  error: `nil` → `nil`, an `error` passes through unchanged, anything
  else becomes an eris error (`"recovered from panic: …"`).

It is the machinery behind `result.OkIfNotPanic` and `async`'s panic
recovery; use it directly when you `recover()` yourself and want the
value as an error.

## The larger components

- **[result.md](result.md)** — the `Result[T]` monad: construct, wrap,
  chain (`Map`/`FlatMap`/…), resolve back to idiomatic Go, plus the
  safe slice/map/pointer accessors and their sentinel errors.
- **[slices.md](slices.md)** — `View[T]`, the fluent slice wrapper:
  filter/project/aggregate/sort, with the in-place vs rebuilt
  distinction spelled out.
- **[maps.md](maps.md)** — `View[K,V]`, the fluent map wrapper: keyed
  access, projection to slices, inversion, merging.
- **[async.md](async.md)** — the futures kit: `Go` starters, panic
  guards, cancellable channel ops, await groups, continuations
  (`after`), and the bounded `workerpool` with keepalive workers.

## Using it in this backend

- Until something imports it, `goods` stays `// indirect` in
  `backend/go.mod`; the first import promotes it and a later
  `go mod tidy` would otherwise drop it.
- Error style is eris end to end: `Wrap`/`Wrapf` on results, `errors.Is`
  against the sentinels each package defines, `%+v` for stacks.
- Panic discipline: goroutines started outside `async.Go*` (bare `go`,
  library callbacks) need their own `recover`; inside the kit, panics
  become task state — guard the await, don't let it escape a handler.
