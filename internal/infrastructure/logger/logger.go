package logger

import (
	"github.com/sirupsen/logrus"
	"github.com/stetsd/monk-logger"
)

var Log = monk_logger.NewLogger(logrus.New())
