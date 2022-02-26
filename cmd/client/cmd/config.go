package cmd

type Config struct {
	Server ServerConf
}

type ServerConf struct {
	Host string
	Port string
}

func NewConfig() Config {
	return Config{}
}
