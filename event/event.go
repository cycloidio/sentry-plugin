package event

import "time"

//go:generate go tool enumer -type=Severity -transform=snake -json -output=severity_enumer_gen.go

// Severity is the severity associated to an event
type Severity uint8

// Available and only accepted severities
const (
	Info Severity = iota + 1
	Warn
	Err
	Crit
)

//go:generate go tool enumer -type=Type -json -output=type_enumer_gen.go

// Type is the type of an event
type Type uint8

// Available and only accepted types
const (
	Cycloid Type = iota + 1
	AWS
	Monitoring
	Custom
)

//go:generate go tool enumer -type=Color -linecomment -transform=lower -json -output=color_enumer_gen.go

// Color is the color associated to an event
type Color uint8

// List of colors from https://developer.mozilla.org/en-US/docs/Web/CSS/color_value.
// As specified by CSS 2: https://www.w3.org/TR/CSS2/syndata.html#value-def-color.
const (
	// NoColor maps to the empty string, this means use the  default  color
	// for the event.
	NoColor Color = iota //
	Black
	Silver
	Gray
	White
	Maroon
	Red
	Purple
	Fuchsia
	Green
	Lime
	Olive
	Yellow
	Navy
	Blue
	Teal
	Aqua
	Orange
)

type Event struct {
	ID        uint32
	Timestamp time.Time
	Title     string
	Message   string
	Icon      string
	Tags      []map[string]string
	Severity  Severity
	Type      Type
	Color     Color
}
