package config

import (
	"log"
	"sync"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		App  App
		HTTP HTTP
		Log  Log
		Ping PingerService
	}

	App struct {
	}

	HTTP struct {
		Port            string        `env:"HTTP_PORT" env-default:"8088"`
		ReadTimeout     time.Duration `env:"HTTP_READ_TIMEOUT" env-default:"60s"`
		WriteTimeout    time.Duration `env:"HTTP_WRITE_TIMEOUT" env-default:"60s"`
		ShutdownTimeout time.Duration `env:"HTTP_SHUTDOWN_TIMEOUT" env-default:"30s"`
	}

	Log struct {
		Level string `env:"QC_LOG_LEVEL" env-default:"debug"`
	}

	PingerService struct {
		WorkerPoolSize     int           `env:"WORKER_POOL_SIZE" env-default:"10"`
		DefaultPingTimeout time.Duration `env:"DEFAULT_PING_TIMEOUT" env-default:"5s"`
	}
)

var (
	config Config
	once   sync.Once
)

func Get() *Config {
	once.Do(func() {
		err := cleanenv.ReadEnv(&config)
		if err != nil {
			log.Fatal("failed to read env", err)
		}
	})

	return &config
}
