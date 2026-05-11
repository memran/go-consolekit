# Validation

Declarative field validation with chainable rule functions.

## Validator

```go
func NewValidator() *Validator
func (v *Validator) Data(data map[string]any) *Validator
func (v *Validator) Rule(field string, rules ...ruleFunc) *Validator
func (v *Validator) Validate() bool
func (v *Validator) Passes() bool
func (v *Validator) Fails() bool
func (v *Validator) Errors() map[string][]string
func (v *Validator) ErrorsFor(field string) []string
```

## Rule Functions

```go
func Required() ruleFunc
func Email() ruleFunc
func MinInt(min int) ruleFunc
func MaxInt(max int) ruleFunc
func MinLen(min int) ruleFunc
func MaxLen(max int) ruleFunc
func BetweenInt(min, max int) ruleFunc
func In(values ...any) ruleFunc
func NotIn(values ...any) ruleFunc
func Match(pattern string) ruleFunc
func Numeric() ruleFunc
func String() ruleFunc
func Bool() ruleFunc
func Integer() ruleFunc
func Confirmed() ruleFunc
func URL() ruleFunc
func NewRuleSet(rules ...ruleFunc) []ruleFunc
```

## Examples

```go
data := map[string]any{
    "name":             "Alice",
    "email":            "alice@example.com",
    "age":              25,
    "password":         "secret123",
    "password_confirmation": "secret123",
    "role":             "admin",
    "website":          "https://example.com",
    "count":            "42",
}

v := console.NewValidator().
    Data(data).
    Rule("name", console.Required(), console.MinLen(2), console.MaxLen(50)).
    Rule("email", console.Required(), console.Email()).
    Rule("age", console.Required(), console.MinInt(18), console.MaxInt(120)).
    Rule("password", console.Required(), console.MinLen(8), console.Confirmed()).
    Rule("role", console.In("admin", "user", "viewer")).
    Rule("website", console.URL()).
    Rule("count", console.Numeric()).
    Rule("active", console.Bool())

if v.Passes() {
    fmt.Println("All validations passed")
} else {
    for field, errors := range v.Errors() {
        fmt.Printf("%s: %v\n", field, errors)
    }
}

// Check specific field errors
errs := v.ErrorsFor("email")
```

### Using NewRuleSet

```go
nameRules := console.NewRuleSet(
    console.Required(),
    console.MinLen(2),
    console.MaxLen(100),
)

v := console.NewValidator().
    Data(data).
    Rule("name", nameRules...)
```

### Integer and Numeric

```go
console.NewValidator().
    Data(map[string]any{"age": 25}).
    Rule("age", console.Integer()).
    Passes() // true

console.NewValidator().
    Data(map[string]any{"score": 95.5}).
    Rule("score", console.Numeric()).
    Passes() // true (float64 is numeric)
```

### BetweenInt

```go
console.NewValidator().
    Data(map[string]any{"level": 5}).
    Rule("level", console.BetweenInt(1, 10)).
    Passes() // true
```

### NotIn

```go
console.NewValidator().
    Data(map[string]any{"status": "banned"}).
    Rule("status", console.NotIn("banned", "deleted")).
    Passes() // false
```

### Match (regex)

```go
console.NewValidator().
    Data(map[string]any{"phone": "123-456-7890"}).
    Rule("phone", console.Match(`^\d{3}-\d{3}-\d{4}$`)).
    Passes() // true
```
