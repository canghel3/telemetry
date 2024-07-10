package log

import (
	"fmt"
	"telemetry/drivers"
	"telemetry/level"
	"time"
)

type Log struct {
	driver   drivers.Generic
	metadata map[any]any
	buf      []byte
	level    uint8
}

func (l *Log) NoLevel() *Log {
	l.level = level.NoLevel
	return l
}

func (l *Log) Error() *Log {
	l.level = level.LevelError
	return l
}

func (l *Log) Warn() *Log {
	l.level = level.LevelWarn
	return l
}

func (l *Log) Info() *Log {
	l.level = level.LevelInfo
	return l
}

func (l *Log) Debug() *Log {
	l.level = level.LevelDebug
	return l
}

func (l *Log) Meta(data map[any]any) *Log {
	l.metadata = data
	return l
}

func File(name string) *Log {
	return &Log{
		driver: drivers.ToFileWithName(name),
	}
}

func Stdout() *Log {
	return &Log{
		driver: drivers.ToStdout(),
	}
}

func Custom(driver drivers.Generic) *Log {
	return &Log{
		driver: driver,
	}
}

func (l *Log) Write(b []byte) {
	l.buf = b
	err := l.driver.Write(formatLogOutput(l))
	if err != nil {
		return
	}
}

func formatLogOutput(l *Log) []byte {
	//TIMESTAMP LEVEL METDATA BUFFER
	timestamp := time.Now()

	var out = make([]byte, 0)
	out = append(out, []byte(timestamp.String())...)
	out = append(out, byte(' '))
	out = append(out, []byte(level.LevelToText[l.level])...)
	out = append(out, byte(' '))

	var meta2bytes = make([]byte, 0)
	for k, v := range l.metadata {
		meta2bytes = append(meta2bytes, []byte(fmt.Sprintf("%v:%v ", k, v))...)
	}
	out = append(out, meta2bytes...)
	out = append(out, l.buf...)
	out = append(out, []byte("\n")...)

	return out
}
