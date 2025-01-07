package loggers

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func InitLogger() {
	Log = logrus.New()

	// Set log level
	Log.SetLevel(logrus.InfoLevel)

	// Set log output (e.g., file or standard output)
	Log.SetOutput(os.Stdout)

	// Set log format (e.g., JSON format)
	Log.SetFormatter(&logrus.JSONFormatter{})
}

func GetLogger() *logrus.Logger {
	if Log == nil {
		InitLogger()
	}
	return Log
}
