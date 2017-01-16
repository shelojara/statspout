package log

import (
	"log"
	"os"
)

var (
	Info  *log.Logger
	Error *log.Logger
	Debug *log.Logger
)

func init() {
	Info = log.New(os.Stdout, "INFO: ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
	Error = log.New(os.Stderr, "ERROR: ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
	Debug = log.New(os.Stdout, "DEBUG: ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
}
