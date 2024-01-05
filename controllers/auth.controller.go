package controllers

import (
	"github.com/AkifhanIlgaz/foody-api/services"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	userService *services.UserService
}

func NewAuthController(userService *services.UserService) *AuthController {
	return &AuthController{
		userService: userService,
	}
}

func (controller *AuthController) SignIn(ctx *gin.Context) {
}

func (controller *AuthController) SignOut(ctx *gin.Context) {
}

func (controller *AuthController) SignUp(ctx *gin.Context) {
	
}
