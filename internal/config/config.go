package config

type Config struct {
	Server Server
}

type Server struct {
	ListenAddr string `yaml:"listen_addr"`
	Port       int    `yaml:"port"`
	Secret     string `yaml:"secret"`
}

func GetConfig() Config {
	return Config{
		Server: Server{
			ListenAddr: "0.0.0.0",
			Port:       80,
			Secret:     "secret",
		},
	}
}
