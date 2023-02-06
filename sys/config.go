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

	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	DEV = "dev"
	STG = "stg"
	PRD = "prd"
)

var ( // config
	env    string
	cfn    string
	typ    string
	dir    string
	base   string
	config *viper.Viper
)

var ( // logger
	path    string
	file    string
	maxNum  int
	maxAge  int
	maxSize int
	logger  *logrus.Logger
)

func init() { // config
	flag.StringVar(&cfn, "cfn", ".env", "config name")
	flag.StringVar(&dir, "dir", ".", "config path")
	flag.StringVar(&env, "env", "dev", "config model")
	flag.StringVar(&typ, "typ", "yaml", "config type")
	if !flag.Parsed() {
		// testing.Init()
		flag.Parse()
	}

	config = viper.New()
	config.AddConfigPath(dir)

	if output, err := exec.Command("go", "env", "GOMOD").Output(); err == nil {
		base = filepath.Dir(string(output))
		config.AddConfigPath(base)
	}

	cpt := cfn + "." + env
	if _, err := os.Stat(cpt); !(err == nil || os.IsExist(err)) {
		if _, err := os.Create(cpt); err != nil {
			log.Println(err)
		}
	}

	config.SetConfigName(cpt)
	config.SetConfigType(typ)
	err := config.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func init() { // logger
	path = String("app.log.path")
	file = String("app.log.file")
	maxNum = Int("app.log.max_num")
	maxAge = Int("app.log.max_age")
	maxSize = Int("app.log.max_size")

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

// config
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
