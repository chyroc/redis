package redis

import (
	"log"
	"os"
)

// Log ...
var Log *log.Logger

func init() {
	Log = log.New(os.Stdout, "redis", log.LstdFlags)
}

// SetLogPath ...
func SetLogPath(path string) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	Log.SetOutput(f)
	return nil
}
