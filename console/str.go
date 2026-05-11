package console

import (
	"crypto/rand"
	"math/big"
	"strings"
	"unicode"
)

type Str struct {
	value string
}

func NewStr(s string) *Str {
	return &Str{value: s}
}

func (s *Str) String() string {
	return s.value
}

func (s *Str) Contains(substr string) bool {
	return strings.Contains(s.value, substr)
}

func (s *Str) ContainsAll(substrs ...string) bool {
	for _, sub := range substrs {
		if !strings.Contains(s.value, sub) {
			return false
		}
	}
	return true
}

func (s *Str) ContainsAny(substrs ...string) bool {
	for _, sub := range substrs {
		if strings.Contains(s.value, sub) {
			return true
		}
	}
	return false
}

func (s *Str) StartsWith(prefix string) bool {
	return strings.HasPrefix(s.value, prefix)
}

func (s *Str) EndsWith(suffix string) bool {
	return strings.HasSuffix(s.value, suffix)
}

func (s *Str) Upper() *Str {
	return &Str{value: strings.ToUpper(s.value)}
}

func (s *Str) Lower() *Str {
	return &Str{value: strings.ToLower(s.value)}
}

func (s *Str) Title() *Str {
	return &Str{value: strings.Title(s.value)}
}

func (s *Str) Ucfirst() *Str {
	if s.value == "" {
		return &Str{value: ""}
	}
	r := []rune(s.value)
	return &Str{value: string(unicode.ToUpper(r[0])) + string(r[1:])}
}

func (s *Str) Lcfirst() *Str {
	if s.value == "" {
		return &Str{value: ""}
	}
	r := []rune(s.value)
	return &Str{value: string(unicode.ToLower(r[0])) + string(r[1:])}
}

func (s *Str) Length() int {
	return len([]rune(s.value))
}

func (s *Str) Limit(length int, append ...string) *Str {
	suffix := "..."
	if len(append) > 0 {
		suffix = append[0]
	}
	r := []rune(s.value)
	if len(r) <= length {
		return &Str{value: s.value}
	}
	return &Str{value: string(r[:length]) + suffix}
}

func (s *Str) Substr(start, length int) *Str {
	r := []rune(s.value)
	if start < 0 {
		start = len(r) + start
	}
	if start < 0 {
		start = 0
	}
	if start >= len(r) {
		return &Str{value: ""}
	}
	end := start + length
	if end > len(r) {
		end = len(r)
	}
	return &Str{value: string(r[start:end])}
}

func (s *Str) Replace(search, replace string) *Str {
	return &Str{value: strings.ReplaceAll(s.value, search, replace)}
}

func (s *Str) ReplaceFirst(search, replace string) *Str {
	return &Str{value: strings.Replace(s.value, search, replace, 1)}
}

func (s *Str) ReplaceLast(search, replace string) *Str {
	i := strings.LastIndex(s.value, search)
	if i == -1 {
		return &Str{value: s.value}
	}
	return &Str{value: s.value[:i] + replace + s.value[i+len(search):]}
}

func (s *Str) Trim(cutset ...string) *Str {
	if len(cutset) > 0 {
		return &Str{value: strings.Trim(s.value, cutset[0])}
	}
	return &Str{value: strings.TrimSpace(s.value)}
}

func (s *Str) Ltrim(cutset ...string) *Str {
	if len(cutset) > 0 {
		return &Str{value: strings.TrimLeft(s.value, cutset[0])}
	}
	return &Str{value: strings.TrimLeftFunc(s.value, unicode.IsSpace)}
}

func (s *Str) Rtrim(cutset ...string) *Str {
	if len(cutset) > 0 {
		return &Str{value: strings.TrimRight(s.value, cutset[0])}
	}
	return &Str{value: strings.TrimRightFunc(s.value, unicode.IsSpace)}
}

func (s *Str) PadLeft(length int, pad ...string) *Str {
	ch := " "
	if len(pad) > 0 {
		ch = pad[0]
	}
	r := []rune(s.value)
	if len(r) >= length {
		return &Str{value: s.value}
	}
	padding := strings.Repeat(ch, length-len(r))
	return &Str{value: padding + s.value}
}

