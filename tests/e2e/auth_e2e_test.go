package e2e

import (
	"auth-golang/internal/database"
	"auth-golang/internal/middleware"
	"auth-golang/internal/repositories"
	"auth-golang/internal/server"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var testServer *httptest.Server

func TestMain(m *testing.M) {
	// Configurar variáveis de ambiente para banco de teste
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5435")
	os.Setenv("DB_USER", "admin")
	os.Setenv("DB_PASSWORD", "123")
	os.Setenv("DB_NAME", "auth_golang")

	// Conectar ao banco
	db := database.New()
	userRepo := repositories.NewUsersRepository(db)

	// Criar handlers
	authHandler := &server.AuthHandler{
		UserRepository: userRepo,
	}
	userHandler := &server.UserHandler{
		UserRepository: userRepo,
	}

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Authorization", "Content-Type", "Accept"},
		AllowCredentials: true,
	}))

	authRoutes := r.Group("/api/auth")
	{
		authRoutes.POST("/register", authHandler.Register)
		authRoutes.POST("/login", authHandler.Login)
		authRoutes.POST("/refresh", authHandler.RefreshToken)
		authRoutes.POST("/logout", authHandler.Logout)
	}

	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware(userRepo))
	{
		protected.GET("/user", userHandler.GetUser)
	}

	// Limpar tabela users antes dos testes
	_, err := db.DB.Exec("DELETE FROM users")
	if err != nil {
		log.Printf("Failed to clean users table: %v", err)
	}
	testServer = httptest.NewServer(r)
	defer testServer.Close()

	os.Exit(m.Run())
}

func TestRegisterEndpoint(t *testing.T) {
	t.Parallel()

	payload := map[string]string{
		"name":     "Test User",
		"username": "testuser",
		"email":    "test@example.com",
		"password": "password123",
	}
	body, _ := json.Marshal(payload)

	resp, err := http.Post(testServer.URL+"/api/auth/register", "application/json", bytes.NewBuffer(body))
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}

func TestLoginEndpoint(t *testing.T) {
	t.Parallel()

	// Primeiro registrar
	payload := map[string]string{
		"name":     "Test User",
		"username": "testuser2",
		"email":    "test2@example.com",
		"password": "password123",
	}
	body, _ := json.Marshal(payload)
	http.Post(testServer.URL+"/api/auth/register", "application/json", bytes.NewBuffer(body))

	// Agora login
	loginPayload := map[string]string{
		"email":    "test2@example.com",
		"password": "password123",
	}
	loginBody, _ := json.Marshal(loginPayload)

	resp, err := http.Post(testServer.URL+"/api/auth/login", "application/json", bytes.NewBuffer(loginBody))
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGetUserEndpoint(t *testing.T) {
	t.Parallel()

	// Registrar e logar para obter token
	payload := map[string]string{
		"name":     "Test User",
		"username": "testuser3",
		"email":    "test3@example.com",
		"password": "password123",
	}
	body, _ := json.Marshal(payload)
	http.Post(testServer.URL+"/api/auth/register", "application/json", bytes.NewBuffer(body))

	loginPayload := map[string]string{
		"email":    "test3@example.com",
		"password": "password123",
	}
	loginBody, _ := json.Marshal(loginPayload)

	client := &http.Client{}
	resp, err := client.Post(testServer.URL+"/api/auth/login", "application/json", bytes.NewBuffer(loginBody))
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Agora fazer request para /api/user com os cookies
	req, _ := http.NewRequest("GET", testServer.URL+"/api/user", nil)
	for _, cookie := range resp.Cookies() {
		req.AddCookie(cookie)
	}

	userResp, err := client.Do(req)
	assert.NoError(t, err)
	defer userResp.Body.Close()

	assert.Equal(t, http.StatusOK, userResp.StatusCode)
}

func TestRefreshTokenEndpoint(t *testing.T) {
	t.Parallel()

	// Registrar e logar
	payload := map[string]string{
		"name":     "Test User",
		"username": "testuser4",
		"email":    "test4@example.com",
		"password": "password123",
	}
	body, _ := json.Marshal(payload)
	http.Post(testServer.URL+"/api/auth/register", "application/json", bytes.NewBuffer(body))

	loginPayload := map[string]string{
		"email":    "test4@example.com",
		"password": "password123",
	}
	loginBody, _ := json.Marshal(loginPayload)

	client := &http.Client{}
	resp, err := client.Post(testServer.URL+"/api/auth/login", "application/json", bytes.NewBuffer(loginBody))
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Refresh
	req, _ := http.NewRequest("POST", testServer.URL+"/api/auth/refresh", bytes.NewBuffer([]byte("{}")))
	for _, cookie := range resp.Cookies() {
		req.AddCookie(cookie)
	}

	refreshResp, err := client.Do(req)
	assert.NoError(t, err)
	defer refreshResp.Body.Close()

	assert.Equal(t, http.StatusOK, refreshResp.StatusCode)
}

func TestLogoutEndpoint(t *testing.T) {
	t.Parallel()

	resp, err := http.Post(testServer.URL+"/api/auth/logout", "application/json", bytes.NewBuffer([]byte("{}")))
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestCompleteFlow(t *testing.T) {
	t.Parallel()

	client := &http.Client{}

	// 1. Register
	payload := map[string]string{
		"name":     "Fluxo Completo",
		"username": "fluxocompleto",
		"email":    "fluxo@completo.com",
		"password": "password123",
	}
	body, _ := json.Marshal(payload)
	resp, err := client.Post(testServer.URL+"/api/auth/register", "application/json", bytes.NewBuffer(body))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	resp.Body.Close()

	// 2. Login
	loginPayload := map[string]string{
		"email":    "fluxo@completo.com",
		"password": "password123",
	}
	loginBody, _ := json.Marshal(loginPayload)
	resp, err = client.Post(testServer.URL+"/api/auth/login", "application/json", bytes.NewBuffer(loginBody))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()

	// 3. Get User
	req, _ := http.NewRequest("GET", testServer.URL+"/api/user", nil)
	for _, cookie := range resp.Cookies() {
		req.AddCookie(cookie)
	}
	userResp, err := client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, userResp.StatusCode)
	userResp.Body.Close()

	// 4. Refresh
	req, _ = http.NewRequest("POST", testServer.URL+"/api/auth/refresh", bytes.NewBuffer([]byte("{}")))
	for _, cookie := range resp.Cookies() {
		req.AddCookie(cookie)
	}
	refreshResp, err := client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, refreshResp.StatusCode)
	refreshResp.Body.Close()

	// 5. Logout
	req, _ = http.NewRequest("POST", testServer.URL+"/api/auth/logout", bytes.NewBuffer([]byte("{}")))
	for _, cookie := range resp.Cookies() {
		req.AddCookie(cookie)
	}
	logoutResp, err := client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, logoutResp.StatusCode)
	logoutResp.Body.Close()

	// 6. Try to access user after logout (should fail)
	req, _ = http.NewRequest("GET", testServer.URL+"/api/user", nil)
	for _, cookie := range resp.Cookies() {
		req.AddCookie(cookie)
	}
	finalResp, err := client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, finalResp.StatusCode)
	finalResp.Body.Close()
}
