package logger

import (
	"bytes"
	"fmt"
	"os"
	"sort"
	"sync"
	"time"
)

type LogLevel string

const (
	LevelInfo  LogLevel = "INFO"
	LevelError LogLevel = "ERROR"
)

type LogEntry struct {
	Time    time.Time
	Level   LogLevel
	Message string
	Fields  map[string]interface{}
}

type Logger struct {
	logChan chan LogEntry
	wg      sync.WaitGroup
	done    chan struct{}
}

func NewLogger(bufferSize int) *Logger {
	if bufferSize <= 0 {
		bufferSize = 128
	}
	l := &Logger{
		logChan: make(chan LogEntry, bufferSize),
		done:    make(chan struct{}),
	}
	l.wg.Add(1)
	go l.processLogs()
	return l
}

func NewNop() *Logger { return &Logger{} }

func (l *Logger) processLogs() {
	defer l.wg.Done()
	for {
		select {
		case entry := <-l.logChan:
			writeEntry(entry)
		case <-l.done:
			for {
				select {
				case entry := <-l.logChan:
					writeEntry(entry)
				default:
					return
				}
			}
		}
	}
}

func writeEntry(entry LogEntry) {
	ts := entry.Time.Format("2006-01-02 15:04:05")
	var buf bytes.Buffer
	if len(entry.Fields) > 0 {
		keys := make([]string, 0, len(entry.Fields))
		for k := range entry.Fields {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			fmt.Fprintf(&buf, ", %s=%v", k, entry.Fields[k])
		}
	}
	out := fmt.Sprintf("[%s] %s: %s%s\n", ts, entry.Level, entry.Message, buf.String())
	if entry.Level == LevelError {
		_, _ = os.Stderr.WriteString(out)
	} else {
		_, _ = os.Stdout.WriteString(out)
	}
}

func (l *Logger) Info(message string, fields map[string]interface{}) {
	l.enqueue(LogEntry{
		Time:    time.Now(),
		Level:   LevelInfo,
		Message: message,
		Fields:  cloneOrEmpty(fields),
	})
}

func (l *Logger) Error(message string, fields map[string]interface{}) {
	l.enqueue(LogEntry{
		Time:    time.Now(),
		Level:   LevelError,
		Message: message,
		Fields:  cloneOrEmpty(fields),
	})
}

func (l *Logger) enqueue(e LogEntry) {
	if l == nil || l.logChan == nil {
		return
	}
	select {
	case l.logChan <- e:
	default:
	}
}

func (l *Logger) Shutdown() {
	if l == nil || l.done == nil {
		return
	}
	close(l.done)
	l.wg.Wait()
}

func cloneOrEmpty(m map[string]interface{}) map[string]interface{} {
	if m == nil {
		return map[string]interface{}{}
	}
	cp := make(map[string]interface{}, len(m))
	for k, v := range m {
		cp[k] = v
	}
	return cp
}
