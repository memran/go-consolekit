package console

import (
	"context"
	"fmt"
	"os"
	"sync"
)

type Context struct {
	ctx           context.Context
	cancel        context.CancelFunc
	shutdownOnce  sync.Once
	shutdownFns   []func() error
	mu            sync.Mutex
	argsMap       map[string]string
	optionsMap    map[string]string
	input         *Input
	output        *Output
}

func newContext(parent context.Context) *Context {
	ctx, cancel := context.WithCancel(parent)
	return &Context{
		ctx:        ctx,
		cancel:     cancel,
		argsMap:    make(map[string]string),
		optionsMap: make(map[string]string),
	}
}

func (c *Context) Arg(name string) string {
	return c.argsMap[name]
}

func (c *Context) Option(name string) string {
	return c.optionsMap[name]
}

func (c *Context) Input() *Input {
	return c.input
}

func (c *Context) Output() *Output {
	return c.output
}

func (c *Context) Line(text string) {
	c.output.Line(text)
}

func (c *Context) Info(text string) {
	c.output.Info(text)
}

func (c *Context) Success(text string) {
	c.output.Success(text)
}

func (c *Context) Warning(text string) {
	c.output.Warning(text)
}

func (c *Context) Error(text string) {
	c.output.Error(text)
}

func (c *Context) Title(text string) {
	c.output.Title(text)
}

func (c *Context) Done() <-chan struct{} {
	return c.ctx.Done()
}

func (c *Context) IsCancelled() bool {
	select {
	case <-c.ctx.Done():
		return true
	default:
		return false
	}
}

func (c *Context) OnShutdown(fn func() error) {
	c.mu.Lock()
	c.shutdownFns = append(c.shutdownFns, fn)
	c.mu.Unlock()
}

func (c *Context) Shutdown(message string) {
	if message != "" {
		c.output.Warning(message)
	}
	c.cancel()
	c.runShutdown()
}

func (c *Context) runShutdown() {
	c.shutdownOnce.Do(func() {
		c.mu.Lock()
		fns := make([]func() error, len(c.shutdownFns))
		copy(fns, c.shutdownFns)
		c.mu.Unlock()

		for _, fn := range fns {
			fn()
		}
	})
}

func (c *Context) IsDaemon() bool {
	return IsDaemonChild()
}

func (c *Context) PID() int {
	return os.Getpid()
}

func (c *Context) WritePID(path string) error {
	return writePID(path)
}

func (c *Context) StopDaemon() error {
	c.cancel()
	c.runShutdown()
	return nil
}

func (c *Context) RunWorker(fn func() error) {
	for {
		select {
		case <-c.Done():
			return
		default:
			func() {
				defer func() {
					if r := recover(); r != nil {
						c.output.Error(fmt.Sprintf("Worker panic: %v", r))
					}
				}()
				fn()
			}()
		}
	}
}
