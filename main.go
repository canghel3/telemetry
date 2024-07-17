package main

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"telemetry/config"
	"telemetry/log"
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

// TODO: what remains to be done is configurable output (timestamp, field order, field label visiblity, disable implicit formatting)
// finish tests
// write examples
// solve TODOs
func main() {
	const logfile = "./telemetry.log"
	const configFile = "./config.json"
	log.Stdout().Error().Log([]byte("HELLO"))

	stdout := log.Stdout()

	stdout.Error().Log([]byte("HELLO"))
	stdout.Info().Log([]byte("WORLD"))

	toFile := log.File(logfile)
	tx := log.BeginTx()
	tx.Append(toFile.Error().Msg([]byte("hallelujah")))
	tx.Append(toFile.Info().Msg([]byte("marcele, la covrigarie!")))
	tx.Append(log.Stdout().Msg([]byte("TO STDOUT!")))
	tx.Log()

	stdoutWithSettings := log.Stdout()
	stdoutWithSettings.Settings("./overwriter.json")
	stdoutWithSettings.Log([]byte("salutare flacai"))

	os.WriteFile(logfile, nil, 0644)
}
