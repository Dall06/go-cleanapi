package utils

import (
	"dall06/go-cleanapi/config"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Initialize() error
	Warn(message string, args ...interface{})
	Info(message string, args ...interface{})
	Error(message string, args ...interface{})
}

type logger struct {
	loggers map[zapcore.Level]*zap.SugaredLogger
}

func NewLogger() Logger {
	return &logger{
		loggers: make(map[zapcore.Level]*zap.SugaredLogger),
	}
}

func (l *logger) Initialize() error {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.MessageKey = "message"
	encoderConfig.LevelKey = "severity"

	for _, level := range []zapcore.Level{zap.WarnLevel, zap.InfoLevel, zap.ErrorLevel} {
		logFilePath, err := getLogFilePath(config.Stage, level.String())
		if err != nil {
			return err
		}
	
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
	
		l.loggers[level] = zl.Sugar()
	}
	return nil
}

func (l *logger) Warn(message string, args ...interface{}) {
	l.loggers[zapcore.WarnLevel].Warnw(message, args...)
}

func (l *logger) Info(message string, args ...interface{}) {
	l.loggers[zapcore.InfoLevel].Infow(message, args...)
}

func (l *logger) Error(message string, args ...interface{}) {
	l.loggers[zapcore.ErrorLevel].Errorw(message, args...)
}

func getLogFilePath(stage string, level string) (string, error) {
	dirName := "logs"
	dir := filepath.Join(config.ProyectPath, dirName)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return "", err
	}
	return dir + "/" + stage + "_" + level + ".log", nil
}
