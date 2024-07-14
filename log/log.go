package log

import (
	"fmt"
	"telemetry/drivers"
	"telemetry/level"
	"time"
)

type Logger struct {
	buf          []byte
	level        level.Level
	outputDriver drivers.OutputDriver
	metadata     map[any]any
}

func Default() *Logger {
	return &Logger{
		outputDriver: drivers.ToStdout(),
		metadata:     nil,
		buf:          nil,
		level:        level.None(),
	}
}

func (l *Logger) NoLevel() *Logger {
	l.level = level.None()
	return l
}

// Error sets the level type to ERROR.
func (l *Logger) Error() *Logger {
	l.level = level.Error()
	return l
}

// Warn sets the level type to WARN.
func (l *Logger) Warn() *Logger {
	l.level = level.Warn()
	return l
}

func (l *Logger) Info() *Logger {
	l.level = level.Info()
	return l
}

func (l *Logger) Debug() *Logger {
	l.level = level.Debug()
	return l
}

func (l *Logger) Level(lvl level.Level) *Logger {
	l.level = lvl
	return l
}

func (l *Logger) Metadata(data map[any]any) *Logger {
	l.metadata = data
	return l
}

func File(name string) *Logger {
	l := Default()
	l.outputDriver = drivers.ToFileWithName(name)
	return l
}

func Stdout() *Logger {
	l := Default()
	l.outputDriver = drivers.ToStdout()
	return l
}

func CustomOutputDriver(driver drivers.OutputDriver) *Logger {
	l := Default()
	l.outputDriver = driver
	return l
}

// Msg returns a copy of the Logger receiver, except for the buf, which is overwritten.
// Useful for logging transactions.
func (l *Logger) Msg(b []byte) *Logger {
	cpy := &Logger{
		buf:          b,
		level:        l.level,
		outputDriver: l.outputDriver,
		metadata:     l.metadata,
	}

	cpy.buf = formatLogOutput(cpy)
	return cpy
}

// Log sends the (current log buffer + received "b") to the output driver for further handling.
// The buffer is emptied and can be reused.
// If an error occurs during writing, it panics.
func (l *Logger) Log(b []byte) {
	l.buf = append(l.buf, b...)
	err := l.outputDriver.Log(formatLogOutput(l))
	if err != nil {
		panic(err)
	}

	l.buf = []byte{}
}

// TODO: enable config formatting
func formatLogOutput(l *Logger) []byte {
	//TIMESTAMP LEVEL METADATA BUFFER
	timestamp := time.Now()

	var out = make([]byte, 0)
	out = append(out, []byte(timestamp.String())...)
	out = append(out, byte(' '))
	out = append(out, []byte(l.level.Type())...)
	out = append(out, byte(' '))

	var meta2bytes = make([]byte, 0)
	for k, v := range l.metadata {
		meta2bytes = append(meta2bytes, []byte(fmt.Sprintf("%v:%v ", k, v))...) //careful
	}
	//very careful (whitespace)
	out = append(out, meta2bytes...)
	out = append(out, l.buf...)
	out = append(out, '\n')

	return out
}
