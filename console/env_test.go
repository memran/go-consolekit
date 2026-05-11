package console

import (
	"os"
	"path/filepath"
	"testing"
)

func TestEnvGetWithDefault(t *testing.T) {
	e := NewEnv()
	val := e.Get("NONEXISTENT_KEY", "fallback")
	if val != "fallback" {
		t.Fatalf("expected 'fallback', got '%s'", val)
	}
}

func TestEnvGetMissingNoDefault(t *testing.T) {
	e := NewEnv()
	val := e.Get("NONEXISTENT_KEY")
	if val != "" {
		t.Fatalf("expected empty string, got '%s'", val)
	}
}

func TestEnvSetAndGet(t *testing.T) {
	e := NewEnv()
	e.Set("APP_NAME", "ConsoleKit")

	val := e.Get("APP_NAME")
	if val != "ConsoleKit" {
		t.Fatalf("expected 'ConsoleKit', got '%s'", val)
	}
}

func TestEnvLoadFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	os.WriteFile(path, []byte("APP_NAME=ConsoleKit\nAPP_DEBUG=true\n"), 0644)

	e := NewEnv().Load(path)

	if e.Get("APP_NAME") != "ConsoleKit" {
		t.Fatalf("expected 'ConsoleKit', got '%s'", e.Get("APP_NAME"))
	}
	if e.Get("APP_DEBUG") != "true" {
		t.Fatalf("expected 'true', got '%s'", e.Get("APP_DEBUG"))
	}
}

func TestEnvLoadFileNotFound(t *testing.T) {
	e := NewEnv().Load("/nonexistent/.env")
	if e.Get("ANYTHING") != "" {
		t.Fatal("expected empty for non-existent file")
	}
}

func TestEnvGetString(t *testing.T) {
	e := NewEnv()
	e.Set("APP_NAME", "ConsoleKit")

	if e.GetString("APP_NAME") != "ConsoleKit" {
		t.Fatalf("expected 'ConsoleKit', got '%s'", e.GetString("APP_NAME"))
	}
	if e.GetString("MISSING", "def") != "def" {
		t.Fatalf("expected 'def', got '%s'", e.GetString("MISSING", "def"))
	}
}

func TestEnvGetInt(t *testing.T) {
	e := NewEnv()
	e.Set("APP_PORT", "8080")

	if e.GetInt("APP_PORT") != 8080 {
		t.Fatalf("expected 8080, got %d", e.GetInt("APP_PORT"))
	}
	if e.GetInt("MISSING", 3000) != 3000 {
		t.Fatalf("expected 3000, got %d", e.GetInt("MISSING", 3000))
	}
	if e.GetInt("MISSING") != 0 {
		t.Fatalf("expected 0, got %d", e.GetInt("MISSING"))
	}
}

func TestEnvGetIntInvalid(t *testing.T) {
	e := NewEnv()
	e.Set("INVALID", "notanumber")

	if e.GetInt("INVALID", 42) != 42 {
		t.Fatalf("expected 42 default, got %d", e.GetInt("INVALID", 42))
	}
}

func TestEnvGetBool(t *testing.T) {
	e := NewEnv()
	e.Set("APP_DEBUG", "true")
	e.Set("APP_CACHE", "false")
	e.Set("FLAG_ON", "1")
	e.Set("FLAG_OFF", "0")
	e.Set("FLAG_YES", "yes")
	e.Set("FLAG_NO", "no")
	e.Set("FLAG_ON2", "on")
	e.Set("FLAG_OFF2", "off")

	tests := []struct {
		key      string
		expected bool
	}{
		{"APP_DEBUG", true},
		{"APP_CACHE", false},
		{"FLAG_ON", true},
		{"FLAG_OFF", false},
		{"FLAG_YES", true},
		{"FLAG_NO", false},
		{"FLAG_ON2", true},
		{"FLAG_OFF2", false},
	}
	for _, tt := range tests {
		got := e.GetBool(tt.key)
		if got != tt.expected {
			t.Fatalf("GetBool(%s) expected %v, got %v", tt.key, tt.expected, got)
		}
	}
}

func TestEnvGetBoolDefault(t *testing.T) {
	e := NewEnv()
	if e.GetBool("MISSING") != false {
		t.Fatal("expected false default")
	}
	if e.GetBool("MISSING", true) != true {
		t.Fatal("expected true default")
	}
}

func TestEnvGetFloat(t *testing.T) {
	e := NewEnv()
	e.Set("APP_RATE", "3.5")

	if e.GetFloat("APP_RATE") != 3.5 {
		t.Fatalf("expected 3.5, got %f", e.GetFloat("APP_RATE"))
	}
	if e.GetFloat("MISSING", 1.5) != 1.5 {
		t.Fatalf("expected 1.5, got %f", e.GetFloat("MISSING", 1.5))
	}
}

func TestEnvHas(t *testing.T) {
	e := NewEnv()
	e.Set("EXISTS_KEY", "value")

	if !e.Has("EXISTS_KEY") {
		t.Fatal("expected Has to be true for set key")
	}
	if e.Has("MISSING_KEY") {
		t.Fatal("expected Has to be false for missing key")
	}
}

func TestEnvHasOSVar(t *testing.T) {
	os.Setenv("TEST_CONSOLEKIT_ENV", "testval")
	defer os.Unsetenv("TEST_CONSOLEKIT_ENV")

	e := NewEnv()
	if !e.Has("TEST_CONSOLEKIT_ENV") {
		t.Fatal("expected Has to find OS env var")
	}
}

