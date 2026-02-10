package service

//go:generate go tool enumer -type=Status -transform=snake -output=status_string.go
type Status int

const (
	Ok Status = iota
	Syncthing
	Error
)
