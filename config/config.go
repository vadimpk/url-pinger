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
	}

	App struct {
		BaseURL     string `env:"QC_APP_BASE_URL"    env-default:"http://localhost:8080/api/v2"`
		CheckoutURL string `env:"QC_CHECKOUT_URL"    env-default:"https://checkout.stag.qcpg.cc"`
		AdminURL    string `env:"QC_ADMIN_URL"       env-default:"https://admin.stag.qcpg.cc"`
		Env         string `env:"QC_APP_ENV"`
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
