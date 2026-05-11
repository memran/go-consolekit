package console

import (
	"testing"
)

func TestStrContains(t *testing.T) {
	if !NewStr("hello world").Contains("world") {
		t.Fatal("expected contains world")
	}
	if NewStr("hello").Contains("xyz") {
		t.Fatal("expected not contains")
	}
}

func TestStrContainsAll(t *testing.T) {
	if !NewStr("hello world foo").ContainsAll("hello", "world") {
		t.Fatal("expected contains all")
	}
	if NewStr("hello").ContainsAll("hello", "world") {
		t.Fatal("expected not contains all")
	}
}

func TestStrContainsAny(t *testing.T) {
	if !NewStr("hello world").ContainsAny("xyz", "world") {
		t.Fatal("expected contains any")
	}
	if NewStr("hello").ContainsAny("xyz") {
		t.Fatal("expected not contains any")
	}
}

func TestStrStartsWith(t *testing.T) {
	if !NewStr("hello world").StartsWith("hello") {
		t.Fatal("expected starts with hello")
	}
	if NewStr("hello").StartsWith("world") {
		t.Fatal("expected not starts with")
	}
}

func TestStrEndsWith(t *testing.T) {
	if !NewStr("hello world").EndsWith("world") {
		t.Fatal("expected ends with world")
	}
	if NewStr("hello").EndsWith("world") {
		t.Fatal("expected not ends with")
	}
}

func TestStrUpper(t *testing.T) {
	if NewStr("hello").Upper().String() != "HELLO" {
		t.Fatal("upper failed")
	}
}

func TestStrLower(t *testing.T) {
	if NewStr("HELLO").Lower().String() != "hello" {
		t.Fatal("lower failed")
	}
}

func TestStrUcfirst(t *testing.T) {
	if NewStr("hello").Ucfirst().String() != "Hello" {
		t.Fatal("ucfirst failed")
	}
}

func TestStrLcfirst(t *testing.T) {
	if NewStr("Hello").Lcfirst().String() != "hello" {
		t.Fatal("lcfirst failed")
	}
}

func TestStrLength(t *testing.T) {
	if NewStr("hello").Length() != 5 {
		t.Fatalf("expected 5, got %d", NewStr("hello").Length())
	}
}

func TestStrLimit(t *testing.T) {
	r := NewStr("hello world").Limit(5)
	if r.String() != "hello..." {
		t.Fatalf("expected 'hello...', got '%s'", r.String())
	}
	r = NewStr("hi").Limit(5)
	if r.String() != "hi" {
		t.Fatalf("expected 'hi', got '%s'", r.String())
	}
}

func TestStrLimitCustom(t *testing.T) {
	r := NewStr("hello world").Limit(5, "***")
	if r.String() != "hello***" {
		t.Fatalf("expected 'hello***', got '%s'", r.String())
	}
}

func TestStrSubstr(t *testing.T) {
	r := NewStr("hello").Substr(1, 3)
	if r.String() != "ell" {
		t.Fatalf("expected 'ell', got '%s'", r.String())
	}
}

func TestStrReplace(t *testing.T) {
	r := NewStr("hello world").Replace("world", "there")
	if r.String() != "hello there" {
		t.Fatalf("expected 'hello there', got '%s'", r.String())
	}
}

func TestStrReplaceFirst(t *testing.T) {
	r := NewStr("foo bar foo").ReplaceFirst("foo", "baz")
	if r.String() != "baz bar foo" {
		t.Fatalf("expected 'baz bar foo', got '%s'", r.String())
	}
}

func TestStrReplaceLast(t *testing.T) {
	r := NewStr("foo bar foo").ReplaceLast("foo", "baz")
	if r.String() != "foo bar baz" {
		t.Fatalf("expected 'foo bar baz', got '%s'", r.String())
	}
}

func TestStrTrim(t *testing.T) {
	if NewStr("  hello  ").Trim().String() != "hello" {
		t.Fatal("trim failed")
	}
	if NewStr("xxhelloxx").Trim("x").String() != "hello" {
		t.Fatal("trim cutset failed")
	}
}

func TestStrLtrim(t *testing.T) {
	if NewStr("  hello").Ltrim().String() != "hello" {
		t.Fatal("ltrim failed")
	}
}

func TestStrRtrim(t *testing.T) {
	if NewStr("hello  ").Rtrim().String() != "hello" {
		t.Fatal("rtrim failed")
	}
}

func TestStrPadLeft(t *testing.T) {
	if NewStr("hi").PadLeft(5).String() != "   hi" {
		t.Fatalf("expected '   hi', got '%s'", NewStr("hi").PadLeft(5).String())
	}
	if NewStr("hi").PadLeft(5, ".").String() != "...hi" {
		t.Fatalf("expected '...hi', got '%s'", NewStr("hi").PadLeft(5, ".").String())
	}
}

