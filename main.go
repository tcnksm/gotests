package main

import (
	"log"
	"os"
)

func main() {
	cli := &CLI{outStream: os.Stdout, errStream: os.Stderr}
	os.Exit(cli.Run(os.Args))
}

const EnvDebug = "DEBUG"

// Debugf
func Debugf(format string, v ...interface{}) {
	if os.Getenv(EnvDebug) != "" {
		log.Printf("[DEBUG] "+format+"\n", v...)
	}
}
