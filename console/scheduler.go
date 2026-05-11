package console

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type ScheduleEntry struct {
	ID         string          `json:"id"`
	Expression string          `json:"expression"`
	TaskName   string          `json:"task_name"`
	Payload    json.RawMessage `json:"payload"`
	Enabled    bool            `json:"enabled"`
	LastRun    *time.Time      `json:"last_run,omitempty"`
	CreatedAt  time.Time       `json:"created_at"`
}

type ScheduleBuilder struct {
	scheduler  *Scheduler
	expression string
}

func (b *ScheduleBuilder) Call(name string, payload any) (*ScheduleEntry, error) {
	pdata, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal payload: %w", err)
	}

	entry := &ScheduleEntry{
		ID:         uuid(),
		Expression: b.expression,
		TaskName:   name,
		Payload:    pdata,
		Enabled:    true,
		CreatedAt:  time.Now(),
	}

	if err := b.scheduler.saveEntry(entry); err != nil {
		return nil, err
	}

	b.scheduler.mu.Lock()
	b.scheduler.entries = append(b.scheduler.entries, entry)
	b.scheduler.mu.Unlock()

	return entry, nil
}

type Scheduler struct {
	basePath string
	queue    *Queue
	entries  []*ScheduleEntry
	done     chan struct{}
	running  bool
	mu       sync.Mutex
	logFn    func(string, ...any)
}

func NewScheduler(basePath string) *Scheduler {
	return &Scheduler{
		basePath: basePath,
	}
}

func (s *Scheduler) Queue(queue *Queue) *Scheduler {
	s.queue = queue
	return s
}

func (s *Scheduler) Logger(fn func(string, ...any)) *Scheduler {
	s.logFn = fn
	return s
}

func (s *Scheduler) Every(dur string) *ScheduleBuilder {
	_, err := time.ParseDuration(dur)
	if err != nil && !strings.HasPrefix(dur, "@") {
		dur = ""
	}
	return &ScheduleBuilder{scheduler: s, expression: dur}
}

func (s *Scheduler) Daily() *ScheduleBuilder {
	return &ScheduleBuilder{scheduler: s, expression: "@daily"}
}

func (s *Scheduler) Hourly() *ScheduleBuilder {
	return &ScheduleBuilder{scheduler: s, expression: "@hourly"}
}

func (s *Scheduler) Run() error {
	if s.queue == nil {
		return fmt.Errorf("scheduler: no queue set")
	}

	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return fmt.Errorf("scheduler already running")
	}
	s.done = make(chan struct{})
	s.running = true
	s.mu.Unlock()

	s.loadEntries()
	s.log("Scheduler started (%d entries)", len(s.entries))

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	s.evaluate()

	for {
		select {
		case <-ticker.C:
			s.evaluate()
		case <-s.done:
			s.log("Scheduler stopped")
			return nil
		}
	}
}

func (s *Scheduler) Start() {
	go func() {
		_ = s.Run()
	}()
}

func (s *Scheduler) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.running {
		close(s.done)
		s.running = false
	}
}

func (s *Scheduler) IsRunning() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.running
}

func (s *Scheduler) Entries() []*ScheduleEntry {
	s.mu.Lock()
	defer s.mu.Unlock()
	result := make([]*ScheduleEntry, len(s.entries))
	copy(result, s.entries)
	return result
}

func (s *Scheduler) Remove(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	idx := -1
	for i, e := range s.entries {
		if e.ID == id {
			idx = i
			break
		}
	}
	if idx == -1 {
		return fmt.Errorf("schedule %s not found", id)
	}

	s.entries = append(s.entries[:idx], s.entries[idx+1:]...)

	path := filepath.Join(s.basePath, id+".json")
	os.Remove(path)
	return nil
}

func (s *Scheduler) loadEntries() {
	os.MkdirAll(s.basePath, 0755)

	entries, err := os.ReadDir(s.basePath)
	if err != nil {
		return
	}

	s.mu.Lock()
	s.entries = nil
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(s.basePath, entry.Name()))
		if err != nil {
			continue
		}
		var se ScheduleEntry
		if err := json.Unmarshal(data, &se); err != nil {
			continue
		}
		if se.Enabled {
			s.entries = append(s.entries, &se)
		}
	}
	s.mu.Unlock()
}

func (s *Scheduler) evaluate() {
	s.mu.Lock()
	entries := make([]*ScheduleEntry, len(s.entries))
	copy(entries, s.entries)
	s.mu.Unlock()

	now := time.Now()

	for _, entry := range entries {
		if !entry.Enabled {
			continue
		}
		if s.isDue(entry, now) {
			s.log("Triggering schedule %s: %s -> %s", entry.ID[:8], entry.Expression, entry.TaskName)
			var pushPayload any
			if len(entry.Payload) > 0 {
				json.Unmarshal(entry.Payload, &pushPayload)
			}
			if _, err := s.queue.Push(entry.TaskName, pushPayload); err != nil {
				s.log("Error pushing task for schedule %s: %v", entry.ID[:8], err)
				continue
			}
			entry.LastRun = &now
			s.saveEntry(entry)
		}
	}
}

func (s *Scheduler) isDue(entry *ScheduleEntry, now time.Time) bool {
	if entry.LastRun == nil {
		return true
	}

	switch {
	case entry.Expression == "@daily":
		last := entry.LastRun.Truncate(24 * time.Hour)
		today := now.Truncate(24 * time.Hour)
		return today.After(last)

	case entry.Expression == "@hourly":
		last := entry.LastRun.Truncate(time.Hour)
		thisHour := now.Truncate(time.Hour)
		return thisHour.After(last)

	default:
		dur, err := time.ParseDuration(entry.Expression)
		if err != nil {
			return false
		}
		return now.Sub(*entry.LastRun) >= dur
	}
}

func (s *Scheduler) saveEntry(entry *ScheduleEntry) error {
	os.MkdirAll(s.basePath, 0755)

	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}

	path := filepath.Join(s.basePath, entry.ID+".json")
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0644); err != nil {
		return fmt.Errorf("write temp: %w", err)
	}
	if err := os.Rename(tmp, path); err != nil {
		return fmt.Errorf("rename: %w", err)
	}
	return nil
}

func (s *Scheduler) log(format string, args ...any) {
	if s.logFn != nil {
		s.logFn("[scheduler] "+format, args...)
	}
}
