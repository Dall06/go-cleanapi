// Package utils is a package that provides general method for the api usage
package utils

import (
	"dall06/go-cleanapi/config"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger refers to the repository as interface of the logger
type Logger interface {
	Initialize() error
	Warn(message string, args ...interface{})
	Info(message string, args ...interface{})
	Error(message string, args ...interface{})
}

var _ Logger = (*logger)(nil)

type logger struct {
	loggers map[zapcore.Level]*zap.SugaredLogger
	config  config.Vars
}

// NewLogger is a function constructor for Logger
func NewLogger(v config.Vars) Logger {
	return logger{
		loggers: make(map[zapcore.Level]*zap.SugaredLogger),
		config:  v,
	}
}

func (l logger) Initialize() error {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.MessageKey = "message"
	encoderConfig.LevelKey = "severity"

	for _, level := range []zapcore.Level{zap.WarnLevel, zap.InfoLevel, zap.ErrorLevel} {
		logFilePath, err := l.getLogFilePath(l.config.Stage, level.String())
		if err != nil {
			return err
		}

		fmt.Printf("logFilePath: %s, level: %s\n", logFilePath, level.String())

		cfg := zap.Config{
			Level:             zap.NewAtomicLevelAt(level),
			Development:       false,
			DisableStacktrace: true,
			Encoding:          "json",
			EncoderConfig:     encoderConfig,
			OutputPaths:       []string{"stdout", logFilePath},
			ErrorOutputPaths:  []string{"stderr"},
		}

		zl, err := cfg.Build()
		if err != nil {
			return err
		}

		logger := zl.Sugar()

		fmt.Printf("logger: %+v\n", logger)

		l.loggers[level] = logger
	}

	if len(l.loggers) < 3 {
		return errors.New("no loggers configured")
	}

	return nil
}

func (l logger) Warn(message string, args ...interface{}) {
	l.loggers[zapcore.WarnLevel].Warnf(message, args...)
}

func (l logger) Info(message string, args ...interface{}) {
	l.loggers[zapcore.InfoLevel].Infof(message, args...)
}

func (l logger) Error(message string, args ...interface{}) {
	l.loggers[zapcore.ErrorLevel].Errorf(message, args...)
}

func (l logger) getLogFilePath(stage string, level string) (string, error) {
	dirName := "logs"
	dir := filepath.Join(l.config.ProyectPath, dirName)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return "", err
	}
	return dir + "/" + stage + "_" + level + ".log", nil
}
