package log

import (
	"fmt"
	"github.com/Ginger955/telemetry/config"
	"github.com/Ginger955/telemetry/drivers"
	"github.com/Ginger955/telemetry/level"
	"github.com/spf13/viper"
	"os"
	"sync"
	"time"
)

type Logger struct {
	buf          []byte
	lock         sync.Mutex
	level        level.Level
	outputDriver drivers.OutputDriver
	metadata     map[any]any
	config       config.PkgConfig
}

// Settings enables overwriting the logger instance configuration.
func (l *Logger) Settings(file string) *Logger {
	v := viper.New()
	v.SetConfigFile(file)
	err := v.ReadInConfig()

	if err != nil {
		Stdout().Error().Log([]byte(fmt.Sprintf("failed to read config: %s", err.Error())))
	}

	err = v.Unmarshal(&l.config)
	if err != nil {
		Stdout().Error().Log([]byte(fmt.Sprintf("failed to unmarshal config: %s", err.Error())))
	}

	return l
}

// Default initiates a Logger instance with a stdout driver and no log level.
func Default() *Logger {
	var cpy = config.PkgConfiguration
	return &Logger{
		outputDriver: drivers.ToStdout(),
		lock:         sync.Mutex{},
		metadata:     nil,
		buf:          nil,
		level:        level.None(),
		config:       cpy,
	}
}

// NoLevel empties the Logger level type.
func (l *Logger) NoLevel() *Logger {
	l.level = level.None()
	return l
}

// Error sets the Logger level type to ERROR.
func (l *Logger) Error() *Logger {
	l.level = level.Error()
	return l
}

// Warn sets the Logger level type to WARN.
func (l *Logger) Warn() *Logger {
	l.level = level.Warn()
	return l
}

// Info sets the Logger level type to INFO.
func (l *Logger) Info() *Logger {
	l.level = level.Info()
	return l
}

// Debug sets the Logger level type to DEBUG.
func (l *Logger) Debug() *Logger {
	l.level = level.Debug()
	return l
}

// Level sets the Logger level type to the specified level.
// Use level.Custom for generating a custom Logger level.
func (l *Logger) Level(lvl level.Level) *Logger {
	l.level = lvl
	return l
}

func (l *Logger) Metadata(data map[any]any) *Logger {
	l.metadata = data
	return l
}

// File initiates a Logger instance for logging to the specified file.
func File(name string) *Logger {
	l := Default()
	l.outputDriver = drivers.ToFileWithName(name)
	return l
}

// Stdout initiates a Logger instance for logging to stdout.
func Stdout() *Logger {
	l := Default()
	l.outputDriver = drivers.ToStdout()
	return l
}

// OutputDriver initiates a Logger instance for logging to a custom output driver.
func OutputDriver(driver drivers.OutputDriver) *Logger {
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

// Log sends the current log buffer and the received "b" to the output driver for further handling.
// The buffer is emptied and can be reused.
// If an error occurs during writing, it is logged to os.Stderr.
func (l *Logger) Log(b []byte) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.buf = append(l.buf, b...)
	var output = b
	if !l.config.Formatting.LogConfig.FormattingDisabled {
		output = formatLogOutput(l)
	}

	_, err := l.outputDriver.Write(append(output, '\n'))
	if err != nil {
		//write the error encountered during writing to os.Stderr
		//we could write to the log output driver because it implements the io.Writer,
		//but if the output driver is fatally broken, the writing failure will be lost as well
		//and debugging becomes more difficult.
		fmt.Fprintf(os.Stderr, "failed to write log %s: %s\n", output, err.Error())
	}

	l.buf = []byte{}
}

// TODO: eventually implement field ordering from config
func formatLogOutput(l *Logger) []byte {
	//TIMESTAMP LEVEL METADATA BUFFER
	var timestamp string
	if len(l.config.Formatting.LogConfig.Timestamp) > 0 {
		timestamp = time.Now().Format(l.config.Formatting.LogConfig.Timestamp)
	} else {
		timestamp = time.Now().Format("2006-01-02 15:04:05")
	}

	var out = make([]byte, 0)
	out = append(out, []byte(timestamp)...)
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

	return out
}
