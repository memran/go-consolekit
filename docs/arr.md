# Array Utilities

`Arr` provides slice operations on `[]any`.

## Constructor

```go
func NewArr(items ...any) Arr
```

## Methods

```go
func (a Arr) Get(index int) any
func (a Arr) First() any
func (a Arr) Last() any
func (a Arr) Has(value any) bool
func (a Arr) Contains(value any) bool
func (a Arr) Where(fn func(any) bool) Arr
func (a Arr) Pluck(key string) Arr
func (a Arr) Flatten() Arr
func (a Arr) Collapse() Arr
func (a Arr) Chunk(size int) Arr
func (a Arr) Unique() Arr
func (a Arr) Wrap(value any) Arr
func (a Arr) Join(separator string) string
func (a Arr) ToJSON() (string, error)
```

## Examples

```go
arr := console.NewArr(1, 2, 3, 4, 5)

arr.First()                    // 1
arr.Last()                     // 5
arr.Get(2)                     // 3
arr.Get(10)                    // nil
arr.Has(3)                     // true
arr.Contains(6)                // false

arr.Where(func(v any) bool {
    return v.(int) > 3
})                             // [4, 5]

// Pluck from map items
items := console.NewArr(
    map[string]any{"id": 1, "name": "Alice"},
    map[string]any{"id": 2, "name": "Bob"},
)
items.Pluck("name")            // ["Alice", "Bob"]

// Flatten nested arrays
nested := console.NewArr(
    console.NewArr(1, 2),
    console.NewArr(3, console.NewArr(4, 5)),
)
nested.Flatten()               // [1, 2, 3, 4, 5]

// Collapse one level
nested2 := console.NewArr(
    console.NewArr(1, 2),
    console.NewArr(3, 4),
)
nested2.Collapse()             // [1, 2, 3, 4]

// Chunk
console.NewArr(1, 2, 3, 4, 5).Chunk(2)
// [[1, 2], [3, 4], [5]]

// Unique
console.NewArr(1, 2, 2, 3, 3, 3).Unique()
// [1, 2, 3]

// Wrap
console.Arr.Wrap("hello")      // ["hello"]
console.Arr.Wrap([]any{1, 2})  // [1, 2]

// Join
console.NewArr("a", "b", "c").Join(", ")
// "a, b, c"

// ToJSON
json, _ := console.NewArr(1, "two", true).ToJSON()
// `[1,"two",true]`
```
