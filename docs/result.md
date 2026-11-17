# result — the Result[T] monad

[← goods index](README.md)

`Result[T]` holds either a value or an error — exactly one. Construct
with `Ok`, error out with `Err`; chain transforms; resolve at the edge.

```go
import "github.com/eitanha/goods/result"

r := result.AccessMap(m, key).          // Result[V]
        Wrapf("loading user %v", key).  // eris wrap on error only
        Map(parseRole).                 // Result[Role]
        Validate(role.NotBanned)        // predicate may error

role := r.OrZero()                      // or OrError(), OrElse(…)
if err := r.Err(); err != nil { … }
```

## Construct

- `Ok(v)` — wrap a value.
- `OkIfNotErr(v, err)` — lift a `(T, error)` pair (new in v0.1.5):
  non-nil `err` → error result with it, otherwise `Ok(v)`.
- `OkIfNotNil(ptr)` — wrap a pointer; nil → error state (`NoValueErr`).
- `OkIfNotNilWithErr(ptr, err)` — pointer-or-error fallback (new in
  v0.1.5): non-nil `ptr` wins and any `err` is ignored; nil `ptr` →
  error result carrying `err` (`NoValueErr` when that is nil too).
- `OkIfNotPanic(f)` — run `f`; a panic becomes an eris error
  ("recovered from panic: …") instead of unwinding.
- `Err(err)` — error result; a nil `err` becomes `NoValueErr`.
- `Errm(msg)` — error from a message (`eris.New`).
- `Errf(fmt, …)` — error from a format string (`eris.Errorf`).

## Check

- `IsOk()` — holds a value.
- `IsNotOk()` — holds an error.
- `IsErr()` — alias of `IsNotOk`.
- `Err() error` — the error, nil when Ok.

## Wrap — all no-ops on Ok results

- `Wrap(msg)` — eris wrap.
- `Wrapf(fmt, …)` — eris wrapf.
- `WithErr(err)` — replace the error outright.
- `ErrJoin(errs…)` — `errors.Join` onto the existing error.

## Flow

- `Do(f)` — run if Ok.
- `Validate(pred func(T) error)` — a predicate error poisons the result.
- `Filter(pred func(T) bool)` — false → `ErrFalsePredicate`.
- `Map(transform)` — `T → U`, keeps the error.
- `FlatMap(transform)` — `T → Result[U]`, keeps the error.
- `MapErr(f)` — remap the error.
- `If(cond)` — false → `ErrFalseCondition`.
- `IfNotErr(err)` — a non-nil error poisons the result.
- `Cast[U]()` — unchecked type assertion via `any`
  (`ErrTypeMismatch`).
- `Recover(f)` — supply a replacement value from the error.
- `TryRecover(f)` — supply a replacement result from the error.
- `Match(okFn, errFn) U` — both branches, one value out.
- `Flatten(Result[Result[T]])` — package-level.

## Resolve — back to idiomatic Go

- `OrError() (T, error)` — the main bridge.
- `Or(def)` — the value, or `def`.
- `OrZero()` — the value, or the zero `T`.
- `OrNil() *T` — pointer to the value, nil on error.
- `OrPanic()` — the value, or panics with the error.
- `OrElse(f)` — the value, or `f(err)`.

## Safe accessors

Each returns a `Result` carrying a sentinel error:

| Accessor | Fails with |
|---|---|
| `AccessSlice(s, i)` | `ErrIndexOutOfBounds` |
| `AccessMap(m, k)` | `ErrKeyNotFound` |
| `AccessPtr(p)` | `ErrNilPointer` |
| `AccessAs[T](v)` | `ErrTypeMismatch` |

## Sentinel errors

`NoValueErr`, `ErrFalsePredicate`, `ErrFalseCondition`,
`ErrIndexOutOfBounds`, `ErrKeyNotFound`, `ErrNilPointer`,
`ErrTypeMismatch` — all eris values, comparable with `errors.Is`.
