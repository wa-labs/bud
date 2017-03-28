package service

import (
	"os"
	"strings"

	"github.com/go-kit/kit/log"
)

// NewLogger ...
func NewLogger(env string) Logger {
	debugLogger := log.NewNopLogger()
	prodLogger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))

	if strings.Compare("Debug", env) == 0 {
		debugLogger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	}
	return &logger{
		debug: debugLogger,
		prod:  prodLogger,
	}
}

// Logger ...
type Logger interface {
	Debug(keyvals ...interface{}) error
	Prod(keyvals ...interface{}) error
	With(keyvals ...interface{})
}

type logger struct {
	debug log.Logger
	prod  log.Logger
}

// Debug can be used in development to print out messages but they won't be
// printed in production
func (l *logger) Debug(keyvals ...interface{}) error {
	return l.debug.Log(keyvals...)
}

// Prod is used to print messages that will be logged in production
func (l *logger) Prod(keyvals ...interface{}) error {
	return l.prod.Log(keyvals...)
}

// With ...
func (l *logger) With(keyvals ...interface{}) {
	l.debug = log.With(l.debug, keyvals...)
	l.prod = log.With(l.prod, keyvals...)
}

// DefaultCaller is a Valuer that returns the file and line where the Log
// method was invoked. It can only be used with Logger.With.
var DefaultCaller = log.Caller(4)
