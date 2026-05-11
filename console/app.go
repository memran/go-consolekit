package console

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
)

type App struct {
	name          string
	version       string
	description   string
	registry      *Registry
	rootCmd       *cobra.Command
	rootCtx       context.Context
	rootCancel    context.CancelFunc
	renderer      Renderer
	daemonEnabled bool
	daemonArgs    []string
	pidFile       string
	logFile       string
}

func New(name string) *App {
	return &App{
		name:     name,
		registry: NewRegistry(),
		renderer: NewCLIRenderer(),
	}
}

func (a *App) Version(version string) *App {
	a.version = version
	return a
}

func (a *App) Description(description string) *App {
	a.description = description
	return a
}

func (a *App) Register(commands ...Command) *App {
	for _, cmd := range commands {
		a.registry.Add(cmd)
	}
	return a
}

func (a *App) Command(name string) *CommandBuilder {
	return newCommandBuilder(a, name)
}

func (a *App) registerBuilder(b *CommandBuilder) {
	cmd := &builderCommand{
		name:        b.name,
		description: b.desc,
		config:      b.config,
		handler:     b.handler,
	}
	a.registry.Add(cmd)
}

func (a *App) EnableDaemon() *App {
	a.daemonEnabled = true
	return a
}

func (a *App) PIDFile(path string) *App {
	a.pidFile = path
	return a
}

func (a *App) LogFile(path string) *App {
	a.logFile = path
	return a
}

func (a *App) Status() (int, error) {
	if a.pidFile == "" {
		return 0, fmt.Errorf("no PID file configured")
	}
	pid, err := readPID(a.pidFile)
	if err != nil {
		return 0, fmt.Errorf("daemon is not running")
	}
	if !processExists(pid) {
		return 0, fmt.Errorf("daemon is not running (stale PID %d)", pid)
	}
	return pid, nil
}

func (a *App) Stop() error {
	pid, err := readPID(a.pidFile)
	if err != nil {
		removePID(a.pidFile)
		return nil
	}
	proc, err := os.FindProcess(pid)
	if err != nil {
		removePID(a.pidFile)
		return nil
	}
	proc.Signal(syscall.SIGTERM)
	removePID(a.pidFile)
	return nil
}

func (a *App) Restart() error {
	a.Stop()
	return a.startChild()
}

func (a *App) startChild() error {
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("cannot determine executable path: %w", err)
	}
	args := a.daemonArgs
	if args == nil {
		args = os.Args
	}
	childArgs := stripDaemonFlag(args)
	env := append(os.Environ(), daemonEnvVar+"=1")

	nullDev, _ := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
	var files []*os.File
	if nullDev != nil {
		defer nullDev.Close()
		files = []*os.File{nil, nullDev, nullDev}
	} else {
		files = []*os.File{nil, nil, nil}
	}

	attr := &os.ProcAttr{
		Env:   env,
		Files: files,
	}
	proc, err := os.StartProcess(execPath, childArgs, attr)
	if err != nil {
		return fmt.Errorf("failed to start daemon: %w", err)
	}
	if a.pidFile != "" {
		dir := filepath.Dir(a.pidFile)
		os.MkdirAll(dir, 0755)
		os.WriteFile(a.pidFile, []byte(fmt.Sprintf("%d\n", proc.Pid)), 0644)
	}
	fmt.Printf("Daemon started with PID %d\n", proc.Pid)
	return nil
}

