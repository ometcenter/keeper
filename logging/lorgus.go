package logging

import "github.com/sirupsen/logrus"

// Lorgus
type LorgusLog struct {
	Logger *logrus.Logger
}

func NewLorgus() *LorgusLog {
	LorgusLogStruct := &LorgusLog{
		Logger: logrus.New(),
	}

	LorgusLogStruct.Logger.Formatter = &logrus.JSONFormatter{}

	return LorgusLogStruct
}

func (LorgusLog LorgusLog) Error(args ...interface{}) {
	LorgusLog.Logger.Error(args...)
}

func (LorgusLog LorgusLog) Errorf(format string, args ...interface{}) {
	LorgusLog.Logger.Errorf(format, args...)
}

func (LorgusLog LorgusLog) Debugf(format string, args ...interface{}) {
	LorgusLog.Logger.Debugf(format, args...)
}

func (LorgusLog LorgusLog) Warningf(format string, args ...interface{}) {
	LorgusLog.Logger.Warnf(format, args...)
}

func (LorgusLog LorgusLog) Info(args ...interface{}) {
	LorgusLog.Logger.Info(args...)
}

func (LorgusLog LorgusLog) Output(calldepth int, s string) error {

	if calldepth == 4 {
		LorgusLog.Logger.Error(s)
	} else {
		LorgusLog.Logger.Info(s)
	}

	return nil
}

func (LorgusLog LorgusLog) Panic() {

}
