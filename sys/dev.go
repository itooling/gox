package sys

type ConfigSys struct {
	Port int `toml:"port"`
}

type ConfigLog struct {
	Path    string `toml:"path"`
	File    string `toml:"file"`
	MaxNum  int    `toml:"max_num"`
	MaxAge  int    `toml:"max_age"`
	MaxSize int    `toml:"max_size"`
}

type ConfigToken struct {
	Secret  string `toml:"secret"`
	Expires int    `toml:"expires"`
}

type ConfigRdbms struct {
	Kind   string `toml:"kind"`
	Host   string `toml:"host"`
	Port   int    `toml:"port"`
	User   string `toml:"user"`
	Pass   string `toml:"pass"`
	Dbname string `toml:"dbname"`
	Prefix string `toml:"prefix"`
}

type ConfigRedis struct {
	Addr      string `toml:"addr"`
	Pass      string `toml:"pass"`
	Index     int    `toml:"index"`
	Nodes     string `toml:"nodes"`
	NodesPass string `toml:"nodes_pass"`
}

// ConfigDefault default config
type ConfigDefault struct {
	Sys   ConfigSys   `toml:"sys"`
	Log   ConfigLog   `toml:"log"`
	Token ConfigToken `toml:"token"`
	Rdbms ConfigRdbms `toml:"rdbms"`
	Redis ConfigRedis `toml:"redis"`
}

func NewConfigDefault() *ConfigDefault {
	return &ConfigDefault{
		Sys: ConfigSys{
			Port: 8080,
		},
		Log: ConfigLog{
			Path:    "log",
			File:    "out.log",
			MaxNum:  10,
			MaxAge:  30,
			MaxSize: 100,
		},
		Token: ConfigToken{
			Secret:  "secret",
			Expires: 30,
		},
		Rdbms: ConfigRdbms{
			Kind:   "pgsql",
			Host:   "127.0.0.1",
			Port:   5432,
			User:   "postgres",
			Pass:   "123456",
			Dbname: "test",
			Prefix: "xxx_",
		},
		Redis: ConfigRedis{
			Addr:      "127.0.0.1:6379",
			Pass:      "",
			Index:     0,
			Nodes:     "127.0.0.1:30001,127.0.0.1:30002,127.0.0.1:30003",
			NodesPass: "",
		},
	}
}
