package logging

import (
	"fmt"
	"time"

	"github.com/ometcenter/keeper/config"
)

var (
	level Level
)

const (
	// DebugLevel defines debug log level.
	DebugLevel Level = iota
	// InfoLevel defines info log level.
	InfoLevel
	// WarnLevel defines warn log level.
	WarnLevel
	// ErrorLevel defines error log level.
	ErrorLevel
	// FatalLevel defines fatal log level.
)

// Level defines log levels.
type Level int8

// SetLevel метод для установки логера
func SetLevel(l Level) {
	level = l
}

func print(pref, format string) {
	fmt.Println(pref, format)
}

func should(lvl Level) bool {
	if lvl > level || lvl == level {
		return true
	}
	return false
}

func getPrefix(l string) string {
	return fmt.Sprint(l, " ", time.Now().Format("2006/01/02 - 15:04:05"), " |")
}

type LogImpl interface {
	//InitLog()
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Warningf(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	//Output(calldepth int, s string) error
	Panic() // Выполняется перед panic
}

var Impl LogImpl

func SetLogger(LogImpl LogImpl) {
	Impl = LogImpl
}

func InitLog(ServiceConfig *config.ServiceConfig) {

	switch ServiceConfig.LoggerConfig.Level {
	case 0:
		SetLevel(DebugLevel)
	case 1:
		SetLevel(InfoLevel)
	case 3:
		SetLevel(WarnLevel)
	case 4:
		SetLevel(ErrorLevel)
	default:
		SetLevel(DebugLevel)
	}

	switch ServiceConfig.LoggerConfig.Name {
	case "Sentry":
		Logger, err := NewSentryLog(ServiceConfig.SentryUrlDSN)
		if err != nil {
			Logger := NewStandartLog()
			SetLogger(Logger)
		} else {
			SetLogger(Logger)

		}
	case "Lorgus":
		Logger := NewLorgus()
		SetLogger(Logger)
	default:
		Logger := NewStandartLog()
		SetLogger(Logger)
	}

}
