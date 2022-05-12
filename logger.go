package loggo

import (
	"fmt"
	"log"
	"time"
)

type (
	ILogger interface {
		Debug(...interface{})
		Debugf(string, ...interface{})

		Info(...interface{})
		Infof(string, ...interface{})

		Warn(...interface{})
		Warnf(string, ...interface{})

		Error(...interface{})
		Errorf(string, ...interface{})
	}

	Logger struct {
		l     *log.Logger
		level LogLevel
	}

	LogLevel int
)

const (
	DebugLevel LogLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	QuietLevel
)

func (level LogLevel) String() string {
	switch level {
	case DebugLevel:
		return cyan(bold("[DEBUG]"))
	case InfoLevel:
		return blue(bold("[INFO]"))
	case WarnLevel:
		return yellow(bold("[WARN]"))
	case ErrorLevel:
		return red(bold("[ERROR]"))
	default:
		return "[UNKNOWN]"
	}
}

func (logger *Logger) lvlAndTime(level LogLevel) string {
	return level.String() + " | " + time.Now().Format(time.StampMicro) + " | "
}

func (logger *Logger) raw(prefix string, entries ...interface{}) {
	entry := prefix

	entry += fmt.Sprintln(entries...)

	logger.l.Print(entry)
}

func (logger *Logger) Info(entries ...interface{}) {
	if logger.level > InfoLevel {
		return
	}

	logger.raw(logger.lvlAndTime(InfoLevel), entries...)
}

func (logger *Logger) Infof(format string, entries ...interface{}) {
	logger.Info(fmt.Sprintf(format, entries...))
}

func (logger *Logger) Debug(entries ...interface{}) {
	if logger.level > DebugLevel {
		return
	}

	logger.raw(logger.lvlAndTime(DebugLevel), entries...)
}

func (logger *Logger) Debugf(format string, entries ...interface{}) {
	logger.Debug(fmt.Sprintf(format, entries...))
}

func (logger *Logger) Warn(entries ...interface{}) {
	if logger.level > WarnLevel {
		return
	}

	logger.raw(logger.lvlAndTime(WarnLevel), entries...)
}

func (logger *Logger) Warnf(format string, entries ...interface{}) {
	logger.Warn(fmt.Sprintf(format, entries...))
}

func (logger *Logger) Error(entries ...interface{}) {
	if logger.level > ErrorLevel {
		return
	}

	logger.raw(logger.lvlAndTime(ErrorLevel), entries...)
}

func (logger *Logger) Errorf(format string, entries ...interface{}) {
	logger.Error(fmt.Sprintf(format, entries...))
}
