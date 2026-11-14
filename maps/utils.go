package maps

import (
	"iter"
	"maps"
	stdslices "slices"

	"github.com/eitanha/goods/optional"
	"github.com/eitanha/goods/result"
	"github.com/eitanha/goods/slices"
	"github.com/rotisserie/eris"
)

var (
	ErrNoItemFound = eris.New("no item found")
	ErrEmptyMap    = eris.New("map is empty")
)

type Pair[K comparable, V any] struct {
	Key   K
	Value V
}

type View[K comparable, V any] struct {
	m map[K]V
}

func From[K comparable, V any](m map[K]V) View[K, V] {
	return View[K, V]{m: m}
}

func FromSeq[K comparable, V any](seq iter.Seq2[K, V]) View[K, V] {
	m := map[K]V{}
	seq(func(key K, value V) bool {
		m[key] = value
		return true
	})
	return From(m)
}

func FromPairs[K comparable, V any](pairs ...Pair[K, V]) View[K, V] {
	m := make(map[K]V, len(pairs))
	for _, pair := range pairs {
		m[pair.Key] = pair.Value
	}
	return From(m)
}

func Invert[K, V comparable](m map[K]V) View[V, K] {
	inverted := make(map[V]K, len(m))
	for key, value := range m {
		inverted[value] = key
	}
	return From(inverted)
}

func equal[V any](a V, b V) bool {
	return any(a) == any(b)
}

func (v View[K, V]) First(predicate func(K, V) bool) result.Result[Pair[K, V]] {
	for key, value := range v.m {
		if predicate(key, value) {
			return result.Ok(Pair[K, V]{Key: key, Value: value})
		}
	}
	return result.Err[Pair[K, V]](ErrNoItemFound)
}

func (v View[K, V]) Where(predicate func(K, V) bool) View[K, V] {
	found := map[K]V{}
	for key, value := range v.m {
		if predicate(key, value) {
			found[key] = value
		}
	}
	return From(found)
}

func (v View[K, V]) WhereNot(predicate func(K, V) bool) View[K, V] {
	return v.Where(func(key K, value V) bool {
		return !predicate(key, value)
	})
}

func (v View[K, V]) Without(key K) View[K, V] {
	delete(v.m, key)
	return v
}

func (v View[K, V]) WithoutKeys(keys ...K) View[K, V] {
	for _, key := range keys {
		delete(v.m, key)
	}
	return v
}

func (v View[K, V]) Select[To any](selector func(K, V) To) slices.View[To] {
	selected := make([]To, 0, len(v.m))
	for key, value := range v.m {
		selected = append(selected, selector(key, value))
	}
	return slices.From(selected)
}

func (v View[K, V]) SelectWhere[To any](selector func(K, V) To, predicate func(K, V) bool) slices.View[To] {
	selected := make([]To, 0, len(v.m))
	for key, value := range v.m {
		if predicate(key, value) {
			selected = append(selected, selector(key, value))
		}
	}
	return slices.From(selected)
}

func (v View[K, V]) SelectOk[To any](selector func(K, V) result.Result[To]) slices.View[To] {
	selected := make([]To, 0, len(v.m))
	for key, value := range v.m {
		selector(key, value).Do(func(value To) {
			selected = append(selected, value)
		})
	}
	return slices.From(selected)
}

func (v View[K, V]) SelectSpecified[To any](selector func(K, V) optional.Ref[To]) slices.View[To] {
	selected := make([]To, 0, len(v.m))
	for key, value := range v.m {
		selector(key, value).IfSpecified(func(value To) {
			selected = append(selected, value)
		})
	}
	return slices.From(selected)
}

func (v View[K, V]) SelectNotNil[To any](selector func(K, V) *To) slices.View[To] {
	selected := make([]To, 0, len(v.m))
	for key, value := range v.m {
		if selectedValue := selector(key, value); selectedValue != nil {
			selected = append(selected, *selectedValue)
		}
	}
	return slices.From(selected)
}

func (v View[K, V]) SelectValues[NewV any](mapper func(K, V) NewV) View[K, NewV] {
	mapped := make(map[K]NewV, len(v.m))
	for key, value := range v.m {
		mapped[key] = mapper(key, value)
	}
	return From(mapped)
}

func (v View[K, V]) SelectKeys[NewK comparable](mapper func(K, V) NewK) View[NewK, V] {
	mapped := make(map[NewK]V, len(v.m))
	for key, value := range v.m {
		mapped[mapper(key, value)] = value
	}
	return From(mapped)
}

