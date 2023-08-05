package util

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"strings"
)

func InitLogger(logLevel logrus.Level, logFilePath string) {
	// set log level
	logrus.SetLevel(logLevel)
	// set log output
	var output io.Writer
	if logFilePath != "" {
		f, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY, 0777)
		if err != nil {
			panic(err)
		}
		output = io.MultiWriter(os.Stdout, f)
	} else {
		output = os.Stdout
	}
	logrus.SetOutput(output)
	// set log formatter
	logrus.SetFormatter(NewPrettierFormatter())
}

type PrettierFormatter struct {
	logrus.TextFormatter
}

func NewPrettierFormatter() *PrettierFormatter {
	formatter := &PrettierFormatter{}
	formatter.TimestampFormat = "2006-01-02 15:04:05"
	formatter.ForceColors = true
	formatter.FullTimestamp = true
	formatter.DisableLevelTruncation = true
	return formatter
}

func (f *PrettierFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// this whole mess of dealing with ansi color codes is required if you want the colored output otherwise you will lose colors in the log levels
	var levelColor int
	switch entry.Level {
	case logrus.DebugLevel, logrus.TraceLevel:
		levelColor = 31 // gray
	case logrus.InfoLevel:
		return []byte(entry.Message), nil
	case logrus.WarnLevel:
		levelColor = 33 // yellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		levelColor = 31 // red
	default:
		levelColor = 36 // blue
	}
	return []byte(fmt.Sprintf("[%s] \x1b[%dm%s\x1b[0m - %s\n", entry.Time.Format(f.TimestampFormat), levelColor, strings.ToUpper(entry.Level.String()), entry.Message)), nil
}