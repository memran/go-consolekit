package console

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

type TaskStatus string

const (
	TaskPending    TaskStatus = "pending"
	TaskProcessing TaskStatus = "processing"
	TaskCompleted  TaskStatus = "completed"
	TaskFailed     TaskStatus = "failed"
)

type Task struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Payload     json.RawMessage `json:"payload"`
	Status      TaskStatus      `json:"status"`
	Attempts    int             `json:"attempts"`
	MaxAttempts int             `json:"max_attempts"`
	RetryDelay  Duration        `json:"retry_delay"`
	CreatedAt   time.Time       `json:"created_at"`
	ScheduledAt *time.Time      `json:"scheduled_at,omitempty"`
	UpdatedAt   time.Time       `json:"updated_at"`
	Error       string          `json:"error,omitempty"`
	Queue       string          `json:"queue"`
}

type TaskContext struct {
	Task    *Task
	queue   *Queue
	handler TaskHandler
	logger  func(string, ...any)
}

func (c *TaskContext) Bind(v any) error {
	return json.Unmarshal(c.Task.Payload, v)
}

func (c *TaskContext) Log(format string, args ...any) {
	if c.logger != nil {
		c.logger(format, args...)
	}
}

func (c *TaskContext) MarkAsFailed(err error) {
	c.Task.Error = err.Error()
}

type TaskHandler func(ctx *TaskContext) error

type TaskOption func(*Task)

func WithDelay(d time.Duration) TaskOption {
	return func(t *Task) {
		now := time.Now()
		sched := now.Add(d)
		t.ScheduledAt = &sched
	}
}

func WithMaxAttempts(n int) TaskOption {
	return func(t *Task) {
		t.MaxAttempts = n
	}
}

func WithRetryDelay(d time.Duration) TaskOption {
	return func(t *Task) {
		t.RetryDelay = Duration(d)
	}
}

type Duration time.Duration

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Duration(d).String())
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	pd, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	*d = Duration(pd)
	return nil
}

type Queue struct {
	basePath    string
	name        string
	maxAttempts int
	retryDelay  time.Duration
	concurrency int
	handlers    map[string]TaskHandler
	done        chan struct{}
	sem         chan struct{}
	running     bool
	mu          sync.RWMutex
	logFn       func(string, ...any)
}

func NewQueue(basePath string) *Queue {
	return &Queue{
		basePath:    basePath,
		name:        "default",
		maxAttempts: 3,
		retryDelay:  30 * time.Second,
		concurrency: 1,
		handlers:    make(map[string]TaskHandler),
	}
}

func (q *Queue) Name(name string) *Queue {
	q.name = name
	return q
}

func (q *Queue) MaxAttempts(n int) *Queue {
	q.maxAttempts = n
	return q
}

func (q *Queue) RetryDelay(d time.Duration) *Queue {
	q.retryDelay = d
	return q
}

func (q *Queue) Concurrency(n int) *Queue {
	q.concurrency = n
	return q
}

func (q *Queue) Logger(fn func(string, ...any)) *Queue {
	q.logFn = fn
	return q
}

func (q *Queue) Register(name string, handler TaskHandler) *Queue {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.handlers[name] = handler
	return q
}

func (q *Queue) Push(name string, payload any, opts ...TaskOption) (*Task, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	task := &Task{
		ID:          uuid(),
		Name:        name,
		Status:      TaskPending,
		Attempts:    0,
		MaxAttempts: q.maxAttempts,
		RetryDelay:  Duration(q.retryDelay),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Queue:       q.name,
	}

	for _, opt := range opts {
		opt(task)
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal payload: %w", err)
	}
	task.Payload = data

	if err := q.writeTask(task, TaskPending); err != nil {
		return nil, err
	}

	return task, nil
}

func (q *Queue) Pending() ([]*Task, error) {
	return q.listTasks(TaskPending)
}

func (q *Queue) Failed() ([]*Task, error) {
	return q.listTasks(TaskFailed)
}

func (q *Queue) ProcessingTasks() ([]*Task, error) {
	return q.listTasks(TaskProcessing)
}

func (q *Queue) Completed() ([]*Task, error) {
	return q.listTasks(TaskCompleted)
}

func (q *Queue) Count() int {
	files, err := os.ReadDir(q.dir(TaskPending))
	if err != nil {
		return 0
	}
	return len(files)
}

func (q *Queue) Retry(id string) error {
	task, err := q.readTask(id, TaskFailed)
	if err != nil {
		return err
	}
	task.Status = TaskPending
	task.Attempts = 0
	task.Error = ""
	task.RetryDelay = Duration(q.retryDelay)
	task.UpdatedAt = time.Now()

	if err := q.writeTask(task, TaskPending); err != nil {
		return err
	}
	if err := q.deleteTask(id, TaskFailed); err != nil {
		return err
	}
	return nil
}

func (q *Queue) Remove(id string) error {
	for _, status := range []TaskStatus{TaskPending, TaskProcessing, TaskFailed, TaskCompleted} {
		if err := q.deleteTask(id, status); err == nil {
			return nil
		}
	}
	return fmt.Errorf("task %s not found", id)
}

func (q *Queue) Flush() error {
	for _, status := range []TaskStatus{TaskPending, TaskProcessing, TaskFailed, TaskCompleted} {
		dir := q.dir(status)
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, entry := range entries {
			os.Remove(filepath.Join(dir, entry.Name()))
		}
	}
	return nil
}

