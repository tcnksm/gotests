package main

import (
	"log"
	"os"
)

const EnvDebug = "DEBUG"

func main() {
	cli := &CLI{outStream: os.Stdout, errStream: os.Stderr}
	os.Exit(cli.Run(os.Args))
}

// Debugf
func Debugf(format string, v ...interface{}) {
	if os.Getenv(EnvDebug) != "" {
		log.Printf("[DEBUG] "+format+"\n", v...)
	}
}
