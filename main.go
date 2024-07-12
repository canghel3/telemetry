package main

import (
	"os"
	"telemetry/log"
)

type CustomDriver struct {
	msg string
}

func (cd *CustomDriver) Write(p []byte) error {
	cd.msg = string(p)
	return nil
}

func main() {
	const LOGFILE = "./telemetry.log"
	log.Stdout().Error().Write([]byte("HELLO"))

	stdout := log.Stdout()

	stdout.Error().Write([]byte("HELLO"))
	stdout.Info().Write([]byte("WORLD"))

	toFile := log.File(LOGFILE)
	tx := toFile.BeginTx()
	tx.Append(toFile.Error().Msg([]byte("hallelujah")))
	tx.Append(toFile.Info().Msg([]byte("marcele, la covrigarie!")))
	tx.Append(log.Stdout().Msg([]byte("TO STDOUT!")))
	tx.Commit()

	os.WriteFile(LOGFILE, nil, 0644)
}
