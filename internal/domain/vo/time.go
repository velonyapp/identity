package vo

import "time"

type Time struct {
	value time.Time
}

func NewTime(value time.Time) Time {
	return Time{value: value}
}

func NewTimeNow() Time {
	return Time{value: time.Now()}
}

func (t Time) AddSeconds(seconds int) Time {
	return Time{value: t.value.Add(time.Duration(seconds) * time.Second)}
}

func (t Time) AddMinutes(minutes int) Time {
	return Time{value: t.value.Add(time.Duration(minutes) * time.Minute)}
}

func (t Time) AddHours(hours int) Time {
	return Time{value: t.value.Add(time.Duration(hours) * time.Hour)}
}

func (t Time) Value() time.Time {
	return t.value
}

func (t Time) String() string {
	return t.value.Format(time.RFC3339)
}

func (t Time) Before(other Time) bool {
	return t.value.Before(other.value)
}

func (t Time) After(other Time) bool {
	return t.value.After(other.value)
}

func (t Time) Equal(other Time) bool {
	return t.value.Equal(other.value)
}
