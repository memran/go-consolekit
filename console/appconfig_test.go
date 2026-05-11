package console

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfigSetAndGet(t *testing.T) {
	cfg := NewConfig()
	cfg.Set("app.name", "ConsoleKit")
	cfg.Set("app.debug", true)
	cfg.Set("app.port", 8080)

	if cfg.Get("app.name") != "ConsoleKit" {
		t.Fatalf("expected 'ConsoleKit', got %v", cfg.Get("app.name"))
	}
	if cfg.Get("app.debug") != true {
		t.Fatalf("expected true, got %v", cfg.Get("app.debug"))
	}
	if cfg.Get("app.port") != 8080 {
		t.Fatalf("expected 8080, got %v", cfg.Get("app.port"))
	}
}

func TestConfigGetNested(t *testing.T) {
	cfg := NewConfig()
	cfg.Set("database.connections.mysql.host", "localhost")
	cfg.Set("database.connections.mysql.port", 3306)

	host := cfg.Get("database.connections.mysql.host")
	if host != "localhost" {
		t.Fatalf("expected 'localhost', got %v", host)
	}

	connections := cfg.Get("database.connections")
	conns, ok := connections.(map[string]interface{})
	if !ok {
		t.Fatal("expected connections to be a map")
	}
	mysql, ok := conns["mysql"].(map[string]interface{})
	if !ok {
		t.Fatal("expected mysql to be a map")
	}
	if mysql["host"] != "localhost" {
		t.Fatalf("expected 'localhost', got %v", mysql["host"])
	}
}

func TestConfigGetDefault(t *testing.T) {
	cfg := NewConfig()
	cfg.Set("app.name", "MyApp")

	val := cfg.Get("app.nonexistent")
	if val != nil {
		t.Fatalf("expected nil, got %v", val)
	}

	val = cfg.Get("app.nonexistent", "default")
	if val != "default" {
		t.Fatalf("expected 'default', got %v", val)
	}

	val = cfg.Get("app.nonexistent.more", "fallback")
	if val != "fallback" {
		t.Fatalf("expected 'fallback', got %v", val)
	}
}

func TestConfigGetString(t *testing.T) {
	cfg := NewConfig()
	cfg.Set("app.name", "ConsoleKit")

	if cfg.GetString("app.name") != "ConsoleKit" {
		t.Fatalf("expected 'ConsoleKit', got '%s'", cfg.GetString("app.name"))
	}
	if cfg.GetString("app.missing") != "" {
		t.Fatalf("expected empty string, got '%s'", cfg.GetString("app.missing"))
	}
	if cfg.GetString("app.missing", "fallback") != "fallback" {
		t.Fatalf("expected 'fallback', got '%s'", cfg.GetString("app.missing", "fallback"))
	}
}

func TestConfigGetInt(t *testing.T) {
	cfg := NewConfig()
	cfg.Set("app.port", 8080)

	if cfg.GetInt("app.port") != 8080 {
		t.Fatalf("expected 8080, got %d", cfg.GetInt("app.port"))
	}
	if cfg.GetInt("app.missing") != 0 {
		t.Fatalf("expected 0, got %d", cfg.GetInt("app.missing"))
	}
	if cfg.GetInt("app.missing", 3000) != 3000 {
		t.Fatalf("expected 3000, got %d", cfg.GetInt("app.missing", 3000))
	}
}

func TestConfigGetBool(t *testing.T) {
	cfg := NewConfig()
	cfg.Set("app.debug", true)
	cfg.Set("app.cache", false)

	if cfg.GetBool("app.debug") != true {
		t.Fatalf("expected true, got %v", cfg.GetBool("app.debug"))
	}
	if cfg.GetBool("app.cache") != false {
		t.Fatalf("expected false, got %v", cfg.GetBool("app.cache"))
	}
	if cfg.GetBool("app.missing") != false {
		t.Fatalf("expected false, got %v", cfg.GetBool("app.missing"))
	}
	if cfg.GetBool("app.missing", true) != true {
		t.Fatalf("expected true, got %v", cfg.GetBool("app.missing", true))
	}
}

func TestConfigHas(t *testing.T) {
	cfg := NewConfig()
	cfg.Set("app.name", "MyApp")
	cfg.Set("app.database.host", "localhost")

	if !cfg.Has("app.name") {
		t.Fatal("expected app.name to exist")
	}
	if !cfg.Has("app.database.host") {
		t.Fatal("expected app.database.host to exist")
	}
	if cfg.Has("app.database.port") {
		t.Fatal("expected app.database.port to not exist")
	}
	if cfg.Has("nonexistent") {
		t.Fatal("expected nonexistent to not exist")
	}
}

