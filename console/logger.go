package console

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

type LogLevel int

const (
	DebugLevel LogLevel = iota
	InfoLevel
	NoticeLevel
	WarningLevel
	ErrorLevel
	CriticalLevel
	AlertLevel
	EmergencyLevel
)

func (l LogLevel) String() string {
	switch l {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case NoticeLevel:
		return "NOTICE"
	case WarningLevel:
		return "WARNING"
	case ErrorLevel:
		return "ERROR"
	case CriticalLevel:
		return "CRITICAL"
	case AlertLevel:
		return "ALERT"
	case EmergencyLevel:
		return "EMERGENCY"
	default:
		return "INFO"
	}
}

type Logger struct {
	mu       sync.Mutex
	writers  []io.Writer
	channels map[string][]io.Writer
	context  map[string]interface{}
	name     string
}

var defaultLogger = NewLogger()

func NewLogger() *Logger {
	return &Logger{
		writers:  []io.Writer{os.Stdout},
		channels: make(map[string][]io.Writer),
		context:  make(map[string]interface{}),
	}
}

func (l *Logger) With(key string, value interface{}) *Logger {
	nl := l.clone()
	nl.context[key] = value
	return nl
}

func (l *Logger) WithMap(data map[string]interface{}) *Logger {
	nl := l.clone()
	for k, v := range data {
		nl.context[k] = v
	}
	return nl
}

func (l *Logger) WithWriter(writer io.Writer) *Logger {
	nl := l.clone()
	nl.writers = []io.Writer{writer}
	return nl
}

func (l *Logger) WithName(name string) *Logger {
	nl := l.clone()
	nl.name = name
	return nl
}

func (l *Logger) Channel(name string) *Logger {
	nl := l.clone()
	if writers, ok := l.channels[name]; ok {
		nl.writers = writers
	} else {
		nl.writers = l.writers
	}
	return nl
}

func (l *Logger) Stack(names ...string) *Logger {
	nl := l.clone()
	var all []io.Writer
	for _, name := range names {
		if writers, ok := l.channels[name]; ok {
			all = append(all, writers...)
		}
	}
	if len(all) == 0 {
		all = l.writers
	}
	nl.writers = all
	return nl
}

func (l *Logger) AddChannel(name string, writer io.Writer) *Logger {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.channels[name] = append(l.channels[name], writer)
	return l
}

func (l *Logger) Debug(msg string, kv ...interface{}) {
	l.log(DebugLevel, msg, kv...)
}

func (l *Logger) Info(msg string, kv ...interface{}) {
	l.log(InfoLevel, msg, kv...)
}

func (l *Logger) Notice(msg string, kv ...interface{}) {
	l.log(NoticeLevel, msg, kv...)
}

func (l *Logger) Warning(msg string, kv ...interface{}) {
	l.log(WarningLevel, msg, kv...)
}

func (l *Logger) Error(msg string, kv ...interface{}) {
	l.log(ErrorLevel, msg, kv...)
}

func (l *Logger) Critical(msg string, kv ...interface{}) {
	l.log(CriticalLevel, msg, kv...)
}

func (l *Logger) Alert(msg string, kv ...interface{}) {
	l.log(AlertLevel, msg, kv...)
}

func (l *Logger) Emergency(msg string, kv ...interface{}) {
	l.log(EmergencyLevel, msg, kv...)
}

func (l *Logger) Log(level LogLevel, msg string, kv ...interface{}) {
	l.log(level, msg, kv...)
}

func (l *Logger) log(level LogLevel, msg string, kv ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	ctx := make(map[string]interface{})
	for k, v := range l.context {
		ctx[k] = v
	}
	for i := 0; i < len(kv)-1; i += 2 {
		key, ok := kv[i].(string)
		if ok {
			ctx[key] = kv[i+1]
		}
	}

	line := l.format(level, msg, ctx)

	for _, w := range l.writers {
		fmt.Fprintln(w, line)
	}
}

func (l *Logger) format(level LogLevel, msg string, ctx map[string]interface{}) string {
	ts := time.Now().Format("2006-01-02 15:04:05")
	name := l.name
	if name == "" {
		name = "local"
	}
	line := fmt.Sprintf("[%s] %s.%s: %s", ts, name, level, msg)
	if len(ctx) > 0 {
		encoded, err := json.Marshal(ctx)
		if err == nil {
			line += " " + string(encoded)
		}
	}
	return line
}

func (l *Logger) clone() *Logger {
	writers := make([]io.Writer, len(l.writers))
	copy(writers, l.writers)

	channels := make(map[string][]io.Writer, len(l.channels))
	for k, v := range l.channels {
		c := make([]io.Writer, len(v))
		copy(c, v)
		channels[k] = c
	}

	ctx := make(map[string]interface{}, len(l.context))
	for k, v := range l.context {
		ctx[k] = v
	}

	return &Logger{
		writers:  writers,
		channels: channels,
		context:  ctx,
		name:     l.name,
	}
}

func Info(msg string, kv ...interface{}) {
	defaultLogger.Info(msg, kv...)
}

func Debug(msg string, kv ...interface{}) {
	defaultLogger.Debug(msg, kv...)
}

func Notice(msg string, kv ...interface{}) {
	defaultLogger.Notice(msg, kv...)
}

func Warning(msg string, kv ...interface{}) {
	defaultLogger.Warning(msg, kv...)
}

func Error(msg string, kv ...interface{}) {
	defaultLogger.Error(msg, kv...)
}

func Critical(msg string, kv ...interface{}) {
	defaultLogger.Critical(msg, kv...)
}

func Alert(msg string, kv ...interface{}) {
	defaultLogger.Alert(msg, kv...)
}

func Emergency(msg string, kv ...interface{}) {
	defaultLogger.Emergency(msg, kv...)
}

func (l *Logger) formatContext(kv []interface{}) string {
	if len(kv) == 0 {
		return ""
	}
	parts := make([]string, 0, len(kv)/2)
	for i := 0; i < len(kv)-1; i += 2 {
		parts = append(parts, fmt.Sprintf("%v=%v", kv[i], kv[i+1]))
	}
	return strings.Join(parts, " ")
}
