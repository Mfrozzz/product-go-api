package middleware

import (
	"net/http"
	"os"
	"product-go-api/model"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func AuthMiddleware() gin.HandlerFunc {
	godotenv.Load()
	var jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			response := model.Response{
				Message: "Missing or invalid token",
			}
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil || !token.Valid {
			response := model.Response{
				Message: "Invalid Token",
			}
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if role, ok := claims["role"].(string); ok {
				ctx.Set("role", role)
			}
		}

		ctx.Next()
	}
}
