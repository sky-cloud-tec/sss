package syslog

// RFC interface
type RFC interface {
	compileMatcher()
	parse([]byte, *map[string]interface{})
}
