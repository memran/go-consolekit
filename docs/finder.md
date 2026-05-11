# File Finder

Recursive file searching with name, size, type, depth, and date filters.

## Constructor

```go
func NewFinder() *Finder
```

## Methods

```go
func (f *Finder) In(dirs ...string) *Finder
func (f *Finder) Name(patterns ...string) *Finder
func (f *Finder) Size(min, max int64) *Finder
func (f *Finder) Depth(d int) *Finder
func (f *Finder) Type(ftype string) *Finder
func (f *Finder) Date(after, before time.Time) *Finder
func (f *Finder) Modified(since time.Time) *Finder
func (f *Finder) Find() []string
```

### Name

Name matching supports glob patterns (`filepath.Match`). If no glob match is found, falls back to substring matching. When no name patterns are set, all names match.

```go
console.NewFinder().
    In(".").
    Name("*.go", "*.md").
    Find()
```

### Size

File size filtering in bytes.

```go
console.NewFinder().
    In(".").
    Size(1024, 1048576). // between 1KB and 1MB
    Find()
```

### Depth

Limit directory traversal depth. `-1` means unlimited.

```go
console.NewFinder().
    In(".").
    Depth(2). // 2 levels deep
    Find()
```

### Type

Filter by entry type: `"file"` or `"dir"`.

```go
console.NewFinder().
    In(".").
    Type("file").
    Find()

console.NewFinder().
    In(".").
    Type("dir").
    Find()
```

### Date

Filter by modification time range.

```go
since := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
until := time.Now()

console.NewFinder().
    In(".").
    Date(since, until).
    Find()
```

### Modified

Shorthand for filtering files modified since a given time.

```go
console.NewFinder().
    In(".").
    Modified(time.Now().Add(-24 * time.Hour)). // last 24 hours
    Find()
```

## Complete Example

```go
results := console.NewFinder().
    In("/var/log", "/tmp").
    Name("*.log").
    Size(0, 1048576).
    Depth(3).
    Type("file").
    Modified(time.Now().Add(-7 * 24 * time.Hour)).
    Find()

for _, path := range results {
    fmt.Println(path)
}
```
