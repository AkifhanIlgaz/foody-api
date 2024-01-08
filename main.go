package main

import (
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

	databases, err := database.Connect(config)
	if err != nil {
		log.Fatal("Could not connect to databases: ", err)
	}

	defer databases.Postgres.Close()
	defer databases.Redis.Close()

	server := gin.Default()
	utils.SetCors(server)

	userService := services.NewUserService(databases.Postgres)
	sessionService := services.NewSessionService(databases.Postgres)

	authController := controllers.NewAuthController(userService, sessionService)

	authRouteController := routes.NewAuthRouteController(authController)

	router := server.Group("/api")
	router.GET("/health-checker", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "API is healthy"})
	})

	authRouteController.AuthRoute(router)

	log.Fatal(server.Run(":" + config.Port))
}
