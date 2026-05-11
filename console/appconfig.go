package console

import (
	"encoding/json"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	data map[string]interface{}
}

func NewConfig() *Config {
	return &Config{
		data: make(map[string]interface{}),
	}
}

func (c *Config) Set(key string, value interface{}) *Config {
	keys := strings.Split(key, ".")
	current := c.data
	for i, k := range keys {
		if i == len(keys)-1 {
			current[k] = value
		} else {
			if next, ok := current[k].(map[string]interface{}); ok {
				current = next
			} else {
				next := make(map[string]interface{})
				current[k] = next
				current = next
			}
		}
	}
	return c
}

func (c *Config) Get(key string, defaults ...interface{}) interface{} {
	keys := strings.Split(key, ".")
	current := c.data
	for i, k := range keys {
		val, ok := current[k]
		if !ok {
			return defaultOrNil(defaults)
		}
		if i == len(keys)-1 {
			return val
		}
		next, ok := val.(map[string]interface{})
		if !ok {
			return defaultOrNil(defaults)
		}
		current = next
	}
	return defaultOrNil(defaults)
}

func (c *Config) GetString(key string, defaults ...string) string {
	val := c.Get(key)
	if val == nil {
		if len(defaults) > 0 {
			return defaults[0]
		}
		return ""
	}
	s, ok := val.(string)
	if !ok {
		if len(defaults) > 0 {
			return defaults[0]
		}
		return ""
	}
	return s
}

func (c *Config) GetInt(key string, defaults ...int) int {
	val := c.Get(key)
	if val == nil {
		if len(defaults) > 0 {
			return defaults[0]
		}
		return 0
	}
	switch v := val.(type) {
	case int:
		return v
	case float64:
		return int(v)
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return 0
}

func (c *Config) GetBool(key string, defaults ...bool) bool {
	val := c.Get(key)
	if val == nil {
		if len(defaults) > 0 {
			return defaults[0]
		}
		return false
	}
	b, ok := val.(bool)
	if !ok {
		if len(defaults) > 0 {
			return defaults[0]
		}
		return false
	}
	return b
}

func (c *Config) Has(key string) bool {
	keys := strings.Split(key, ".")
	current := c.data
	for i, k := range keys {
		val, ok := current[k]
		if !ok {
			return false
		}
		if i == len(keys)-1 {
			return true
		}
		next, ok := val.(map[string]interface{})
		if !ok {
			return false
		}
		current = next
	}
	return false
}

func (c *Config) All() map[string]interface{} {
	return c.data
}

func (c *Config) Load(data map[string]interface{}) *Config {
	for k, v := range data {
		c.Set(k, v)
	}
	return c
}

func (c *Config) LoadJSON(path string) *Config {
	data, err := os.ReadFile(path)
	if err != nil {
		return c
	}
	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		return c
	}
	c.flatten("", parsed)
	return c
}

func (c *Config) LoadYAML(path string) *Config {
	data, err := os.ReadFile(path)
	if err != nil {
		return c
	}
	var parsed map[string]interface{}
	if err := yaml.Unmarshal(data, &parsed); err != nil {
		return c
	}
	c.flatten("", parsed)
	return c
}

func (c *Config) flatten(prefix string, data map[string]interface{}) {
	for k, v := range data {
		key := k
		if prefix != "" {
			key = prefix + "." + k
		}
		if nested, ok := v.(map[string]interface{}); ok {
			c.flatten(key, nested)
		} else {
			c.Set(key, v)
		}
	}
}

func defaultOrNil(defaults []interface{}) interface{} {
	if len(defaults) > 0 {
		return defaults[0]
	}
	return nil
}
