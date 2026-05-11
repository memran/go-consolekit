package console

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"
)

type Process struct {
	name     string
	args     []string
	dir      string
	env      map[string]string
	stdin    io.Reader
	timeout  time.Duration
	cmd      *exec.Cmd
	running  bool
	captured bool
	stdout   bytes.Buffer
	stderr   bytes.Buffer
}

type Result struct {
	exitCode int
	stdout   string
	stderr   string
	err      error
}

func NewProcess(name string, args ...string) *Process {
	return &Process{
		name: name,
		args: args,
		env:  make(map[string]string),
	}
}

func Run(name string, args ...string) *Process {
	return NewProcess(name, args...)
}

func (p *Process) WithWorkingDir(dir string) *Process {
	p.dir = dir
	return p
}

func (p *Process) WithEnv(key, value string) *Process {
	p.env[key] = value
	return p
}

func (p *Process) WithEnvs(env map[string]string) *Process {
	for k, v := range env {
		p.env[k] = v
	}
	return p
}

func (p *Process) WithInput(input io.Reader) *Process {
	p.stdin = input
	return p
}

func (p *Process) WithInputString(input string) *Process {
	p.stdin = bytes.NewBufferString(input)
	return p
}

func (p *Process) Timeout(timeout time.Duration) *Process {
	p.timeout = timeout
	return p
}

func (p *Process) buildCmd() *exec.Cmd {
	cmd := exec.Command(p.name, p.args...)
	if p.dir != "" {
		cmd.Dir = p.dir
	}
	if len(p.env) > 0 {
		cmd.Env = os.Environ()
		for k, v := range p.env {
			cmd.Env = append(cmd.Env, k+"="+v)
		}
	}
	if p.stdin != nil {
		cmd.Stdin = p.stdin
	}
	return cmd
}

func (p *Process) Run() *Result {
	cmd := p.buildCmd()
	cmd.Stdout = &p.stdout
	cmd.Stderr = &p.stderr
	p.cmd = cmd

	if p.timeout > 0 {
		timer := time.AfterFunc(p.timeout, func() {
			if cmd.Process != nil {
				cmd.Process.Kill()
			}
		})
		defer timer.Stop()
	}

	err := cmd.Run()
	p.running = false

	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = -1
		}
	}

	return &Result{
		exitCode: exitCode,
		stdout:   p.stdout.String(),
		stderr:   p.stderr.String(),
		err:      err,
	}
}

func (p *Process) MustRun() *Result {
	result := p.Run()
	if result.err != nil {
		panic(fmt.Sprintf("process failed: %s\n%s", result.err, result.stderr))
	}
	return result
}

func (p *Process) Start() error {
	cmd := p.buildCmd()
	cmd.Stdout = &p.stdout
	cmd.Stderr = &p.stderr
	p.cmd = cmd

	if p.timeout > 0 {
		timer := time.AfterFunc(p.timeout, func() {
			if cmd.Process != nil {
				cmd.Process.Kill()
			}
		})
		defer func() {
			if !p.running {
				timer.Stop()
			}
		}()
	}

	err := cmd.Start()
	if err != nil {
		return err
	}
	p.running = true
	return nil
}

func (p *Process) Wait() *Result {
	if p.cmd == nil || !p.running {
		return &Result{exitCode: -1, err: fmt.Errorf("process not started")}
	}
	err := p.cmd.Wait()
	p.running = false

	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = -1
		}
	}

	return &Result{
		exitCode: exitCode,
		stdout:   p.stdout.String(),
		stderr:   p.stderr.String(),
		err:      err,
	}
}

func (p *Process) Stop() error {
	if p.cmd == nil || p.cmd.Process == nil {
		return fmt.Errorf("process not started")
	}
	return p.cmd.Process.Kill()
}

func (p *Process) Signal(sig os.Signal) error {
	if p.cmd == nil || p.cmd.Process == nil {
		return fmt.Errorf("process not started")
	}
	return p.cmd.Process.Signal(sig)
}

func (p *Process) PID() int {
	if p.cmd == nil || p.cmd.Process == nil {
		return 0
	}
	return p.cmd.Process.Pid
}

func (p *Process) IsRunning() bool {
	return p.running
}

func (p *Process) String() string {
	return fmt.Sprintf("%s %v", p.name, p.args)
}

func (r *Result) ExitCode() int {
	return r.exitCode
}

func (r *Result) Output() string {
	return r.stdout
}

func (r *Result) Error() string {
	return r.stderr
}

func (r *Result) Err() error {
	return r.err
}

func (r *Result) IsSuccessful() bool {
	return r.exitCode == 0 && r.err == nil
}

func (r *Result) IsFailed() bool {
	return !r.IsSuccessful()
}

func (r *Result) Lines() []string {
	return splitLines(r.stdout)
}

func (r *Result) ErrorLines() []string {
	return splitLines(r.stderr)
}

func splitLines(s string) []string {
	if s == "" {
		return nil
	}
	var lines []string
	buf := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, string(buf))
			buf = buf[:0]
		} else if s[i] != '\r' {
			buf = append(buf, s[i])
		}
	}
	if len(buf) > 0 {
		lines = append(lines, string(buf))
	}
	return lines
}
