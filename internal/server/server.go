package server

import (
	"auth-golang/internal/database"
	"auth-golang/internal/middleware"
	"auth-golang/internal/repositories"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	usersRepository repositories.UsersRepositoryInterface
}

func NewServer() *http.Server {
	port := os.Getenv("PORT")

	db := database.New()
	userRepository := repositories.NewUsersRepository(db)

	server := &Server{
		usersRepository: userRepository,
	}

	return &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      server.RegisterRoutes(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Authorization", "Content-Type", "Accept"},
		AllowCredentials: true,
	}))

	authHandler := newAuthHandler(s.usersRepository)
	userHandler := newUserHandler(s.usersRepository)

	authRoutes := r.Group("/api/auth")
	{
		authRoutes.POST("/register", authHandler.Register)
		authRoutes.POST("/login", authHandler.Login)
		authRoutes.POST("/refresh", authHandler.RefreshToken)
		authRoutes.POST("/logout", authHandler.Logout)
	}

	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware(s.usersRepository))
	{
		protected.GET("/user", userHandler.GetUser)

	}

	return r
}
