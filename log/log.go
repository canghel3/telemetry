package log

import (
	"fmt"
	"github.com/Ginger955/telemetry/config"
	"github.com/Ginger955/telemetry/drivers"
	"github.com/Ginger955/telemetry/level"
	"github.com/spf13/viper"
	"sync"
)

type Output struct {
	lock   sync.Mutex
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

// Settings overwrites the default output instance configuration.
func (o *Output) Settings(file string) *Output {
	var n = new(Output)

	o.lock.Lock()
	n.driver = o.driver
	n.config = o.config
	o.lock.Unlock()

	v := viper.New()
	v.SetConfigFile(file)
	err := v.ReadInConfig()

	if err != nil {
		Stdout().Error().Log(fmt.Sprintf("failed to read config: %s", err.Error()))
		return n
	}

	err = v.Unmarshal(&n.config)
	if err != nil {
		Stdout().Error().Log(fmt.Sprintf("failed to unmarshal config: %s", err.Error()))
		return n
	}

	return n
}
