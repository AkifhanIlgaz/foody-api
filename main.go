package main

import (
	"context"
	"log"
	"net/http"

	"github.com/AkifhanIlgaz/foody-api/cfg"
	"github.com/AkifhanIlgaz/foody-api/controllers"
	"github.com/AkifhanIlgaz/foody-api/database"
	"github.com/AkifhanIlgaz/foody-api/routes"
	"github.com/AkifhanIlgaz/foody-api/services"
	"github.com/AkifhanIlgaz/foody-api/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	config, err := cfg.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not read environment variables", err)
	}

	ctx := context.TODO()

	databases, err := database.Connect(config)
	if err != nil {
		log.Fatal("Could not connect to databases: ", err)
	}

	defer databases.Mongo.Disconnect(ctx)
	defer databases.Redis.Close()

	server := gin.Default()
	utils.SetCors(server)

	userService := services.NewUserService(ctx, databases.Mongo, config)
	sessionService := services.NewSessionService(ctx, databases.Mongo, config)

	authController := controllers.NewAuthController(userService, sessionService)

	authRouteController := routes.NewAuthRouteController(authController)

	router := server.Group("/api")
	router.GET("/health-checker", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "API is healthy"})
	})

	authRouteController.AuthRoute(router)

	log.Fatal(server.Run(":" + config.Port))
}
