package sys

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/natefinch/lumberjack"
	"github.com/pelletier/go-toml/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	DEV = "dev"
	STG = "stg"
	PRD = "prd"
)

var (
	dir    string
	env    string
	base   string
	config *viper.Viper
)

var (
	path    string
	file    string
	maxNum  int
	maxAge  int
	maxSize int
	logger  *logrus.Logger
)

func init() {
	flag.StringVar(&dir, "dir", ".", "config path")
	flag.StringVar(&env, "env", "dev", "the env(dev/stg/prd)")
	if !flag.Parsed() {
		if env == DEV {
			testing.Init()
		}
		flag.Parse()
	}

	format := "toml"
	config = viper.New()
	config.AddConfigPath(dir)

	if output, err := exec.Command("go", "env", "GOMOD").Output(); err == nil {
		base = filepath.Dir(string(output))
		config.AddConfigPath(base)
	}

	cfn := env + "." + format
	ncd := NewConfigDefault()
	if _, err := os.Stat(cfn); !(err == nil || os.IsExist(err)) {
		if cf, err := os.Create(cfn); err != nil {
			log.Println(err)
		} else if out, err := toml.Marshal(ncd); err != nil {
			log.Println(err)
		} else {
			cf.Write(out)
		}
	}

	config.SetConfigName(cfn)
	config.SetConfigType(format)
	err := config.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func init() {
	path = String("log.path")
	file = String("log.file")
	maxNum = Int("log.max_num")
	maxAge = Int("log.max_age")
	maxSize = Int("log.max_size")

	if path == "" {
		path = "./log"
	}
	if file == "" {
		file = "out.log"
	}
	if maxNum == 0 {
		maxNum = 10
	}
	if maxAge == 0 {
		maxAge = 30
	}
	if maxSize == 0 {
		maxSize = 100
	}

	if _, err := os.Stat(path); !(err == nil || os.IsExist(err)) {
		os.MkdirAll(path, os.ModePerm)
	}

	path := filepath.Join(path, file)

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

func Object(k string) interface{} {
	return config.Get(k)
}

func DefaultObject(k string, def interface{}) interface{} {
	res := Object(k)
	if res == nil && def != nil {
		res = def
	}
	return res
}

func String(k string) string {
	return config.GetString(k)
}

func Int(k string) int {
	return config.GetInt(k)
}

func Float(k string) float64 {
	return config.GetFloat64(k)
}

func Bool(k string) bool {
	return config.GetBool(k)
}

func StringSlice(k string) []string {
	return config.GetStringSlice(k)
}

func StringMap(k string) map[string]interface{} {
	return config.GetStringMap(k)
}

func StringMapString(k string) map[string]string {
	return config.GetStringMapString(k)
}

func StringMapStringSlice(k string) map[string][]string {
	return config.GetStringMapStringSlice(k)
}

func IntSlice(k string) []int {
	return config.GetIntSlice(k)
}

func Env() string {
	return env
}

func Config() *viper.Viper {
	return config
}

func Logger() *logrus.Logger {
	return logger
}
