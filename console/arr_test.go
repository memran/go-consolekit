package console

import "testing"

func TestArrGet(t *testing.T) {
	a := NewArr(10, 20, 30)
	if a.Get(0) != 10 || a.Get(2) != 30 {
		t.Fatal("get failed")
	}
	if a.Get(-1) != nil || a.Get(99) != nil {
		t.Fatal("get out of range failed")
	}
}

func TestArrFirstLast(t *testing.T) {
	a := NewArr(1, 2, 3)
	if a.First() != 1 {
		t.Fatal("first failed")
	}
	if a.Last() != 3 {
		t.Fatal("last failed")
	}
	if NewArr().First() != nil {
		t.Fatal("first empty failed")
	}
}

func TestArrHas(t *testing.T) {
	a := NewArr("a", "b", "c")
	if !a.Has("b") {
		t.Fatal("has failed")
	}
	if a.Has("z") {
		t.Fatal("has should be false")
	}
}

func TestArrContains(t *testing.T) {
	a := NewArr(1, 2, 3)
	if !a.Contains(2) {
		t.Fatal("contains failed")
	}
}

func TestArrWhere(t *testing.T) {
	a := NewArr(1, 2, 3, 4, 5)
	r := a.Where(func(v any) bool { return v.(int) > 2 })
	if len(r) != 3 {
		t.Fatalf("expected 3, got %d", len(r))
	}
}

func TestArrPluck(t *testing.T) {
	a := NewArr(
		map[string]any{"id": 1, "name": "Alice"},
		map[string]any{"id": 2, "name": "Bob"},
	)
	names := a.Pluck("name")
	if len(names) != 2 || names[0] != "Alice" || names[1] != "Bob" {
		t.Fatalf("pluck failed: %v", names)
	}
}

func TestArrFlatten(t *testing.T) {
	a := NewArr(1, NewArr(2, 3), NewArr(4, 5, 6))
	r := a.Flatten()
	if len(r) != 6 {
		t.Fatalf("expected 6, got %d %v", len(r), r)
	}
}

func TestArrCollapse(t *testing.T) {
	a := NewArr(NewArr(1, 2), NewArr(3, 4))
	r := a.Collapse()
	if len(r) != 4 {
		t.Fatalf("expected 4, got %d", len(r))
	}
}

func TestArrChunk(t *testing.T) {
	a := NewArr(1, 2, 3, 4, 5)
	chunks := a.Chunk(2)
	if len(chunks) != 3 {
		t.Fatalf("expected 3 chunks, got %d", len(chunks))
	}
}

func TestArrUnique(t *testing.T) {
	a := NewArr(1, 2, 2, 3, 3, 3)
	r := a.Unique()
	if len(r) != 3 {
		t.Fatalf("expected 3, got %d", len(r))
	}
}

func TestArrWrap(t *testing.T) {
	var a Arr
	r := a.Wrap("hello")
	if len(r) != 1 || r[0] != "hello" {
		t.Fatal("wrap failed")
	}
	r = a.Wrap([]any{1, 2, 3})
	if len(r) != 3 {
		t.Fatal("wrap slice failed")
	}
}

func TestArrJoin(t *testing.T) {
	a := NewArr("a", "b", "c")
	if a.Join(",") != "a,b,c" {
		t.Fatalf("expected 'a,b,c', got '%s'", a.Join(","))
	}
}

func TestArrToJSON(t *testing.T) {
	a := NewArr(1, 2, 3)
	json, err := a.ToJSON()
	if err != nil || json != "[1,2,3]" {
		t.Fatalf("expected '[1,2,3]', got '%s'", json)
	}
}
