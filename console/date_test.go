package console

import (
	"testing"
	"time"
)

func TestDateNow(t *testing.T) {
	d := Now()
	if d.Time().IsZero() {
		t.Fatal("expected non-zero time")
	}
}

func TestDateToday(t *testing.T) {
	d := Today()
	if d.Hour() != 0 || d.Minute() != 0 || d.Second() != 0 {
		t.Fatalf("expected start of day, got %d:%d:%d", d.Hour(), d.Minute(), d.Second())
	}
}

func TestDateParse(t *testing.T) {
	d := Parse("2024-01-15")
	if d.Year() != 2024 || d.Month() != time.January || d.Day() != 15 {
		t.Fatalf("expected 2024-01-15, got %d-%d-%d", d.Year(), d.Month(), d.Day())
	}
}

func TestDateParseWithTime(t *testing.T) {
	d := Parse("2024-01-15 14:30:00")
	if d.Hour() != 14 || d.Minute() != 30 {
		t.Fatalf("expected 14:30, got %d:%d", d.Hour(), d.Minute())
	}
}

func TestDateFormat(t *testing.T) {
	d := Parse("2024-03-15")
	if d.Format("2006/01/02") != "2024/03/15" {
		t.Fatalf("expected '2024/03/15', got '%s'", d.Format("2006/01/02"))
	}
}

func TestDateToDateTimeString(t *testing.T) {
	d := Parse("2024-06-20 10:30:00")
	if d.ToDateTimeString() != "2024-06-20 10:30:00" {
		t.Fatalf("expected '2024-06-20 10:30:00', got '%s'", d.ToDateTimeString())
	}
}

func TestDateToDateString(t *testing.T) {
	d := Parse("2024-06-20")
	if d.ToDateString() != "2024-06-20" {
		t.Fatalf("expected '2024-06-20', got '%s'", d.ToDateString())
	}
}

func TestDateAddDays(t *testing.T) {
	d := Parse("2024-01-01")
	r := d.AddDays(5)
	if r.Day() != 6 {
		t.Fatalf("expected day 6, got %d", r.Day())
	}
}

func TestDateSubDays(t *testing.T) {
	d := Parse("2024-01-10")
	r := d.SubDays(5)
	if r.Day() != 5 {
		t.Fatalf("expected day 5, got %d", r.Day())
	}
}

func TestDateAddMonths(t *testing.T) {
	d := Parse("2024-01-01")
	r := d.AddMonths(2)
	if r.Month() != time.March {
		t.Fatalf("expected March, got %s", r.Month())
	}
}

func TestDateAddYears(t *testing.T) {
	d := Parse("2024-01-01")
	r := d.AddYears(1)
	if r.Year() != 2025 {
		t.Fatalf("expected 2025, got %d", r.Year())
	}
}

func TestDateStartOfDay(t *testing.T) {
	d := Parse("2024-06-15 14:30:00")
	s := d.StartOfDay()
	if s.Hour() != 0 || s.Minute() != 0 || s.Second() != 0 {
		t.Fatal("expected start of day")
	}
}

func TestDateEndOfDay(t *testing.T) {
	d := Parse("2024-06-15 14:30:00")
	e := d.EndOfDay()
	if e.Hour() != 23 || e.Minute() != 59 {
		t.Fatalf("expected 23:59, got %d:%d", e.Hour(), e.Minute())
	}
}

func TestDateStartOfMonth(t *testing.T) {
	d := Parse("2024-06-15")
	s := d.StartOfMonth()
	if s.Day() != 1 {
		t.Fatalf("expected day 1, got %d", s.Day())
	}
}

func TestDateEndOfMonth(t *testing.T) {
	d := Parse("2024-06-15")
	e := d.EndOfMonth()
	if e.Day() != 30 {
		t.Fatalf("expected day 30, got %d", e.Day())
	}
}

func TestDateIsWeekend(t *testing.T) {
	sat := Parse("2024-06-15") // Saturday
	if !sat.IsWeekend() {
		t.Fatal("Saturday should be weekend")
	}
	mon := Parse("2024-06-17") // Monday
	if mon.IsWeekend() {
		t.Fatal("Monday should not be weekend")
	}
}

