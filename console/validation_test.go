package console

import "testing"

func TestValidatorRequiredPass(t *testing.T) {
	v := NewValidator().Data(map[string]any{"name": "John"}).Rule("name", Required())
	if !v.Passes() {
		t.Fatal("expected passes")
	}
}

func TestValidatorRequiredFail(t *testing.T) {
	v := NewValidator().Data(map[string]any{"name": ""}).Rule("name", Required())
	if !v.Fails() {
		t.Fatal("empty string should fail required")
	}
}

func TestValidatorRequiredNil(t *testing.T) {
	v := NewValidator().Data(map[string]any{"name": nil}).Rule("name", Required())
	if !v.Fails() {
		t.Fatal("expected fails for nil")
	}
}

func TestValidatorEmailPass(t *testing.T) {
	v := NewValidator().Data(map[string]any{"email": "test@example.com"}).Rule("email", Email())
	if !v.Passes() {
		t.Fatal("valid email should pass")
	}
}

func TestValidatorEmailFail(t *testing.T) {
	v := NewValidator().Data(map[string]any{"email": "invalid"}).Rule("email", Email())
	if !v.Fails() {
		t.Fatal("invalid email should fail")
	}
}

func TestValidatorMinInt(t *testing.T) {
	v := NewValidator().Data(map[string]any{"age": 18}).Rule("age", MinInt(18))
	if !v.Passes() {
		t.Fatal("age 18 should pass min 18")
	}
	v = NewValidator().Data(map[string]any{"age": 15}).Rule("age", MinInt(18))
	if !v.Fails() {
		t.Fatal("age 15 should fail min 18")
	}
}

func TestValidatorMaxInt(t *testing.T) {
	v := NewValidator().Data(map[string]any{"age": 25}).Rule("age", MaxInt(30))
	if !v.Passes() {
		t.Fatal("age 25 should pass max 30")
	}
	v = NewValidator().Data(map[string]any{"age": 35}).Rule("age", MaxInt(30))
	if !v.Fails() {
		t.Fatal("age 35 should fail max 30")
	}
}

func TestValidatorMinLen(t *testing.T) {
	v := NewValidator().Data(map[string]any{"name": "John"}).Rule("name", MinLen(3))
	if !v.Passes() {
		t.Fatal("min len should pass")
	}
	v = NewValidator().Data(map[string]any{"name": "Jo"}).Rule("name", MinLen(3))
	if !v.Fails() {
		t.Fatal("min len should fail")
	}
}

func TestValidatorMaxLen(t *testing.T) {
	v := NewValidator().Data(map[string]any{"name": "Jo"}).Rule("name", MaxLen(10))
	if !v.Passes() {
		t.Fatal("max len should pass")
	}
	v = NewValidator().Data(map[string]any{"name": "VeryLongNameHere"}).Rule("name", MaxLen(10))
	if !v.Fails() {
		t.Fatal("max len should fail")
	}
}

func TestValidatorBetweenInt(t *testing.T) {
	v := NewValidator().Data(map[string]any{"age": 25}).Rule("age", BetweenInt(18, 65))
	if !v.Passes() {
		t.Fatal("between should pass")
	}
	v = NewValidator().Data(map[string]any{"age": 10}).Rule("age", BetweenInt(18, 65))
	if !v.Fails() {
		t.Fatal("below range should fail")
	}
}

func TestValidatorIn(t *testing.T) {
	v := NewValidator().Data(map[string]any{"role": "admin"}).Rule("role", In("admin", "user"))
	if !v.Passes() {
		t.Fatal("in should pass")
	}
	v = NewValidator().Data(map[string]any{"role": "superadmin"}).Rule("role", In("admin", "user"))
	if !v.Fails() {
		t.Fatal("not in should fail")
	}
}

func TestValidatorNotIn(t *testing.T) {
	v := NewValidator().Data(map[string]any{"status": "banned"}).Rule("status", NotIn("banned", "deleted"))
	if !v.Fails() {
		t.Fatal("not in should fail for banned")
	}
}

