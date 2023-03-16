package utils

import (
	"dall06/go-cleanapi/config"
	"io"
	"log"
	"os"
	"regexp"
	"strings"

	"golang.org/x/exp/slices"
)

type LoggerRepository interface {
	Info(message interface{}, args ...interface{}) error
	Warn(message interface{}, args ...interface{}) error
	Error(message interface{}, args ...interface{}) error
}

var _ LoggerRepository = (*Logger)(nil)

type Logger struct {
	stage string
}

func NewLogger(stage string) LoggerRepository {
	stages := []string{"dev", "prod", "test"}
	stage = strings.ToLower(stage)

	if !slices.Contains(stages, stage) {
		return nil
	}

	return Logger{
		stage: stage,
	}
}

func (l Logger) log(level string, args ...interface{}) error {
	levels := []string{"error", "warn", "info"}
	level = strings.ToLower(level)

	if !slices.Contains(levels, level) {
		level = "warn"
	}

	projectName := regexp.MustCompile(`^(.*` + config.ProjectDirName + `)`)
    currentWorkDirectory, _ := os.Getwd()
    rootPath := projectName.Find([]byte(currentWorkDirectory))
	file := string(rootPath) + `/log/` +l.stage+ `/` +level+ `.log`

	f, err := os.OpenFile(file ,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}

	defer f.Close()

	log.SetOutput(io.MultiWriter(f, os.Stdout))
	log.SetFlags(log.Ldate | log.Ltime)

	log.Printf("%s", args)
	return nil
}

func (l Logger) Warn(msg interface{}, args ...interface{}) error {
	err := l.log("warn", msg, args)
	if err != nil {
		return err
	}

	return nil
}

func (l Logger) Info(msg interface{}, args ...interface{}) error {
	err := l.log("info", msg, args)
	if err != nil {
		return err
	}
	
	return nil
}

func (l Logger) Error(msg interface{}, args ...interface{}) error {
	err := l.log("error", msg, args)
	if err != nil {
		return err
	}

	return nil
}
