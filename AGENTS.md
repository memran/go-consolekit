# go-consolekit

A Go framework for building console applications. Single `console` package in `console/`, module path `go-consolekit`, Go 1.25.

## Architecture

- **`console/`** — the library (52 files, all `package console`)
- **`cmd/example/`** — example app with 20+ command demos (entrypoint: `cmd/example/main.go`)
- **Renderer** abstraction: `CLIRenderer` (terminal output via `fatih/color`), `TUIRenderer` (Bubble Tea TUI via `charmbracelet/bubbletea` + `lipgloss`, scrollable output, q/Ctrl+C to quit). Set via `app.Renderer(r)`. Lazy-starts on first output call. `Stop()` and `Wait()` methods for lifecycle.
- **Underlying CLI parser**: `spf13/cobra` — commands registered via `Command` interface or fluent `CommandBuilder`
- **Daemon mode**: `EnableDaemon()`, `--daemon` flag, PID/redirect to file. Platform-specific `processExists` in `daemon_unix.go` / `daemon_windows.go` (build constraints). `Context.IsDaemon()` / `Context.StopDaemon()` / `Context.RunWorker()`

## Non-trivial subsystems (all in `console/`)

| Package area | File(s) | Notes |
|---|---|---|
| Event bus | `event.go` | Priority listeners, wildcard patterns (`*`), `EmitAsync`/`Flush`, `Subscriber` interface, panic recovery |
| Queue | `queue.go` | File-based (JSON on disk under `<basePath>/queues/<name>/<status>/`), configurable concurrency, retry, scheduled tasks |
| Scheduler | `scheduler.go` | File-based JSON entries, `@daily`/`@hourly`/duration expressions, pushes tasks to Queue |
| Config | `appconfig.go` | Dot-notation, loads JSON/YAML files |
| Logger | `logger.go` | Channels, structured context, named loggers, clone-on-write |
| Utilities | `str.go`, `arr.go`, `obj.go`, `collection.go`, `date.go`, `env.go`, `fileops.go`, `finder.go`, `http.go`, `process.go`, `security.go`, `validation.go` | Standalone helpers, no required wiring |

## Dev commands

```sh
# Test (all)
go test ./console/

# Single test
go test -run 'TestCommandBuilderRegistration' ./console/

# Run example commands
go run ./cmd/example hello Emran
go run ./cmd/example install blog --db postgres

# Build
go build ./...
```

**Note**: Full test suite may get `signal: terminated` due to terminal-dependent renderer tests — use `-run` to filter.

## Key behaviors

- **Registry panics** on duplicate command registration (`registry.go:16`)
- **`Context`** is cratead per-command invocation with `newContext(parent)` — exposes `Arg()`, `Option()`, `Input()`, `Output()`, shutdown hooks, daemon helpers
- **Commands** can be namespaced with `:` (e.g. `admin:users:list`), mapped to nested cobra subcommands via `attachCommand`
- **Fluent builder**: `app.Command("name").Argument("x").Required().Description("...")` returns `CommandBuilder`, `.Handle(fn)` finalizes
- **`CollectFrom[T any]`** helper for generic collection construction (`collection.go`)
- **`question` tool prompts** are interactive and may block in non-terminal environments
