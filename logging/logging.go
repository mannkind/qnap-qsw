package logging

import (
	"github.com/go-logr/logr"
	"github.com/go-logr/stdr"
)

func New(verbosity int) logr.Logger {
	stdr.SetVerbosity(verbosity)
	stdLog := stdr.New(nil)

	return stdLog
}
