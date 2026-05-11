# Dot-Notation Object

`Obj` provides nested map access using dot-notation keys.

## Constructor

```go
func NewObj(data ...map[string]any) *Obj
```

With no arguments, creates an empty `Obj`. Optionally accepts an initial map.

## Methods

```go
func (o *Obj) Set(key string, value any) *Obj
func (o *Obj) Get(key string, defaults ...any) any
func (o *Obj) GetString(key string, defaults ...string) string
func (o *Obj) GetInt(key string, defaults ...int) int
func (o *Obj) GetBool(key string, defaults ...bool) bool
func (o *Obj) Has(key string) bool
func (o *Obj) Forget(key string) *Obj
func (o *Obj) All() map[string]any
func (o *Obj) ToJSON() string
func (o *Obj) Count() int
func (o *Obj) IsEmpty() bool
```

## Examples

```go
obj := console.NewObj()

// Set via dot notation
obj.Set("user.name", "Alice")
obj.Set("user.email", "alice@example.com")
obj.Set("user.address.city", "New York")
obj.Set("enabled", true)
obj.Set("count", 42)

// Get with dot notation
obj.Get("user.name")              // "Alice"
obj.Get("user.address.city")      // "New York"
obj.Get("user.address.zip")       // nil
obj.Get("missing", "default")     // "default"

// Typed getters
obj.GetString("user.name")        // "Alice"
obj.GetString("missing", "n/a")   // "n/a"
obj.GetInt("count")               // 42
obj.GetInt("missing", 10)         // 10
obj.GetBool("enabled")            // true

// Check existence
obj.Has("user.name")              // true
obj.Has("user.address.zip")       // false

// Remove keys
obj.Forget("user.address.city")

// All returns a copy
all := obj.All()
// map[string]any{
//     "user": map[string]any{
//         "name": "Alice",
//         "email": "alice@example.com",
//         "address": map[string]any{},
//     },
//     "enabled": true,
//     "count": 42,
// }

// JSON output
obj.ToJSON()
// `{"count":42,"enabled":true,"user":{"address":{},"email":"alice@example.com","name":"Alice"}}`

// Count and empty check
obj.Count()                       // 3 (top-level keys)
obj.IsEmpty()                     // false

// Create from existing map
data := map[string]any{
    "host": "localhost",
    "port": 8080,
}
obj2 := console.NewObj(data)
obj2.GetString("host")           // "localhost"

// Deeply nested keys
obj.Set("a.b.c.d.e", "deep")
obj.Get("a.b.c.d.e")             // "deep"
obj.Has("a.b.c")                 // true
obj.Forget("a.b.c.d")
obj.Get("a.b.c.d")               // nil
```
