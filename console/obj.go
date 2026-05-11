package console

import (
	"encoding/json"
	"strings"
)

type Obj struct {
	data map[string]any
}

func NewObj(data ...map[string]any) *Obj {
	if len(data) > 0 {
		m := make(map[string]any, len(data[0]))
		for k, v := range data[0] {
			m[k] = v
		}
		return &Obj{data: m}
	}
	return &Obj{data: make(map[string]any)}
}

func (o *Obj) Set(key string, value any) *Obj {
	keys := strings.Split(key, ".")
	current := o.data
	for i, k := range keys {
		if i == len(keys)-1 {
			current[k] = value
		} else {
			if next, ok := current[k].(map[string]any); ok {
				current = next
			} else {
				next := make(map[string]any)
				current[k] = next
				current = next
			}
		}
	}
	return o
}

func (o *Obj) Get(key string, defaults ...any) any {
	keys := strings.Split(key, ".")
	current := o.data
	for i, k := range keys {
		val, ok := current[k]
		if !ok {
			return defaultOrNilObj(defaults)
		}
		if i == len(keys)-1 {
			return val
		}
		next, ok := val.(map[string]any)
		if !ok {
			return defaultOrNilObj(defaults)
		}
		current = next
	}
	return defaultOrNilObj(defaults)
}

func (o *Obj) GetString(key string, defaults ...string) string {
	val := o.Get(key)
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

func (o *Obj) GetInt(key string, defaults ...int) int {
	val := o.Get(key)
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

func (o *Obj) GetBool(key string, defaults ...bool) bool {
	val := o.Get(key)
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

func (o *Obj) Has(key string) bool {
	keys := strings.Split(key, ".")
	current := o.data
	for i, k := range keys {
		val, ok := current[k]
		if !ok {
			return false
		}
		if i == len(keys)-1 {
			return true
		}
		next, ok := val.(map[string]any)
		if !ok {
			return false
		}
		current = next
	}
	return false
}

func (o *Obj) Forget(key string) *Obj {
	keys := strings.Split(key, ".")
	current := o.data
	for i, k := range keys {
		if i == len(keys)-1 {
			delete(current, k)
			return o
		}
		next, ok := current[k].(map[string]any)
		if !ok {
			return o
		}
		current = next
	}
	return o
}

func (o *Obj) All() map[string]any {
	result := make(map[string]any, len(o.data))
	for k, v := range o.data {
		result[k] = v
	}
	return result
}

func (o *Obj) ToJSON() string {
	data, err := json.Marshal(o.data)
	if err != nil {
		return "{}"
	}
	return string(data)
}

func (o *Obj) Count() int {
	return len(o.data)
}

func (o *Obj) IsEmpty() bool {
	return len(o.data) == 0
}

func defaultOrNilObj(defaults []any) any {
	if len(defaults) > 0 {
		return defaults[0]
	}
	return nil
}
