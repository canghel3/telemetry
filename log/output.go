package log

import (
	"fmt"
	"github.com/Ginger955/telemetry/config"
	"github.com/Ginger955/telemetry/drivers"
	"github.com/Ginger955/telemetry/level"
	"github.com/spf13/viper"
	"io"
	"sync"
)

type Output struct {
	//the lock is used when calling the Settings method
	// in order to return an exact shallow copy of your Output instance.
	lock sync.Mutex

	driver io.Writer
	config config.PkgConfig

	meta map[any]any
}

// Default initiates an Output instance with a stdout driver.
func Default() *Output {
	return &Output{
		driver: drivers.NewStdoutDriver(),
		config: config.PkgConfiguration,
	}
}

// File initiates an Output instance for logging to the specified file.
func File(name string) *Output {
	l := Default()
	l.driver = drivers.NewFileDriver(name)
	return l
}

// Stdout initiates an Output instance for logging to stdout.
func Stdout() *Output {
	d := Default()
	return d
}

// OutputDriver initiates an Output instance for logging to a custom output driver.
func OutputDriver(driver io.Writer) *Output {
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

// Settings overwrites your current Output instance configuration.
// Returns a shallow copy of your Output instance.
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

// WithMetadata sets the metadata for the output driver.
// All messages generated with this driver will contain the given metadata.
func (o *Output) WithMetadata(meta map[any]any) *Output {
	o.meta = meta
	return o
}
