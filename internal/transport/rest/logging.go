package rest

import "github.com/sirupsen/logrus"

func logFields(handler, problemMsg string) logrus.Fields {
	return logrus.Fields{
		"handler": handler,
		"problem": problemMsg,
	}
}

func logError(handler, problemMsg string, err error) {
	logrus.WithFields(logFields(handler, problemMsg)).Error(err)
}
