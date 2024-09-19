package utils

import (
	"github.com/sirupsen/logrus"
	"os"
)

var Log = logrus.New()

func InitLogger() {
	Log.Out = os.Stdout
	Log.SetFormatter(&logrus.JSONFormatter{})
	Log.SetLevel(logrus.InfoLevel)
}
