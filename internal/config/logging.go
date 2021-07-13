package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

const (
	// DefaultLogLevel is the default log level
	DefaultLogLevel = "info"

	// DefaultLogFile is the default log file
	DefaultLogFile = "newrelic-cli.log"
)

var (
	fileHookConfigured = false
)

func InitLogger(logLevel string) {
	l := log.StandardLogger()

	l.SetFormatter(&log.TextFormatter{
		DisableLevelTruncation:    true,
		DisableTimestamp:          true,
		EnvironmentOverrideColors: true,
	})

	switch level := strings.ToUpper(logLevel); level {
	case "TRACE":
		l.SetLevel(log.TraceLevel)
	case "DEBUG":
		l.SetLevel(log.DebugLevel)
	case "WARN":
		l.SetLevel(log.WarnLevel)
	case "ERROR":
		l.SetLevel(log.ErrorLevel)
	default:
		l.SetLevel(log.InfoLevel)
	}
}

func GetDefaultLogFilePath() string {
	return filepath.Join(BasePath, DefaultLogFile)
}

func InitFileLogger() {
	if fileHookConfigured {
		log.Debug("file logger already configured")
		return
	}

	_, err := os.Stat(BasePath)

	if os.IsNotExist(err) {
		errDir := os.MkdirAll(BasePath, 0750)
		if errDir != nil {
			log.Warnf("Could not create log file folder: %s", err)
		}
	}

	fileHook, err := NewLogrusFileHook(BasePath+"/"+DefaultLogFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0640)
	if err == nil && !fileHookConfigured {
		l := log.StandardLogger()
		l.Hooks.Add(fileHook)
		fileHookConfigured = true
	}
}

type LogrusFileHook struct {
	file      *os.File
	flag      int
	chmod     os.FileMode
	formatter *log.TextFormatter
}

func NewLogrusFileHook(file string, flag int, chmod os.FileMode) (*LogrusFileHook, error) {
	plainFormatter := &log.TextFormatter{DisableColors: true}
	logFile, err := os.OpenFile(file, flag, chmod)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to write file on filehook %v", err)
		return nil, err
	}

	return &LogrusFileHook{logFile, flag, chmod, plainFormatter}, err
}

func (hook *LogrusFileHook) Fire(entry *log.Entry) error {
	plainformat, err := hook.formatter.Format(entry)
	if err != nil {
		return err
	}

	line := string(plainformat)
	_, err = hook.file.WriteString(line)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to write file on filehook(entry.String)%v", err)
		return err
	}

	return nil
}

func (hook *LogrusFileHook) Levels() []log.Level {
	return []log.Level{
		log.PanicLevel,
		log.FatalLevel,
		log.ErrorLevel,
		log.WarnLevel,
		log.InfoLevel,
		log.DebugLevel,
	}
}
