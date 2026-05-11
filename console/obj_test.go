package console

import "testing"

func TestObjSetAndGet(t *testing.T) {
	o := NewObj()
	o.Set("name", "ConsoleKit")
	o.Set("debug", true)
	o.Set("port", 8080)

	if o.GetString("name") != "ConsoleKit" {
		t.Fatalf("expected 'ConsoleKit', got '%s'", o.GetString("name"))
	}
	if o.GetBool("debug") != true {
		t.Fatal("expected true")
	}
	if o.GetInt("port") != 8080 {
		t.Fatalf("expected 8080, got %d", o.GetInt("port"))
	}
}

func TestObjGetNested(t *testing.T) {
	o := NewObj()
	o.Set("database.host", "localhost")
	o.Set("database.port", 3306)

	if o.GetString("database.host") != "localhost" {
		t.Fatalf("expected 'localhost', got '%s'", o.GetString("database.host"))
	}
	if o.GetInt("database.port") != 3306 {
		t.Fatalf("expected 3306, got %d", o.GetInt("database.port"))
	}
}

func TestObjGetDefault(t *testing.T) {
	o := NewObj()
	if o.Get("missing") != nil {
		t.Fatal("expected nil")
	}
	if o.Get("missing", "default") != "default" {
		t.Fatal("expected default")
	}
	if o.GetString("missing", "fallback") != "fallback" {
		t.Fatal("expected fallback")
	}
	if o.GetInt("missing", 42) != 42 {
		t.Fatal("expected 42")
	}
	if o.GetBool("missing", true) != true {
		t.Fatal("expected true")
	}
}

func TestObjHas(t *testing.T) {
	o := NewObj()
	o.Set("app.name", "Test")

	if !o.Has("app.name") {
		t.Fatal("expected has")
	}
	if o.Has("app.missing") {
		t.Fatal("expected not has")
	}
}

func TestObjForget(t *testing.T) {
	o := NewObj(map[string]any{"name": "Test", "version": "1.0"})
	o.Forget("name")
	if o.Has("name") {
		t.Fatal("expected name to be forgotten")
	}
	if !o.Has("version") {
		t.Fatal("expected version to remain")
	}
}

func TestObjForgetNested(t *testing.T) {
	o := NewObj()
	o.Set("database.host", "localhost")
	o.Set("database.port", 3306)
	o.Forget("database.host")
	if o.Has("database.host") {
		t.Fatal("expected host to be forgotten")
	}
	if !o.Has("database.port") {
		t.Fatal("expected port to remain")
	}
}

func TestObjAll(t *testing.T) {
	o := NewObj(map[string]any{"name": "Test", "version": "1.0"})
	all := o.All()
	if all["name"] != "Test" || all["version"] != "1.0" {
		t.Fatal("all failed")
	}
}

func TestObjToJSON(t *testing.T) {
	o := NewObj(map[string]any{"name": "Test"})
	json := o.ToJSON()
	if json != `{"name":"Test"}` {
		t.Fatalf("expected '{\"name\":\"Test\"}', got '%s'", json)
	}
}

func TestObjCount(t *testing.T) {
	o := NewObj()
	if o.Count() != 0 {
		t.Fatal("expected 0")
	}
	o.Set("name", "Test")
	if o.Count() != 1 {
		t.Fatal("expected 1")
	}
}

func TestObjIsEmpty(t *testing.T) {
	if !NewObj().IsEmpty() {
		t.Fatal("expected empty")
	}
	if NewObj(map[string]any{"a": 1}).IsEmpty() {
		t.Fatal("expected not empty")
	}
}

func TestObjChaining(t *testing.T) {
	o := NewObj().
		Set("app.name", "ChainApp").
		Set("app.debug", true)

	if o.GetString("app.name") != "ChainApp" {
		t.Fatal("chaining getstring failed")
	}
	if o.GetBool("app.debug") != true {
		t.Fatal("chaining getbool failed")
	}
}

func TestObjFromMap(t *testing.T) {
	o := NewObj(map[string]any{"name": "FromMap", "nested": map[string]any{"key": "val"}})
	if o.GetString("name") != "FromMap" {
		t.Fatal("from map failed")
	}
}
