package log

import (
	"fmt"
	"sync"
	"telemetry/drivers"
	"telemetry/levels"
	"time"
)

//TODO: add global logger so that we do not have to specify the output driver every time?

var logger = Default()

type Log struct {
	buf          []byte
	levelLock    *sync.Mutex
	level        uint8 //TODO: use an interface so that we can easily add other levels?
	outputDriver drivers.OutputDriver
	metadata     map[any]any
}

func Default() *Log {
	return &Log{
		outputDriver: drivers.ToStdout(),
		levelLock:    &sync.Mutex{},
		metadata:     nil,
		buf:          nil,
		level:        levels.NoLevel,
	}
}

func (l *Log) NoLevel() *Log {
	l.level = levels.NoLevel
	return l
}

func (l *Log) Error() *Log {
	l.level = levels.LevelError
	return l
}

func (l *Log) Warn() *Log {
	l.level = levels.LevelWarn
	return l
}

func (l *Log) Info() *Log {
	l.level = levels.LevelInfo
	return l
}

func (l *Log) Debug() *Log {
	l.level = levels.LevelDebug
	return l
}

func (l *Log) CustomLevel() {

}

func (l *Log) Metadata(data map[any]any) *Log {
	l.metadata = data
	return l
}

func File(name string) *Log {
	return &Log{
		outputDriver: drivers.ToFileWithName(name),
	}
}

func Stdout() *Log {
	return &Log{
		outputDriver: drivers.ToStdout(),
	}
}

func Custom(driver drivers.OutputDriver) *Log {
	return &Log{
		outputDriver: driver,
	}
}

// Write sends the current log buffer to the output driver for further handling.
// The buffer is emptied and can be reused.
func (l *Log) Write(b []byte) {
	l.buf = b
	err := l.outputDriver.Write(formatLogOutput(l))
	if err != nil {
		return
	}

	l.buf = []byte{}
}

func formatLogOutput(l *Log) []byte {
	//TIMESTAMP LEVEL METDATA BUFFER
	timestamp := time.Now()

	var out = make([]byte, 0)
	out = append(out, []byte(timestamp.String())...)
	out = append(out, byte(' '))
	out = append(out, []byte(levels.LevelToText[l.level])...)
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