func TestConfigAll(t *testing.T) {
	cfg := NewConfig()
	cfg.Set("app.name", "MyApp")
	cfg.Set("app.debug", true)

	all := cfg.All()
	if all == nil {
		t.Fatal("expected non-nil All()")
	}
	app, ok := all["app"].(map[string]interface{})
	if !ok {
		t.Fatal("expected app to be a map")
	}
	if app["name"] != "MyApp" {
		t.Fatalf("expected 'MyApp', got %v", app["name"])
	}
}

func TestConfigLoad(t *testing.T) {
	cfg := NewConfig()
	cfg.Load(map[string]interface{}{
		"app.name":  "LoadedApp",
		"app.debug": true,
	})

	if cfg.Get("app.name") != "LoadedApp" {
		t.Fatalf("expected 'LoadedApp', got %v", cfg.Get("app.name"))
	}
	if cfg.Get("app.debug") != true {
		t.Fatalf("expected true, got %v", cfg.Get("app.debug"))
	}
}

func TestConfigFluentChaining(t *testing.T) {
	cfg := NewConfig().
		Set("app.name", "ChainApp").
		Set("app.version", "1.0").
		Set("app.debug", true)

	if cfg.GetString("app.name") != "ChainApp" {
		t.Fatalf("expected 'ChainApp', got '%s'", cfg.GetString("app.name"))
	}
	if cfg.GetString("app.version") != "1.0" {
		t.Fatalf("expected '1.0', got '%s'", cfg.GetString("app.version"))
	}
	if cfg.GetBool("app.debug") != true {
		t.Fatalf("expected true, got %v", cfg.GetBool("app.debug"))
	}
}

func TestConfigOverwrite(t *testing.T) {
	cfg := NewConfig()
	cfg.Set("app.name", "Original")
	cfg.Set("app.name", "Overwritten")

	if cfg.Get("app.name") != "Overwritten" {
		t.Fatalf("expected 'Overwritten', got %v", cfg.Get("app.name"))
	}
}

func TestConfigDeepOverwrite(t *testing.T) {
	cfg := NewConfig()
	cfg.Set("database.default", "mysql")
	cfg.Set("database.connections.mysql.host", "localhost")

	if cfg.GetString("database.default") != "mysql" {
		t.Fatalf("expected 'mysql', got '%s'", cfg.GetString("database.default"))
	}
	if cfg.GetString("database.connections.mysql.host") != "localhost" {
		t.Fatalf("expected 'localhost', got '%s'", cfg.GetString("database.connections.mysql.host"))
	}
}

func TestConfigGetStringWithNonString(t *testing.T) {
	cfg := NewConfig()
	cfg.Set("app.port", 8080)

	if cfg.GetString("app.port") != "" {
		t.Fatalf("expected empty for non-string, got '%s'", cfg.GetString("app.port"))
	}
}

func TestConfigGetIntWithNonInt(t *testing.T) {
	cfg := NewConfig()
	cfg.Set("app.name", "test")

	if cfg.GetInt("app.name") != 0 {
		t.Fatalf("expected 0 for non-int, got %d", cfg.GetInt("app.name"))
	}
	if cfg.GetInt("app.name", 42) != 42 {
		t.Fatalf("expected 42 default, got %d", cfg.GetInt("app.name", 42))
	}
}

func TestConfigGetBoolWithNonBool(t *testing.T) {
	cfg := NewConfig()
	cfg.Set("app.name", "test")

	if cfg.GetBool("app.name") != false {
		t.Fatalf("expected false for non-bool")
	}
	if cfg.GetBool("app.name", true) != true {
		t.Fatalf("expected true default")
	}
}

func TestConfigLoadAndChain(t *testing.T) {
	cfg := NewConfig().
		Load(map[string]interface{}{
			"app.name":    "App",
			"app.version": "2.0",
		}).
		Set("app.debug", true)

	if cfg.GetString("app.name") != "App" {
		t.Fatalf("expected 'App', got '%s'", cfg.GetString("app.name"))
	}
	if cfg.GetBool("app.debug") != true {
		t.Fatalf("expected true")
	}
}

func TestConfigEmpty(t *testing.T) {
	cfg := NewConfig()

	if cfg.Get("anything") != nil {
		t.Fatal("expected nil from empty config")
	}
	if cfg.Has("anything") {
		t.Fatal("expected Has to return false for empty config")
	}
	if cfg.GetString("anything") != "" {
		t.Fatal("expected empty string from empty config")
	}
	if cfg.GetInt("anything") != 0 {
		t.Fatal("expected 0 from empty config")
	}
	if cfg.GetBool("anything") != false {
		t.Fatal("expected false from empty config")
	}
}

func TestConfigLoadJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	content := `{
		"app": { "name": "TestApp", "debug": true, "port": 8080 },
		"database": { "host": "localhost", "port": 3306 }
	}`
	writeFile(t, path, content)

	cfg := NewConfig().LoadJSON(path)

	if cfg.GetString("app.name") != "TestApp" {
		t.Fatalf("expected 'TestApp', got '%s'", cfg.GetString("app.name"))
	}
	if cfg.GetBool("app.debug") != true {
		t.Fatal("expected true")
	}
	if cfg.GetInt("app.port") != 8080 {
		t.Fatalf("expected 8080, got %d", cfg.GetInt("app.port"))
	}
	if cfg.GetString("database.host") != "localhost" {
		t.Fatalf("expected 'localhost', got '%s'", cfg.GetString("database.host"))
	}
	if cfg.GetInt("database.port") != 3306 {
		t.Fatalf("expected 3306, got %d", cfg.GetInt("database.port"))
	}
}

func TestConfigLoadJSONNotFound(t *testing.T) {
	cfg := NewConfig().LoadJSON("/nonexistent/config.json")
	if cfg.Get("anything") != nil {
		t.Fatal("expected empty config after loading nonexistent JSON")
	}
}

func TestConfigLoadJSONChaining(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "cfg.json")
	writeFile(t, path, `{"app": {"name": "LoadChain"}}`)

	cfg := NewConfig().
		LoadJSON(path).
		Set("app.debug", true)

	if cfg.GetString("app.name") != "LoadChain" {
		t.Fatalf("expected 'LoadChain', got '%s'", cfg.GetString("app.name"))
	}
	if cfg.GetBool("app.debug") != true {
		t.Fatal("expected true")
	}
}

func TestConfigLoadYAML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	content := `app:
  name: TestApp
  debug: true
  port: 8080
database:
  host: localhost
  port: 3306
`
	writeFile(t, path, content)

	cfg := NewConfig().LoadYAML(path)

	if cfg.GetString("app.name") != "TestApp" {
		t.Fatalf("expected 'TestApp', got '%s'", cfg.GetString("app.name"))
	}
	if cfg.GetBool("app.debug") != true {
		t.Fatal("expected true")
	}
	if cfg.GetInt("app.port") != 8080 {
		t.Fatalf("expected 8080, got %d", cfg.GetInt("app.port"))
	}
	if cfg.GetString("database.host") != "localhost" {
		t.Fatalf("expected 'localhost', got '%s'", cfg.GetString("database.host"))
	}
	if cfg.GetInt("database.port") != 3306 {
		t.Fatalf("expected 3306, got %d", cfg.GetInt("database.port"))
	}
}

func TestConfigLoadYAMLNotFound(t *testing.T) {
	cfg := NewConfig().LoadYAML("/nonexistent/config.yaml")
	if cfg.Get("anything") != nil {
		t.Fatal("expected empty config after loading nonexistent YAML")
	}
}

func TestConfigLoadYAMLChaining(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "cfg.yaml")
	writeFile(t, path, "app:\n  name: LoadChain\n")

	cfg := NewConfig().
		LoadYAML(path).
		Set("app.debug", true)

	if cfg.GetString("app.name") != "LoadChain" {
		t.Fatalf("expected 'LoadChain', got '%s'", cfg.GetString("app.name"))
	}
	if cfg.GetBool("app.debug") != true {
		t.Fatal("expected true")
	}
}

func TestConfigLoadJSONAndYAMLCombined(t *testing.T) {
	dir := t.TempDir()
	jsonPath := filepath.Join(dir, "db.json")
	yamlPath := filepath.Join(dir, "app.yaml")
	writeFile(t, jsonPath, `{"database": {"host": "json_host", "port": 5432}}`)
	writeFile(t, yamlPath, "app:\n  name: CombinedApp\n")

	cfg := NewConfig().
		LoadJSON(jsonPath).
		LoadYAML(yamlPath)

	if cfg.GetString("app.name") != "CombinedApp" {
		t.Fatalf("expected 'CombinedApp', got '%s'", cfg.GetString("app.name"))
	}
	if cfg.GetString("database.host") != "json_host" {
		t.Fatalf("expected 'json_host', got '%s'", cfg.GetString("database.host"))
	}
	if cfg.GetInt("database.port") != 5432 {
		t.Fatalf("expected 5432, got %d", cfg.GetInt("database.port"))
	}
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("writeFile failed: %v", err)
	}
}

func TestConfigGetIntWithFloat64(t *testing.T) {
	cfg := NewConfig()
	cfg.Set("app.rate", float64(3.5))

	if cfg.GetInt("app.rate") != 3 {
		t.Fatalf("expected 3, got %d", cfg.GetInt("app.rate"))
	}
}
