package routes

import (
	"github.com/AkifhanIlgaz/foody-api/controllers"
	"github.com/gin-gonic/gin"
)

type AuthRouteController struct {
	authController *controllers.AuthController
	userMiddleware *controllers.UserMiddleware
}

func NewAuthRouteController(authController *controllers.AuthController, userMiddleware *controllers.UserMiddleware) *AuthRouteController {
	return &AuthRouteController{
		authController: authController,
		userMiddleware: userMiddleware,
	}
}

func (routeController *AuthRouteController) AuthRoute(rg *gin.RouterGroup) {
	router := rg.Group("/auth", routeController.userMiddleware.SetUser())

	router.POST("/signup", routeController.authController.SignUp)
	router.POST("/signin", routeController.authController.SignIn)
	router.POST("/signout", routeController.userMiddleware.RequireUser(), routeController.authController.SignOut)
}
