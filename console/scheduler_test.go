package console

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func setupScheduler(t *testing.T) (*Scheduler, *Queue) {
	t.Helper()
	qdir := filepath.Join(os.TempDir(), "consolekit-sched-queue-"+uuid())
	sdir := filepath.Join(os.TempDir(), "consolekit-sched-data-"+uuid())

	q := NewQueue(qdir).Name("sched-test").Concurrency(5)
	s := NewScheduler(sdir).Queue(q)

	return s, q
}

func cleanupScheduler(t *testing.T, s *Scheduler, q *Queue) {
	t.Helper()
	q.Flush()
	s.Stop()
	os.RemoveAll(q.basePath)
	os.RemoveAll(s.basePath)
}

func TestSchedulerNew(t *testing.T) {
	dir := filepath.Join(os.TempDir(), "consolekit-sched-"+uuid())
	s := NewScheduler(dir)
	defer os.RemoveAll(dir)

	if s == nil {
		t.Fatal("scheduler should not be nil")
	}
}

func TestSchedulerEvery(t *testing.T) {
	dir := filepath.Join(os.TempDir(), "consolekit-sched-"+uuid())
	s := NewScheduler(dir)
	defer os.RemoveAll(dir)

	b := s.Every("5m")
	if b.expression != "5m" {
		t.Fatalf("expected '5m', got '%s'", b.expression)
	}
}

func TestSchedulerDaily(t *testing.T) {
	dir := filepath.Join(os.TempDir(), "consolekit-sched-"+uuid())
	s := NewScheduler(dir)
	defer os.RemoveAll(dir)

	b := s.Daily()
	if b.expression != "@daily" {
		t.Fatalf("expected '@daily', got '%s'", b.expression)
	}
}

func TestSchedulerHourly(t *testing.T) {
	dir := filepath.Join(os.TempDir(), "consolekit-sched-"+uuid())
	s := NewScheduler(dir)
	defer os.RemoveAll(dir)

	b := s.Hourly()
	if b.expression != "@hourly" {
		t.Fatalf("expected '@hourly', got '%s'", b.expression)
	}
}

func TestSchedulerCallCreatesEntry(t *testing.T) {
	dir := filepath.Join(os.TempDir(), "consolekit-sched-"+uuid())
	s := NewScheduler(dir)
	defer os.RemoveAll(dir)

	entry, err := s.Every("1m").Call("test-task", map[string]string{"foo": "bar"})
	if err != nil {
		t.Fatal(err)
	}

	if entry.ID == "" {
		t.Fatal("entry should have ID")
	}
	if entry.TaskName != "test-task" {
		t.Fatalf("expected 'test-task', got '%s'", entry.TaskName)
	}
	if entry.Expression != "1m" {
		t.Fatalf("expected '1m', got '%s'", entry.Expression)
	}
	if !entry.Enabled {
		t.Fatal("entry should be enabled")
	}
}

func TestSchedulerCallPersistsToDisk(t *testing.T) {
	dir := filepath.Join(os.TempDir(), "consolekit-sched-"+uuid())
	s := NewScheduler(dir)
	defer os.RemoveAll(dir)

	entry, _ := s.Every("5m").Call("persist-test", nil)

	loaded := NewScheduler(dir)
	loaded.loadEntries()

	entries := loaded.Entries()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry on reload, got %d", len(entries))
	}
	if entries[0].ID != entry.ID {
		t.Fatalf("IDs don't match")
	}
}

func TestSchedulerRunWithoutQueue(t *testing.T) {
	dir := filepath.Join(os.TempDir(), "consolekit-sched-"+uuid())
	s := NewScheduler(dir)
	defer os.RemoveAll(dir)

	err := s.Run()
	if err == nil {
		t.Fatal("expected error when no queue set")
	}
}

func TestSchedulerIsDue(t *testing.T) {
	dir := filepath.Join(os.TempDir(), "consolekit-sched-"+uuid())
	s := NewScheduler(dir)
	defer os.RemoveAll(dir)

	now := time.Now()

	t.Run("never run", func(t *testing.T) {
		entry := &ScheduleEntry{Expression: "5m"}
		if !s.isDue(entry, now) {
			t.Fatal("entry with no LastRun should be due")
		}
	})

	t.Run("duration elapsed", func(t *testing.T) {
		last := now.Add(-10 * time.Minute)
		entry := &ScheduleEntry{
			Expression: "5m",
			LastRun:    &last,
		}
		if !s.isDue(entry, now) {
			t.Fatal("entry should be due after 5m has passed")
		}
	})

	t.Run("duration not elapsed", func(t *testing.T) {
		last := now.Add(-1 * time.Minute)
		entry := &ScheduleEntry{
			Expression: "5m",
			LastRun:    &last,
		}
		if s.isDue(entry, now) {
			t.Fatal("entry should not be due before 5m")
		}
	})

	t.Run("@daily yesterday", func(t *testing.T) {
		yesterday := now.Add(-24 * time.Hour)
		entry := &ScheduleEntry{
			Expression: "@daily",
			LastRun:    &yesterday,
		}
		if !s.isDue(entry, now) {
			t.Fatal("@daily should be due after 24h")
		}
	})

	t.Run("@daily today", func(t *testing.T) {
		today := now.Truncate(24 * time.Hour).Add(1 * time.Hour)
		entry := &ScheduleEntry{
			Expression: "@daily",
			LastRun:    &today,
		}
		if s.isDue(entry, now) {
			t.Fatal("@daily should not be due same day")
		}
	})

	t.Run("@hourly this hour", func(t *testing.T) {
		thisHour := now.Truncate(time.Hour).Add(5 * time.Minute)
		entry := &ScheduleEntry{
			Expression: "@hourly",
			LastRun:    &thisHour,
		}
		if s.isDue(entry, now) {
			t.Fatal("@hourly should not be due same hour")
		}
	})

	t.Run("invalid duration", func(t *testing.T) {
		last := now.Add(-1 * time.Hour)
		entry := &ScheduleEntry{
			Expression: "invalid",
			LastRun:    &last,
		}
		if s.isDue(entry, now) {
			t.Fatal("invalid duration should not be due")
		}
	})
}