func TestEnvOSVarTakesPriority(t *testing.T) {
	os.Setenv("APP_NAME", "from_os")
	defer os.Unsetenv("APP_NAME")

	e := NewEnv()
	e.Set("APP_NAME", "from_env")

	val := e.Get("APP_NAME")
	if val != "from_os" {
		t.Fatalf("expected 'from_os', got '%s'", val)
	}
}

func TestEnvAll(t *testing.T) {
	e := NewEnv()
	e.Set("KEY1", "val1")
	e.Set("KEY2", "val2")

	all := e.All()
	if all["KEY1"] != "val1" {
		t.Fatalf("expected 'val1', got '%s'", all["KEY1"])
	}
	if all["KEY2"] != "val2" {
		t.Fatalf("expected 'val2', got '%s'", all["KEY2"])
	}
}

func TestEnvParseBasic(t *testing.T) {
	e := NewEnv()
	e.parse("KEY=value")

	if e.Get("KEY") != "value" {
		t.Fatalf("expected 'value', got '%s'", e.Get("KEY"))
	}
}

func TestEnvParseComments(t *testing.T) {
	e := NewEnv()
	e.parse("# comment\nKEY=value\n# another")

	if e.Get("KEY") != "value" {
		t.Fatalf("expected 'value', got '%s'", e.Get("KEY"))
	}
}

func TestEnvParseEmptyLines(t *testing.T) {
	e := NewEnv()
	e.parse("\n\nKEY=value\n\n")

	if e.Get("KEY") != "value" {
		t.Fatalf("expected 'value', got '%s'", e.Get("KEY"))
	}
}

func TestEnvParseQuotedValue(t *testing.T) {
	e := NewEnv()
	e.parse("KEY=\"quoted value\"")

	if e.Get("KEY") != "quoted value" {
		t.Fatalf("expected 'quoted value', got '%s'", e.Get("KEY"))
	}
}

func TestEnvParseSingleQuoted(t *testing.T) {
	e := NewEnv()
	e.parse("KEY='single quoted'")

	if e.Get("KEY") != "single quoted" {
		t.Fatalf("expected 'single quoted', got '%s'", e.Get("KEY"))
	}
}

func TestEnvParseExport(t *testing.T) {
	e := NewEnv()
	e.parse("export KEY=value")

	if e.Get("KEY") != "value" {
		t.Fatalf("expected 'value', got '%s'", e.Get("KEY"))
	}
}

func TestEnvParseWhitespace(t *testing.T) {
	e := NewEnv()
	e.parse("  KEY  =  value  ")

	if e.Get("KEY") != "value" {
		t.Fatalf("expected 'value', got '%s'", e.Get("KEY"))
	}
}

func TestEnvParseEscapedQuotes(t *testing.T) {
	e := NewEnv()
	e.parse(`KEY="hello \"world\""`)

	if e.Get("KEY") != `hello "world"` {
		t.Fatalf(`expected 'hello "world"', got '%s'`, e.Get("KEY"))
	}
}

func TestEnvParseEscapeSequences(t *testing.T) {
	e := NewEnv()
	e.parse("KEY=\"line1\\nline2\"")

	if e.Get("KEY") != "line1\nline2" {
		t.Fatalf("expected 'line1\\nline2', got '%s'", e.Get("KEY"))
	}
}

func TestEnvLoadFileFormat(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	content := `APP_NAME=MyApp
APP_ENV=local
APP_DEBUG=true
APP_PORT=8080

# Database config
DB_HOST=localhost
DB_PORT=3306
DB_NAME="my_database"
`
	os.WriteFile(path, []byte(content), 0644)

	e := NewEnv().Load(path)

	if e.GetString("APP_NAME") != "MyApp" {
		t.Fatalf("expected 'MyApp', got '%s'", e.GetString("APP_NAME"))
	}
	if e.GetBool("APP_DEBUG") != true {
		t.Fatal("expected true")
	}
	if e.GetInt("APP_PORT") != 8080 {
		t.Fatalf("expected 8080, got %d", e.GetInt("APP_PORT"))
	}
	if e.GetString("DB_HOST") != "localhost" {
		t.Fatalf("expected 'localhost', got '%s'", e.GetString("DB_HOST"))
	}
	if e.GetInt("DB_PORT") != 3306 {
		t.Fatalf("expected 3306, got %d", e.GetInt("DB_PORT"))
	}
	if e.GetString("DB_NAME") != "my_database" {
		t.Fatalf("expected 'my_database', got '%s'", e.GetString("DB_NAME"))
	}
}

func TestEnvFluentChaining(t *testing.T) {
	e := NewEnv().
		Set("APP_NAME", "ChainApp").
		Set("APP_DEBUG", "true").
		Set("APP_PORT", "9090")

	if e.GetString("APP_NAME") != "ChainApp" {
		t.Fatalf("expected 'ChainApp', got '%s'", e.GetString("APP_NAME"))
	}
	if e.GetBool("APP_DEBUG") != true {
		t.Fatal("expected true")
	}
	if e.GetInt("APP_PORT") != 9090 {
		t.Fatalf("expected 9090, got %d", e.GetInt("APP_PORT"))
	}
}
