# Environment

`Env` loads `.env` files and provides typed access to environment variables. Merges loaded values with system environment variables.

## Constructor

```go
func NewEnv() *Env
```

## Methods

```go
func (e *Env) Load(path string) *Env
func (e *Env) LoadFile(path string) *Env
func (e *Env) Set(key, value string) *Env
func (e *Env) Get(key string, defaults ...string) string
func (e *Env) GetString(key string, defaults ...string) string
func (e *Env) GetInt(key string, defaults ...int) int
func (e *Env) GetBool(key string, defaults ...bool) bool
func (e *Env) GetFloat(key string, defaults ...float64) float64
func (e *Env) Has(key string) bool
func (e *Env) All() map[string]string
func (e *Env) Dump() *Env
```

- `Get` checks the system environment first, then the loaded `.env` data
- `Load` and `LoadFile` are aliases; errors silently return the current `Env`
- `Dump` prints all loaded key-value pairs to stdout

## .env File Format

```env
# Comments start with #
APP_NAME=myapp
APP_PORT=8080
APP_DEBUG=true
DB_HOST=localhost
DB_PASSWORD="secret123"
export EXPORTED_VAR=value
MULTI_LINE="line1\nline2"
```

Supports:
- `#` comments
- `export` prefix stripping
- Quoted values (single and double quotes)
- `\n`, `\r`, `\t` escape sequences in double-quoted strings

## Examples

```go
e := console.NewEnv()

// Load .env file
e.Load(".env")

// Set values programmatically
e.Set("APP_NAME", "myapp")

// Get with defaults
e.Get("APP_NAME")              // "myapp"
e.Get("DB_HOST")               // "localhost"
e.Get("MISSING", "fallback")   // "fallback"

// Typed getters
e.GetString("APP_NAME")        // "myapp"
e.GetInt("APP_PORT")           // 8080
e.GetBool("APP_DEBUG")         // true
e.GetFloat("SCORE")            // 0.0 (or parsed float)

// Check existence
e.Has("APP_NAME")              // true
e.Has("MISSING")               // false

// All variables (system + loaded)
all := e.All()
for k, v := range all {
    fmt.Printf("%s=%s\n", k, v)
}

// Dump to stdout
e.Dump()
```

## Chaining

```go
dbHost := console.NewEnv().
    Load(".env").
    Load(".env.local").
    Get("DB_HOST", "localhost")
```
