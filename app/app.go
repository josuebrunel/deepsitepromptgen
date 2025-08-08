package app

import (
	"log"
	"os"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

const (
	ListenAddr = "127.0.0.1:8888"
)

func envGetOrDflt(key, dflt string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return dflt
}

type App struct {
	ListenAddr string
}

func New() App {
	return App{
		ListenAddr: envGetOrDflt("LISTEN_ADDR", ListenAddr),
	}
}

func (a App) Run() {
	e := echo.New()
	e.Debug = true
	e.Use(middleware.Logger())
	e.Use(middleware.CORS())
	e.Use(middleware.Recover())
	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup: "form:csrf,header:csrf",
	}))
	e.GET("/", HandlerIndex())
	e.POST("/", HandlerIndex())
	log.Fatal(e.Start(a.ListenAddr))

}
