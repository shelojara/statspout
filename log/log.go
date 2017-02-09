package log

import (
	"log"
	"os"
)

var (
	Info    *log.Logger
	Error   *log.Logger
	Debug   *log.Logger
	Warning *log.Logger
)

func init() {
	Info    = log.New(os.Stdout, "INFO: ", log.LstdFlags|log.Lmicroseconds)
	Error   = log.New(os.Stderr, "ERROR: ", log.LstdFlags|log.Lmicroseconds)
	Debug   = log.New(os.Stdout, "DEBUG: ", log.LstdFlags|log.Lmicroseconds)
	Warning = log.New(os.Stdout, "WARNING: ", log.LstdFlags|log.Lmicroseconds)
}
