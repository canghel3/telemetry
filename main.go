package main

import (
	"github.com/Ginger955/telemetry/config"
	"github.com/Ginger955/telemetry/log"
	"github.com/spf13/viper"
	"os"
	"time"
)

var Stdout = log.Stdout()

type CustomDriver struct {
	msg string
}

func (cd *CustomDriver) Write(p []byte) error {
	cd.msg = string(p)
	return nil
}

func init() {
	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvPrefix("TELEMETRY")
	v.AddConfigPath("./")
	v.SetConfigName("telemetry")
	v.SetConfigType("json")
	err := v.ReadInConfig()
	if err != nil {
		log.Stdout().Error().Logf("failed to read config: %s", err.Error())
	}

	err = v.Unmarshal(&config.PkgConfiguration)
	if err != nil {
		log.Stdout().Error().Logf("failed to unmarshal config: %s", err.Error())
	}
}

// formatting is very basic, but not vital for proper functionality
// solve TODOs
func main() {
	const logfile = "./telemetry.log"
	const configFile = "./config.json"
	log.Stdout().Error().Log("HELLO")

	std := log.Stdout()
	time.Sleep(time.Millisecond * 50)
	o := std.Settings("./telemetry.json")
	std = log.File("./xyz.log")

	o.Info().Log("t1")
	std.Warn().Log("t2")

	//settings are not persistent unless the output is assigned to a variable
	std.Settings("some config file")

	s := std.Settings("another config file.json")
	s.Info().Log("persistent custom settings")

	toFile := log.File(logfile)
	tx := log.BeginTx()
	tx.Append(toFile.Error().Msg("something is going on"))
	tx.Append(toFile.Info().Msg("marcele, la covrigarie!"))
	tx.Append(log.Stdout().Warn().Msg("TO STDOUT!"))
	tx.Log()

	os.WriteFile(logfile, nil, 0644)

	s = log.Stdout().WithMetadata(map[any]any{"5": 6})
	s.Info().Metadata(map[any]any{"1": 2}).Log("good mornin'")

	k := log.Stdout().WithMetadata(map[any]any{"2": 3})
	k.Info().Log("good night")
}
