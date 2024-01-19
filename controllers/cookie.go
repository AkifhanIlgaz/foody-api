package controllers

import (
	"math"

	"github.com/gin-gonic/gin"
)

const cookieSession = "session"

func setCookie(ctx *gin.Context, name, value string) {
	ctx.SetCookie(name, value, math.MaxInt, "/", "localhost:8000", false, true)
}

func deleteCookie(ctx *gin.Context, name string) {
	ctx.SetCookie(name, "", -1, "/", "localhost:8000", false, true)
}
