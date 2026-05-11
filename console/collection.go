package console

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sort"
	"strings"
)

type Collection struct {
	items []any
}

func Collect(items []any) *Collection {
	c := make([]any, len(items))
	copy(c, items)
	return &Collection{items: c}
}

func CollectFrom[T any](items []T) *Collection {
	c := make([]any, len(items))
	for i, v := range items {
		c[i] = v
	}
	return &Collection{items: c}
}

func (c *Collection) All() []any {
	result := make([]any, len(c.items))
	copy(result, c.items)
	return result
}

func (c *Collection) Count() int {
	return len(c.items)
}

func (c *Collection) IsEmpty() bool {
	return len(c.items) == 0
}

func (c *Collection) IsNotEmpty() bool {
	return len(c.items) > 0
}

func (c *Collection) First() any {
	if len(c.items) == 0 {
		return nil
	}
	return c.items[0]
}

func (c *Collection) Last() any {
	if len(c.items) == 0 {
		return nil
	}
	return c.items[len(c.items)-1]
}

func (c *Collection) Get(index int) any {
	if index < 0 || index >= len(c.items) {
		return nil
	}
	return c.items[index]
}

func (c *Collection) Map(fn func(any) any) *Collection {
	result := make([]any, len(c.items))
	for i, v := range c.items {
		result[i] = fn(v)
	}
	return &Collection{items: result}
}

func (c *Collection) Filter(fn func(any) bool) *Collection {
	result := make([]any, 0, len(c.items))
	for _, v := range c.items {
		if fn(v) {
			result = append(result, v)
		}
	}
	return &Collection{items: result}
}

func (c *Collection) Reject(fn func(any) bool) *Collection {
	result := make([]any, 0, len(c.items))
	for _, v := range c.items {
		if !fn(v) {
			result = append(result, v)
		}
	}
	return &Collection{items: result}
}

func Tap[T any](v T, fn func(T)) T {
	fn(v)
	return v
}

func (c *Collection) Tap(fn func(*Collection)) *Collection {
	fn(c)
	return c
}

func (c *Collection) Each(fn func(any)) *Collection {
	for _, v := range c.items {
		fn(v)
	}
	return c
}

func (c *Collection) Reduce(fn func(any, any) any, initial any) any {
	carry := initial
	for _, v := range c.items {
		carry = fn(carry, v)
	}
	return carry
}

func (c *Collection) Slice(offset, length int) *Collection {
	if offset >= len(c.items) {
		return &Collection{items: []any{}}
	}
	end := offset + length
	if end > len(c.items) {
		end = len(c.items)
	}
	if offset < 0 {
		offset = 0
	}
	return &Collection{items: append([]any{}, c.items[offset:end]...)}
}

func (c *Collection) Sort(fn func(any, any) bool) *Collection {
	result := make([]any, len(c.items))
	copy(result, c.items)
	sort.Slice(result, func(i, j int) bool {
		return fn(result[i], result[j])
	})
	return &Collection{items: result}
}

func (c *Collection) Reverse() *Collection {
	result := make([]any, len(c.items))
	for i, v := range c.items {
		result[len(c.items)-1-i] = v
	}
	return &Collection{items: result}
}

func (c *Collection) Shuffle() *Collection {
	result := make([]any, len(c.items))
	perm := rand.Perm(len(c.items))
	for i, v := range perm {
		result[v] = c.items[i]
	}
	return &Collection{items: result}
}

func (c *Collection) Unique() *Collection {
	seen := make(map[any]struct{}, len(c.items))
	result := make([]any, 0, len(c.items))
	for _, v := range c.items {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			result = append(result, v)
		}
	}
	return &Collection{items: result}
}

func (c *Collection) Values() *Collection {
	values := make([]any, len(c.items))
	copy(values, c.items)
	return &Collection{items: values}
}

func (c *Collection) Keys() *Collection {
	keys := make([]any, len(c.items))
	for i := range c.items {
		keys[i] = i
	}
	return &Collection{items: keys}
}

func (c *Collection) Chunk(size int) *Collection {
	if size <= 0 {
		return &Collection{items: []any{}}
	}
	var chunks []any
	for i := 0; i < len(c.items); i += size {
		end := i + size
		if end > len(c.items) {
			end = len(c.items)
		}
		chunk := make([]any, end-i)
		copy(chunk, c.items[i:end])
		chunks = append(chunks, chunk)
	}
	return &Collection{items: chunks}
}

func (c *Collection) Collapse() *Collection {
	var result []any
	for _, item := range c.items {
		if slice, ok := item.([]any); ok {
			result = append(result, slice...)
		}
	}
	return &Collection{items: result}
}

func (c *Collection) Merge(items []any) *Collection {
	result := make([]any, len(c.items)+len(items))
	copy(result, c.items)
	copy(result[len(c.items):], items)
	return &Collection{items: result}
}

