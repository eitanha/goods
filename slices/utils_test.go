package slices_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/eitanha/goods/optional"
	"github.com/eitanha/goods/result"
	"github.com/eitanha/goods/slices"
)

func assert(t *testing.T, name string, got, want any) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("%s: got %v, want %v", name, got, want)
	}
}

func TestFluentChain(t *testing.T) {
	numbers := []int{5, 2, 8, 2, 1, 9, 4}

	got := slices.From(numbers).
		WhereNot(func(n int) bool { return n%2 == 0 }).
		Distinct().
		Sort(func(i, j int) bool { return i > j }).
		Select(func(n int) int { return n * 10 }).
		Take(2).
		Slice()
	assert(t, "chain", got, []int{90, 10})

	words := slices.From([]string{"banana", "apple", "cherry", "avocado"}).
		Where(func(s string) bool { return strings.HasPrefix(s, "a") }).
		Select(strings.ToUpper).
		Distinct().
		Slice()
	assert(t, "type-changing chain", words, []string{"APPLE", "AVOCADO"})
}

func TestOperations(t *testing.T) {
	numbers := []int{5, 2, 8, 2, 1, 9, 4}

	assert(t, "Where", slices.From(numbers).Where(func(n int) bool { return n > 4 }).Slice(), []int{5, 8, 9})
	assert(t, "WhereNot", slices.From(numbers).WhereNot(func(n int) bool { return n > 4 }).Slice(), []int{2, 2, 1, 4})
	assert(t, "Without", slices.From(numbers).Without(2).Slice(), []int{5, 8, 1, 9, 4})
	assert(t, "Distinct", slices.From(numbers).Distinct().Slice(), []int{5, 2, 8, 1, 9, 4})
	assert(t, "DistinctBy", slices.From(numbers).DistinctBy(func(a, b int) bool { return a%3 == b%3 }).Slice(), []int{5, 1, 9})

	assert(t, "Select", slices.From(numbers).Select(func(n int) int { return n / 2 }).Slice(), []int{2, 1, 4, 1, 0, 4, 2})
	assert(t, "SelectWhere", slices.From(numbers).SelectWhere(func(n int) string { return string(rune('A' + n)) }, func(n int) bool { return n > 8 }).Slice(), []string{"J"})
	assert(t, "SelectOk", slices.From(numbers).SelectOk(func(n int) result.Result[int] {
		if n%2 == 0 {
			return result.Ok(n)
		}
		return result.Err[int](slices.ErrNoItemFound)
	}).Slice(), []int{2, 8, 2, 4})
	assert(t, "SelectSpecified", slices.From(numbers).SelectSpecified(func(n int) optional.Ref[int] {
		if n > 4 {
			return optional.Make(n)
		}
		return optional.Empty[int]()
	}).Slice(), []int{5, 8, 9})
	assert(t, "SelectNotNil", slices.From(numbers).SelectNotNil(func(n int) *int {
		if n > 8 {
			return &n
		}
		return nil
	}).Slice(), []int{9})

	first := slices.From(numbers).First(func(n int) bool { return n > 7 })
	if !first.IsOk() || *first.OrNil() != 8 {
		t.Fatalf("First: %v", first)
	}
	if slices.From(numbers).First(func(n int) bool { return n > 100 }).IsOk() {
		t.Fatal("First miss should not be Ok")
	}

	sum := slices.From(numbers).Aggregate(func(a, b int) int { return a + b })
	if !sum.IsOk() || *sum.OrNil() != 31 {
		t.Fatalf("Aggregate: %v", sum)
	}
	if slices.From([]int{}).Aggregate(func(a, b int) int { return a + b }).IsOk() {
		t.Fatal("Aggregate on empty should not be Ok")
	}
	assert(t, "AggregateWith", slices.From(numbers).AggregateWith(func(item, acc int) int { return item + acc }, 100), 131)

	nums := []int{1, 2, 3}
	slices.From(nums).Map(func(n int) int { return n * 2 })
	assert(t, "Map in place", nums, []int{2, 4, 6})

	nums2 := []int{1, 2, 3}
	slices.From(nums2).MapWhere(func(n int) int { return 0 }, func(n int) bool { return n == 2 })
	assert(t, "MapWhere in place", nums2, []int{1, 0, 3})

	nums3 := []int{3, 1, 2}
	slices.From(nums3).SortFunc(func(a, b int) int { return a - b })
	assert(t, "SortFunc in place", nums3, []int{1, 2, 3})

	assert(t, "Take", slices.From(numbers).Take(3).Slice(), []int{5, 2, 8})
	if slices.From(numbers).Take(0).OrNil() != nil {
		t.Fatal("Take(0) should be nil")
	}
	assert(t, "Take over", slices.From([]int{1}).Take(5).Slice(), []int{1})

	if slices.From([]int{}).OrNil() != nil {
		t.Fatal("OrNil empty should be nil")
	}
	if slices.From(numbers).OrNil() == nil {
		t.Fatal("OrNil non-empty should not be nil")
	}

	assert(t, "Sort", slices.From([]int{3, 1, 2}).Sort(func(i, j int) bool { return i > j }).Slice(), []int{2, 1, 3})
	assert(t, "Reverse", slices.From([]int{1, 2, 3}).Reverse().Slice(), []int{3, 2, 1})

	assert(t, "Flatten", slices.Flatten(slices.From([][]int{{1, 2}, {3}})).Slice(), []int{1, 2, 3})

	seq := func(yield func(int) bool) {
		for i := 1; i <= 3; i++ {
			if !yield(i) {
				return
			}
		}
	}
	assert(t, "FromSeq", slices.FromSeq(seq).Select(func(n int) int { return n * 2 }).Slice(), []int{2, 4, 6})
	assert(t, "SelectSeq", slices.SelectSeq(seq, func(n int) int { return n + 1 }).Slice(), []int{2, 3, 4})
	assert(t, "Range", slices.Range(2, 5), []int{2, 3, 4})

	seen := []int{}
	slices.From(numbers).Where(func(n int) bool { return n > 4 }).ForEach(func(n int) { seen = append(seen, n) })
	assert(t, "ForEach", seen, []int{5, 8, 9})

	assert(t, "Exists", slices.From(numbers).Exists(func(n int) bool { return n == 9 }), true)
	assert(t, "Exists miss", slices.From(numbers).Exists(func(n int) bool { return n == 100 }), false)
	assert(t, "Contains", slices.From(numbers).Contains(4), true)
	assert(t, "Contains miss", slices.From(numbers).Contains(100), false)
}
