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
				"error": "Não autorizado",
			})
			return
		}

		secret := os.Getenv("ACCESS_TOKEN_SECRET")
		_, claims, err := auth.ValidateToken(accessToken, secret)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Token inválido",
			})
			return
		}

		userId, ok := claims["sub"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "ID de usuário inválido no token",
			})
			return
		}

		user, err := userRepository.FindByID(c.Request.Context(), userId)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Erro interno ao buscar usuário",
			})
			return
		}

		if user == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Usuário não encontrado",
			})
			return
		}

		// Verificar se o refresh token ainda é válido (não foi limpo no logout)
		refreshToken, err := c.Cookie("refresh_token")
		if err != nil || refreshToken == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Refresh token ausente",
			})
			return
		}

		// Verificar se o refresh token do usuário corresponde
		if user.RefreshToken != refreshToken {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Refresh token inválido",
			})
			return
		}

		c.Set("userId", userId)
		c.Next()
	}

}