func (q *Queue) Work() error {
	q.mu.Lock()
	if q.running {
		q.mu.Unlock()
		return fmt.Errorf("queue worker already running")
	}
	q.done = make(chan struct{})
	q.sem = make(chan struct{}, q.concurrency)
	q.running = true
	q.mu.Unlock()

	q.log("Worker started (concurrency: %d)", q.concurrency)

	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			q.processPending()
		case <-q.done:
			q.log("Worker stopped")
			return nil
		}
	}
}

func (q *Queue) Run() {
	go func() {
		_ = q.Work()
	}()
}

func (q *Queue) Stop() {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.running {
		close(q.done)
		q.running = false
	}
}

func (q *Queue) IsRunning() bool {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return q.running
}

func (q *Queue) processPending() {
	entries, err := os.ReadDir(q.dir(TaskPending))
	if err != nil {
		return
	}

	sort.Slice(entries, func(i, j int) bool {
		infoI, _ := entries[i].Info()
		infoJ, _ := entries[j].Info()
		if infoI == nil || infoJ == nil {
			return false
		}
		return infoI.ModTime().Before(infoJ.ModTime())
	})

	for _, entry := range entries {
		select {
		case q.sem <- struct{}{}:
			go func(name string) {
				defer func() { <-q.sem }()
				q.executeTask(name)
			}(entry.Name())
		default:
			return
		}
	}
}

func (q *Queue) executeTask(filename string) {
	id := strings.TrimSuffix(filename, ".json")

	task, err := q.readTask(id, TaskPending)
	if err != nil {
		q.log("Error reading task %s: %v", id, err)
		return
	}

	if task.ScheduledAt != nil && time.Now().Before(*task.ScheduledAt) {
		return
	}

	if err := q.moveTask(id, TaskPending, TaskProcessing); err != nil {
		q.log("Error moving task %s to processing: %v", id, err)
		return
	}
	task.Status = TaskProcessing
	task.UpdatedAt = time.Now()

	q.mu.RLock()
	handler, ok := q.handlers[task.Name]
	q.mu.RUnlock()

	if !ok {
		task.Error = fmt.Sprintf("no handler registered for task: %s", task.Name)
		task.Status = TaskFailed
		task.UpdatedAt = time.Now()
		q.writeTask(task, TaskFailed)
		q.deleteTask(id, TaskProcessing)
		q.log("Task %s failed: %s", id, task.Error)
		return
	}

	ctx := &TaskContext{
		Task:    task,
		queue:   q,
		handler: handler,
		logger:  q.logFn,
	}

	err = handler(ctx)

	task.Attempts++
	task.UpdatedAt = time.Now()

	if err != nil {
		task.Error = err.Error()
		q.log("Task %s attempt %d/%d failed: %v", id, task.Attempts, task.MaxAttempts, err)

		if task.Attempts >= task.MaxAttempts {
			task.Status = TaskFailed
			q.writeTask(task, TaskFailed)
			q.deleteTask(id, TaskProcessing)
			q.log("Task %s permanently failed after %d attempts", id, task.Attempts)
		} else {
			sched := time.Now().Add(time.Duration(task.RetryDelay))
			task.ScheduledAt = &sched
			task.Status = TaskPending
			q.writeTask(task, TaskPending)
			q.deleteTask(id, TaskProcessing)
			q.log("Task %s scheduled for retry at %v", id, sched)
		}
	} else {
		task.Status = TaskCompleted
		task.Error = ""
		q.writeTask(task, TaskCompleted)
		q.deleteTask(id, TaskProcessing)
		q.log("Task %s completed successfully", id)
	}
}

func (q *Queue) dir(status TaskStatus) string {
	return filepath.Join(q.basePath, "queues", q.name, string(status))
}

func (q *Queue) taskPath(id string, status TaskStatus) string {
	return filepath.Join(q.dir(status), id+".json")
}

func (q *Queue) ensureDirs() {
	for _, status := range []TaskStatus{TaskPending, TaskProcessing, TaskFailed, TaskCompleted} {
		os.MkdirAll(q.dir(status), 0755)
	}
}

func (q *Queue) writeTask(task *Task, status TaskStatus) error {
	q.ensureDirs()
	data, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("marshal task: %w", err)
	}
	path := q.taskPath(task.ID, status)
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return fmt.Errorf("write temp: %w", err)
	}
	if err := os.Rename(tmp, path); err != nil {
		return fmt.Errorf("rename: %w", err)
	}
	return nil
}

func (q *Queue) readTask(id string, status TaskStatus) (*Task, error) {
	path := q.taskPath(id, status)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var task Task
	if err := json.Unmarshal(data, &task); err != nil {
		return nil, fmt.Errorf("unmarshal task: %w", err)
	}
	return &task, nil
}

func (q *Queue) deleteTask(id string, status TaskStatus) error {
	path := q.taskPath(id, status)
	return os.Remove(path)
}

func (q *Queue) moveTask(id string, from, to TaskStatus) error {
	src := q.taskPath(id, from)
	dst := q.taskPath(id, to)
	if err := os.Rename(src, dst); err != nil {
		return err
	}
	return nil
}

func (q *Queue) listTasks(status TaskStatus) ([]*Task, error) {
	dir := q.dir(status)
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return []*Task{}, nil
		}
		return nil, err
	}
	var tasks []*Task
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		id := strings.TrimSuffix(entry.Name(), ".json")
		task, err := q.readTask(id, status)
		if err != nil {
			continue
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (q *Queue) log(format string, args ...any) {
	if q.logFn != nil {
		q.logFn("[queue:"+q.name+"] "+format, args...)
	}
}

func uuid() string {
	b := make([]byte, 16)
	rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