func (c *Collection) Diff(items []any) *Collection {
	exclude := make(map[any]struct{}, len(items))
	for _, v := range items {
		exclude[v] = struct{}{}
	}
	result := make([]any, 0, len(c.items))
	for _, v := range c.items {
		if _, ok := exclude[v]; !ok {
			result = append(result, v)
		}
	}
	return &Collection{items: result}
}

func (c *Collection) Intersect(items []any) *Collection {
	lookup := make(map[any]struct{}, len(items))
	for _, v := range items {
		lookup[v] = struct{}{}
	}
	result := make([]any, 0, len(c.items))
	for _, v := range c.items {
		if _, ok := lookup[v]; ok {
			result = append(result, v)
		}
	}
	return &Collection{items: result}
}

func (c *Collection) Contains(value any) bool {
	for _, v := range c.items {
		if v == value {
			return true
		}
	}
	return false
}

func (c *Collection) Search(value any) int {
	for i, v := range c.items {
		if v == value {
			return i
		}
	}
	return -1
}

func (c *Collection) Where(key string, value any) *Collection {
	result := make([]any, 0, len(c.items))
	for _, v := range c.items {
		if m, ok := v.(map[string]any); ok {
			if m[key] == value {
				result = append(result, v)
			}
		}
	}
	return &Collection{items: result}
}

func (c *Collection) Pluck(key string) *Collection {
	result := make([]any, 0, len(c.items))
	for _, v := range c.items {
		if m, ok := v.(map[string]any); ok {
			if val, exists := m[key]; exists {
				result = append(result, val)
			}
		}
	}
	return &Collection{items: result}
}

func (c *Collection) GroupBy(key string) *Collection {
	groups := make(map[any][]any)
	for _, v := range c.items {
		if m, ok := v.(map[string]any); ok {
			if val, exists := m[key]; exists {
				groups[val] = append(groups[val], v)
			}
		}
	}
	result := make([]any, 0, len(groups))
	for _, group := range groups {
		result = append(result, &Collection{items: group})
	}
	return &Collection{items: result}
}

func (c *Collection) KeyBy(key string) *Collection {
	result := make(map[any]any, len(c.items))
	for _, v := range c.items {
		if m, ok := v.(map[string]any); ok {
			if val, exists := m[key]; exists {
				result[val] = v
			}
		}
	}
	items := make([]any, 0, len(result))
	for k, v := range result {
		items = append(items, map[string]any{"key": k, "value": v})
	}
	return &Collection{items: items}
}

func (c *Collection) Sum(key ...string) float64 {
	var total float64
	for _, v := range c.items {
		if len(key) > 0 {
			if m, ok := v.(map[string]any); ok {
				if val, exists := m[key[0]]; exists {
					total += toFloat64(val)
				}
			}
		} else {
			total += toFloat64(v)
		}
	}
	return total
}

func (c *Collection) Avg(key ...string) float64 {
	if len(c.items) == 0 {
		return 0
	}
	return c.Sum(key...) / float64(len(c.items))
}

func (c *Collection) Min(key ...string) float64 {
	if len(c.items) == 0 {
		return 0
	}
	min := toFloat64(c.items[0], key...)
	for _, v := range c.items[1:] {
		val := toFloat64(v, key...)
		if val < min {
			min = val
		}
	}
	return min
}

func (c *Collection) Max(key ...string) float64 {
	if len(c.items) == 0 {
		return 0
	}
	max := toFloat64(c.items[0], key...)
	for _, v := range c.items[1:] {
		val := toFloat64(v, key...)
		if val > max {
			max = val
		}
	}
	return max
}

func (c *Collection) Implode(separator string) string {
	parts := make([]string, len(c.items))
	for i, v := range c.items {
		parts[i] = toString(v)
	}
	return strings.Join(parts, separator)
}

func (c *Collection) Join(separator, lastSeparator string) string {
	if len(c.items) == 0 {
		return ""
	}
	if len(c.items) == 1 {
		return toString(c.items[0])
	}
	parts := make([]string, len(c.items)-1)
	for i, v := range c.items[:len(c.items)-1] {
		parts[i] = toString(v)
	}
	return strings.Join(parts, separator) + lastSeparator + toString(c.items[len(c.items)-1])
}

func (c *Collection) ToJSON() (string, error) {
	data, err := json.Marshal(c.items)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func toFloat64(v any, key ...string) float64 {
	switch val := v.(type) {
	case float64:
		return val
	case float32:
		return float64(val)
	case int:
		return float64(val)
	case int64:
		return float64(val)
	case int32:
		return float64(val)
	case uint:
		return float64(val)
	case uint64:
		return float64(val)
	case uint32:
		return float64(val)
	default:
		return 0
	}
}

func toString(v any) string {
	switch val := v.(type) {
	case string:
		return val
	default:
		return fmt.Sprintf("%v", val)
	}
}
