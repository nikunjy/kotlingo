package kotlingo

import "log"

type Logger interface {
	Info(format string, vals ...interface{})
}

type EmptyLogger struct {
}

func (el EmptyLogger) Info(format string, vals ...interface{}) {
}

type ConsoleLogger struct {
}

func (cl ConsoleLogger) Info(format string, vals ...interface{}) {
	log.Printf(format, vals...)
}
