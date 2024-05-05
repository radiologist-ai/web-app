package config

type Config struct {
	Server   Server
	Database Database
}

type (
	Server struct {
		ListenAddr string `yaml:"listen_addr"`
		Port       int    `yaml:"port"`
		Secret     string `yaml:"secret"`
	}
	Database struct {
		Host     string `yaml:"db_host"`
		Port     int    `yaml:"db_port"`
		Username string `yaml:"db_username"`
		Password string `yaml:"db_password"`
		Database string `yaml:"db_database"`
	}
)

func GetConfig() Config {
	return Config{
		Server: Server{
			ListenAddr: "0.0.0.0",
			Port:       5000,
			Secret:     "secret",
		},
		Database: Database{
			Host:     "localhost",
			Port:     5432,
			Username: "postgres",
			Password: "password",
			Database: "aidb",
		},
	}
}
