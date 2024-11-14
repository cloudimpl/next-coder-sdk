package polycode

import (
	"sync"
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
	msg       *LogMsg
	publisher LogPublisher
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
	_ = entry.publisher.Publish(*entry.msg)
}

func (entry *LogEntry) Msg(msg string) {
	entry.msg.Message = msg
	entry.msg.Timestamp = time.Now().UnixMicro()
	_ = entry.publisher.Publish(*entry.msg)
}

type Logger struct {
	section   string
	publisher LogPublisher
}

func (logger *Logger) Log(level LogLevel) *LogEntry {
	return &LogEntry{
		msg: &LogMsg{
			Level:   level,
			Section: logger.section,
			Tags:    make(map[string]interface{}),
		},
		publisher: logger.publisher,
	}
}

func (logger *Logger) Debug() *LogEntry {
	return logger.Log(DebugLevel)
}

func (logger *Logger) Info() *LogEntry {
	return logger.Log(InfoLevel)
}

func (logger *Logger) Warn() *LogEntry {
	return logger.Log(WarnLevel)
}

func (logger *Logger) Error() *LogEntry {
	return logger.Log(ErrorLevel)
}

type LogPublisher interface {
	Publish(msg LogMsg) error
	PublishList(msgList []LogMsg) error
	Flush(process func(messages []LogMsg) error) error
}

type LogAggregator struct {
	messages []LogMsg
	mu       sync.Mutex
}

func (l *LogAggregator) PublishList(msgList []LogMsg) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.messages = append(l.messages, msgList...)
	return nil
}

func (l *LogAggregator) Publish(msg LogMsg) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.messages = append(l.messages, msg)
	return nil
}

func (l *LogAggregator) Flush(process func(messages []LogMsg) error) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	err := process(l.messages)
	if err != nil {
		return err
	}

	l.messages = make([]LogMsg, 0)
	return nil
}

func CreateLogger(section string, publisher LogPublisher) *Logger {
	return &Logger{
		section:   section,
		publisher: publisher,
	}
}

func CreateLogAggregator() *LogAggregator {
	return &LogAggregator{messages: make([]LogMsg, 0)}
}
