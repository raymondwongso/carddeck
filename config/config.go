// Package config stores usecases for fetching configuration files
package config

import (
	"github.com/joeshaw/envdecode"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Server   server
	Postgres postgres
}

type server struct {
	Port              string `env:"SERVER_PORT,default=8080"`
	ShutdownTimeout   int    `env:"SERVER_SHUTDOWN_TIMEOUT,default=15"`
	ReadTimeout       int    `env:"SERVER_READ_TIMEOUT,default=5"`
	ReadHeaderTimeout int    `env:"SERVER_READ_HEADER_TIMEOUT,default=5"`
	WriteTimeout      int    `env:"SERVER_WRITE_TIMEOUT,default=5"`
}

type postgres struct {
	Host         string `env:"POSTGRES_HOST,default=localhost"`
	Port         string `env:"POSTGRES_PORT,default=5432"`
	User         string `env:"POSTGRES_USER,default=defaultuser"`
	Pass         string `env:"POSTGRESS_PASS,default=defaultpass"`
	DatabaseName string `env:"POSTGRESS_DBNAME,default=carddeck_dev"`
	// in production, typically we will have more configuration parameters
	// such as max_open_conn, max_conn_idle, etc.
}

// Load returns config object that is populated from env variables and additional
// env provided in filepath.
// filepath is the path to additional env files
// return error if decoding failed, in which it is advised to abort further operation
func Load(filepath ...string) (*Config, error) {
	// load additional env from specified filepath
	if len(filepath) > 0 {
		_ = godotenv.Load(filepath...)
	}

	var config Config
	if err := envdecode.Decode(&config); err != nil {
		log.Error().Err(err).Msg("failed decoding env")
		return nil, err
	}

	return &config, nil
}
