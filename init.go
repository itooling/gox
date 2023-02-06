package gox

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	DEV = "dev"
	STG = "stg"
	PRD = "prd"
)

var (
	env     string
	cfn     string
	dir     string
	out     string
	base    string
	maxNum  int
	maxAge  int
	maxSize int
	config  *viper.Viper
	logger  *logrus.Logger
)

func init() { //config
	flag.StringVar(&env, "env", "dev", "the env (dev, stg, prd)")
	flag.StringVar(&cfn, "config", "config", "the config file")
	if !flag.Parsed() {
		testing.Init()
		flag.Parse()
	}

	config = viper.New()
	config.AddConfigPath(".")
	config.AddConfigPath("./config")

	if output, err := exec.Command("go", "env", "GOMOD").Output(); err == nil {
		base = filepath.Dir(string(output))
		config.AddConfigPath(base)
	}

	config.SetConfigName(cfn)
	config.SetConfigType("yaml")
	err := config.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func init() { //log
	dir = String("app.log.dir")
	out = String("app.log.out")
	maxNum = Int("app.log.max_num")
	maxAge = Int("app.log.max_age")
	maxSize = Int("app.log.max_size")

	if dir == "" {
		dir = "./log"
	}
	if out == "" {
		out = "out.log"
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

	if _, err := os.Stat(dir); !(err == nil || os.IsExist(err)) {
		os.MkdirAll(dir, os.ModePerm)
	}

	path := filepath.Join(dir, out)

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
	return config.Get(env + "." + k)
}

func DefaultObject(k string, def interface{}) interface{} {
	res := Object(k)
	if res == nil && def != nil {
		res = def
	}
	return res
}

func String(k string) string {
	return config.GetString(env + "." + k)
}

func Int(k string) int {
	return config.GetInt(env + "." + k)
}

func Float(k string) float64 {
	return config.GetFloat64(env + "." + k)
}

func Bool(k string) bool {
	return config.GetBool(env + "." + k)
}

func StringSlice(k string) []string {
	return config.GetStringSlice(env + "." + k)
}

func StringMap(k string) map[string]interface{} {
	return config.GetStringMap(env + "." + k)
}

func StringMapString(k string) map[string]string {
	return config.GetStringMapString(env + "." + k)
}

func StringMapStringSlice(k string) map[string][]string {
	return config.GetStringMapStringSlice(env + "." + k)
}

func IntSlice(k string) []int {
	return config.GetIntSlice(env + "." + k)
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
