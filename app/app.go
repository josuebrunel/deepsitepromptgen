package app

import (
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

const (
	ListenAddr = ":8080"
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
	e.Use(middleware.RequestID())
	e.Use(middleware.RequestLogger())
	e.Use(middleware.CORSWithConfig(
		middleware.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
			AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		},
	))
	e.Use(middleware.Recover())
	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup: "form:csrf,header:csrf",
	}))
	e.GET("/", HandlerIndex())
	e.POST("/", HandlerIndex())
	log.Fatal(e.Start(a.ListenAddr))

}
