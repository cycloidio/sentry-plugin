package issue

import "time"

type Issue struct {
	ID        string
	Title     string
	Permalink string

	HasSeen   bool
	FirstSeen time.Time
	LastSeen  time.Time
	UserCount int

	Level  string
	Status string
	Type   string
}
