// @title           Authentication Service API
// @version         1.0
// @description     API for user registration, login, logout, email validation, OTP, and password flows.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

package main

import (
	"Authentication_Service/internal/config"
	"Authentication_Service/pkg/dbFactory"
	"Authentication_Service/pkg/redis"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {
	// Use gin.New() instead of gin.Default() to avoid duplicate middleware
	// gin.Default() includes Logger and Recovery which we add explicitly in app.Run()
	r := gin.New()

	cfg := config.LoadConfig()

	db := dbFactory.GetDatabase[*gorm.DB](cfg)
	if db == nil {
		log.Fatal("Database not found")
	}

	dbInstance, errCn := db.Connect(dbFactory.GormStrategy{})

	if errCn != nil {
		log.Fatal(errCn)
	}

	redisInstance, errRedis := redis.InitRedis(cfg)
	_ = redisInstance

	if errRedis != nil {
		log.Fatal("Failed to connect to Redis:", errRedis)
	}

	app := NewApp(r, redisInstance, dbInstance, cfg)
	errRunApp := app.Run()

	if errRunApp != nil {
		log.Fatal("Failed to run the app:", errRunApp)
	}
}
