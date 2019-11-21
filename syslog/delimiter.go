package syslog

// Delimiter interface
type Delimiter interface {
	Push(b byte) (string, bool)
	Vestige() (string, bool)
}
