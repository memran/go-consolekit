# String Utilities

`Str` provides fluent string manipulation via method chaining.

## Constructor

```go
func NewStr(s string) *Str
```

## Methods

### Inspection

```go
func (s *Str) String() string
func (s *Str) Length() int
func (s *Str) IsEmpty() bool
func (s *Str) IsNotEmpty() bool
func (s *Str) Equals(other string) bool
func (s *Str) EqualsIgnoreCase(other string) bool
```

### Contains

```go
func (s *Str) Contains(substr string) bool
func (s *Str) ContainsAll(substrs ...string) bool
func (s *Str) ContainsAny(substrs ...string) bool
```

### Prefix/Suffix

```go
func (s *Str) StartsWith(prefix string) bool
func (s *Str) EndsWith(suffix string) bool
```

### Case Conversion

```go
func (s *Str) Upper() *Str
func (s *Str) Lower() *Str
func (s *Str) Title() *Str
func (s *Str) Ucfirst() *Str
func (s *Str) Lcfirst() *Str
func (s *Str) Slug(separator ...string) *Str
func (s *Str) Snake() *Str
func (s *Str) Kebab() *Str
func (s *Str) Studly() *Str
func (s *Str) Camel() *Str
```

### Substring

```go
func (s *Str) Limit(length int, append ...string) *Str
func (s *Str) Substr(start, length int) *Str
```

### Replace

```go
func (s *Str) Replace(search, replace string) *Str
func (s *Str) ReplaceFirst(search, replace string) *Str
func (s *Str) ReplaceLast(search, replace string) *Str
```

### Trim

```go
func (s *Str) Trim(cutset ...string) *Str
func (s *Str) Ltrim(cutset ...string) *Str
func (s *Str) Rtrim(cutset ...string) *Str
```

### Pad

```go
func (s *Str) PadLeft(length int, pad ...string) *Str
func (s *Str) PadRight(length int, pad ...string) *Str
```

### Extract

```go
func (s *Str) Before(delimiter string) *Str
func (s *Str) After(delimiter string) *Str
func (s *Str) Between(from, to string) *Str
```

### Pattern Matching

```go
func (s *Str) Is(pattern string) bool
```

Supports `*` wildcards. `NewStr("hello.go").Is("*.go")` returns `true`.

### Masking

```go
func (s *Str) Mask(ch rune, start, length int) *Str
```

### Repeat

```go
func (s *Str) Repeat(count int) *Str
```

### Random

```go
func Random(length int) string
```

Package-level function for random alphanumeric strings.

## Examples

```go
s := console.NewStr("hello world")

// Chaining
result := s.Trim().Upper().Replace("HELLO", "HI")
// "HI WORLD"

// Case conversion
console.NewStr("helloWorld").Snake()     // "hello_world"
console.NewStr("hello_world").Studly()   // "HelloWorld"
console.NewStr("HelloWorld").Camel()     // "helloWorld"
console.NewStr("hello world").Kebab()    // "hello-world"
console.NewStr("Hello World").Slug()     // "hello-world"
console.NewStr("Hello World").Slug("_") // "hello_world"

// Limiting
console.NewStr("hello.txt").Limit(5)        // "hello..."
console.NewStr("hello").Limit(10)           // "hello"

// Substring
console.NewStr("hello").Substr(1, 3)        // "ell"
console.NewStr("hello").Substr(-3, 2)       // "ll"

// Before/After/Between
console.NewStr("user@example.com").Before("@") // "user"
console.NewStr("user@example.com").After("@")  // "example.com"
console.NewStr("{{name}}").Between("{{", "}}") // "name"

// Padding
console.NewStr("5").PadLeft(3, "0")   // "005"
console.NewStr("hi").PadRight(5)      // "hi   "

// Masking
console.NewStr("hello@example.com").Mask('*', 0, 5)
// "*****@example.com"

// Pattern matching
console.NewStr("readme.md").Is("*.md")      // true
console.NewStr("config.json").Is("*.json")  // true

// Trimming
console.NewStr("  spaced  ").Trim()         // "spaced"
console.NewStr("...hello...").Trim(".")     // "hello"

// ReplaceLast
console.NewStr("a/b/c").ReplaceLast("/", ":") // "a/b:c"

// Random
rand := console.Random(16) // e.g. "a3Bx9K2mQ7wE5nR1"
```
