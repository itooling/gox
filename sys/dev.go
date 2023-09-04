package sys

type ConfigNet struct {
	Port int `yaml:"port"`
}

type ConfigLog struct {
	Path    string `yaml:"path"`
	File    string `yaml:"file"`
	MaxNum  int    `yaml:"max_num"`
	MaxAge  int    `yaml:"max_age"`
	MaxSize int    `yaml:"max_size"`
}

type ConfigDbs struct {
	Rdbms ConfigDbsRdbms `yaml:"rdbms"`
	Redis ConfigDbsRedis `yaml:"redis"`
}

type ConfigDbsRdbms struct {
	Kind   string `yaml:"kind"`
	Host   string `yaml:"host"`
	Port   int    `yaml:"port"`
	User   string `yaml:"user"`
	Pass   string `yaml:"pass"`
	Dbname string `yaml:"dbname"`
	Prefix string `yaml:"prefix"`
}

type ConfigDbsRedis struct {
	Addr      string `yaml:"addr"`
	Pass      string `yaml:"pass"`
	Index     int    `yaml:"index"`
	Nodes     string `yaml:"nodes"`
	NodesPass string `yaml:"nodes_pass"`
}

type ConfigApp struct {
	Net ConfigNet `yaml:"net"`
	Log ConfigLog `yaml:"log"`
	Dbs ConfigDbs `yaml:"dbs"`
}

type ConfigDefault struct {
	App ConfigApp `yaml:"app"`
}

func NewConfigDefault() *ConfigDefault {
	return &ConfigDefault{
		App: ConfigApp{
			Net: ConfigNet{
				Port: 8080,
			},
			Log: ConfigLog{
				Path:    "log",
				File:    "out.log",
				MaxNum:  10,
				MaxAge:  30,
				MaxSize: 100,
			},
			Dbs: ConfigDbs{
				Rdbms: ConfigDbsRdbms{
					Kind:   "pgsql",
					Host:   "127.0.0.1",
					Port:   5432,
					User:   "postgres",
					Pass:   "123456",
					Dbname: "test",
					Prefix: "xxx_",
				},
				Redis: ConfigDbsRedis{
					Addr:      "127.0.0.1:6379",
					Pass:      "",
					Index:     0,
					Nodes:     "127.0.0.1:30001,127.0.0.1:30002,127.0.0.1:30003",
					NodesPass: "",
				},
			},
		},
	}
}
