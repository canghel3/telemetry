package log

import (
	"fmt"
	"telemetry/drivers"
	"telemetry/levels"
	"time"
)

type Log struct {
	outputDriver drivers.OutputDriver
	metadata     map[any]any
	buf          []byte
	level        uint8
}

type Tx struct {
	id  string
	log *Log
}

func (l *Log) MakeTx() *Tx {
	return &Tx{
		id:  "", //set id to uuid v4
		log: l,
	}
}

func (tx *Tx) Write(p []byte) {
	tx.log.buf = append(tx.log.buf, p...)
}

func (tx *Tx) End() error {
	err := tx.log.outputDriver.Write(tx.log.buf)
	if err != nil {
		return err

	}

	tx.log.buf = []byte{}
	return nil
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

func (l *Log) Meta(data map[any]any) *Log {
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

func (l *Log) Write(b []byte) {
	l.buf = b
	err := l.outputDriver.Write(formatLogOutput(l))
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
