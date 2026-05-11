# Collections

`Collection` wraps `[]any` with a fluent pipeline API inspired by Laravel collections.

## Constructors

```go
func Collect(items []any) *Collection
func CollectFrom[T any](items []T) *Collection
```

`CollectFrom` accepts any typed slice and converts elements to `any`.

## Methods

### Accessors

```go
func (c *Collection) All() []any
func (c *Collection) Count() int
func (c *Collection) IsEmpty() bool
func (c *Collection) IsNotEmpty() bool
func (c *Collection) First() any
func (c *Collection) Last() any
func (c *Collection) Get(index int) any
func (c *Collection) Keys() *Collection
func (c *Collection) Values() *Collection
```

### Transformation

```go
func (c *Collection) Map(fn func(any) any) *Collection
func (c *Collection) Filter(fn func(any) bool) *Collection
func (c *Collection) Reject(fn func(any) bool) *Collection
func (c *Collection) Reduce(fn func(any, any) any, initial any) any
func (c *Collection) Each(fn func(any)) *Collection
func (c *Collection) Slice(offset, length int) *Collection
```

### Sorting & Ordering

```go
func (c *Collection) Sort(fn func(any, any) bool) *Collection
func (c *Collection) Reverse() *Collection
func (c *Collection) Shuffle() *Collection
```

### Set Operations

```go
func (c *Collection) Unique() *Collection
func (c *Collection) Chunk(size int) *Collection
func (c *Collection) Collapse() *Collection
func (c *Collection) Merge(items []any) *Collection
func (c *Collection) Diff(items []any) *Collection
func (c *Collection) Intersect(items []any) *Collection
```

### Search

```go
func (c *Collection) Contains(value any) bool
func (c *Collection) Search(value any) int
```

### Map Operations

```go
func (c *Collection) Where(key string, value any) *Collection
func (c *Collection) Pluck(key string) *Collection
func (c *Collection) GroupBy(key string) *Collection
func (c *Collection) KeyBy(key string) *Collection
```

### Aggregation

```go
func (c *Collection) Sum(key ...string) float64
func (c *Collection) Avg(key ...string) float64
func (c *Collection) Min(key ...string) float64
func (c *Collection) Max(key ...string) float64
```

### Output

```go
func (c *Collection) Implode(separator string) string
func (c *Collection) Join(separator, lastSeparator string) string
func (c *Collection) ToJSON() (string, error)
```

### Utility

```go
func (c *Collection) Tap(fn func(*Collection)) *Collection
```

## Examples

```go
items := []any{1, 2, 3, 4, 5, 6}
c := console.Collect(items)

// Map and filter
result := c.
    Filter(func(v any) bool { return v.(int) > 2 }).
    Map(func(v any) any { return v.(int) * 2 })
// [6, 8, 10, 12]

// Reduce
sum := c.Reduce(func(carry, v any) any {
    return carry.(int) + v.(int)
}, 0)
// 15

// Sort descending
sorted := c.Sort(func(a, b any) bool {
    return a.(int) > b.(int)
})
// [6, 5, 4, 3, 2, 1]

// GroupBy maps
users := console.Collect([]any{
    map[string]any{"role": "admin", "name": "Alice"},
    map[string]any{"role": "user", "name": "Bob"},
    map[string]any{"role": "admin", "name": "Charlie"},
})
groups := users.GroupBy("role")
// [[{role:admin,name:Alice},{role:admin,name:Charlie}], [{role:user,name:Bob}]]

// Pluck
users.Pluck("name") // ["Alice", "Bob", "Charlie"]

// Join with last separator
console.Collect([]any{"a", "b", "c"}).Join(", ", " and ")
// "a, b and c"

// Implode
console.Collect([]any{1, 2, 3}).Implode("-")
// "1-2-3"

// Sum/Avg/Min/Max
nums := console.Collect([]any{10, 20, 30, 40})
nums.Sum()   // 100
nums.Avg()   // 25
nums.Min()   // 10
nums.Max()   // 40

// Tap for side effects
nums.Tap(func(c *console.Collection) {
    fmt.Println("Count:", c.Count())
}).Map(...)

// KeyBy
users.KeyBy("name")
// [{key:Alice,value:{...}}, {key:Bob,value:{...}}]

// Chunk
c.Chunk(2) // [[1,2], [3,4], [5,6]]

// Diff / Intersect
a := console.Collect([]any{1, 2, 3, 4})
b := []any{2, 4, 6}
a.Diff(b)       // [1, 3]
a.Intersect(b)  // [2, 4]
```

## Typed Collections

```go
ints := []int{10, 20, 30, 40, 50}
c := console.CollectFrom(ints)
avg := c.Avg() // 30
```