func (v View[K, V]) SelectEntries[NewK comparable, NewV any](mapper func(K, V) (NewK, NewV)) View[NewK, NewV] {
	mapped := make(map[NewK]NewV, len(v.m))
	for key, value := range v.m {
		newKey, newValue := mapper(key, value)
		mapped[newKey] = newValue
	}
	return From(mapped)
}

func (v View[K, V]) Keys() slices.View[K] {
	keys := make([]K, 0, len(v.m))
	for key := range v.m {
		keys = append(keys, key)
	}
	return slices.From(keys)
}

func (v View[K, V]) Values() slices.View[V] {
	values := make([]V, 0, len(v.m))
	for _, value := range v.m {
		values = append(values, value)
	}
	return slices.From(values)
}

func (v View[K, V]) Pairs() slices.View[Pair[K, V]] {
	pairs := make([]Pair[K, V], 0, len(v.m))
	for key, value := range v.m {
		pairs = append(pairs, Pair[K, V]{Key: key, Value: value})
	}
	return slices.From(pairs)
}

func (v View[K, V]) Get(key K) result.Result[V] {
	value, ok := v.m[key]
	if !ok {
		return result.Err[V](result.ErrKeyNotFound).Wrapf("key %v not found in map", key)
	}
	return result.Ok(value)
}

func (v View[K, V]) GetOr(key K, fallback V) V {
	if value, ok := v.m[key]; ok {
		return value
	}
	return fallback
}

func (v View[K, V]) Map(mapper func(K, V) V) View[K, V] {
	for key, value := range v.m {
		v.m[key] = mapper(key, value)
	}
	return v
}

func (v View[K, V]) MapWhere(mapper func(K, V) V, predicate func(K, V) bool) View[K, V] {
	for key, value := range v.m {
		if predicate(key, value) {
			v.m[key] = mapper(key, value)
		}
	}
	return v
}

func (v View[K, V]) Aggregate(aggregator func(V, V) V) result.Result[V] {
	if len(v.m) == 0 {
		return result.Err[V](ErrEmptyMap)
	}
	var aggregated V
	first := true
	for _, value := range v.m {
		if first {
			aggregated = value
			first = false
			continue
		}
		aggregated = aggregator(value, aggregated)
	}
	return result.Ok(aggregated)
}

func (v View[K, V]) AggregateWith[Accumulator any](aggregator func(K, V, Accumulator) Accumulator, initial Accumulator) Accumulator {
	aggregated := initial
	for key, value := range v.m {
		aggregated = aggregator(key, value, aggregated)
	}
	return aggregated
}

func (v View[K, V]) ForEach(consumer func(K, V)) View[K, V] {
	for key, value := range v.m {
		consumer(key, value)
	}
	return v
}

func (v View[K, V]) Exists(predicate func(K, V) bool) bool {
	return v.First(predicate).IsOk()
}

func (v View[K, V]) HasKey(key K) bool {
	_, ok := v.m[key]
	return ok
}

func (v View[K, V]) ContainsValue(value V) bool {
	return v.Exists(func(_ K, candidate V) bool {
		return equal(candidate, value)
	})
}

func (v View[K, V]) OrNil() map[K]V {
	if len(v.m) == 0 {
		return nil
	}
	return v.m
}

func (v View[K, V]) Take(n int) View[K, V] {
	if n <= 0 {
		return View[K, V]{}
	}
	if n > len(v.m) {
		return v
	}
	taken := make(map[K]V, n)
	for key, value := range v.m {
		if len(taken) == n {
			break
		}
		taken[key] = value
	}
	return From(taken)
}

func (v View[K, V]) Merge(others ...map[K]V) View[K, V] {
	for _, other := range others {
		for key, value := range other {
			v.m[key] = value
		}
	}
	return v
}

func (v View[K, V]) Clone() View[K, V] {
	return From(maps.Clone(v.m))
}

func (v View[K, V]) Sort(cmp func(a Pair[K, V], b Pair[K, V]) int) slices.View[Pair[K, V]] {
	pairs := v.Pairs().Slice()
	stdslices.SortFunc(pairs, cmp)
	return slices.From(pairs)
}

func (v View[K, V]) Seq() iter.Seq2[K, V] {
	return maps.All(v.m)
}

func (v View[K, V]) ToMap() map[K]V {
	return v.m
}
