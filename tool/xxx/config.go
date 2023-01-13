package xxx

import (
	"flag"
	"github.com/spf13/viper"
	"os/exec"
	"path/filepath"
	"testing"
)

const (
	DEV = "dev"
	STG = "stg"
	PRD = "prd"
)

var (
	env    string
	cfn    string
	base   string
	config *viper.Viper
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
	config.SetConfigType("json")
	err := config.ReadInConfig()
	if err != nil {
		panic(err)
	}
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

func CONFIG() *viper.Viper {
	return config
}
