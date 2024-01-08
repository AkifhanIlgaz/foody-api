package controllers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/AkifhanIlgaz/foody-api/models"
	"github.com/AkifhanIlgaz/foody-api/services"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	userService    *services.UserService
	sessionService *services.SessionService
}

func NewAuthController(userService *services.UserService, sessionService *services.SessionService) *AuthController {
	return &AuthController{
		userService:    userService,
		sessionService: sessionService,
	}
}

func (controller *AuthController) SignUp(ctx *gin.Context) {
	var credentials models.SignUpCredentials

	if err := ctx.ShouldBindJSON(&credentials); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Missing required fields",
		})
		return
	}

	user, err := controller.userService.Create(credentials.Email, credentials.Password)
	if err != nil {
		if errors.Is(err, services.ErrEmailTaken) {
			ctx.AbortWithStatusJSON(http.StatusConflict, gin.H{
				"status":  "fail",
				"message": "This email address is already associated with an account.",
			})
			return
		}
		fmt.Println(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "Something went wrong",
		})
		return
	}

	session, err := controller.sessionService.Create(user.Id)
	if err != nil {
		fmt.Println(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "Something went wrong",
		})
		return
	}

	fmt.Println(session)
}

func (controller *AuthController) SignIn(ctx *gin.Context) {
}

func (controller *AuthController) SignOut(ctx *gin.Context) {
}
