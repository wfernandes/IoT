package logging

import (
	"bytes"
	"fmt"
)

type LogLevel int

const (
	SILENT LogLevel = iota
	FATAL
	ERROR
	INFO
	DEBUG
)

func (l LogLevel) String() string {
	switch l {
	case FATAL:
		return "FATAL"
	case ERROR:
		return "ERROR"
	case INFO:
		return "INFO"
	case DEBUG:
		return "DEBUG"
	default:
		return "INVALID"
	}
}

func (l *LogLevel) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, l.String())), nil
}

func (l *LogLevel) UnmarshalJSON(data []byte) error {
	data = bytes.Trim(data, `"`)
	strData := string(data)
	switch strData {
	case "FATAL":
		*l = FATAL
	case "ERROR":
		*l = ERROR
	case "INFO":
		*l = INFO
	case "DEBUG":
		*l = DEBUG
	default:
		return fmt.Errorf("Unknown LogLevel: %s", strData)
	}
	return nil
}
