# Configuration

`Config` provides a dot-notation configuration store with JSON and YAML loading.

## Constructor

```go
func NewConfig() *Config
```

## Methods

```go
func (c *Config) Set(key string, value interface{}) *Config
func (c *Config) Get(key string, defaults ...interface{}) interface{}
func (c *Config) GetString(key string, defaults ...string) string
func (c *Config) GetInt(key string, defaults ...int) int
func (c *Config) GetBool(key string, defaults ...bool) bool
func (c *Config) Has(key string) bool
func (c *Config) All() map[string]interface{}
func (c *Config) Load(data map[string]interface{}) *Config
func (c *Config) LoadJSON(path string) *Config
func (c *Config) LoadYAML(path string) *Config
```

## Dot Notation

Keys use `.` as a path separator for nested access:

```go
c := console.NewConfig()
c.Set("database.host", "localhost")
c.Set("database.port", 5432)
c.Set("database.credentials.user", "admin")
c.Set("debug", true)

c.Get("database.host")                // "localhost"
c.Get("database.port")                // 5432
c.Get("database.credentials.user")    // "admin"
c.Get("database.credentials.pass")    // nil
c.Get("missing", "default")           // "default"

// Typed getters
c.GetString("database.host")          // "localhost"
c.GetString("missing", "fallback")    // "fallback"
c.GetInt("database.port")             // 5432
c.GetInt("missing")                   // 0
c.GetBool("debug")                    // true

// Check existence
c.Has("database.host")                // true
c.Has("database.timeout")             // false

// All returns the underlying map
all := c.All()
```

## Loading from Files

### JSON

```go
c := console.NewConfig()
c.LoadJSON("config.json")
// Reads: { "app": { "name": "myapp", "port": 8080 } }
c.GetString("app.name") // "myapp"
c.GetInt("app.port")    // 8080
```

### YAML

```go
c := console.NewConfig()
c.LoadYAML("config.yaml")
// Reads:
// app:
//   name: myapp
//   port: 8080
c.GetString("app.name") // "myapp"
c.GetInt("app.port")    // 8080
```

### Load from map

```go
c := console.NewConfig()
c.Load(map[string]interface{}{
    "name": "myapp",
    "features": map[string]interface{}{
        "auth": true,
        "logs": false,
    },
})
c.GetBool("features.auth") // true
```
