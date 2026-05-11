package console

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

type TestEvent struct {
	name string
	data string
}

func (e TestEvent) Name() string { return e.name }

type OrderPlaced struct {
	OrderID int
}

func TestEventBusOnEmit(t *testing.T) {
	bus := NewEventBus()
	var called bool
	bus.On("user.registered", func(e any) error {
		called = true
		return nil
	})
	bus.Emit("user.registered", nil)
	if !called {
		t.Fatal("handler should have been called")
	}
}

func TestEventBusMultipleListeners(t *testing.T) {
	bus := NewEventBus()
	var mu sync.Mutex
	order := []string{}

	bus.On("event", func(e any) error {
		mu.Lock()
		order = append(order, "first")
		mu.Unlock()
		return nil
	})
	bus.On("event", func(e any) error {
		mu.Lock()
		order = append(order, "second")
		mu.Unlock()
		return nil
	})
	bus.Emit("event", nil)
	if len(order) != 2 {
		t.Fatalf("expected 2 calls, got %d", len(order))
	}
}

func TestEventBusPriority(t *testing.T) {
	bus := NewEventBus()
	order := []string{}

	bus.On("e", func(e any) error {
		order = append(order, "low")
		return nil
	}, 0)
	bus.On("e", func(e any) error {
		order = append(order, "high")
		return nil
	}, 10)
	bus.Emit("e", nil)
	if order[0] != "high" || order[1] != "low" {
		t.Fatalf("expected [high low], got %v", order)
	}
}

func TestEventBusWildcard(t *testing.T) {
	bus := NewEventBus()
	count := 0

	bus.On("user.*", func(e any) error {
		count++
		return nil
	})
	bus.Emit("user.registered", nil)
	bus.Emit("user.deleted", nil)
	if count != 2 {
		t.Fatalf("expected 2, got %d", count)
	}
}

func TestEventBusWildcardNested(t *testing.T) {
	bus := NewEventBus()
	var events []string
	bus.On("app.*", func(e any) error {
		events = append(events, e.(string))
		return nil
	})
	bus.Emit("app.user.created", "app.user.created")
	bus.Emit("app.user.deleted", "app.user.deleted")
	if len(events) != 2 {
		t.Fatalf("expected 2, got %d: %v", len(events), events)
	}
}

func TestEventBusOff(t *testing.T) {
	bus := NewEventBus()
	count := 0
	h := func(e any) error {
		count++
		return nil
	}
	bus.On("e", h)
	bus.Emit("e", nil)
	if count != 1 {
		t.Fatal("expected 1")
	}
	bus.Off("e")
	bus.Emit("e", nil)
	if count != 1 {
		t.Fatal("should not be called after Off")
	}
}

func TestEventBusHasListeners(t *testing.T) {
	bus := NewEventBus()
	if bus.HasListeners("nonexistent") {
		t.Fatal("expected no listeners")
	}
	bus.On("my.event", func(e any) error { return nil })
	if !bus.HasListeners("my.event") {
		t.Fatal("expected listeners")
	}
}

func TestEventBusHasListenersWildcard(t *testing.T) {
	bus := NewEventBus()
	bus.On("user.*", func(e any) error { return nil })
	if !bus.HasListeners("user.created") {
		t.Fatal("wildcard should match")
	}
}

func TestEventBusEmitAsync(t *testing.T) {
	bus := NewEventBus()
	var mu sync.Mutex
	results := []string{}

	bus.On("task", func(e any) error {
		mu.Lock()
		results = append(results, "done")
		mu.Unlock()
		return nil
	})
	bus.EmitAsync("task", nil)
	errs := bus.Flush()
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
}

func TestEventBusEmitAsyncMultiple(t *testing.T) {
	bus := NewEventBus()
	var mu sync.Mutex
	count := 0

	bus.On("task", func(e any) error {
		time.Sleep(10 * time.Millisecond)
		mu.Lock()
		count++
		mu.Unlock()
		return nil
	})
	bus.On("task", func(e any) error {
		time.Sleep(5 * time.Millisecond)
		mu.Lock()
		count++
		mu.Unlock()
		return nil
	})
	bus.EmitAsync("task", nil)
	errs := bus.Flush()
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if count != 2 {
		t.Fatalf("expected 2, got %d", count)
	}
}

func TestEventBusErrorHandling(t *testing.T) {
	bus := NewEventBus()
	bus.On("failing", func(e any) error {
		return fmt.Errorf("something went wrong")
	})
	bus.Emit("failing", nil)
	errs := bus.Flush()
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
}

func TestEventBusPanicRecovery(t *testing.T) {
	bus := NewEventBus()
	bus.On("panic", func(e any) error {
		panic("oops")
	})
	bus.Emit("panic", nil)
	errs := bus.Flush()
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
}

func TestEventBusFlushMultipleCalls(t *testing.T) {
	bus := NewEventBus()
	bus.On("e", func(e any) error { return fmt.Errorf("err") })
	bus.Emit("e", nil)
	errs := bus.Flush()
	if len(errs) != 1 {
		t.Fatalf("expected 1, got %d", len(errs))
	}
	errs = bus.Flush()
	if len(errs) != 0 {
		t.Fatalf("expected 0 after second flush, got %d", len(errs))
	}
}

func TestEventBusOrderAcrossListeners(t *testing.T) {
	bus := NewEventBus()
	order := []string{}
	var mu sync.Mutex

	bus.On("e", func(e any) error {
		mu.Lock()
		order = append(order, "a")
		mu.Unlock()
		return nil
	}, 5)
	bus.On("e", func(e any) error {
		mu.Lock()
		order = append(order, "b")
		mu.Unlock()
		return nil
	}, 10)

	bus.Emit("e", nil)
	mu.Lock()
	if len(order) != 2 || order[0] != "b" || order[1] != "a" {
		t.Fatalf("expected [b a], got %v", order)
	}
	mu.Unlock()
}

func TestEventBusSubscribe(t *testing.T) {
	bus := NewEventBus()
	called := false

	sub := &testSubscriber{fn: func(b *EventBus) {
		b.On("sub.event", func(e any) error {
			called = true
			return nil
		})
	}}
	bus.Subscribe(sub)
	bus.Emit("sub.event", nil)
	if !called {
		t.Fatal("subscriber handler should have been called")
	}
}

func TestEventBusChaining(t *testing.T) {
	bus := NewEventBus()
	count := 0

	bus.On("e1", func(e any) error { count++; return nil }).
		On("e2", func(e any) error { count++; return nil })

	bus.Emit("e1", nil).Emit("e2", nil)
	if count != 2 {
		t.Fatalf("expected 2, got %d", count)
	}
}

func TestEventBusEmitAsyncNoFlushLeak(t *testing.T) {
	bus := NewEventBus()
	bus.On("leak", func(e any) error { return nil })
	bus.EmitAsync("leak", nil)
	bus.EmitAsync("leak", nil)
	errs := bus.Flush()
	if len(errs) != 0 {
		t.Fatalf("expected 0 errors, got %d", len(errs))
	}
}

func TestEventBusEmitterPattern(t *testing.T) {
	bus := NewEventBus()
	var captured string
	bus.On("order.placed", func(e any) error {
		captured = e.(string)
		return nil
	})
	bus.Emit("order.placed", "order-123")
	if captured != "order-123" {
		t.Fatalf("expected 'order-123', got '%s'", captured)
	}
}

type testSubscriber struct {
	fn func(*EventBus)
}

func (s *testSubscriber) Subscribe(bus *EventBus) {
	s.fn(bus)
}
