package console

import (
	"regexp"
	"strconv"
	"strings"
)

type ValidationRule struct {
	field   string
	rules   []ruleFunc
	value   any
	errors  []string
}

type ruleFunc func(value any, data map[string]any, field string) string

type Validator struct {
	data   map[string]any
	rules  map[string][]ruleFunc
	errors map[string][]string
}

func NewValidator() *Validator {
	return &Validator{
		rules:  make(map[string][]ruleFunc),
		errors: make(map[string][]string),
	}
}

func (v *Validator) Data(data map[string]any) *Validator {
	v.data = data
	return v
}

func (v *Validator) Rule(field string, rules ...ruleFunc) *Validator {
	v.rules[field] = append(v.rules[field], rules...)
	return v
}

func (v *Validator) Validate() bool {
	v.errors = make(map[string][]string)
	for field, rules := range v.rules {
		val := v.data[field]
		for _, rule := range rules {
			if err := rule(val, v.data, field); err != "" {
				v.errors[field] = append(v.errors[field], err)
			}
		}
	}
	return len(v.errors) == 0
}

func (v *Validator) Passes() bool {
	return v.Validate()
}

func (v *Validator) Fails() bool {
	return !v.Validate()
}

func (v *Validator) Errors() map[string][]string {
	return v.errors
}

func (v *Validator) ErrorsFor(field string) []string {
	return v.errors[field]
}

func Required() ruleFunc {
	return func(value any, data map[string]any, field string) string {
		if value == nil {
			return "The " + field + " field is required."
		}
		s, ok := value.(string)
		if ok && strings.TrimSpace(s) == "" {
			return "The " + field + " field is required."
		}
		return ""
	}
}

func Email() ruleFunc {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return func(value any, data map[string]any, field string) string {
		s, ok := value.(string)
		if !ok || !re.MatchString(s) {
			return "The " + field + " must be a valid email address."
		}
		return ""
	}
}

func MinInt(min int) ruleFunc {
	return func(value any, data map[string]any, field string) string {
		n, ok := value.(int)
		if !ok {
			f, ok := value.(float64)
			if ok {
				n = int(f)
			} else {
				return "The " + field + " must be a number."
			}
		}
		if n < min {
			return "The " + field + " must be at least " + strconv.Itoa(min) + "."
		}
		return ""
	}
}

func MaxInt(max int) ruleFunc {
	return func(value any, data map[string]any, field string) string {
		n, ok := value.(int)
		if !ok {
			f, ok := value.(float64)
			if ok {
				n = int(f)
			} else {
				return "The " + field + " must be a number."
			}
		}
		if n > max {
			return "The " + field + " must not exceed " + strconv.Itoa(max) + "."
		}
		return ""
	}
}

func MinLen(min int) ruleFunc {
	return func(value any, data map[string]any, field string) string {
		s, ok := value.(string)
		if !ok {
			return "The " + field + " must be a string."
		}
		if len([]rune(s)) < min {
			return "The " + field + " must be at least " + strconv.Itoa(min) + " characters."
		}
		return ""
	}
}

func MaxLen(max int) ruleFunc {
	return func(value any, data map[string]any, field string) string {
		s, ok := value.(string)
		if !ok {
			return "The " + field + " must be a string."
		}
		if len([]rune(s)) > max {
			return "The " + field + " must not exceed " + strconv.Itoa(max) + " characters."
		}
		return ""
	}
}

func BetweenInt(min, max int) ruleFunc {
	return func(value any, data map[string]any, field string) string {
		n, ok := value.(int)
		if !ok {
			f, ok := value.(float64)
			if ok {
				n = int(f)
			} else {
				return "The " + field + " must be a number."
			}
		}
		if n < min || n > max {
			return "The " + field + " must be between " + strconv.Itoa(min) + " and " + strconv.Itoa(max) + "."
		}
		return ""
	}
}

func In(values ...any) ruleFunc {
	return func(value any, data map[string]any, field string) string {
		for _, v := range values {
			if v == value {
				return ""
			}
		}
		return "The " + field + " must be one of: " + strings.Join(stringsSplit(values), ", ") + "."
	}
}

func NotIn(values ...any) ruleFunc {
	return func(value any, data map[string]any, field string) string {
		for _, v := range values {
			if v == value {
				return "The " + field + " is invalid."
			}
		}
		return ""
	}
}

func Match(pattern string) ruleFunc {
	re := regexp.MustCompile(pattern)
	return func(value any, data map[string]any, field string) string {
		s, ok := value.(string)
		if !ok || !re.MatchString(s) {
			return "The " + field + " format is invalid."
		}
		return ""
	}
}

func Numeric() ruleFunc {
	return func(value any, data map[string]any, field string) string {
		switch value.(type) {
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
			return ""
		}
		s, ok := value.(string)
		if ok {
			_, err := strconv.ParseFloat(s, 64)
			if err == nil {
				return ""
			}
		}
		return "The " + field + " must be numeric."
	}
}

func String() ruleFunc {
	return func(value any, data map[string]any, field string) string {
		_, ok := value.(string)
		if !ok {
			return "The " + field + " must be a string."
		}
		return ""
	}
}

func Bool() ruleFunc {
	return func(value any, data map[string]any, field string) string {
		_, ok := value.(bool)
		if !ok {
			return "The " + field + " must be a boolean."
		}
		return ""
	}
}

func Integer() ruleFunc {
	return func(value any, data map[string]any, field string) string {
		_, ok := value.(int)
		if !ok {
			f, ok := value.(float64)
			if ok && f == float64(int(f)) {
				return ""
			}
			return "The " + field + " must be an integer."
		}
		_ = ok
		return ""
	}
}

func Confirmed() ruleFunc {
	return func(value any, data map[string]any, field string) string {
		confField := field + "_confirmation"
		conf, ok := data[confField]
		if !ok {
			return "The " + field + " confirmation does not match."
		}
		if value != conf {
			return "The " + field + " confirmation does not match."
		}
		return ""
	}
}

func URL() ruleFunc {
	re := regexp.MustCompile(`^https?://\S+$`)
	return func(value any, data map[string]any, field string) string {
		s, ok := value.(string)
		if !ok || !re.MatchString(s) {
			return "The " + field + " must be a valid URL."
		}
		return ""
	}
}

func stringsSplit(values []any) []string {
	result := make([]string, len(values))
	for i, v := range values {
		result[i] = toString(v)
	}
	return result
}

func NewRuleSet(rules ...ruleFunc) []ruleFunc {
	return rules
}
