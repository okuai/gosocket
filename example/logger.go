package example

import (
	"fmt"
	"io"
	"log"
)

type Logger struct {
	*log.Logger
}

func NewLogger(out io.Writer) *Logger {
	return &Logger{log.New(out, "", log.LstdFlags|log.Lshortfile)}
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.Printf(fmt.Sprintf("[DEBUG]: %s ", format), v)
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.Printf(fmt.Sprintf("[Infof]: %s ", format), v)
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	l.Printf(fmt.Sprintf("[Warnf]: %s ", format), v)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.Printf(fmt.Sprintf("[Errorf]: %s ", format), v)
}
