package main

import (
	"github.com/op/go-logging"
	"os"
)

var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} [%{level}]%{color:reset} %{message}`,
)

func InitLogging(verbose bool) {
	logBackend := logging.NewLogBackend(os.Stderr, "", 0)
	logFormatter := logging.NewBackendFormatter(logBackend, format)
	// Set the logging severity
	logLeveled := logging.AddModuleLevel(logFormatter)
	if verbose {
		logLeveled.SetLevel(logging.DEBUG, "")
	} else {
		logLeveled.SetLevel(logging.INFO, "")
	}

	// Set the backends to be used.
	logging.SetBackend(logLeveled)
}
