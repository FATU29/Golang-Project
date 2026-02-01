package main

import (
	"Authentication_Service/internal/config"
	"Authentication_Service/internal/handler"
	"Authentication_Service/pkg/dbFactory"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {
	r := gin.Default()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	db := dbFactory.GetDatabase[*gorm.DB](cfg)
	if db == nil {
		log.Fatal("Database not found")
	}

	_, errCn := db.Connect(dbFactory.GormStrategy{})

	if errCn != nil {
		log.Fatal(err)
	}

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	handler.DefineRouters(r)

	server := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	errRun := r.Run(server)
	if errRun != nil {
		return
	}
}
