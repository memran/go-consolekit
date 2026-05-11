package console

import (
	"testing"
)

func TestCollectFromSlice(t *testing.T) {
	c := Collect([]any{1, 2, 3})
	if c.Count() != 3 {
		t.Fatalf("expected 3, got %d", c.Count())
	}
}

func TestCollectFromTyped(t *testing.T) {
	c := CollectFrom([]int{1, 2, 3})
	if c.Count() != 3 {
		t.Fatalf("expected 3, got %d", c.Count())
	}
}

func TestAll(t *testing.T) {
	c := Collect([]any{1, 2, 3})
	result := c.All()
	if len(result) != 3 || result[0] != 1 {
		t.Fatalf("expected [1 2 3], got %v", result)
	}
}

func TestAllIsolate(t *testing.T) {
	original := []any{1, 2, 3}
	c := Collect(original)
	original[0] = 99
	if c.All()[0] == 99 {
		t.Fatal("Collect should copy the slice")
	}
}

func TestCount(t *testing.T) {
	if Collect([]any{}).Count() != 0 {
		t.Fatal("expected 0 for empty")
	}
	if Collect([]any{1, 2}).Count() != 2 {
		t.Fatal("expected 2")
	}
}

func TestIsEmpty(t *testing.T) {
	if !Collect([]any{}).IsEmpty() {
		t.Fatal("expected empty")
	}
	if Collect([]any{1}).IsEmpty() {
		t.Fatal("expected not empty")
	}
}

func TestIsNotEmpty(t *testing.T) {
	if Collect([]any{}).IsNotEmpty() {
		t.Fatal("expected empty")
	}
	if !Collect([]any{1}).IsNotEmpty() {
		t.Fatal("expected not empty")
	}
}

func TestFirst(t *testing.T) {
	if Collect([]any{}).First() != nil {
		t.Fatal("expected nil for empty")
	}
	if Collect([]any{10, 20}).First() != 10 {
		t.Fatal("expected 10")
	}
}

func TestLast(t *testing.T) {
	if Collect([]any{}).Last() != nil {
		t.Fatal("expected nil for empty")
	}
	if Collect([]any{10, 20}).Last() != 20 {
		t.Fatal("expected 20")
	}
}

func TestGet(t *testing.T) {
	c := Collect([]any{"a", "b", "c"})
	if c.Get(0) != "a" || c.Get(1) != "b" || c.Get(2) != "c" {
		t.Fatal("Get failed")
	}
	if c.Get(-1) != nil || c.Get(99) != nil {
		t.Fatal("expected nil for out of range")
	}
}

func TestMap(t *testing.T) {
	c := Collect([]any{1, 2, 3})
	result := c.Map(func(v any) any { return v.(int) * 2 })
	all := result.All()
	if all[0] != 2 || all[1] != 4 || all[2] != 6 {
		t.Fatalf("expected [2 4 6], got %v", all)
	}
}

func TestMapImmutability(t *testing.T) {
	c := Collect([]any{1, 2, 3})
	c.Map(func(v any) any { return v.(int) * 2 })
	if c.All()[0] != 1 {
		t.Fatal("Map should not modify original")
	}
}

func TestFilter(t *testing.T) {
	c := Collect([]any{1, 2, 3, 4, 5})
	result := c.Filter(func(v any) bool { return v.(int) > 2 })
	all := result.All()
	if len(all) != 3 || all[0] != 3 || all[1] != 4 || all[2] != 5 {
		t.Fatalf("expected [3 4 5], got %v", all)
	}
}

func TestFilterEmpty(t *testing.T) {
	result := Collect([]any{1, 2}).Filter(func(v any) bool { return false })
	if result.Count() != 0 {
		t.Fatal("expected empty after filter all")
	}
}

func TestReject(t *testing.T) {
	c := Collect([]any{1, 2, 3, 4, 5})
	result := c.Reject(func(v any) bool { return v.(int) <= 2 })
	all := result.All()
	if len(all) != 3 || all[0] != 3 || all[1] != 4 || all[2] != 5 {
		t.Fatalf("expected [3 4 5], got %v", all)
	}
}

func TestEach(t *testing.T) {
	c := Collect([]any{1, 2, 3})
	var sum int
	c.Each(func(v any) { sum += v.(int) })
	if sum != 6 {
		t.Fatalf("expected 6, got %d", sum)
	}
}

func TestReduce(t *testing.T) {
	c := Collect([]any{1, 2, 3, 4, 5})
	sum := c.Reduce(func(carry, v any) any { return carry.(int) + v.(int) }, 0)
	if sum != 15 {
		t.Fatalf("expected 15, got %d", sum)
	}
}

func TestReduceEmpty(t *testing.T) {
	sum := Collect([]any{}).Reduce(func(carry, v any) any { return carry.(int) + v.(int) }, 42)
	if sum != 42 {
		t.Fatalf("expected 42, got %d", sum)
	}
}

func TestSlice(t *testing.T) {
	c := Collect([]any{1, 2, 3, 4, 5})
	all := c.Slice(1, 3).All()
	if len(all) != 3 || all[0] != 2 || all[1] != 3 || all[2] != 4 {
		t.Fatalf("expected [2 3 4], got %v", all)
	}
}

