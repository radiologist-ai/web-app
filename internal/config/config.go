package config

import "strconv"

type Config struct {
	Server   Server
	Database Database
	GRPC     GRPC
	Minio    Minio
}

type (
	Server struct {
		ListenAddr string `env:"listen_addr"`
		Port       int    `env:"port"`
		Secret     string `env:"secret"`
	}
	Database struct {
		Host     string `env:"db_host"`
		Port     int    `env:"db_port"`
		Username string `env:"db_username"`
		Password string `env:"db_password"`
		Database string `env:"db_database"`
	}
	GRPC struct {
		host string `env:"grpc_host"`
		port int    `env:"grpc_port"`
	}
	Minio struct {
		ServerURL  string `env:"minio_url"`
		AccessKey  string `env:"minio_access"`
		SecretKey  string `env:"minio_secret"`
		BucketName string `env:"minio_bucket"`
	}
)

func (g *GRPC) Addr() string {
	return g.host + ":" + strconv.Itoa(g.port)
}

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
		GRPC: GRPC{
			host: "0.0.0.0",
			port: 8080,
		},
		Minio: Minio{
			ServerURL:  "minio:9000",
			AccessKey:  "radiologist",
			SecretKey:  "password",
			BucketName: "public",
		},
	}
}
