package slices

import (
	"iter"
	"slices"
	"sort"

	"github.com/eitanha/goods/optional"
	"github.com/eitanha/goods/result"
	"github.com/rotisserie/eris"
)

var (
	ErrNoItemFound = eris.New("no item found")
	ErrEmptySlice  = eris.New("slice is empty")
)

type View[T any] struct {
	slice []T
}

func From[T any](slice []T) View[T] {
	return View[T]{slice: slice}
}

func FromSeq[T any](seq iter.Seq[T]) View[T] {
	slice := []T{}
	seq(func(item T) bool {
		slice = append(slice, item)
		return true
	})
	return From(slice)
}

func SelectSeq[Item, To any](seq iter.Seq[Item], selector func(Item) To) View[To] {
	var selected []To
	seq(func(item Item) bool {
		selected = append(selected, selector(item))
		return true
	})
	return From(selected)
}

func Range(start int, end int) []int {
	if start >= end {
		return nil
	}
	numbers := make([]int, end-start)
	for i := range end - start {
		numbers[i] = i + start
	}
	return numbers
}

func Flatten[T any](view View[[]T]) View[T] {
	flattened := []T{}
	for _, slice := range view.slice {
		flattened = append(flattened, slice...)
	}
	return From(flattened)
}

func equal[T any](a T, b T) bool {
	return any(a) == any(b)
}

func contains[T any](slice []T, item T) bool {
	return slices.ContainsFunc(slice, func(candidate T) bool {
		return equal(candidate, item)
	})
}

func (v View[T]) First(predicate func(T) bool) result.Result[T] {
	index := slices.IndexFunc(v.slice, predicate)
	if index == -1 {
		return result.Err[T](ErrNoItemFound)
	}
	return result.Ok(v.slice[index])
}

func (v View[T]) Where(predicate func(T) bool) View[T] {
	found := []T{}
	for _, item := range v.slice {
		if predicate(item) {
			found = append(found, item)
		}
	}
	return From(found)
}

func (v View[T]) WhereNot(predicate func(T) bool) View[T] {
	return v.Where(func(item T) bool {
		return !predicate(item)
	})
}

func (v View[T]) Without(item T) View[T] {
	return v.WhereNot(func(candidate T) bool {
		return equal(candidate, item)
	})
}

func (v View[T]) Select[To any](selector func(T) To) View[To] {
	selected := make([]To, 0, len(v.slice))
	for _, item := range v.slice {
		selected = append(selected, selector(item))
	}
	return From(selected)
}

func (v View[T]) SelectWhere[To any](selector func(T) To, predicate func(T) bool) View[To] {
	selected := make([]To, 0, len(v.slice))
	for _, item := range v.slice {
		if predicate(item) {
			selected = append(selected, selector(item))
		}
	}
	return From(selected)
}

func (v View[T]) SelectOk[To any](selector func(T) result.Result[To]) View[To] {
	selected := make([]To, 0, len(v.slice))
	for _, item := range v.slice {
		selector(item).Do(func(value To) {
			selected = append(selected, value)
		})
	}
	return From(selected)
}

func (v View[T]) SelectSpecified[To any](selector func(T) optional.Ref[To]) View[To] {
	selected := make([]To, 0, len(v.slice))
	for _, item := range v.slice {
		selector(item).IfSpecified(func(value To) {
			selected = append(selected, value)
		})
	}
	return From(selected)
}

func (v View[T]) SelectNotNil[To any](selector func(T) *To) View[To] {
	selected := make([]To, 0, len(v.slice))
	for _, item := range v.slice {
		if value := selector(item); value != nil {
			selected = append(selected, *value)
		}
	}
	return From(selected)
}

func (v View[T]) Map(mapper func(T) T) View[T] {
	for i, item := range v.slice {
		v.slice[i] = mapper(item)
	}
	return v
}

func (v View[T]) MapWhere(mapper func(T) T, predicate func(T) bool) View[T] {
	for i, item := range v.slice {
		if predicate(item) {
			v.slice[i] = mapper(item)
		}
	}
	return v
}

func (v View[T]) Aggregate(aggregator func(T, T) T) result.Result[T] {
	if len(v.slice) == 0 {
		return result.Err[T](ErrEmptySlice)
	}
	return result.Ok(From(v.slice[1:]).AggregateWith(aggregator, v.slice[0]))
}

func (v View[T]) AggregateWith[Accumulator any](aggregator func(T, Accumulator) Accumulator, initial Accumulator) Accumulator {
	aggregated := initial
	for _, item := range v.slice {
		aggregated = aggregator(item, aggregated)
	}
	return aggregated
}

func (v View[T]) ForEach(consumer func(T)) View[T] {
	for _, item := range v.slice {
		consumer(item)
	}
	return v
}

func (v View[T]) Exists(predicate func(T) bool) bool {
	return v.First(predicate).IsOk()
}

func (v View[T]) Contains(item T) bool {
	return contains(v.slice, item)
}

func (v View[T]) Distinct() View[T] {
	unique := make([]T, 0, len(v.slice))
	for _, item := range v.slice {
		if !contains(unique, item) {
			unique = append(unique, item)
		}
	}
	return From(unique)
}

func (v View[T]) DistinctBy(equals func(T, T) bool) View[T] {
	unique := make([]T, 0, len(v.slice))
	for _, item := range v.slice {
		if !slices.ContainsFunc(unique, func(u T) bool {
			return equals(item, u)
		}) {
			unique = append(unique, item)
		}
	}
	return From(unique)
}

func (v View[T]) OrNil() []T {
	if len(v.slice) == 0 {
		return nil
	}
	return v.slice
}

func (v View[T]) Take(n int) View[T] {
	if n <= 0 {
		return View[T]{}
	}
	if n > len(v.slice) {
		return From(v.slice)
	}
	return From(v.slice[:n])
}

func (v View[T]) Sort(less func(int, int) bool) View[T] {
	indices := Range(0, len(v.slice))
	sort.Slice(indices, func(i, j int) bool {
		return less(indices[i], indices[j])
	})
	return From(indices).Select(func(index int) T {
		return v.slice[index]
	})
}

func (v View[T]) SortFunc(cmp func(a T, b T) int) View[T] {
	slices.SortFunc(v.slice, cmp)
	return v
}

func (v View[T]) Reverse() View[T] {
	reversed := make([]T, len(v.slice))
	copy(reversed, v.slice)
	slices.Reverse(reversed)
	return From(reversed)
}

func (v View[T]) Slice() []T {
	return v.slice
}