func TestStrPadRight(t *testing.T) {
	if NewStr("hi").PadRight(5).String() != "hi   " {
		t.Fatalf("expected 'hi   ', got '%s'", NewStr("hi").PadRight(5).String())
	}
	if NewStr("hi").PadRight(5, ".").String() != "hi..." {
		t.Fatalf("expected 'hi...', got '%s'", NewStr("hi").PadRight(5, ".").String())
	}
}

func TestStrBefore(t *testing.T) {
	r := NewStr("hello world").Before(" ")
	if r.String() != "hello" {
		t.Fatalf("expected 'hello', got '%s'", r.String())
	}
}

func TestStrAfter(t *testing.T) {
	r := NewStr("hello world").After(" ")
	if r.String() != "world" {
		t.Fatalf("expected 'world', got '%s'", r.String())
	}
}

func TestStrBetween(t *testing.T) {
	r := NewStr("foo[bar]baz").Between("[", "]")
	if r.String() != "bar" {
		t.Fatalf("expected 'bar', got '%s'", r.String())
	}
}

func TestStrSlug(t *testing.T) {
	r := NewStr("Hello World").Slug()
	if r.String() != "hello-world" {
		t.Fatalf("expected 'hello-world', got '%s'", r.String())
	}
	r = NewStr("Foo Bar Baz").Slug("_")
	if r.String() != "foo_bar_baz" {
		t.Fatalf("expected 'foo_bar_baz', got '%s'", r.String())
	}
}

func TestStrSnake(t *testing.T) {
	r := NewStr("helloWorld").Snake()
	if r.String() != "hello_world" {
		t.Fatalf("expected 'hello_world', got '%s'", r.String())
	}
	r = NewStr("HelloWorld").Snake()
	if r.String() != "hello_world" {
		t.Fatalf("expected 'hello_world', got '%s'", r.String())
	}
}

func TestStrKebab(t *testing.T) {
	r := NewStr("helloWorld").Kebab()
	if r.String() != "hello-world" {
		t.Fatalf("expected 'hello-world', got '%s'", r.String())
	}
}

func TestStrStudly(t *testing.T) {
	r := NewStr("hello_world").Studly()
	if r.String() != "HelloWorld" {
		t.Fatalf("expected 'HelloWorld', got '%s'", r.String())
	}
}

func TestStrCamel(t *testing.T) {
	r := NewStr("hello_world").Camel()
	if r.String() != "helloWorld" {
		t.Fatalf("expected 'helloWorld', got '%s'", r.String())
	}
}

func TestStrRepeat(t *testing.T) {
	r := NewStr("ab").Repeat(3)
	if r.String() != "ababab" {
		t.Fatalf("expected 'ababab', got '%s'", r.String())
	}
}

func TestStrMask(t *testing.T) {
	r := NewStr("1234567890").Mask('*', 2, 4)
	if r.String() != "12****7890" {
		t.Fatalf("expected '12****7890', got '%s'", r.String())
	}
}

func TestStrIsEmpty(t *testing.T) {
	if !NewStr("").IsEmpty() {
		t.Fatal("expected empty")
	}
	if NewStr("a").IsEmpty() {
		t.Fatal("expected not empty")
	}
}

func TestStrIsNotEmpty(t *testing.T) {
	if !NewStr("a").IsNotEmpty() {
		t.Fatal("expected not empty")
	}
	if NewStr("").IsNotEmpty() {
		t.Fatal("expected empty")
	}
}

func TestStrEquals(t *testing.T) {
	if !NewStr("hello").Equals("hello") {
		t.Fatal("expected equals")
	}
	if NewStr("hello").Equals("world") {
		t.Fatal("expected not equals")
	}
}

func TestStrEqualsIgnoreCase(t *testing.T) {
	if !NewStr("Hello").EqualsIgnoreCase("hello") {
		t.Fatal("expected equals ignore case")
	}
}

func TestStrIs(t *testing.T) {
	tests := []struct {
		pattern string
		value   string
		want    bool
	}{
		{"hello", "hello", true},
		{"hello", "world", false},
		{"h*o", "hello", true},
		{"h*", "hello", true},
		{"*o", "hello", true},
		{"h*d", "hello", false},
		{"h*x*o", "hello", false},
	}
	for _, tt := range tests {
		got := NewStr(tt.value).Is(tt.pattern)
		if got != tt.want {
			t.Fatalf("Is(%q, %q) = %v, want %v", tt.pattern, tt.value, got, tt.want)
		}
	}
}

func TestRandom(t *testing.T) {
	r1 := Random(16)
	r2 := Random(16)
	if len(r1) != 16 || len(r2) != 16 {
		t.Fatal("expected length 16")
	}
	if r1 == r2 {
		t.Fatal("expected different random strings")
	}
}

func TestStrChaining(t *testing.T) {
	r := NewStr("  hello world  ").
		Trim().
		Upper().
		Replace("WORLD", "THERE").
		Slug()
	if r.String() != "hello-there" {
		t.Fatalf("expected 'hello-there', got '%s'", r.String())
	}
}
