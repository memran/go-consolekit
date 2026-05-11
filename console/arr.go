package console

type Arr []any

func NewArr(items ...any) Arr {
	return Arr(items)
}

func (a Arr) Get(index int) any {
	if index < 0 || index >= len(a) {
		return nil
	}
	return a[index]
}

func (a Arr) First() any {
	if len(a) == 0 {
		return nil
	}
	return a[0]
}

func (a Arr) Last() any {
	if len(a) == 0 {
		return nil
	}
	return a[len(a)-1]
}

func (a Arr) Has(value any) bool {
	for _, v := range a {
		if v == value {
			return true
		}
	}
	return false
}

func (a Arr) Contains(value any) bool {
	return a.Has(value)
}

func (a Arr) Where(fn func(any) bool) Arr {
	result := make(Arr, 0, len(a))
	for _, v := range a {
		if fn(v) {
			result = append(result, v)
		}
	}
	return result
}

func (a Arr) Pluck(key string) Arr {
	result := make(Arr, 0, len(a))
	for _, v := range a {
		if m, ok := v.(map[string]any); ok {
			if val, exists := m[key]; exists {
				result = append(result, val)
			}
		}
	}
	return result
}

func (a Arr) Flatten() Arr {
	result := make(Arr, 0)
	for _, v := range a {
		if nested, ok := v.(Arr); ok {
			result = append(result, nested.Flatten()...)
		} else if nested, ok := v.([]any); ok {
			result = append(result, Arr(nested).Flatten()...)
		} else {
			result = append(result, v)
		}
	}
	return result
}

func (a Arr) Collapse() Arr {
	result := make(Arr, 0)
	for _, v := range a {
		if nested, ok := v.(Arr); ok {
			result = append(result, nested...)
		} else if nested, ok := v.([]any); ok {
			result = append(result, nested...)
		} else {
			result = append(result, v)
		}
	}
	return result
}

func (a Arr) Chunk(size int) Arr {
	if size <= 0 {
		return Arr{}
	}
	var chunks Arr
	for i := 0; i < len(a); i += size {
		end := i + size
		if end > len(a) {
			end = len(a)
		}
		chunk := make(Arr, end-i)
		copy(chunk, a[i:end])
		chunks = append(chunks, chunk)
	}
	return chunks
}

func (a Arr) Unique() Arr {
	seen := make(map[any]struct{}, len(a))
	result := make(Arr, 0, len(a))
	for _, v := range a {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			result = append(result, v)
		}
	}
	return result
}

func (a Arr) Wrap(value any) Arr {
	if arr, ok := value.(Arr); ok {
		return arr
	}
	if slice, ok := value.([]any); ok {
		return Arr(slice)
	}
	return Arr{value}
}

func (a Arr) Join(separator string) string {
	parts := make([]string, len(a))
	for i, v := range a {
		parts[i] = toString(v)
	}
	return stringsJoin(parts, separator)
}

func (a Arr) ToJSON() (string, error) {
	return Collect(a).ToJSON()
}

func stringsJoin(parts []string, sep string) string {
	result := ""
	for i, p := range parts {
		if i > 0 {
			result += sep
		}
		result += p
	}
	return result
}
