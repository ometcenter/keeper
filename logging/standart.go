package logging

import (
	"fmt"
	"log"
)

type StandartLog struct{}

func NewStandartLog() *StandartLog {
	StandartLogStruct := &StandartLog{}

	return StandartLogStruct
}

func (StandartLog StandartLog) Error(args ...interface{}) {
	if ok := should(DebugLevel); ok {
		s := fmt.Sprint(args...)
		print(getPrefix("[ERR]"), s)
	}
}

func (StandartLog StandartLog) Errorf(format string, args ...interface{}) {
	if ok := should(ErrorLevel); ok {
		s := fmt.Sprintf(format, args...)
		print(getPrefix("[ERR]"), s)
	}
}

func (StandartLog StandartLog) Info(args ...interface{}) {
	if ok := should(DebugLevel); ok {
		s := fmt.Sprint(args...)
		print(getPrefix("[INF]"), s)
	}
}

func (StandartLog StandartLog) Debugf(format string, args ...interface{}) {
	if ok := should(DebugLevel); ok {
		s := fmt.Sprintf(format, args...)
		print(getPrefix("[DEB]"), s)
	}
}

func (StandartLog StandartLog) Warningf(format string, args ...interface{}) {
	if ok := should(WarnLevel); ok {
		s := fmt.Sprintf(format, args...)
		print(getPrefix("[WAR]"), s)
	}
}

// TODO: При установке в NSQ logger всегда приходит уровень 2, понять почему
func (StandartLog StandartLog) Output(calldepth int, s string) error {

	if calldepth == 4 {
		log.Print(s)
	} else {
		log.Print(s)
	}

	return nil
}

func (StandartLog StandartLog) Panic() {

}
