package middleware

import (
	"auth-golang/internal/auth"
	"auth-golang/internal/repositories"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(
	userRepository repositories.UsersRepositoryInterface,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken, err := c.Cookie("access_token")

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized",
			})
			return
		}

		secret := os.Getenv("ACCESS_TOKEN_SECRET")
		_, claims, err := auth.ValidateToken(accessToken, secret)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token",
			})
			return
		}

		userId, ok := claims["sub"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid userID in token",
			})
			return
		}

		user, err := userRepository.FindByID(c.Request.Context(), userId)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		if user == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "User not found",
			})
			return
		}

		// Verificar se o refresh token ainda é válido (não foi limpo no logout)
		refreshToken, err := c.Cookie("refresh_token")
		if err != nil || refreshToken == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Refresh token missing",
			})
			return
		}

		// Verificar se o refresh token do usuário corresponde
		if user.RefreshToken != refreshToken {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid refresh token",
			})
			return
		}

		c.Set("userId", userId)
		c.Next()
	}

}
