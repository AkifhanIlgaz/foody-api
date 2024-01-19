package controllers

import (
	"errors"
	"log"
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
	var credentials models.AuthCredentials

	if err := ctx.ShouldBindJSON(&credentials); err != nil {
		log.Println(err)
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
		log.Println(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "Something went wrong",
		})
		return
	}

	session, err := controller.sessionService.Create(user.Id)
	if err != nil {
		log.Println(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "Something went wrong",
		})
		return
	}

	setCookie(ctx, cookieSession, session.Token)

	ctx.JSON(http.StatusOK, nil)
}

func (controller *AuthController) SignIn(ctx *gin.Context) {
	var credentials models.AuthCredentials

	if err := ctx.ShouldBindJSON(&credentials); err != nil {
		log.Println(err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Missing required fields",
		})
		return
	}

	user, err := controller.userService.Authenticate(credentials.Email, credentials.Password)
	if err != nil {
		log.Println(err)

		if errors.Is(err, services.ErrUserNotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"status":  "fail",
				"message": "User with this email doesn't exist",
			})
			return
		}

		if errors.Is(err, services.ErrWrongPassword) {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"status":  "fail",
				"message": "Wrong password",
			})
			return
		}

		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "Something went wrong",
		})
		return
	}

	session, err := controller.sessionService.Create(user.Id)
	if err != nil {
		log.Println(err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "Something went wrong",
		})
		return
	}

	setCookie(ctx, cookieSession, session.Token)

	ctx.JSON(http.StatusOK, nil)
}

func (controller *AuthController) SignOut(ctx *gin.Context) {
	sessionToken := ctx.GetHeader("Session")
	if sessionToken == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "Session header is missing",
		})
		return
	}

	err := controller.sessionService.Delete(sessionToken)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "Something went wrong",
		})
		return
	}

	deleteCookie(ctx, cookieSession)

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Signed out successfully",
	})
}
