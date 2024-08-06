package log

import (
	"fmt"
	"github.com/Ginger955/telemetry/config"
	"github.com/Ginger955/telemetry/drivers"
	"github.com/Ginger955/telemetry/level"
	"github.com/spf13/viper"
)

//TODO: refactor package such that it is concurrency safe.
// potential solution is to restructure code as such:
// create LogDriver struct that can be created by File, Stdout and/or OutputDriver
// the level setting (and other) methods have a LogDriver receiver and return a new Message struct
// the Message struct methods are simply Msg and Log
// structuring the code in this way allows me (the creator) to impose the precise flow i want for the operations
// and ensures (hopefully) that the package is concurrency safe

type Output struct {
	driver drivers.OutputDriver
	config config.PkgConfig
}

// Default initiates an Output instance with a stdout driver.
func Default() *Output {
	return &Output{
		driver: drivers.ToStdout(),
		config: config.PkgConfiguration,
	}
}

// File initiates an Output instance for logging to the specified file.
func File(name string) *Output {
	l := Default()
	l.driver = drivers.ToFileWithName(name)
	return l
}

// Stdout initiates an Output instance for logging to stdout.
func Stdout() *Output {
	d := Default()
	d.driver = drivers.ToStdout()
	return d
}

// OutputDriver initiates an Output instance for logging to a custom output driver.
func OutputDriver(driver drivers.OutputDriver) *Output {
	l := Default()
	l.driver = driver
	return l
}

func (o *Output) Info() *Message {
	return newMessage(o, level.Info())
}

func (o *Output) Error() *Message {
	return newMessage(o, level.Error())
}

func (o *Output) Warn() *Message {
	return newMessage(o, level.Warn())
}

func (o *Output) Debug() *Message {
	return newMessage(o, level.Debug())
}

func (o *Output) Level(custom level.Level) *Message {
	return newMessage(o, custom)
}

// Settings enables overwriting the output instance configuration.
func (o *Output) Settings(file string) *Output {
	v := viper.New()
	v.SetConfigFile(file)
	err := v.ReadInConfig()

	if err != nil {
		Stdout().Error().Log(fmt.Sprintf("failed to read config: %s", err.Error()))
	}

	err = v.Unmarshal(&o.config)
	if err != nil {
		Stdout().Error().Log(fmt.Sprintf("failed to unmarshal config: %s", err.Error()))
	}

	return o
}