func (a *App) Run() error {
	a.rootCmd = &cobra.Command{
		Use:              a.name,
		Version:          a.version,
		Short:            a.description,
		SilenceUsage:     true,
		SilenceErrors:    true,
		TraverseChildren: true,
	}
	a.rootCmd.SetVersionTemplate(fmt.Sprintf("%s version {{.Version}}\n", a.name))

	if a.daemonEnabled {
		a.rootCmd.PersistentFlags().Bool("daemon", false, "Run as daemon")
		a.rootCmd.PersistentFlags().String("pid-file", "", "Path to PID file")
		a.rootCmd.PersistentFlags().String("log-file", "", "Path to log file")
	}

	for _, cmd := range a.registry.All() {
		cobraCmd := a.buildCobraCommand(cmd)
		a.attachCommand(a.rootCmd, strings.Split(cmd.Name(), ":"), cobraCmd)
	}

	if a.daemonEnabled && !IsDaemonChild() && hasDaemonFlag(os.Args) {
		a.daemonArgs = os.Args
		if pf := extractFlagValue(os.Args, "--pid-file"); pf != "" {
			a.pidFile = pf
		}
		if lf := extractFlagValue(os.Args, "--log-file"); lf != "" {
			a.logFile = lf
		}
		return a.startChild()
	}

	if IsDaemonChild() && a.daemonEnabled {
		if pf := extractFlagValue(os.Args, "--pid-file"); pf != "" {
			a.pidFile = pf
		}
		if lf := extractFlagValue(os.Args, "--log-file"); lf != "" {
			a.logFile = lf
		}
		if a.pidFile != "" {
			writePID(a.pidFile)
		}
		if a.logFile != "" {
			redirectOutput(a.logFile)
		}
	}

	a.rootCtx, a.rootCancel = context.WithCancel(context.Background())
	defer a.rootCancel()

	sigChan := make(chan os.Signal, 2)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigChan)

	var cmdErr error
	cmdDone := make(chan struct{})

	go func() {
		cmdErr = a.rootCmd.Execute()
		close(cmdDone)
	}()

	select {
	case <-cmdDone:
		return cmdErr
	case <-sigChan:
		a.renderer.Warning("Received interrupt. Shutting down gracefully (press again to force)...")
		a.rootCancel()

		select {
		case <-cmdDone:
			return fmt.Errorf("shutdown")
		case <-sigChan:
			a.renderer.Error("Forced exit.")
			os.Exit(1)
			return nil
		}
	}
}

func (a *App) attachCommand(parent *cobra.Command, parts []string, cmd *cobra.Command) {
	if len(parts) == 1 {
		parent.AddCommand(cmd)
		return
	}
	childName := parts[0]
	var child *cobra.Command
	for _, c := range parent.Commands() {
		if c.Name() == childName {
			child = c
			break
		}
	}
	if child == nil {
		child = &cobra.Command{
			Use:   childName,
			Short: childName,
		}
		parent.AddCommand(child)
	}
	a.attachCommand(child, parts[1:], cmd)
}

func (a *App) buildCobraCommand(cmd Command) *cobra.Command {
	config := NewCommandConfig(cmd.Name())
	cmd.Configure(config)

	parts := strings.Split(config.Name, ":")
	leafName := parts[len(parts)-1]

	useLine := leafName
	for _, arg := range config.Arguments {
		if arg.required {
			useLine += " <" + arg.Name + ">"
		} else {
			useLine += " [" + arg.Name + "]"
		}
	}

	cobraCmd := &cobra.Command{
		Use:   useLine,
		Short: cmd.Description(),
		RunE: func(cobraCmd *cobra.Command, args []string) error {
			parentCtx := a.rootCtx
			if parentCtx == nil {
				parentCtx = context.Background()
			}
			ctx := newContext(parentCtx)
			ctx.input = NewInput()
			ctx.output = NewOutput(a.renderer)

			if a.pidFile != "" && IsDaemonChild() {
				ctx.OnShutdown(func() error {
					return removePID(a.pidFile)
				})
			}

			for i, arg := range config.Arguments {
				if i < len(args) {
					ctx.argsMap[arg.Name] = args[i]
				} else if arg.defaultVal != "" {
					ctx.argsMap[arg.Name] = arg.defaultVal
				}
			}

			for _, opt := range config.Options {
				val, _ := cobraCmd.Flags().GetString(opt.Name)
				ctx.optionsMap[opt.Name] = val
			}

			err := cmd.Handle(ctx)
			ctx.runShutdown()
			return err
		},
	}

	for _, opt := range config.Options {
		if opt.shortcut != "" {
			cobraCmd.Flags().StringP(opt.Name, opt.shortcut, opt.defaultVal, opt.description)
		} else {
			cobraCmd.Flags().String(opt.Name, opt.defaultVal, opt.description)
		}
	}

	if len(config.Arguments) > 0 {
		required := 0
		for _, arg := range config.Arguments {
			if arg.required {
				required++
			}
		}
		cobraCmd.Args = func(cmd *cobra.Command, args []string) error {
			if len(args) < required {
				return fmt.Errorf("requires at least %d argument(s), received %d", required, len(args))
			}
			return nil
		}
	} else {
		cobraCmd.Args = cobra.NoArgs
	}

	return cobraCmd
}
