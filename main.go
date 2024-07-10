package main

import "telemetry/log"

func main() {
	log.Stdout().Error().Write([]byte("HELLO"))
}