func (s *Str) PadRight(length int, pad ...string) *Str {
	ch := " "
	if len(pad) > 0 {
		ch = pad[0]
	}
	r := []rune(s.value)
	if len(r) >= length {
		return &Str{value: s.value}
	}
	padding := strings.Repeat(ch, length-len(r))
	return &Str{value: s.value + padding}
}

func (s *Str) Before(delimiter string) *Str {
	i := strings.Index(s.value, delimiter)
	if i == -1 {
		return &Str{value: s.value}
	}
	return &Str{value: s.value[:i]}
}

func (s *Str) After(delimiter string) *Str {
	i := strings.Index(s.value, delimiter)
	if i == -1 {
		return &Str{value: ""}
	}
	return &Str{value: s.value[i+len(delimiter):]}
}

func (s *Str) Between(from, to string) *Str {
	after := s.After(from)
	return after.Before(to)
}

func (s *Str) Is(pattern string) bool {
	parts := strings.Split(pattern, "*")
	if len(parts) == 1 {
		return s.value == pattern
	}
	if !strings.HasPrefix(s.value, parts[0]) {
		return false
	}
	remain := s.value[len(parts[0]):]
	for i := 1; i < len(parts); i++ {
		idx := strings.Index(remain, parts[i])
		if idx == -1 {
			return false
		}
		if i == len(parts)-1 && parts[i] == "" {
			return true
		}
		remain = remain[idx+len(parts[i]):]
	}
	if len(parts[len(parts)-1]) > 0 {
		return strings.HasSuffix(s.value, parts[len(parts)-1])
	}
	return remain == ""
}

func (s *Str) Slug(separator ...string) *Str {
	sep := "-"
	if len(separator) > 0 {
		sep = separator[0]
	}
	var result strings.Builder
	for _, r := range strings.ToLower(s.value) {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			result.WriteRune(r)
		} else if unicode.IsSpace(r) || r == '-' || r == '_' {
			result.WriteString(sep)
		}
	}
	return &Str{value: result.String()}
}

func (s *Str) Snake() *Str {
	return &Str{value: toSnakeCase(s.value)}
}

func (s *Str) Kebab() *Str {
	return &Str{value: strings.ReplaceAll(toSnakeCase(s.value), "_", "-")}
}

func (s *Str) Studly() *Str {
	return &Str{value: toStudlyCase(s.value)}
}

func (s *Str) Camel() *Str {
	return s.Studly().Lcfirst()
}

func (s *Str) Repeat(count int) *Str {
	return &Str{value: strings.Repeat(s.value, count)}
}

func (s *Str) Mask(ch rune, start, length int) *Str {
	r := []rune(s.value)
	if start >= len(r) {
		return &Str{value: s.value}
	}
	end := start + length
	if end > len(r) {
		end = len(r)
	}
	result := make([]rune, len(r))
	copy(result, r)
	for i := start; i < end; i++ {
		result[i] = ch
	}
	return &Str{value: string(result)}
}

func (s *Str) IsEmpty() bool {
	return s.value == ""
}

func (s *Str) IsNotEmpty() bool {
	return s.value != ""
}

func (s *Str) Equals(other string) bool {
	return s.value == other
}

func (s *Str) EqualsIgnoreCase(other string) bool {
	return strings.EqualFold(s.value, other)
}

func Random(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		result[i] = charset[n.Int64()]
	}
	return string(result)
}

func toSnakeCase(s string) string {
	var result strings.Builder
	r := []rune(s)
	for i, ch := range r {
		if unicode.IsUpper(ch) {
			if i > 0 && (unicode.IsLower(r[i-1]) || (i+1 < len(r) && unicode.IsLower(r[i+1]))) {
				result.WriteRune('_')
			}
			result.WriteRune(unicode.ToLower(ch))
		} else if ch == '-' || ch == ' ' {
			result.WriteRune('_')
		} else {
			result.WriteRune(ch)
		}
	}
	return result.String()
}

func toStudlyCase(s string) string {
	var result strings.Builder
	nextUpper := true
	for _, ch := range s {
		if ch == '_' || ch == '-' || ch == ' ' {
			nextUpper = true
		} else if nextUpper {
			result.WriteRune(unicode.ToUpper(ch))
			nextUpper = false
		} else {
			result.WriteRune(ch)
		}
	}
	return result.String()
}
