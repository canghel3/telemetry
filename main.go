package main

import (
	"fmt"
	"github.com/Ginger955/telemetry/config"
	"github.com/Ginger955/telemetry/log"
	"github.com/spf13/viper"
	"os"
)

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
		log.Stdout().Error().Log([]byte(fmt.Sprintf("failed to read config: %s", err.Error())))
	}

	err = v.Unmarshal(&config.PkgConfiguration)
	if err != nil {
		log.Stdout().Error().Log([]byte(fmt.Sprintf("failed to unmarshal config: %s", err.Error())))
	}
}

// formatting is very basic, but not vital for proper functionality
// solve TODOs
func main() {
	const logfile = "./telemetry.log"
	const configFile = "./config.json"
	log.Stdout().Error().Log([]byte("HELLO"))

	stdout := log.Stdout()

	stdout.Metadata(map[any]any{"something": "clean"})
	stdout.Info().Log(nil)

	stdout.Error().Log([]byte("HELLO"))
	stdout.Info().Log([]byte("WORLD"))

	toFile := log.File(logfile)
	tx := log.BeginTx()
	tx.Append(toFile.Error().Msg([]byte("something is going on")))
	tx.Append(toFile.Info().Msg([]byte("marcele, la covrigarie!")))
	tx.Append(log.Stdout().Msg([]byte("TO STDOUT!")))
	tx.Log()

	stdoutWithSettings := log.Stdout()
	stdoutWithSettings.Settings("./overwriter.json")
	stdoutWithSettings.Log([]byte("inghetata de fistic"))

	os.WriteFile(logfile, nil, 0644)
}
