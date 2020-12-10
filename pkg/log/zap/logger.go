package zap

import (
	"go.uber.org/zap"

	"github.com/hb-chen/docker-sync/pkg/log"
)

var defaultLogger *zap.Logger

func init() {
	zapConf := zap.NewProductionConfig()
	zapEncoderConf := zap.NewProductionEncoderConfig()
	zapConf.EncoderConfig = zapEncoderConf
	l, err := zapConf.Build(zap.AddCallerSkip(2))
	if err != nil {
		log.Fatal(err)
	}
	defaultLogger = l
}

func DefaultLogger() *zap.Logger {
	return defaultLogger
}

type sugaredLogger struct {
	*zap.SugaredLogger
}

func (l *sugaredLogger) Named(name string) log.Logger {
	sl := l.SugaredLogger.Desugar().Named(name).WithOptions(zap.AddCallerSkip(-1)).Sugar()
	return &sugaredLogger{SugaredLogger: sl}
}

func ReplaceLogger(logger *zap.Logger) {
	defaultLogger = logger

	l := &sugaredLogger{SugaredLogger: logger.Sugar()}
	log.SetLogger(l)
}
