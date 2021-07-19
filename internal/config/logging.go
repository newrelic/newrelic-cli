package config

import (
	"fmt"
	"io/ioutil"
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
	logger             = log.StandardLogger()
)

func InitLogger(logger *log.Logger, logLevel string) {
	logger.SetFormatter(&log.TextFormatter{
		DisableLevelTruncation:    true,
		DisableTimestamp:          true,
		EnvironmentOverrideColors: true,
	})

	level := getLevelFromString(logLevel, log.InfoLevel)

	logger.SetLevel(level)
}

func GetDefaultLogFilePath() string {
	return filepath.Join(BasePath, DefaultLogFile)
}

func InitFileLogger(terminalLogLevel string) {
	logger = log.New()
	InitLogger(logger, terminalLogLevel)

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

	fileLoggerLevel := log.DebugLevel
	if l := os.Getenv("NEW_RELIC_CLI_FILE_LOG_LEVEL"); l != "" {
		fileLoggerLevel = getLevelFromString(l, log.DebugLevel)
	}

	fileHook, err := NewLogrusFileHook(BasePath+"/"+DefaultLogFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0640)
	if err == nil && !fileHookConfigured {
		l := log.StandardLogger()
		l.SetOutput(ioutil.Discard)
		l.SetLevel(fileLoggerLevel)
		l.Hooks.Add(fileHook)
		fileHookConfigured = true
	}
}

type LogrusFileHook struct {
	file      *os.File
	flag      int
	chmod     os.FileMode
	formatter *log.JSONFormatter
}

func NewLogrusFileHook(file string, flag int, chmod os.FileMode) (*LogrusFileHook, error) {
	formatter := &log.JSONFormatter{}
	logFile, err := os.OpenFile(file, flag, chmod)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to write file on filehook %v", err)
		return nil, err
	}

	return &LogrusFileHook{logFile, flag, chmod, formatter}, err
}

func (hook *LogrusFileHook) Fire(entry *log.Entry) error {

	logger.Log(entry.Level, entry.Message)

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
		log.TraceLevel,
	}
}

func getLevelFromString(logLevel string, defaultLevel log.Level) log.Level {
	switch level := strings.ToUpper(logLevel); level {
	case "TRACE":
		return log.TraceLevel
	case "DEBUG":
		return log.DebugLevel
	case "WARN":
		return log.WarnLevel
	case "ERROR":
		return log.ErrorLevel
	default:
		log.Debugf("log level %s could not be parsed, using default level %s", level, defaultLevel)
		return defaultLevel
	}
}
