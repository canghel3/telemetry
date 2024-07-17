package main

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"telemetry/log"
	"telemetry/models"
)

type CustomDriver struct {
	msg string
}

func (cd *CustomDriver) Write(p []byte) error {
	cd.msg = string(p)
	return nil
}

func init() {
	viper.AddConfigPath("./")
	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	if err != nil {
		log.Stdout().Error().Log([]byte(fmt.Sprintf("failed to read config: %s", err.Error())))
	}

	err = viper.Unmarshal(&models.PkgConfig)
	if err != nil {
		log.Stdout().Error().Log([]byte(fmt.Sprintf("failed to unmarshal config: %s", err.Error())))
	}
}

// TODO: what remains to be done is configurable output (timestamp, field order, field label visiblity, disable implicit formatting)
// finish tests
// write examples
// solve TODOs
func main() {
	const LOGFILE = "./telemetry.log"
	log.Stdout().Error().Log([]byte("HELLO"))

	stdout := log.Stdout()

	stdout.Error().Log([]byte("HELLO"))
	stdout.Info().Log([]byte("WORLD"))

	toFile := log.File(LOGFILE)
	tx := log.BeginTx()
	tx.Append(toFile.Error().Msg([]byte("hallelujah")))
	tx.Append(toFile.Info().Msg([]byte("marcele, la covrigarie!")))
	tx.Append(log.Stdout().Msg([]byte("TO STDOUT!")))
	tx.Log()

	os.WriteFile(LOGFILE, nil, 0644)
}
