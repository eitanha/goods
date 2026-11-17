# slices — the fluent slice View[T]

[← goods index](README.md)

`View[T]` wraps a slice for chaining. Builders: `From(slice)`,
`FromSeq(seq)`, `SelectSeq(seq, selector)`, `Range(start, end)`
(end-exclusive), plus package-level `Flatten(View[[]T])`.

```go
import "github.com/eitanha/goods/slices"

names := slices.From(users).
        Where(user.Active).
        Select(user.Name).
        SortFunc(cmp.Compare[string])
```

## Queries

- `First(pred) result.Result[T]` — first match, `ErrNoItemFound`.
- `Exists(pred) bool` — any match.
- `Contains(item) bool` — equality via `any(a) == any(b)` (panics on
  non-comparable element types).
- `Aggregate(fold) result.Result[T]` — fold without seed,
  `ErrEmptySlice` when empty.
- `AggregateWith(fold, initial) Acc` — fold with seed.

## Filters — each builds a new slice

- `Where(pred)`.
- `WhereNot(pred)`.
- `Without(item)`.
- `Distinct()` — same `any` equality as `Contains`.
- `DistinctBy(equals)`.
- `Take(n)`.

## Projects — each builds a new slice

- `Select(selector)`.
- `SelectWhere(selector, pred)`.
- `SelectOk(selector)` — unwraps `result.Result`, skips errors.
- `SelectSpecified(selector)` — unwraps `optional.Ref`, skips
  unspecified.
- `SelectNotNil(selector)` — unwraps pointers, skips nil.

## Sorting & ordering

- `Sort(less func(i, j int) bool)` — index-comparator like
  `sort.Slice`; unstable; rebuilds the view.
- `SortFunc(cmp func(a, b T) int)` — sorts in place.
- `Reverse()` — rebuilds the view.

## Mutators — careful, in place

- `Map(mapper)` — changes the **underlying slice**; anything else
  holding that slice sees the change.
- `MapWhere(mapper, pred)` — same, only where `pred` holds.

## Exits

- `Slice() []T` — the underlying slice.
- `OrNil()` — nil when empty.
- `ForEach(f)` — side effects; returns the view for chaining.

## Sentinel errors

`ErrNoItemFound`, `ErrEmptySlice`.
