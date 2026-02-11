package app

import (
	"dppg/internal/db"
	"dppg/internal/db/migrations"
	"dppg/internal/db/repository"
	"dppg/internal/prompt"
	"log"
	"net/http"

	"github.com/josuebrunel/ezauth"
	"github.com/josuebrunel/ezauth/pkg/config"
	"github.com/josuebrunel/gopkg/xlog"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

const (
	ListenAddr = ":8080"
)

type App struct {
	ListenAddr string
}

func New() App {
	return App{
		ListenAddr: ":8080", // Default, will be overridden by config in Run if we pass it, but Run loads config.
		// Let's rely on the config loaded in Run.
	}
}

func (a App) Run() {
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Database connection
	dbConn, err := db.New(cfg.DBDSN)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer dbConn.Close()

	// EzAuth initialization
	authCfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load auth config: %v", err)
	}
	authCfg.Debug = true
	xlog.Info("auth config loaded", "config", authCfg)
	// auth, err := ezauth.NewWithDB(&authCfg, dbConn, "auth")
	auth, err := ezauth.New(&authCfg, "auth")
	if err != nil {
		log.Fatalf("failed to init ezauth: %v", err)
	}

	// Auth Migrations
	if err := auth.Migrate(); err != nil {
		log.Printf("failed to migrate auth tables: %v", err)
	}

	// App Migrations
	if err := migrations.Up(dbConn); err != nil {
		log.Printf("failed to migrate app tables: %v", err)
	}

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

	// EzAuth Routes
	// auth.Handler implements http.Handler
	e.Any("/auth/*", echo.WrapHandler(auth.Handler))

	// Protected Routes
	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup: "form:csrf,header:csrf",
		// Skipper: func(c echo.Context) bool {
		// 	// Skip CSRF for auth routes as EzAuth handles its own CSRF or tokens
		// 	path := c.Path()
		// 	return len(path) >= 5 && path[:5] == "/auth"
		// },
	}))

	// Handlers
	repo := repository.NewRepository(dbConn)
	h := NewHandler(auth, repo)

	// WrapMiddleware fits echo.MiddlewareFunc
	e.GET("/", h.Index()) // Load user if logged in
	e.GET("/signup", h.SignUp())

	// Prompt Routes
	prompt.Mount(e, auth, repo)

	log.Fatal(e.Start(a.ListenAddr))

}
