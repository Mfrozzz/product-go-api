package middleware

import (
	"net/http"
	"product-go-api/model"

	"github.com/gin-gonic/gin"
)

func RequireAdmin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		role, exists := ctx.Get("role")
		if !exists || (role != "admin" && role != "super_admin") {
			response := model.Response{
				Message: "Only Admins are allowed here.",
			}
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		ctx.Next()
	}
}