func TestSchedulerEvaluateTriggersTask(t *testing.T) {
	s, q := setupScheduler(t)
	defer cleanupScheduler(t, s, q)

	q.Register("scheduled-task", func(ctx *TaskContext) error {
		return nil
	})

	s.Every("1s").Call("scheduled-task", nil)

	s.evaluate()

	time.Sleep(100 * time.Millisecond)

	pending, _ := q.Pending()
	if len(pending) != 1 {
		t.Fatalf("expected 1 pending task after evaluate, got %d", len(pending))
	}
}

func TestSchedulerEvaluateDoesNotDoubleTrigger(t *testing.T) {
	s, q := setupScheduler(t)
	defer cleanupScheduler(t, s, q)

	q.Register("no-double", func(ctx *TaskContext) error {
		return nil
	})

	s.Every("1m").Call("no-double", nil)

	s.evaluate()
	s.evaluate()

	time.Sleep(100 * time.Millisecond)

	pending, _ := q.Pending()
	if len(pending) != 1 {
		t.Fatalf("expected 1 pending (not double-triggered), got %d", len(pending))
	}
}

func TestSchedulerStartStop(t *testing.T) {
	s, q := setupScheduler(t)
	defer cleanupScheduler(t, s, q)

	wait := make(chan struct{})
	q.Register("start-stop", func(ctx *TaskContext) error {
		close(wait)
		return nil
	})

	s.Every("1s").Call("start-stop", nil)

	s.Start()
	time.Sleep(50 * time.Millisecond)
	if !s.IsRunning() {
		t.Fatal("expected running")
	}

	s.Stop()
	if s.IsRunning() {
		t.Fatal("expected stopped")
	}
}

func TestSchedulerEntries(t *testing.T) {
	dir := filepath.Join(os.TempDir(), "consolekit-sched-"+uuid())
	s := NewScheduler(dir)
	defer os.RemoveAll(dir)

	s.Every("5m").Call("task-1", nil)
	s.Every("10m").Call("task-2", nil)

	entries := s.Entries()
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
}

func TestSchedulerRemove(t *testing.T) {
	dir := filepath.Join(os.TempDir(), "consolekit-sched-"+uuid())
	s := NewScheduler(dir)
	defer os.RemoveAll(dir)

	entry, _ := s.Every("5m").Call("remove-me", nil)

	if err := s.Remove(entry.ID); err != nil {
		t.Fatal(err)
	}

	entries := s.Entries()
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries after remove, got %d", len(entries))
	}
}

func TestSchedulerLoadEntriesOnlyEnabled(t *testing.T) {
	dir := filepath.Join(os.TempDir(), "consolekit-sched-"+uuid())
	s := NewScheduler(dir)
	defer os.RemoveAll(dir)

	entry, _ := s.Every("5m").Call("disabled-task", nil)
	entry.Enabled = false
	s.saveEntry(entry)

	s.loadEntries()
	entries := s.Entries()
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries (disabled), got %d", len(entries))
	}
}

func TestSchedulerLogger(t *testing.T) {
	dir := filepath.Join(os.TempDir(), "consolekit-sched-"+uuid())
	s := NewScheduler(dir)
	defer os.RemoveAll(dir)

	var mu sync.Mutex
	var logged string
	s.Logger(func(format string, args ...any) {
		mu.Lock()
		logged = format
		mu.Unlock()
	})

	s.log("test %s", "log")
	mu.Lock()
	if logged == "" {
		t.Fatal("expected log output")
	}
	mu.Unlock()
}

func TestSchedulerChaining(t *testing.T) {
	dir := filepath.Join(os.TempDir(), "consolekit-sched-"+uuid())
	qdir := filepath.Join(os.TempDir(), "consolekit-chain-queue-"+uuid())

	q := NewQueue(qdir)
	s := NewScheduler(dir).
		Queue(q).
		Logger(nil)
	defer os.RemoveAll(dir)
	defer os.RemoveAll(qdir)

	if s.queue != q {
		t.Fatal("queue should be set")
	}
}

func TestSchedulerEntryWithPayload(t *testing.T) {
	s, q := setupScheduler(t)
	defer cleanupScheduler(t, s, q)

	type Data struct {
		Key string
	}

	var result string
	var mu sync.Mutex

	q.Register("with-data", func(ctx *TaskContext) error {
		var d Data
		ctx.Bind(&d)
		mu.Lock()
		result = d.Key
		mu.Unlock()
		return nil
	})

	s.Every("1m").Call("with-data", Data{Key: "scheduled-value"})

	s.evaluate()

	done := make(chan struct{})
	q.Concurrency(1)
	go func() {
		q.Work()
		close(done)
	}()
	time.Sleep(500 * time.Millisecond)
	q.Stop()
	<-done

	mu.Lock()
	if result != "scheduled-value" {
		t.Fatalf("expected 'scheduled-value', got '%s'", result)
	}
	mu.Unlock()
}