func TestSliceOutOfRange(t *testing.T) {
	if Collect([]any{1}).Slice(10, 5).Count() != 0 {
		t.Fatal("expected empty for out-of-range slice")
	}
}

func TestSort(t *testing.T) {
	c := Collect([]any{3, 1, 2})
	result := c.Sort(func(a, b any) bool { return a.(int) < b.(int) })
	all := result.All()
	if all[0] != 1 || all[1] != 2 || all[2] != 3 {
		t.Fatalf("expected [1 2 3], got %v", all)
	}
}

func TestReverse(t *testing.T) {
	c := Collect([]any{1, 2, 3})
	all := c.Reverse().All()
	if all[0] != 3 || all[1] != 2 || all[2] != 1 {
		t.Fatalf("expected [3 2 1], got %v", all)
	}
}

func TestReverseEmpty(t *testing.T) {
	if Collect([]any{}).Reverse().Count() != 0 {
		t.Fatal("expected empty")
	}
}

func TestUnique(t *testing.T) {
	c := Collect([]any{1, 2, 2, 3, 3, 3})
	all := c.Unique().All()
	if len(all) != 3 {
		t.Fatalf("expected 3 unique, got %d: %v", len(all), all)
	}
}

func TestValues(t *testing.T) {
	c := Collect([]any{10, 20, 30})
	all := c.Values().All()
	if all[0] != 10 || all[1] != 20 || all[2] != 30 {
		t.Fatal("Values should return items")
	}
}

func TestKeys(t *testing.T) {
	c := Collect([]any{"a", "b", "c"})
	all := c.Keys().All()
	if all[0] != 0 || all[1] != 1 || all[2] != 2 {
		t.Fatalf("expected [0 1 2], got %v", all)
	}
}

func TestChunk(t *testing.T) {
	c := Collect([]any{1, 2, 3, 4, 5})
	chunks := c.Chunk(2).All()
	if len(chunks) != 3 {
		t.Fatalf("expected 3 chunks, got %d", len(chunks))
	}
	first := chunks[0].([]any)
	if first[0] != 1 || first[1] != 2 {
		t.Fatalf("first chunk [1 2], got %v", first)
	}
}

func TestChunkEmpty(t *testing.T) {
	if Collect([]any{}).Chunk(3).Count() != 0 {
		t.Fatal("expected empty")
	}
}

func TestCollapse(t *testing.T) {
	c := Collect([]any{[]any{1, 2}, []any{3, 4}})
	all := c.Collapse().All()
	if len(all) != 4 || all[0] != 1 || all[3] != 4 {
		t.Fatalf("expected [1 2 3 4], got %v", all)
	}
}

func TestMerge(t *testing.T) {
	c := Collect([]any{1, 2}).Merge([]any{3, 4})
	all := c.All()
	if len(all) != 4 || all[2] != 3 || all[3] != 4 {
		t.Fatalf("expected [1 2 3 4], got %v", all)
	}
}

func TestMergeEmpty(t *testing.T) {
	c := Collect([]any{1, 2}).Merge([]any{})
	if c.Count() != 2 {
		t.Fatal("expected 2")
	}
}

func TestDiff(t *testing.T) {
	c := Collect([]any{1, 2, 3, 4}).Diff([]any{2, 4})
	all := c.All()
	if len(all) != 2 || all[0] != 1 || all[1] != 3 {
		t.Fatalf("expected [1 3], got %v", all)
	}
}

func TestIntersect(t *testing.T) {
	c := Collect([]any{1, 2, 3, 4}).Intersect([]any{2, 4, 6})
	all := c.All()
	if len(all) != 2 || all[0] != 2 || all[1] != 4 {
		t.Fatalf("expected [2 4], got %v", all)
	}
}

func TestContains(t *testing.T) {
	c := Collect([]any{"a", "b", "c"})
	if !c.Contains("b") {
		t.Fatal("expected to contain 'b'")
	}
	if c.Contains("z") {
		t.Fatal("expected not to contain 'z'")
	}
}

func TestSearch(t *testing.T) {
	c := Collect([]any{"a", "b", "c"})
	if c.Search("b") != 1 {
		t.Fatalf("expected 1, got %d", c.Search("b"))
	}
	if c.Search("z") != -1 {
		t.Fatalf("expected -1, got %d", c.Search("z"))
	}
}

func TestWhere(t *testing.T) {
	c := Collect([]any{
		map[string]any{"name": "Alice", "role": "admin"},
		map[string]any{"name": "Bob", "role": "user"},
		map[string]any{"name": "Charlie", "role": "admin"},
	})
	result := c.Where("role", "admin")
	if result.Count() != 2 {
		t.Fatalf("expected 2, got %d", result.Count())
	}
}

func TestPluck(t *testing.T) {
	c := Collect([]any{
		map[string]any{"id": 1, "name": "Alice"},
		map[string]any{"id": 2, "name": "Bob"},
	})
	result := c.Pluck("name")
	all := result.All()
	if len(all) != 2 || all[0] != "Alice" || all[1] != "Bob" {
		t.Fatalf("expected [Alice Bob], got %v", all)
	}
}

