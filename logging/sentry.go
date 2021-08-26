package logging

import (
	"errors"
	"fmt"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/ometcenter/keeper/version"
)

type SentryLog struct{}

func NewSentryLog(urlDNS string) (*SentryLog, error) {
	StandartLogStruct := &SentryLog{}

	err := sentry.Init(sentry.ClientOptions{
		Dsn:     urlDNS,
		Release: version.Commit + " - " + version.BuildTime,
		Debug:   true,
	})

	fmt.Println(urlDNS)

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return StandartLogStruct, nil
}

func (SentryLog SentryLog) Error(args ...interface{}) {

	StringError := fmt.Sprint(args...)
	// TODO: Создавать ошибку из текста как-то много лишнего.
	sentry.CaptureException(errors.New(StringError))

	print(getPrefix("[ERR]"), StringError)

}

func (SentryLog SentryLog) Info(args ...interface{}) {

	if ok := should(InfoLevel); ok {
		s := fmt.Sprint(args...)
		print(getPrefix("[INF]"), s)
	}
}

func (SentryLog SentryLog) Infof(format string, args ...interface{}) {
	if ok := should(InfoLevel); ok {
		s := fmt.Sprintf(format, args...)
		print(getPrefix("[INF]"), s)
	}
}

func (SentryLog SentryLog) Errorf(format string, args ...interface{}) {

	StringError := fmt.Sprintf(format, args...)
	// TODO: Создавать ошибку из текста как-то много лишнего.
	sentry.CaptureException(errors.New(StringError))
	print(getPrefix("[ERR]"), StringError)

}

func (SentryLog SentryLog) Debugf(format string, args ...interface{}) {
	if ok := should(DebugLevel); ok {
		s := fmt.Sprintf(format, args...)
		print(getPrefix("[INF]"), s)
	}
}

func (SentryLog SentryLog) Warningf(format string, args ...interface{}) {
	if ok := should(WarnLevel); ok {
		s := fmt.Sprintf(format, args...)
		print(getPrefix("[WAR]"), s)
	}
}

func (SentryLog SentryLog) Output(calldepth int, s string) error {

	if calldepth == 4 {
		sentry.CaptureException(errors.New(s))
		print(getPrefix("[ERR]"), s)
	} else {

		if ok := should(InfoLevel); ok {
			//s := fmt.Sprint(args...)
			print(getPrefix(""), s)
		}
	}

	return nil
}

func (SentryLog SentryLog) Panic() {
	sentry.Flush(time.Second * 5)
}
