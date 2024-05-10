package logging

import (
	"fmt"
	"log/slog"
	"os"
)

// Lorgus
type SlogLog struct {
	logger *slog.Logger
}

func NewSlog() *SlogLog {

	logHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	})
	logger := slog.New(logHandler)

	logger.Debug("slog test - debug level")
	logger.Info("slog test - info level")
	logger.Warn("slog test - warn level")
	logger.Error("slog test - error level")

	SlogLogStruct := &SlogLog{
		logger: logger,
	}

	return SlogLogStruct

	// logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
	// 	Level: slog.LevelDebug,
	// 	AddSource: true,

	//    }).WithAttrs([]slog.Attr{
	// 	slog.Int("What's the meaning of life?", 42),
	//    })
}

func (s SlogLog) Error(args ...interface{}) {
	s.logger.Error(fmt.Sprint(args...))
	//s.logger.Error("", args...)
}

func (s SlogLog) Errorf(msg string, args ...interface{}) {
	s.logger.Error(msg, args...)
}

func (s SlogLog) Debugf(msg string, args ...interface{}) {
	s.logger.Debug(msg, args...)
}

func (s SlogLog) Warningf(msg string, args ...interface{}) {
	s.logger.Warn(msg, args...)
}

func (s SlogLog) Info(args ...interface{}) {
	s.logger.Info(fmt.Sprint(args...))
	//s.logger.Error("", args...)
}

func (s SlogLog) Infof(msg string, args ...interface{}) {
	s.logger.Info(msg, args...)
}

// func (LorgusLog LorgusLog) Output(calldepth int, s string) error {

// 	if calldepth == 4 {
// 		LorgusLog.Logger.Error(s)
// 	} else {
// 		LorgusLog.Logger.Info(s)
// 	}

// 	return nil
// }

func (s SlogLog) Panic() {

}
