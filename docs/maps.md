# maps — the fluent map View[K,V]

[← goods index](README.md)

`View[K,V]` wraps a map for chaining. Builders: `From(m)`,
`FromSeq(iter.Seq2)`, `FromPairs(Pair… )`, `Invert(m)` (swaps keys and
values). `Pair[K,V]` is `{Key K; Value V}`.

```go
import "github.com/eitanha/goods/maps"

active := maps.From(sessions).
        WhereNot(session.Expired).
        SelectValues(func(_ string, s Session) string { return s.UserID })

role := maps.From(roles).Get(userID).Or(RoleGuest)
```

## Queries

- `Get(k) result.Result[V]` — `ErrKeyNotFound` when absent.
- `GetOr(k, def) V` — fallback when absent.
- `HasKey(k) bool`.
- `ContainsValue(v) bool` — equality via `any(a) == any(b)` (panics on
  non-comparable value types).
- `Exists(pred) bool` — any entry matches.
- `First(pred) result.Result[Pair]` — first match, `ErrNoItemFound`.
- `Aggregate(fold) result.Result[V]` — fold without seed,
  `ErrEmptyMap` when empty.
- `AggregateWith(fold, initial) Acc` — fold with seed.

## Filters — builds a new map

- `Where(pred)`.
- `WhereNot(pred)`.

## Removals — in place

- `Without(k)`.
- `WithoutKeys(k… )`.

## Projects

To slices:

- `Select(selector)` — entries to values.
- `SelectWhere(selector, pred)`.
- `SelectOk(selector)` — unwraps `result.Result`, skips errors.
- `SelectSpecified(selector)` — unwraps `optional.Ref`, skips
  unspecified.
- `SelectNotNil(selector)` — unwraps pointers, skips nil.
- `Keys()` — as a slice view.
- `Values()` — as a slice view.
- `Pairs()` — as a slice view of `Pair`.

To maps (each builds a new map):

- `SelectValues(mapper)` — remap values.
- `SelectKeys(mapper)` — remap keys.
- `SelectEntries(mapper)` — remap both.

## Mutators — careful, in place

- `Map(mapper)` — changes the **underlying map**; anything else
  holding that map sees the change.
- `MapWhere(mapper, pred)` — same, only where `pred` holds.
- `Merge(others… )` — copy every entry from `others` in place.

## Exits

- `ToMap()` — the underlying map.
- `OrNil()` — nil when empty.
- `Clone()` — shallow copy (new map, same values; `maps.Clone`).
- `Seq()` — `iter.Seq2[K, V]`.
- `Take(n)` — up to `n` entries, new map.
- `Sort(cmp) → slices.View[Pair[K,V]]`.

## Sentinel errors

`ErrNoItemFound`, `ErrEmptyMap`.