func TestDateIsWeekday(t *testing.T) {
	mon := Parse("2024-06-17") // Monday
	if !mon.IsWeekday() {
		t.Fatal("Monday should be weekday")
	}
}

func TestDateDiffInDays(t *testing.T) {
	d1 := Parse("2024-01-10")
	d2 := Parse("2024-01-01")
	diff := d1.DiffInDays(d2)
	if diff != 9 {
		t.Fatalf("expected 9 days, got %d", diff)
	}
}

func TestDateDiffInHours(t *testing.T) {
	d1 := Parse("2024-01-02 12:00:00")
	d2 := Parse("2024-01-01 12:00:00")
	diff := d1.DiffInHours(d2)
	if diff != 24 {
		t.Fatalf("expected 24 hours, got %f", diff)
	}
}

func TestDateGt(t *testing.T) {
	d1 := Parse("2024-01-10")
	d2 := Parse("2024-01-01")
	if !d1.Gt(d2) {
		t.Fatal("expected d1 > d2")
	}
	if d2.Gt(d1) {
		t.Fatal("expected d2 < d1")
	}
}

func TestDateLt(t *testing.T) {
	d1 := Parse("2024-01-01")
	d2 := Parse("2024-01-10")
	if !d1.Lt(d2) {
		t.Fatal("expected d1 < d2")
	}
}

func TestDateEq(t *testing.T) {
	d1 := Parse("2024-01-01")
	d2 := Parse("2024-01-01")
	if !d1.Eq(d2) {
		t.Fatal("expected equal")
	}
}

func TestYesterday(t *testing.T) {
	y := Yesterday()
	if y.ToDateString() != Now().SubDays(1).ToDateString() {
		t.Fatal("yesterday failed")
	}
}

func TestTomorrow(t *testing.T) {
	tm := Tomorrow()
	if tm.ToDateString() != Now().AddDays(1).ToDateString() {
		t.Fatal("tomorrow failed")
	}
}

func TestDateTimestamp(t *testing.T) {
	d := Parse("2024-01-01")
	if d.Timestamp() <= 0 {
		t.Fatal("expected positive timestamp")
	}
}

func TestDateYearMonthDay(t *testing.T) {
	d := Parse("2024-06-15")
	if d.Year() != 2024 || d.Month() != time.June || d.Day() != 15 {
		t.Fatal("year/month/day failed")
	}
}

func TestDateCopy(t *testing.T) {
	d1 := Parse("2024-01-01")
	d2 := d1.Copy()
	d2.AddDays(1)
	if d1.Day() != 1 {
		t.Fatal("copy should be independent")
	}
}

func TestDateHumanDiff(t *testing.T) {
	d := Now().SubHours(2)
	diff := d.HumanDiff()
	if diff == "" {
		t.Fatal("expected non-empty human diff")
	}
}

func TestParseDateAlias(t *testing.T) {
	d := ParseDate("2024-12-25")
	if d.Month() != time.December || d.Day() != 25 {
		t.Fatal("parse date alias failed")
	}
}

func TestDateRfc3339(t *testing.T) {
	d := Parse("2024-01-01")
	rfc := d.Rfc3339()
	if rfc == "" {
		t.Fatal("expected non-empty RFC3339")
	}
}

func TestDateISO8601(t *testing.T) {
	d := Parse("2024-01-01")
	iso := d.ISO8601()
	if iso == "" {
		t.Fatal("expected non-empty ISO8601")
	}
}

func TestDateAge(t *testing.T) {
	d := Parse("2000-01-01")
	if d.Age() < 20 {
		t.Fatalf("expected age >= 20, got %d", d.Age())
	}
}

func TestDateChaining(t *testing.T) {
	d := Parse("2024-01-15").
		AddDays(5).
		StartOfMonth()
	if d.Day() != 1 {
		t.Fatalf("chaining expected day 1, got %d", d.Day())
	}
}
