package log

const (
	DEBUG LogLevel = iota
	INFO
	WARNING
	ERROR
	SILENT
)

var LogLevelMapping = map[string]LogLevel{
	DEBUG.String():   DEBUG,
	INFO.String():    INFO,
	WARNING.String(): WARNING,
	ERROR.String():   ERROR,
	SILENT.String():  SILENT,
}

type LogLevel int

func (l LogLevel) String() string {
	switch l {
	case INFO:
		return "info"
	case WARNING:
		return "warning"
	case ERROR:
		return "error"
	case DEBUG:
		return "debug"
	case SILENT:
		return "silent"
	default:
		return "unknown"
	}
}
