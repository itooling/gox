package xxx

import (
	"fmt"
	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var (
	dir     string
	maxNum  int
	maxAge  int
	maxSize int
	logger  *logrus.Logger
)

func init() { //log
	dir = String("app.log.dir")
	maxNum = Int("app.log.max_num")
	maxAge = Int("app.log.max_age")
	maxSize = Int("app.log.max_size")

	if _, err := os.Stat(dir); !(err == nil || os.IsExist(err)) {
		os.MkdirAll(dir, os.ModePerm)
	}

	path := filepath.Join(dir, "app.log")

	lf := &lumberjack.Logger{
		Filename:   path,
		MaxSize:    maxSize, //mb
		MaxAge:     maxAge,  //day
		MaxBackups: maxNum,  //number
		Compress:   false,
	}

	mw := io.MultiWriter(lf, os.Stdout)

	logger = logrus.New()
	logger.SetOutput(mw)
	logger.SetReportCaller(true)
	logger.SetFormatter(&LineFormatter{TimestampFormat: "2006/01/02 15:04:05"})
	logger.SetLevel(logrus.InfoLevel)
}

type LineFormatter struct {
	TimestampFormat string //2006/01/02 15:04:05.000
}

func (f *LineFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var file string
	var line int
	if entry.Caller != nil {
		file = filepath.Base(entry.Caller.File)
		line = entry.Caller.Line
	}

	timestamp := entry.Time.Format(f.TimestampFormat)

	content := fmt.Sprintf("%s [%s] [%s:%d] %s\n",
		timestamp, strings.ToUpper(entry.Level.String()), file, line, entry.Message)

	return []byte(content), nil
}

func LOGGER() *logrus.Logger {
	return logger
}
