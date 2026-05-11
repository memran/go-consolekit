# Date/Time

`Date` wraps `time.Time` with fluent manipulation, formatting, and diff methods.

## Constructors

```go
func Now() *Date
func Today() *Date
func Yesterday() *Date
func Tomorrow() *Date
func Parse(value string) *Date
func ParseDate(value string) *Date
func NewDate(t ...time.Time) *Date
```

- `NewDate()` with no arguments uses `time.Now()`
- `Parse` and `ParseDate` auto-detect multiple formats (RFC3339, ISO8601, `2006-01-02`, etc.)
- `Today` returns the current date at midnight
- `Yesterday` / `Tomorrow` are relative to `Now()`

## Conversion

```go
func (d *Date) Time() time.Time
func (d *Date) Timestamp() int64
func (d *Date) String() string              // "2006-01-02 15:04:05"
func (d *Date) Format(layout string) string
func (d *Date) ToDateTimeString() string    // "2006-01-02 15:04:05"
func (d *Date) ToDateString() string        // "2006-01-02"
func (d *Date) ToTimeString() string        // "15:04:05"
func (d *Date) Rfc3339() string
func (d *Date) ISO8601() string
```

## Component Getters

```go
func (d *Date) Year() int
func (d *Date) Month() time.Month
func (d *Date) Day() int
func (d *Date) Hour() int
func (d *Date) Minute() int
func (d *Date) Second() int
func (d *Date) Weekday() time.Weekday
```

## Arithmetic

```go
func (d *Date) AddDays(n int) *Date
func (d *Date) SubDays(n int) *Date
func (d *Date) AddHours(n int) *Date
func (d *Date) SubHours(n int) *Date
func (d *Date) AddMinutes(n int) *Date
func (d *Date) AddMonths(n int) *Date
func (d *Date) AddYears(n int) *Date
```

## Start/End of Period

```go
func (d *Date) StartOfDay() *Date
func (d *Date) EndOfDay() *Date
func (d *Date) StartOfMonth() *Date
func (d *Date) EndOfMonth() *Date
func (d *Date) StartOfYear() *Date
func (d *Date) EndOfYear() *Date
```

## Comparison

```go
func (d *Date) Gt(other *Date) bool
func (d *Date) Lt(other *Date) bool
func (d *Date) Eq(other *Date) bool
func (d *Date) IsWeekend() bool
func (d *Date) IsWeekday() bool
func (d *Date) IsPast() bool
func (d *Date) IsFuture() bool
func (d *Date) IsToday() bool
```

## Difference

```go
func (d *Date) DiffInDays(other *Date) int
func (d *Date) DiffInHours(other *Date) float64
func (d *Date) DiffInMinutes(other *Date) float64
func (d *Date) DiffInSeconds(other *Date) float64
func (d *Date) Age() int
func (d *Date) HumanDiff() string
```

## Utility

```go
func (d *Date) Copy() *Date
```

## Examples

```go
// Create dates
now := console.Now()
today := console.Today()
yesterday := console.Yesterday()
tomorrow := console.Tomorrow()
parsed := console.Parse("2024-01-15")
fromTime := console.NewDate(time.Now())

// Formatting
now.Format("2006-01-02")       // "2025-05-11"
now.ToDateTimeString()         // "2025-05-11 14:30:00"
now.ToDateString()             // "2025-05-11"
now.ToTimeString()             // "14:30:00"
now.ISO8601()                  // "2025-05-11T14:30:00Z07:00"

// Components
now.Year()                     // 2025
now.Month()                    // May
now.Day()                      // 11
now.Hour()                     // 14
now.Weekday()                  // Sunday

// Arithmetic
now.AddDays(7)                 // +7 days
now.SubDays(3)                 // -3 days
now.AddMonths(2)               // +2 months
now.AddYears(-1)               // -1 year
now.AddHours(5)                // +5 hours

// Start/End of period
now.StartOfDay()               // 2025-05-11 00:00:00
now.EndOfDay()                 // 2025-05-11 23:59:59
now.StartOfMonth()             // 2025-05-01 00:00:00
now.EndOfMonth()               // 2025-05-31 23:59:59
now.StartOfYear()              // 2025-01-01 00:00:00

// Comparison
now.Gt(console.Yesterday())    // true
now.Lt(tomorrow)               // true
now.Eq(console.Now())          // true
now.IsPast()                   // false
now.IsFuture()                 // false
now.IsToday()                  // true
now.IsWeekend()                // true (Sunday)
now.IsWeekday()                // false

// Difference
d1 := console.Parse("2024-01-01")
d2 := console.Parse("2024-01-10")
d1.DiffInDays(d2)              // -9
d1.DiffInHours(d2)             // -216
console.Parse("2024-01-10").DiffInDays(d1) // 9

// Human diff
console.Now().SubDays(2).HumanDiff()  // "2 days ago"
console.Now().AddHours(5).HumanDiff() // "5 hours from now"
console.Parse("2020-01-01").Age()     // 5

// Parse multiple formats
console.Parse("2024-01-15T10:30:00Z")
console.Parse("2024/01/15")
console.Parse("01/15/2024")

// Copy
copy := now.Copy()
```
