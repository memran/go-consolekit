# go-consolekit

A Go CLI framework with utilities for building console applications.

## Features

- **App Framework** - Build CLI apps with commands, arguments, options, and subcommands using `App`, `CommandBuilder`, `Registry`, and `Context`
- **Daemon Mode** - Run commands as background daemons with PID files, log redirection, and graceful shutdown
- **TUI Renderer** - Full-screen Bubble Tea terminal UI with scrollable styled output, set via `app.Renderer(NewTUIRenderer())`
- **Event Bus** - Priority-based pub/sub with wildcard patterns, async emission, and panic recovery
- **Task Queue** - File-based queue with configurable concurrency, retry, and scheduled tasks
- **Scheduler** - File-based cron scheduler (`@daily`, `@hourly`, duration expressions) that pushes tasks to the Queue
- **String Utilities** (`Str`) - Fluent string manipulation: case conversion, substring, padding, slug, masking, pattern matching
- **Array Utilities** (`Arr`) - Slice operations: pluck, flatten, collapse, chunk, unique, wrap
- **Collections** (`Collection`) - Map, filter, reduce, sort, group, aggregate with fluent API
- **Date/Time** (`Date`) - Parsing, formatting, diff, human diff, start/end of period, comparison
- **Dot-Notation Object** (`Obj`) - Nested map access with dot notation, typed getters
- **Security** - SHA256 hashing, AES-GCM encryption/decryption, random bytes/strings
- **Validation** - Declarative field validation with rule chaining
- **File Finder** - Recursive file search with name/size/type/date filters
- **File Operations** (`FileOps`) - Read, write, append, prepend, copy, move, info
- **HTTP Client** (`HttpClient`) - Fluent HTTP client with retries, auth, JSON support
- **Configuration** (`Config`) - Dot-notation config with JSON/YAML loading
- **Environment** (`Env`) - `.env` file loading, typed getters, system env merging
- **Logger** (`Logger`) - Structured logging with levels, channels, context, key-value pairs
- **Interactive Input** - Ask, confirm, select, secret prompts using survey
- **Styled Output** - Colored text, progress bars, spinners, tables via Output

## Package

```
go get github.com/user/go-consolekit
```

```go
import "github.com/user/go-consolekit/console"
```

## Quick Start

```go
package main

import "github.com/user/go-consolekit/console"

func main() {
    app := console.New("myapp").
        Version("1.0.0").
        Description("My CLI application")

    app.Command("greet").
        Description("Say hello").
        Argument("name").Required().
        Option("lang").Default("en").
        Handle(func(ctx *console.Context) error {
            ctx.Success("Hello, " + ctx.Arg("name"))
            return nil
        })

    app.Run()
}
```

## Documentation

- [Console Framework](console.md) - App, CommandBuilder, Registry, Context, Renderer, Errors
- [Daemon Mode](daemon.md) - Background daemon process management
- [String Utilities](str.md) - Str type and string manipulation
- [Array Utilities](arr.md) - Arr type and array operations
- [Collections](collection.md) - Collection type with map/filter/reduce
- [Date/Time](date.md) - Date type for date manipulation
- [Dot-Notation Object](obj.md) - Obj type for nested map access
- [Security](security.md) - Hashing, encryption, random generation
- [Validation](validation.md) - Field validation rules
- [File Finder](finder.md) - Recursive file searching
- [File Operations](fileops.md) - File read/write/manage
- [HTTP Client](http.md) - HTTP requests with fluent API
- [Configuration](config.md) - Config with JSON/YAML loading
- [Environment](env.md) - Environment variable management
- [Logger](logger.md) - Structured logging

## Types Overview

| Type | Constructor | Purpose |
|------|-------------|---------|
| `App` | `New(name)` | CLI application |
| `Str` | `NewStr(s)` | String manipulation |
| `Arr` | `NewArr(items...)` | Array operations |
| `Collection` | `Collect(items)` / `CollectFrom[T](items)` | Collection pipeline |
| `Date` | `Now()` / `NewDate(t)` / `Parse(s)` | Date/time |
| `Obj` | `NewObj(data)` | Dot-notation object |
| `Finder` | `NewFinder()` | File search |
| `FileOps` | `File(path)` | File operations |
| `HttpClient` | `Http()` | HTTP client |
| `Config` | `NewConfig()` | Configuration store |
| `Env` | `NewEnv()` | Environment variables |
| `Logger` | `NewLogger()` | Structured logger |
| `Validator` | `NewValidator()` | Data validation |
| `Input` | `NewInput()` | Interactive prompts |
| `Output` | `NewOutput(renderer)` | Styled output |
| `Renderer` | `NewCLIRenderer()` / `NewTUIRenderer()` | Output rendering |