func TestSum(t *testing.T) {
	c := Collect([]any{1.0, 2.0, 3.0})
	if c.Sum() != 6.0 {
		t.Fatalf("expected 6, got %f", c.Sum())
	}
}

func TestSumIntegers(t *testing.T) {
	c := Collect([]any{1, 2, 3})
	if c.Sum() != 6.0 {
		t.Fatalf("expected 6, got %f", c.Sum())
	}
}

func TestSumWithKey(t *testing.T) {
	c := Collect([]any{
		map[string]any{"price": 10.0},
		map[string]any{"price": 20.0},
	})
	if c.Sum("price") != 30.0 {
		t.Fatalf("expected 30, got %f", c.Sum("price"))
	}
}

func TestAvg(t *testing.T) {
	c := Collect([]any{2.0, 4.0, 6.0})
	if c.Avg() != 4.0 {
		t.Fatalf("expected 4, got %f", c.Avg())
	}
}

func TestAvgEmpty(t *testing.T) {
	if Collect([]any{}).Avg() != 0 {
		t.Fatal("expected 0 for empty")
	}
}

func TestMin(t *testing.T) {
	c := Collect([]any{3.0, 1.0, 2.0})
	if c.Min() != 1.0 {
		t.Fatalf("expected 1, got %f", c.Min())
	}
}

func TestMax(t *testing.T) {
	c := Collect([]any{3.0, 1.0, 2.0})
	if c.Max() != 3.0 {
		t.Fatalf("expected 3, got %f", c.Max())
	}
}

func TestImplode(t *testing.T) {
	c := Collect([]any{"a", "b", "c"})
	if c.Implode(",") != "a,b,c" {
		t.Fatalf("expected 'a,b,c', got '%s'", c.Implode(","))
	}
}

func TestJoin(t *testing.T) {
	c := Collect([]any{"a", "b", "c"})
	if c.Join(", ", " and ") != "a, b and c" {
		t.Fatalf("expected 'a, b and c', got '%s'", c.Join(", ", " and "))
	}
}

func TestJoinSingle(t *testing.T) {
	if Collect([]any{"only"}).Join(",", " and ") != "only" {
		t.Fatal("expected 'only'")
	}
}

func TestJoinEmpty(t *testing.T) {
	if Collect([]any{}).Join(",", " and ") != "" {
		t.Fatal("expected ''")
	}
}

func TestToJSON(t *testing.T) {
	c := Collect([]any{1, 2, 3})
	json, err := c.ToJSON()
	if err != nil || json != "[1,2,3]" {
		t.Fatalf("expected '[1,2,3]', got '%s' err=%v", json, err)
	}
}

func TestFluentChaining(t *testing.T) {
	c := Collect([]any{1, 2, 3, 4, 5, 6}).
		Filter(func(v any) bool { return v.(int) > 2 }).
		Map(func(v any) any { return v.(int) * 2 })

	all := c.All()
	if len(all) != 4 || all[0] != 6 || all[1] != 8 || all[2] != 10 || all[3] != 12 {
		t.Fatalf("expected [6 8 10 12], got %v", all)
	}
}

func TestChainingWithReduce(t *testing.T) {
	sum := Collect([]any{1, 2, 3, 4, 5}).
		Filter(func(v any) bool { return v.(int) > 2 }).
		Map(func(v any) any { return v.(int) * 10 }).
		Reduce(func(carry, v any) any { return carry.(int) + v.(int) }, 0)

	if sum != 120 {
		t.Fatalf("expected 120, got %d", sum)
	}
}

func TestCollectionTap(t *testing.T) {
	var sideEffect []int
	c := Collect([]any{1, 2, 3}).
		Tap(func(col *Collection) {
			sideEffect = append(sideEffect, col.Count())
		}).
		Map(func(v any) any { return v.(int) * 2 })

	all := c.All()
	if len(all) != 3 || all[0] != 2 || all[1] != 4 || all[2] != 6 {
		t.Fatalf("Tap should forward chain, got %v", all)
	}
	if len(sideEffect) != 1 || sideEffect[0] != 3 {
		t.Fatalf("Tap side effect should run, got %v", sideEffect)
	}
}

func TestStandaloneTap(t *testing.T) {
	var sideEffect int
	result := Tap(42, func(v int) {
		sideEffect = v * 2
	})
	if result != 42 {
		t.Fatalf("Tap should return original value, got %d", result)
	}
	if sideEffect != 84 {
		t.Fatalf("Tap side effect should run, got %d", sideEffect)
	}
}

func TestStandaloneTapWithSlice(t *testing.T) {
	items := []int{1, 2, 3}
	result := Tap(items, func(v []int) {
		v[0] = 99
	})
	if result[0] != 99 {
		t.Fatalf("Tap should allow mutation, got %d", result[0])
	}
}

func TestCollectFromStrings(t *testing.T) {
	c := CollectFrom([]string{"a", "b", "c"})
	if c.Implode(",") != "a,b,c" {
		t.Fatalf("expected 'a,b,c', got '%s'", c.Implode(","))
	}
}
