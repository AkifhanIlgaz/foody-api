package controllers

import (
	"errors"
	"log"
	"net/http"

	"github.com/AkifhanIlgaz/foody-api/models"
	"github.com/AkifhanIlgaz/foody-api/services"
	"github.com/gin-gonic/gin"
)

type UserMiddleware struct {
	sessionService *services.SessionService
}

func NewUserMiddleware(sessionService *services.SessionService) *UserMiddleware {
	return &UserMiddleware{
		sessionService: sessionService,
	}
}

func (middleware *UserMiddleware) SetUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sessionToken, err := ctx.Cookie(cookieSession)
		if err != nil {
			log.Println(err)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "session cookie missing"})
			return
		}

		user, err := middleware.sessionService.User(sessionToken)
		if err != nil {
			log.Println(err)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "something went wrong"})
			return
		}

		ctx.Set("currentUser", user)
		ctx.Next()
	}
}

func (middleware *UserMiddleware) RequireUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_, err := getUserFromCtx(ctx)
		if err != nil {
			log.Println(err)
			ctx.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		ctx.Next()
	}
}

func getUserFromCtx(ctx *gin.Context) (*models.User, error) {
	raw, exists := ctx.Get("currentUser")
	if !exists {
		return nil, errors.New("user not logged in")
	}

	user, ok := raw.(*models.User)
	if !ok {
		return nil, errors.New("wrong type for user")
	}

	return user, nil
}
