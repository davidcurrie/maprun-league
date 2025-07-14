package config

import "time"

type ClosingDate struct {
	time.Time
}

func (d *ClosingDate) UnmarshalText(text []byte) error {
	t, err := time.Parse("2006-01-02", string(text))
	if err != nil {
		return err
	}
	d.Time = t
	return nil
}

func (d ClosingDate) MarshalText() ([]byte, error) {
	return []byte(d.Time.Format("2006-01-02")), nil
}

func (d ClosingDate) After(t time.Time) bool {
	nextDay := d.AddDate(0, 0, 1).Truncate(24 * time.Hour)
	return nextDay.After(t)
}
