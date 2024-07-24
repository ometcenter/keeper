package logging

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/ometcenter/keeper/version"
)

//	{
//		"timestamp" : "2022-11-21 18:09:12.714",
//		"level" : "INFO",
//		"thread" : "main",
//		"logger" : "org.springframework.boot.web.embedded.jetty.JettyWebServer",
//		"message" : "Jetty started on port(s) 8080 (http/1.1) with context path '/'",
//		"application" : "<your-application-name>"
//	  }
type MeshSpecificLog struct {
	Timestamp   string `json:"timestamp"`
	Level       string `json:"level"`
	Thread      string `json:"thread"`
	Logger      string `json:"logger"`
	Message     string `json:"message"`
	Application string `json:"application"`
}

type MeshSpecificSentryLog struct{}

func NewMeshSpecificSentryLog(urlDNS string) (*MeshSpecificSentryLog, error) {
	StandartLogStruct := &MeshSpecificSentryLog{}

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

func (SentryLog MeshSpecificSentryLog) Error(args ...interface{}) {

	StringError := fmt.Sprint(args...)
	// TODO: Создавать ошибку из текста как-то много лишнего.
	sentry.CaptureException(errors.New(StringError))

	MeshSpecificLog := MeshSpecificLog{Timestamp: time.Now().Format("2006-01-02 15:04:05"),
		Level:       "ERR",
		Thread:      "main",
		Logger:      "fmt.Println",
		Message:     StringError,
		Application: "keeper"}

	jsonB, _ := json.Marshal(MeshSpecificLog)
	fmt.Println(string(jsonB))

}

func (SentryLog MeshSpecificSentryLog) Errorf(format string, args ...interface{}) {

	StringError := fmt.Sprintf(format, args...)
	// TODO: Создавать ошибку из текста как-то много лишнего.
	sentry.CaptureException(errors.New(StringError))
	//print(getPrefix("[ERR]"), StringError)
	MeshSpecificLog := MeshSpecificLog{Timestamp: time.Now().Format("2006-01-02 15:04:05"),
		Level:       "ERR",
		Thread:      "main",
		Logger:      "fmt.Println",
		Message:     StringError,
		Application: "keeper"}

	jsonB, _ := json.Marshal(MeshSpecificLog)
	fmt.Println(string(jsonB))

}

func (SentryLog MeshSpecificSentryLog) Info(args ...interface{}) {

	// if ok := should(InfoLevel); ok {
	s := fmt.Sprint(args...)
	// 	print(getPrefix("[INF]"), s)
	// }
	MeshSpecificLog := MeshSpecificLog{Timestamp: time.Now().Format("2006-01-02 15:04:05"),
		Level:       "INF",
		Thread:      "main",
		Logger:      "fmt.Println",
		Message:     s,
		Application: "keeper"}

	jsonB, _ := json.Marshal(MeshSpecificLog)
	fmt.Println(string(jsonB))
}

func (SentryLog MeshSpecificSentryLog) Infof(format string, args ...interface{}) {
	// if ok := should(InfoLevel); ok {
	s := fmt.Sprintf(format, args...)
	// 	print(getPrefix("[INF]"), s)
	// }
	MeshSpecificLog := MeshSpecificLog{Timestamp: time.Now().Format("2006-01-02 15:04:05"),
		Level:       "INF",
		Thread:      "main",
		Logger:      "fmt.Println",
		Message:     s,
		Application: "keeper"}

	jsonB, _ := json.Marshal(MeshSpecificLog)
	fmt.Println(string(jsonB))
}

func (SentryLog MeshSpecificSentryLog) Debugf(format string, args ...interface{}) {
	if ok := should(DebugLevel); ok {
		s := fmt.Sprintf(format, args...)
		print(getPrefix("[INF]"), s)
	}
}

func (SentryLog MeshSpecificSentryLog) Warningf(format string, args ...interface{}) {
	if ok := should(WarnLevel); ok {
		s := fmt.Sprintf(format, args...)
		print(getPrefix("[WAR]"), s)
	}
}

// func (SentryLog SentryLog) Output(calldepth int, s string) error {

// 	if calldepth == 4 {
// 		sentry.CaptureException(errors.New(s))
// 		print(getPrefix("[ERR]"), s)
// 	} else {

// 		if ok := should(InfoLevel); ok {
// 			//s := fmt.Sprint(args...)
// 			print(getPrefix(""), s)
// 		}
// 	}

// 	return nil
// }

func (SentryLog MeshSpecificSentryLog) Panic() {
	sentry.Flush(time.Second * 5)
}
