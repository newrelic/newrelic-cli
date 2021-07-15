package logging

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/newrelic/newrelic-cli/internal/config"
	log "github.com/sirupsen/logrus"
)

const (
	// DefaultLogFile is the default log file
	DefaultLogFile = "newrelic-cli.log"
)

var (
	fileHook *LogrusFileHook = nil
)

func LogTrace(args ...interface{}) {
	log.Trace(args)
	if fileHook != nil {
		fileHook.logger.Trace(args)
	}
}

func Warn(args ...interface{}) {
	log.Warn(args)
	if fileHook != nil {
		fileHook.logger.Warn(args)
	}
}

func Error(args ...interface{}) {
	log.Error(args)
	if fileHook != nil {
		fileHook.logger.Error(args)
	}
}

func Fatal(args ...interface{}) {
	log.Fatal(args)
	if fileHook != nil {
		fileHook.logger.Fatal(args)
	}
}

func Tracef(format string, args ...interface{}) {
	log.Tracef(format, args)
	if fileHook != nil {
		fileHook.logger.Tracef(format, args)
	}
}

func Debugf(format string, args ...interface{}) {
	log.Debugf(format, args)
	if fileHook != nil {
		fileHook.logger.Debugf(format, args)
	}
}

func Warnf(format string, args ...interface{}) {
	log.Warnf(format, args)
	if fileHook != nil {
		fileHook.logger.Warnf(format, args)
	}
}

func Errorf(format string, args ...interface{}) {
	log.Errorf(format, args)
	if fileHook != nil {
		fileHook.logger.Errorf(format, args)
	}
}

func Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args)
	if fileHook != nil {
		fileHook.logger.Fatalf(format, args)
	}
}

func Printf(format string, args ...interface{}) {
	log.Printf(format, args)
	if fileHook != nil {
		fileHook.logger.Printf(format, args)
	}
}

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
	return filepath.Join(config.BasePath, DefaultLogFile)
}

func InitFileLogger() {
	if fileHook != nil {
		log.Debug("file logger already configured")
		return
	}

	_, err := os.Stat(config.BasePath)

	if os.IsNotExist(err) {
		errDir := os.MkdirAll(config.BasePath, 0750)
		if errDir != nil {
			log.Warnf("Could not create log file folder: %s", err)
		}
	}

	filepath := GetDefaultLogFilePath()
	fileHook, err := NewLogrusFileHook(filepath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	fileHook.logger.SetLevel(log.DebugLevel)
	fileHook.logger.Hooks.Add(fileHook)
}

type LogrusFileHook struct {
	file      *os.File
	flag      int
	chmod     os.FileMode
	formatter *log.TextFormatter
	logger    *log.Logger
}

func NewLogrusFileHook(file string, flag int, chmod os.FileMode) (*LogrusFileHook, error) {
	plainFormatter := &log.TextFormatter{DisableColors: true}
	logFile, err := os.OpenFile(file, flag, chmod)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to write file on filehook %v", err)
		return nil, err
	}

	return &LogrusFileHook{logFile, flag, chmod, plainFormatter, log.New()}, err
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
