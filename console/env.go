package console

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Env struct {
	data map[string]string
}

func NewEnv() *Env {
	return &Env{
		data: make(map[string]string),
	}
}

func (e *Env) Load(path string) *Env {
	data, err := os.ReadFile(path)
	if err != nil {
		return e
	}
	e.parse(string(data))
	return e
}

func (e *Env) LoadFile(path string) *Env {
	return e.Load(path)
}

func (e *Env) Set(key, value string) *Env {
	e.data[key] = value
	return e
}

func (e *Env) Get(key string, defaults ...string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	if val, ok := e.data[key]; ok {
		return val
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return ""
}

func (e *Env) GetString(key string, defaults ...string) string {
	return e.Get(key, defaults...)
}

func (e *Env) GetInt(key string, defaults ...int) int {
	val := e.Get(key)
	if val == "" {
		if len(defaults) > 0 {
			return defaults[0]
		}
		return 0
	}
	n, err := strconv.Atoi(val)
	if err != nil {
		if len(defaults) > 0 {
			return defaults[0]
		}
		return 0
	}
	return n
}

func (e *Env) GetBool(key string, defaults ...bool) bool {
	val := e.Get(key)
	if val == "" {
		if len(defaults) > 0 {
			return defaults[0]
		}
		return false
	}
	switch strings.ToLower(val) {
	case "1", "true", "yes", "on":
		return true
	case "0", "false", "no", "off":
		return false
	default:
		if len(defaults) > 0 {
			return defaults[0]
		}
		return false
	}
}

func (e *Env) GetFloat(key string, defaults ...float64) float64 {
	val := e.Get(key)
	if val == "" {
		if len(defaults) > 0 {
			return defaults[0]
		}
		return 0
	}
	f, err := strconv.ParseFloat(val, 64)
	if err != nil {
		if len(defaults) > 0 {
			return defaults[0]
		}
		return 0
	}
	return f
}

func (e *Env) Has(key string) bool {
	_, ok := os.LookupEnv(key)
	if ok {
		return true
	}
	_, ok = e.data[key]
	return ok
}

func (e *Env) All() map[string]string {
	result := make(map[string]string, len(e.data))
	for k, v := range e.data {
		result[k] = v
	}
	for _, kv := range os.Environ() {
		parts := strings.SplitN(kv, "=", 2)
		if len(parts) == 2 {
			result[parts[0]] = parts[1]
		}
	}
	return result
}

func (e *Env) parse(content string) {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.TrimPrefix(line, "export ")
		eqIdx := strings.Index(line, "=")
		if eqIdx < 0 {
			continue
		}
		key := strings.TrimSpace(line[:eqIdx])
		if key == "" {
			continue
		}
		value := strings.TrimSpace(line[eqIdx+1:])
		value = e.unquote(value)
		e.data[key] = value
	}
}

func (e *Env) unquote(value string) string {
	if len(value) < 2 {
		return value
	}
	quote := value[0]
	if quote != '"' && quote != '\'' {
		return value
	}
	if value[len(value)-1] != byte(quote) {
		return value
	}
	inner := value[1 : len(value)-1]
	if quote == '"' {
		inner = strings.NewReplacer(
			`\n`, "\n",
			`\r`, "\r",
			`\t`, "\t",
			`\\`, "\\",
			`\"`, `"`,
			`\'`, `'`,
		).Replace(inner)
	}
	return inner
}

func (e *Env) Dump() *Env {
	for k, v := range e.data {
		fmt.Printf("%s=%s\n", k, v)
	}
	return e
}
