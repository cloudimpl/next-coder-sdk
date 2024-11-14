package polycode

import (
	"encoding/json"
	"time"
)

type LogLevel string

const (
	DebugLevel LogLevel = "DEBUG"
	InfoLevel  LogLevel = "INFO"
	WarnLevel  LogLevel = "WARN"
	ErrorLevel LogLevel = "ERROR"
)

type LogMsg struct {
	Level     LogLevel               `json:"level"`
	Section   string                 `json:"section"`
	Tags      map[string]interface{} `json:"tags"`
	Timestamp int64                  `json:"timestamp"`
	Message   string                 `json:"message"`
}

type LogEntry struct {
	msg *LogMsg
}

func (entry *LogEntry) Str(key string, val string) *LogEntry {
	entry.msg.Tags[key] = val
	return entry
}

func (entry *LogEntry) Int64(key string, val int64) *LogEntry {
	entry.msg.Tags[key] = val
	return entry
}

func (entry *LogEntry) Float64(key string, val float64) *LogEntry {
	entry.msg.Tags[key] = val
	return entry
}

func (entry *LogEntry) Bool(key string, val bool) *LogEntry {
	entry.msg.Tags[key] = val
	return entry
}

func (entry *LogEntry) Done() {
	entry.msg.Timestamp = time.Now().UnixMicro()
	logJson, err := json.Marshal(entry.msg)
	if err == nil {
		println(string(logJson))
	}
}

func (entry *LogEntry) Msg(msg string) {
	entry.msg.Message = msg
	entry.Done()
}

type Logger interface {
	Log(level LogLevel) *LogEntry
	Debug() *LogEntry
	Info() *LogEntry
	Warn() *LogEntry
	Error() *LogEntry
}

type JsonLogger struct {
	section string
}

func (logger *JsonLogger) Log(level LogLevel) *LogEntry {
	return &LogEntry{
		msg: &LogMsg{
			Level:   level,
			Section: logger.section,
			Tags:    make(map[string]interface{}),
		},
	}
}

func (logger *JsonLogger) Debug() *LogEntry {
	return logger.Log(DebugLevel)
}

func (logger *JsonLogger) Info() *LogEntry {
	return logger.Log(InfoLevel)
}

func (logger *JsonLogger) Warn() *LogEntry {
	return logger.Log(WarnLevel)
}

func (logger *JsonLogger) Error() *LogEntry {
	return logger.Log(ErrorLevel)
}

func CreateLogger(section string) Logger {
	return &JsonLogger{
		section: section,
	}
}
