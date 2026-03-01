package main

import (
	"Authentication_Service/internal/config"
	"Authentication_Service/internal/config/customValidation"
	"Authentication_Service/internal/middleware"
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type App struct {
	Re  *gin.Engine
	Rd  *redis.Client
	Db  *gorm.DB
	Cfg *config.Config
}

func NewApp(re *gin.Engine, rd *redis.Client, db *gorm.DB, cfg *config.Config) *App {
	return &App{
		Re:  re,
		Rd:  rd,
		Db:  db,
		Cfg: cfg,
	}
}

func (app *App) Run() error {

	customValidation.RegisterCustomValidations()

	// Request ID middleware (must be first to track all requests)
	app.Re.Use(middleware.RequestIDMiddleware())

	app.Re.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "X-Request-ID", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length", "X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Logger middleware with request ID tracking
	app.Re.Use(middleware.RequestLogger())

	// Custom recovery middleware
	app.Re.Use(gin.CustomRecovery(middleware.GlobalPanicHandler()))

	instance := NewInitialInstance(app)
	instance.DefineRouter()

	server := fmt.Sprintf("%s:%s", app.Cfg.Host, app.Cfg.Port)

	errRun := app.Re.Run(server)

	if errRun != nil {
		return errRun
	}

	return nil
}
