package console

import (
	"fmt"
	"time"
)

type Date struct {
	t time.Time
}

var dateLayous = []string{
	time.RFC3339,
	"2006-01-02T15:04:05",
	"2006-01-02 15:04:05",
	"2006-01-02",
	"2006/01/02",
	"01/02/2006",
	"02-01-2006",
	"2006-01-02T15:04:05Z07:00",
	time.RFC822,
	time.RFC1123,
}

func NewDate(t ...time.Time) *Date {
	if len(t) > 0 {
		return &Date{t: t[0]}
	}
	return &Date{t: time.Now()}
}

func Now() *Date {
	return &Date{t: time.Now()}
}

func Today() *Date {
	now := time.Now()
	return &Date{t: time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())}
}

func Parse(value string) *Date {
	for _, layout := range dateLayous {
		t, err := time.Parse(layout, value)
		if err == nil {
			return &Date{t: t}
		}
	}
	return &Date{t: time.Now()}
}

func (d *Date) Time() time.Time {
	return d.t
}

func (d *Date) Format(layout string) string {
	return d.t.Format(layout)
}

func (d *Date) ToDateTimeString() string {
	return d.t.Format("2006-01-02 15:04:05")
}

func (d *Date) ToDateString() string {
	return d.t.Format("2006-01-02")
}

func (d *Date) ToTimeString() string {
	return d.t.Format("15:04:05")
}

func (d *Date) Timestamp() int64 {
	return d.t.Unix()
}

func (d *Date) Year() int {
	return d.t.Year()
}

func (d *Date) Month() time.Month {
	return d.t.Month()
}

func (d *Date) Day() int {
	return d.t.Day()
}

func (d *Date) Hour() int {
	return d.t.Hour()
}

func (d *Date) Minute() int {
	return d.t.Minute()
}

func (d *Date) Second() int {
	return d.t.Second()
}

func (d *Date) Weekday() time.Weekday {
	return d.t.Weekday()
}

func (d *Date) AddDays(n int) *Date {
	return &Date{t: d.t.AddDate(0, 0, n)}
}

func (d *Date) SubDays(n int) *Date {
	return &Date{t: d.t.AddDate(0, 0, -n)}
}

func (d *Date) AddHours(n int) *Date {
	return &Date{t: d.t.Add(time.Duration(n) * time.Hour)}
}

func (d *Date) SubHours(n int) *Date {
	return &Date{t: d.t.Add(time.Duration(-n) * time.Hour)}
}

func (d *Date) AddMinutes(n int) *Date {
	return &Date{t: d.t.Add(time.Duration(n) * time.Minute)}
}

func (d *Date) AddMonths(n int) *Date {
	return &Date{t: d.t.AddDate(0, n, 0)}
}

func (d *Date) AddYears(n int) *Date {
	return &Date{t: d.t.AddDate(n, 0, 0)}
}

func (d *Date) StartOfDay() *Date {
	y, m, day := d.t.Date()
	return &Date{t: time.Date(y, m, day, 0, 0, 0, 0, d.t.Location())}
}

func (d *Date) EndOfDay() *Date {
	y, m, day := d.t.Date()
	return &Date{t: time.Date(y, m, day, 23, 59, 59, 999999999, d.t.Location())}
}

func (d *Date) StartOfMonth() *Date {
	return &Date{t: time.Date(d.t.Year(), d.t.Month(), 1, 0, 0, 0, 0, d.t.Location())}
}

func (d *Date) EndOfMonth() *Date {
	return d.StartOfMonth().AddMonths(1).SubDays(1).EndOfDay()
}

func (d *Date) StartOfYear() *Date {
	return &Date{t: time.Date(d.t.Year(), 1, 1, 0, 0, 0, 0, d.t.Location())}
}

func (d *Date) EndOfYear() *Date {
	return &Date{t: time.Date(d.t.Year(), 12, 31, 23, 59, 59, 999999999, d.t.Location())}
}

func (d *Date) IsWeekend() bool {
	switch d.t.Weekday() {
	case time.Saturday, time.Sunday:
		return true
	}
	return false
}

func (d *Date) IsWeekday() bool {
	return !d.IsWeekend()
}

func (d *Date) IsPast() bool {
	return d.t.Before(time.Now())
}

func (d *Date) IsFuture() bool {
	return d.t.After(time.Now())
}

func (d *Date) IsToday() bool {
	y1, m1, d1 := d.t.Date()
	y2, m2, d2 := time.Now().Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func (d *Date) DiffInDays(other *Date) int {
	return int(d.t.Sub(other.t).Hours() / 24)
}

func (d *Date) DiffInHours(other *Date) float64 {
	return d.t.Sub(other.t).Hours()
}

func (d *Date) DiffInMinutes(other *Date) float64 {
	return d.t.Sub(other.t).Minutes()
}

func (d *Date) DiffInSeconds(other *Date) float64 {
	return d.t.Sub(other.t).Seconds()
}

func (d *Date) Gt(other *Date) bool {
	return d.t.After(other.t)
}

func (d *Date) Lt(other *Date) bool {
	return d.t.Before(other.t)
}

func (d *Date) Eq(other *Date) bool {
	return d.t.Equal(other.t)
}

func (d *Date) String() string {
	return d.ToDateTimeString()
}

func (d *Date) Age() int {
	return int(time.Since(d.t).Hours() / 8760)
}

func (d *Date) Copy() *Date {
	return &Date{t: d.t}
}

func Yesterday() *Date {
	return Now().SubDays(1)
}

func Tomorrow() *Date {
	return Now().AddDays(1)
}

func ParseDate(value string) *Date {
	return Parse(value)
}

func (d *Date) Rfc3339() string {
	return d.t.Format(time.RFC3339)
}

func (d *Date) ISO8601() string {
	return d.t.Format("2006-01-02T15:04:05Z07:00")
}

func (d *Date) HumanDiff() string {
	diff := time.Since(d.t)
	abs := diff
	if abs < 0 {
		abs = -abs
	}

	ago := "ago"
	if diff < 0 {
		ago = "from now"
	}

	switch {
	case abs < time.Minute:
		return fmt.Sprintf("%d seconds %s", int(abs.Seconds()), ago)
	case abs < time.Hour:
		return fmt.Sprintf("%d minutes %s", int(abs.Minutes()), ago)
	case abs < 24*time.Hour:
		return fmt.Sprintf("%d hours %s", int(abs.Hours()), ago)
	case abs < 30*24*time.Hour:
		return fmt.Sprintf("%d days %s", int(abs.Hours()/24), ago)
	case abs < 365*24*time.Hour:
		return fmt.Sprintf("%d months %s", int(abs.Hours()/(24*30)), ago)
	default:
		return fmt.Sprintf("%d years %s", int(abs.Hours()/(24*365)), ago)
	}
}
