package httpx

import (
	"strings"
	"time"
)

// Date marshals/unmarshals as "YYYY-MM-DD", matching what date-only inputs
// (HTML <input type="date">, due dates, inspection dates) actually send —
// unlike time.Time's default RFC3339, which rejects a bare date string.
type Date struct {
	time.Time
}

func (d Date) MarshalJSON() ([]byte, error) {
	if d.IsZero() {
		return []byte("null"), nil
	}
	return []byte(`"` + d.Format("2006-01-02") + `"`), nil
}

func (d *Date) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)
	if s == "" || s == "null" {
		return nil
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	d.Time = t
	return nil
}

// DateToTime unwraps a *Date into a *time.Time for the repository layer,
// which stores dates as time.Time under the hood.
func DateToTime(d *Date) *time.Time {
	if d == nil {
		return nil
	}
	t := d.Time
	return &t
}

// TimeToDate wraps a *time.Time coming back from the repository layer into
// a *Date for JSON responses.
func TimeToDate(t *time.Time) *Date {
	if t == nil {
		return nil
	}
	return &Date{Time: *t}
}
