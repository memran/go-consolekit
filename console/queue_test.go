package console

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func setupQueue(t *testing.T) *Queue {
	t.Helper()
	dir := filepath.Join(os.TempDir(), "consolekit-queue-test-"+uuid())
	q := NewQueue(dir).Name("test").Concurrency(5)
	return q
}

func cleanupQueue(t *testing.T, q *Queue) {
	t.Helper()
	q.Flush()
	os.RemoveAll(q.basePath)
}

func TestQueuePush(t *testing.T) {
	q := setupQueue(t)
	defer cleanupQueue(t, q)

	task, err := q.Push("test-task", map[string]string{"foo": "bar"})
	if err != nil {
		t.Fatal(err)
	}
	if task.ID == "" {
		t.Fatal("task should have an ID")
	}
	if task.Name != "test-task" {
		t.Fatalf("expected 'test-task', got '%s'", task.Name)
	}
	if task.Status != TaskPending {
		t.Fatalf("expected pending, got %s", task.Status)
	}
	if task.MaxAttempts != 3 {
		t.Fatalf("expected 3 max attempts, got %d", task.MaxAttempts)
	}
}

func TestQueuePushWithDelay(t *testing.T) {
	q := setupQueue(t)
	defer cleanupQueue(t, q)

	task, err := q.Push("delayed", "data", WithDelay(1*time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	if task.ScheduledAt == nil {
		t.Fatal("expected ScheduledAt to be set")
	}
	if task.ScheduledAt.Before(time.Now()) {
		t.Fatal("ScheduledAt should be in the future")
	}
}

func TestQueuePushWithMaxAttempts(t *testing.T) {
	q := setupQueue(t)
	defer cleanupQueue(t, q)

	task, err := q.Push("limited", "data", WithMaxAttempts(1))
	if err != nil {
		t.Fatal(err)
	}
	if task.MaxAttempts != 1 {
		t.Fatalf("expected 1, got %d", task.MaxAttempts)
	}
}

func TestQueuePushWithRetryDelay(t *testing.T) {
	q := setupQueue(t)
	defer cleanupQueue(t, q)

	task, err := q.Push("retry", "data", WithRetryDelay(10*time.Second))
	if err != nil {
		t.Fatal(err)
	}
	if time.Duration(task.RetryDelay) != 10*time.Second {
		t.Fatalf("expected 10s, got %v", task.RetryDelay)
	}
}

func TestQueuePending(t *testing.T) {
	q := setupQueue(t)
	defer cleanupQueue(t, q)

	q.Push("task-1", "data")
	q.Push("task-2", "data")

	pending, err := q.Pending()
	if err != nil {
		t.Fatal(err)
	}
	if len(pending) != 2 {
		t.Fatalf("expected 2 pending, got %d", len(pending))
	}
}

func TestQueueCount(t *testing.T) {
	q := setupQueue(t)
	defer cleanupQueue(t, q)

	q.Push("t1", "data")
	q.Push("t2", "data")
	q.Push("t3", "data")

	if count := q.Count(); count != 3 {
		t.Fatalf("expected 3, got %d", count)
	}
}

func TestQueueFailed(t *testing.T) {
	q := setupQueue(t)
	defer cleanupQueue(t, q)

	task, _ := q.Push("failing-task", "data")
	q.writeTask(task, TaskFailed)

	failed, err := q.Failed()
	if err != nil {
		t.Fatal(err)
	}
	if len(failed) != 1 {
		t.Fatalf("expected 1 failed, got %d", len(failed))
	}
}

func TestQueueRetry(t *testing.T) {
	q := setupQueue(t)
	defer cleanupQueue(t, q)

	task, _ := q.Push("retry-me", "data")
	q.writeTask(task, TaskFailed)

	if err := q.Retry(task.ID); err != nil {
		t.Fatal(err)
	}

	pending, _ := q.Pending()
	if len(pending) != 1 {
		t.Fatalf("expected 1 pending after retry, got %d", len(pending))
	}
	failed, _ := q.Failed()
	if len(failed) != 0 {
		t.Fatalf("expected 0 failed after retry, got %d", len(failed))
	}
}

func TestQueueRemove(t *testing.T) {
	q := setupQueue(t)
	defer cleanupQueue(t, q)

	task, _ := q.Push("remove-me", "data")
	if err := q.Remove(task.ID); err != nil {
		t.Fatal(err)
	}
	if q.Count() != 0 {
		t.Fatal("expected 0 after remove")
	}
}

func TestQueueFlush(t *testing.T) {
	q := setupQueue(t)
	defer cleanupQueue(t, q)

	q.Push("t1", "data")
	q.Push("t2", "data")
	q.Flush()

	if q.Count() != 0 {
		t.Fatal("expected 0 after flush")
	}
}

func TestQueueWorkProcessesTask(t *testing.T) {
	q := setupQueue(t)
	defer cleanupQueue(t, q)

	var mu sync.Mutex
	processed := false

	q.Register("echo", func(ctx *TaskContext) error {
		mu.Lock()
		processed = true
		mu.Unlock()
		return nil
	})

	q.Push("echo", "hello")

	done := make(chan struct{})
	go func() {
		q.Work()
		close(done)
	}()

	time.Sleep(500 * time.Millisecond)
	q.Stop()
	<-done

	mu.Lock()
	if !processed {
		t.Fatal("task should have been processed")
	}
	mu.Unlock()
}

func TestQueueWorkFailedTaskMovesToFailed(t *testing.T) {
	q := setupQueue(t)
	defer cleanupQueue(t, q)

	q.Register("failing", func(ctx *TaskContext) error {
		return &ValidationError{Message: "always fails"}
	})

	q.Push("failing", nil, WithMaxAttempts(1))

	done := make(chan struct{})
	go func() {
		q.Work()
		close(done)
	}()
	time.Sleep(500 * time.Millisecond)
	q.Stop()
	<-done

	failed, _ := q.Failed()
	if len(failed) == 0 {
		t.Fatal("expected failed task")
	}
}

func TestQueueWorkWithPayloadBinding(t *testing.T) {
	q := setupQueue(t)
	defer cleanupQueue(t, q)

	type Payload struct {
		Message string
	}

	var result string
	var mu sync.Mutex

	q.Register("with-payload", func(ctx *TaskContext) error {
		var p Payload
		ctx.Bind(&p)
		mu.Lock()
		result = p.Message
		mu.Unlock()
		return nil
	})

	q.Push("with-payload", Payload{Message: "hello world"})

	done := make(chan struct{})
	go func() {
		q.Work()
		close(done)
	}()
	time.Sleep(500 * time.Millisecond)
	q.Stop()
	<-done

	mu.Lock()
	if result != "hello world" {
		t.Fatalf("expected 'hello world', got '%s'", result)
	}
	mu.Unlock()
}

func TestQueueWorkRetryOnFailure(t *testing.T) {
	q := setupQueue(t)
	defer cleanupQueue(t, q)

	q.MaxAttempts(3).RetryDelay(50 * time.Millisecond)

	var attempts []int
	var mu sync.Mutex

	q.Register("flaky", func(ctx *TaskContext) error {
		mu.Lock()
		attempts = append(attempts, ctx.Task.Attempts+1)
		mu.Unlock()
		return &ValidationError{Message: "not ready yet"}
	})

	q.Push("flaky", nil)

	done := make(chan struct{})
	go func() {
		q.Work()
		close(done)
	}()
	time.Sleep(800 * time.Millisecond)
	q.Stop()
	<-done

	mu.Lock()
	if len(attempts) != 3 {
		t.Fatalf("expected 3 attempts, got %d: %v", len(attempts), attempts)
	}
	mu.Unlock()

	failed, _ := q.Failed()
	if len(failed) != 1 {
		t.Fatal("task should be failed after max attempts")
	}
}

func TestQueueWorkCompletedTask(t *testing.T) {
	q := setupQueue(t)
	defer cleanupQueue(t, q)

	q.Register("quick", func(ctx *TaskContext) error {
		return nil
	})

	q.Push("quick", nil)

	done := make(chan struct{})
	go func() {
		q.Work()
		close(done)
	}()
	time.Sleep(500 * time.Millisecond)
	q.Stop()
	<-done

	completed, _ := q.Completed()
	if len(completed) != 1 {
		t.Fatalf("expected 1 completed, got %d", len(completed))
	}
}

func TestQueueWorkMultipleConcurrent(t *testing.T) {
	q := setupQueue(t)
	defer cleanupQueue(t, q)

	q.Concurrency(10)

	var mu sync.Mutex
	var order []int

	q.Register("worker", func(ctx *TaskContext) error {
		var n int
		ctx.Bind(&n)
		time.Sleep(50 * time.Millisecond)
		mu.Lock()
		order = append(order, n)
		mu.Unlock()
		return nil
	})

	for i := 0; i < 5; i++ {
		q.Push("worker", i)
	}

	done := make(chan struct{})
	go func() {
		q.Work()
		close(done)
	}()
	time.Sleep(600 * time.Millisecond)
	q.Stop()
	<-done

	mu.Lock()
	if len(order) != 5 {
		t.Fatalf("expected 5 results, got %d", len(order))
	}
	mu.Unlock()
}

func TestQueueRunStop(t *testing.T) {
	q := setupQueue(t)
	defer cleanupQueue(t, q)

	q.Register("quick", func(ctx *TaskContext) error { return nil })
	q.Push("quick", nil)

	q.Run()
	time.Sleep(200 * time.Millisecond)
	if !q.IsRunning() {
		t.Fatal("expected running")
	}
	q.Stop()
	if q.IsRunning() {
		t.Fatal("expected stopped")
	}
}

func TestQueueNoHandler(t *testing.T) {
	q := setupQueue(t)
	defer cleanupQueue(t, q)

	q.Push("no-handler", "data")

	done := make(chan struct{})
	go func() {
		q.Work()
		close(done)
	}()
	time.Sleep(500 * time.Millisecond)
	q.Stop()
	<-done

	failed, _ := q.Failed()
	if len(failed) != 1 {
		t.Fatal("expected failed task with no handler")
	}
}

func TestQueueLogger(t *testing.T) {
	q := setupQueue(t)
	defer cleanupQueue(t, q)

	var mu sync.Mutex
	var logged string
	q.Logger(func(format string, args ...any) {
		mu.Lock()
		logged = format
		mu.Unlock()
	}).Register("log-test", func(ctx *TaskContext) error {
		return nil
	})

	q.Push("log-test", nil)

	done := make(chan struct{})
	go func() {
		q.Work()
		close(done)
	}()
	time.Sleep(500 * time.Millisecond)
	q.Stop()
	<-done

	mu.Lock()
	if logged == "" {
		t.Fatal("expected log output")
	}
	mu.Unlock()
}

func TestTaskContextMarkAsFailed(t *testing.T) {
	q := setupQueue(t)
	defer cleanupQueue(t, q)

	q.Register("mark-fail", func(ctx *TaskContext) error {
		ctx.MarkAsFailed(&ValidationError{Message: "manual fail"})
		return &ValidationError{Message: "manual fail"}
	})

	q.Push("mark-fail", nil, WithMaxAttempts(1))

	done := make(chan struct{})
	go func() {
		q.Work()
		close(done)
	}()
	time.Sleep(500 * time.Millisecond)
	q.Stop()
	<-done

	failed, _ := q.Failed()
	if len(failed) != 1 {
		t.Fatal("expected failed task")
	}
}

func TestQueueChaining(t *testing.T) {
	q := setupQueue(t)
	defer cleanupQueue(t, q)

	q.Name("chain").MaxAttempts(5).Concurrency(3)

	task, err := q.Push("chain-task", "value")
	if err != nil {
		t.Fatal(err)
	}
	if task.Queue != "chain" {
		t.Fatalf("expected 'chain', got '%s'", task.Queue)
	}
}

func TestQueueTaskOrderPreserved(t *testing.T) {
	q := setupQueue(t)
	defer cleanupQueue(t, q)

	q.Concurrency(1)

	var order []int
	var mu sync.Mutex

	q.Register("ordered", func(ctx *TaskContext) error {
		var n int
		ctx.Bind(&n)
		mu.Lock()
		order = append(order, n)
		mu.Unlock()
		return nil
	})

	q.Push("ordered", 1)
	time.Sleep(10 * time.Millisecond)
	q.Push("ordered", 2)
	time.Sleep(10 * time.Millisecond)
	q.Push("ordered", 3)

	done := make(chan struct{})
	go func() {
		q.Work()
		close(done)
	}()
	time.Sleep(800 * time.Millisecond)
	q.Stop()
	<-done

	mu.Lock()
	if len(order) != 3 || order[0] != 1 || order[1] != 2 || order[2] != 3 {
		t.Fatalf("expected [1 2 3], got %v", order)
	}
	mu.Unlock()
}

func TestQueueScheduledTaskNotProcessedEarly(t *testing.T) {
	q := setupQueue(t)
	defer cleanupQueue(t, q)

	q.Register("later", func(ctx *TaskContext) error { return nil })

	q.Push("later", nil, WithDelay(1*time.Hour))

	done := make(chan struct{})
	go func() {
		q.Work()
		close(done)
	}()
	time.Sleep(500 * time.Millisecond)
	q.Stop()
	<-done

	pending, _ := q.Pending()
	if len(pending) != 1 {
		t.Fatalf("expected 1 pending (scheduled), got %d", len(pending))
	}
}

func TestQueue(t *testing.T) {
	q := setupQueue(t)
	defer cleanupQueue(t, q)

	q.Register("simple", func(ctx *TaskContext) error { return nil })
	q.Push("simple", nil)

	done := make(chan struct{})
	go func() {
		q.Work()
		close(done)
	}()
	time.Sleep(500 * time.Millisecond)
	q.Stop()
	<-done

	completed, _ := q.Completed()
	if len(completed) != 1 {
		t.Fatalf("expected 1 completed, got %d", len(completed))
	}
}
