package server

import (
	"auth-golang/internal/auth"
	"auth-golang/internal/models"
	"auth-golang/internal/repositories"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	userRepository repositories.UsersRepositoryInterface
}

func newAuthHandler(userRepository repositories.UsersRepositoryInterface) *AuthHandler {
	return &AuthHandler{
		userRepository: userRepository,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		Name     string `json:"name" binding:"required"`
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	user, err := h.userRepository.FindByEmail(c.Request.Context(), req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if user != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	// hash password
	hashedPassword, err := auth.HashPassword(req.Password)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//create a new user
	user, err = models.NewUser(req.Name, req.Username, req.Email, hashedPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("error creating user: %s", err.Error()))
		return
	}

	//save the user to the database
	err = h.userRepository.Create(c.Request.Context(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
	})

}

func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	user, err := h.userRepository.FindByEmail(c.Request.Context(), req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Login failed. Please check your credentials."})
		return
	}

	if !auth.CheckPasswordHash(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Login failed. Please check your credentials."})
		return
	}

	accessToken, err := auth.GenerateAccessToken(user.ID.Hex())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	refreshToken, err := auth.GenerateRefreshToken(user.ID.Hex())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	c.SetCookie(
		"access_token",
		accessToken,
		15*60,
		"/",
		"localhost",
		false,
		true,
	)

	user.AddRefreshToken(refreshToken)

	_, err = h.userRepository.Update(c.Request.Context(), user)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user with refresh token"})
		return
	}

	c.SetCookie(
		"refresh_token",
		refreshToken,
		7*24*60*60,
		"/api/auth/refresh",
		"localhost",
		false,
		true,
	)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token not found"})
		return
	}

	user, err := h.userRepository.FindByRefreshToken(c.Request.Context(), refreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	newAccessToken, err := auth.GenerateAccessToken(user.ID.Hex())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate new access token"})
		return
	}

	newRefreshToken, err := auth.GenerateRefreshToken(user.ID.Hex())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate new refresh token"})
		return
	}

	c.SetCookie(
		"access_token",
		newAccessToken,
		15*60,
		"/",
		"localhost",
		false,
		true,
	)

	user.AddRefreshToken(newRefreshToken)
	_, err = h.userRepository.Update(c.Request.Context(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user with new refresh token"})
		return
	}

	c.SetCookie(
		"refresh_token",
		newRefreshToken,
		7*24*60*60,
		"/api/auth/refresh",
		"localhost",
		false,
		true,
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "Token refreshed successfully",
	})

}

func (h *AuthHandler) Logout(c *gin.Context) {
	// Clear the cookies
	c.SetCookie("access_token", "", -1, "/", "localhost", false, true)
	c.SetCookie("refresh_token", "", -1, "/api/auth/refresh", "localhost", false, true)

	// Optionally, you can also remove the refresh token from the user in the database
	refreshToken, err := c.Cookie("refresh_token")
	if err == nil {
		user, err := h.userRepository.FindByRefreshToken(c.Request.Context(), refreshToken)
		if err == nil && user != nil {
			user.AddRefreshToken("")                                  // Clear the refresh token
			_, _ = h.userRepository.Update(c.Request.Context(), user) // Ignore error for logout
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}
