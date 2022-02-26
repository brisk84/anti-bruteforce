package main

type Config struct {
	Logger LoggerConf
	Server ServerConf
	App    AppConf
}

type LoggerConf struct {
	Level string
	Path  string
}

type ServerConf struct {
	Host string
	Port string
}

type AppConf struct {
	LoginLimit int `mapstructure:"login_limit"`
	PassLimit  int `mapstructure:"pass_limit"`
	IPLimit    int `mapstructure:"ip_limit"`
}

func NewConfig() Config {
	return Config{}
}