func TestValidatorMatch(t *testing.T) {
	v := NewValidator().Data(map[string]any{"phone": "123-456-7890"}).Rule("phone", Match(`^\d{3}-\d{3}-\d{4}$`))
	if !v.Passes() {
		t.Fatal("match should pass")
	}
	v = NewValidator().Data(map[string]any{"phone": "invalid"}).Rule("phone", Match(`^\d{3}-\d{3}-\d{4}$`))
	if !v.Fails() {
		t.Fatal("match should fail")
	}
}

func TestValidatorNumeric(t *testing.T) {
	v := NewValidator().Data(map[string]any{"price": 42.5}).Rule("price", Numeric())
	if !v.Passes() {
		t.Fatal("numeric should pass for float")
	}
	v = NewValidator().Data(map[string]any{"price": "notanumber"}).Rule("price", Numeric())
	if !v.Fails() {
		t.Fatal("numeric should fail for string")
	}
}

func TestValidatorString(t *testing.T) {
	v := NewValidator().Data(map[string]any{"name": "John"}).Rule("name", String())
	if !v.Passes() {
		t.Fatal("string should pass")
	}
	v = NewValidator().Data(map[string]any{"name": 42}).Rule("name", String())
	if !v.Fails() {
		t.Fatal("string should fail for int")
	}
}

func TestValidatorBool(t *testing.T) {
	v := NewValidator().Data(map[string]any{"active": true}).Rule("active", Bool())
	if !v.Passes() {
		t.Fatal("bool should pass")
	}
	v = NewValidator().Data(map[string]any{"active": "yes"}).Rule("active", Bool())
	if !v.Fails() {
		t.Fatal("bool should fail for string")
	}
}

func TestValidatorConfirmed(t *testing.T) {
	v := NewValidator().Data(map[string]any{"password": "secret", "password_confirmation": "secret"}).Rule("password", Confirmed())
	if !v.Passes() {
		t.Fatal("confirmed should pass")
	}
	v = NewValidator().Data(map[string]any{"password": "secret", "password_confirmation": "wrong"}).Rule("password", Confirmed())
	if !v.Fails() {
		t.Fatal("confirmed should fail")
	}
}

func TestValidatorURL(t *testing.T) {
	v := NewValidator().Data(map[string]any{"url": "https://example.com"}).Rule("url", URL())
	if !v.Passes() {
		t.Fatal("url should pass")
	}
	v = NewValidator().Data(map[string]any{"url": "not-a-url"}).Rule("url", URL())
	if !v.Fails() {
		t.Fatal("url should fail")
	}
}

func TestValidatorMultipleRules(t *testing.T) {
	v := NewValidator().
		Data(map[string]any{"name": "John", "email": "john@example.com", "age": 25}).
		Rule("name", Required(), MinLen(2)).
		Rule("email", Required(), Email()).
		Rule("age", Required(), MinInt(18), MaxInt(65))

	if !v.Passes() {
		t.Fatal("all rules should pass")
	}
}

func TestValidatorErrors(t *testing.T) {
	v := NewValidator().
		Data(map[string]any{"name": "", "email": "invalid"}).
		Rule("name", Required()).
		Rule("email", Email())

	if v.Passes() {
		t.Fatal("should fail")
	}
	errs := v.Errors()
	if len(errs) != 2 {
		t.Fatalf("expected 2 fields with errors, got %d", len(errs))
	}
	if len(errs["name"]) == 0 {
		t.Fatal("expected name errors")
	}
	if len(errs["email"]) == 0 {
		t.Fatal("expected email errors")
	}
}

func TestValidatorErrorsFor(t *testing.T) {
	v := NewValidator().
		Data(map[string]any{"name": ""}).
		Rule("name", Required())

	v.Validate()
	if len(v.ErrorsFor("name")) == 0 {
		t.Fatal("expected errors for name")
	}
	if v.ErrorsFor("nonexistent") != nil {
		t.Fatal("expected nil for nonexistent field")
	}
}

func TestNewRuleSet(t *testing.T) {
	rules := NewRuleSet(Required(), Email())
	if len(rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(rules))
	}
}
