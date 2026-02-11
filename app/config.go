package app

import (
	"github.com/josuebrunel/gopkg/xenv"
	"github.com/josuebrunel/gopkg/xlog"
)

type Config struct {
	ListenAddr string `env:"LISTEN_ADDR" default:":8080"`
	BaseURL    string `env:"BASE_URL" default:"http://localhost:8080"`
	DBDSN      string `env:"DB_DSN" required:"true"`
}

func LoadConfig() (*Config, error) {
	var cfg Config
	if err := xenv.Load(&cfg); err != nil {
		xlog.Error("failed to load config", "error", err)
		return nil, err
	}
	return &cfg, nil
}
