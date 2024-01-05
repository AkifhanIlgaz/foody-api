package routes

import (
	"github.com/AkifhanIlgaz/foody-api/controllers"
	"github.com/gin-gonic/gin"
)

type AuthRouteController struct {
	authController *controllers.AuthController
}

func NewAuthRouteController(authController *controllers.AuthController) *AuthRouteController {
	return &AuthRouteController{
		authController: authController,
	}
}

func (routeController *AuthRouteController) AuthRoute(rg *gin.RouterGroup) {
	router := rg.Group("/auth")

	router.POST("/signup", routeController.authController.SignUp)
	router.POST("/signin", routeController.authController.SignIn)
	router.POST("/signout", routeController.authController.SignOut)
}
