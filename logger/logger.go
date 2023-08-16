package logger

import (
	"os"

	"github.com/op/go-logging"
)

var Log = logging.MustGetLogger("Main")
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

var backend = logging.NewLogBackend(os.Stderr, "", 0)
var backendFormatter = logging.NewBackendFormatter(backend, format)
var backendLeveled = logging.AddModuleLevel(backend)

func Init() {
	backendLeveled.SetLevel(logging.ERROR, "")
	logging.SetBackend(backendFormatter)
}
