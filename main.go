package main

import (
	"io/ioutil"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v2"

	"github.com/hb-chen/docker-sync/config"
	"github.com/hb-chen/docker-sync/pkg/log"
	zapLog "github.com/hb-chen/docker-sync/pkg/log/zap"
	"github.com/hb-chen/docker-sync/sync"
)

const (
	logCallerSkip = 1
)

func initLogger(level string, debug, e bool) error {
	logLevel := zapcore.WarnLevel
	err := logLevel.UnmarshalText([]byte(level))
	if err != nil {
		return err
	}

	stderr, close, err := zap.Open("stderr")
	if err != nil {
		close()
		return err
	}
	writer := stderr

	encoder := getLogEncoder(debug)
	core := zapcore.NewCore(encoder, writer, logLevel)
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(logCallerSkip))
	zapLog.ReplaceLogger(logger)

	return nil
}

func getLogEncoder(debug bool) zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	if debug {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	}

	return zapcore.NewConsoleEncoder(encoderConfig)
}

func loadConfig(conf *config.Config) error {
	file, err := os.Open("config.yml")
	if err != nil {
		return err
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(b, conf)
}

func main() {
	initLogger("DEBUG", true, true)

	conf := config.Config{}
	err := loadConfig(&conf)
	if err != nil {
		log.Fatal(err)
	}

	sync.Run(&conf)
}
