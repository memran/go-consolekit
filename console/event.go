package console

import (
	"fmt"
	"strings"
	"sync"
)

type Event interface {
	Name() string
}

type EventHandler func(any) error

type eventListener struct {
	handler  EventHandler
	priority int
}

type EventBus struct {
	mu       sync.RWMutex
	handlers map[string][]eventListener
	wg       sync.WaitGroup
	errors   []error
	errMu    sync.Mutex
	async    bool
}

func NewEventBus() *EventBus {
	return &EventBus{
		handlers: make(map[string][]eventListener),
	}
}

func (b *EventBus) On(event string, handler EventHandler, priority ...int) *EventBus {
	b.mu.Lock()
	defer b.mu.Unlock()

	prio := 0
	if len(priority) > 0 {
		prio = priority[0]
	}

	listener := eventListener{handler: handler, priority: prio}
	insert := func(listeners []eventListener, l eventListener) []eventListener {
		idx := len(listeners)
		for i, el := range listeners {
			if l.priority > el.priority {
				idx = i
				break
			}
		}
		return append(listeners[:idx], append([]eventListener{l}, listeners[idx:]...)...)
	}
	b.handlers[event] = insert(b.handlers[event], listener)
	return b
}

func (b *EventBus) Off(event string, handler ...EventHandler) *EventBus {
	b.mu.Lock()
	defer b.mu.Unlock()

	if len(handler) == 0 {
		delete(b.handlers, event)
		return b
	}

	listeners := b.handlers[event]
	filtered := make([]eventListener, 0, len(listeners))
	for _, l := range listeners {
		keep := true
		for _, h := range handler {
			if equalHandler(l.handler, h) {
				keep = false
				break
			}
		}
		if keep {
			filtered = append(filtered, l)
		}
	}
	if len(filtered) == 0 {
		delete(b.handlers, event)
	} else {
		b.handlers[event] = filtered
	}
	return b
}

func (b *EventBus) HasListeners(event string) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()

	if _, ok := b.handlers[event]; ok {
		return true
	}
	for pattern := range b.handlers {
		if matchWildcard(pattern, event) {
			return true
		}
	}
	return false
}

func (b *EventBus) Emit(event string, payload any) *EventBus {
	b.dispatch(event, payload, false)
	return b
}

func (b *EventBus) EmitAsync(event string, payload any) *EventBus {
	b.async = true
	b.dispatch(event, payload, true)
	return b
}

func (b *EventBus) Flush() []error {
	b.wg.Wait()
	b.errMu.Lock()
	errs := make([]error, len(b.errors))
	copy(errs, b.errors)
	b.errors = b.errors[:0]
	b.errMu.Unlock()
	return errs
}

func (b *EventBus) Subscribe(subscriber Subscriber) *EventBus {
	subscriber.Subscribe(b)
	return b
}

func (b *EventBus) dispatch(event string, payload any, async bool) {
	b.mu.RLock()
	listeners := b.collectListeners(event)
	b.mu.RUnlock()

	if async {
		for _, l := range listeners {
			b.wg.Add(1)
			go func(h EventHandler) {
				defer b.wg.Done()
				if err := b.safeCall(h, payload); err != nil {
					b.errMu.Lock()
					b.errors = append(b.errors, fmt.Errorf("[%s] %w", event, err))
					b.errMu.Unlock()
				}
			}(l.handler)
		}
	} else {
		for _, l := range listeners {
			if err := b.safeCall(l.handler, payload); err != nil {
				b.errMu.Lock()
				b.errors = append(b.errors, fmt.Errorf("[%s] %w", event, err))
				b.errMu.Unlock()
			}
		}
	}
}

func (b *EventBus) collectListeners(event string) []eventListener {
	var result []eventListener
	for pattern, listeners := range b.handlers {
		if pattern == event || matchWildcard(pattern, event) {
			result = append(result, listeners...)
		}
	}
	return result
}

func (b *EventBus) safeCall(h EventHandler, payload any) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic: %v", r)
		}
	}()
	return h(payload)
}

type Subscriber interface {
	Subscribe(*EventBus)
}

func matchWildcard(pattern, event string) bool {
	if pattern == event {
		return true
	}
	if !strings.Contains(pattern, "*") {
		return false
	}
	parts := strings.Split(pattern, "*")
	if len(parts) == 1 {
		return pattern == event
	}
	if !strings.HasPrefix(event, parts[0]) {
		return false
	}
	remain := event[len(parts[0]):]
	for i := 1; i < len(parts); i++ {
		idx := strings.Index(remain, parts[i])
		if idx == -1 {
			return false
		}
		if i == len(parts)-1 && parts[i] == "" {
			return true
		}
		remain = remain[idx+len(parts[i]):]
	}
	if len(parts[len(parts)-1]) > 0 {
		return strings.HasSuffix(event, parts[len(parts)-1])
	}
	return true
}

func equalHandler(a, b EventHandler) bool {
	return fmt.Sprintf("%p", a) == fmt.Sprintf("%p", b)
}
