package server

import (
	"auth-golang/internal/models"
	"auth-golang/internal/repositories/mocks"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAuthHandler_Register(t *testing.T) {
	// Como diria o Gandalf: "You shall not pass... sem testes!" 🧙‍♂️
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name               string
		requestBody        interface{}
		setupMock          func(*mocks.MockUsersRepository)
		expectedStatusCode int
		expectedBody       map[string]interface{}
		expectUserCreated  bool
	}{
		{
			name: "successful registration - como um Gopher bem-comportado 🐹",
			requestBody: map[string]string{
				"name":     "João Silva",
				"username": "joaosilva",
				"email":    "joao@example.com",
				"password": "senha123456",
			},
			setupMock: func(mockRepo *mocks.MockUsersRepository) {
				// User doesn't exist yet
				mockRepo.FindByEmailFunc = func(ctx context.Context, email string) (*models.User, error) {
					return nil, nil
				}
				// Create succeeds
				mockRepo.CreateFunc = func(ctx context.Context, user *models.User) error {
					return nil
				}
			},
			expectedStatusCode: http.StatusCreated,
			expectedBody: map[string]interface{}{
				"message": "User registered successfully",
			},
			expectUserCreated: true,
		},
		{
			name: "invalid request body - JSON mais quebrado que promessa de político 💔",
			requestBody: map[string]interface{}{
				"name":     "",
				"username": "js",
				"email":    "not-an-email",
				"password": "123",
			},
			setupMock: func(mockRepo *mocks.MockUsersRepository) {
				// Nenhuma função será chamada
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "Invalid request data",
			},
			expectUserCreated: false,
		},
		{
			name: "user already exists - como ex que volta, não queremos! 😅",
			requestBody: map[string]string{
				"name":     "João Silva",
				"username": "joaosilva",
				"email":    "joao@example.com",
				"password": "senha123456",
			},
			setupMock: func(mockRepo *mocks.MockUsersRepository) {
				existingUser := &models.User{
					ID:       uuid.New(),
					Name:     "João Silva",
					Username: "joaosilva",
					Email:    "joao@example.com",
				}
				mockRepo.FindByEmailFunc = func(ctx context.Context, email string) (*models.User, error) {
					return existingUser, nil
				}
			},
			expectedStatusCode: http.StatusConflict,
			expectedBody: map[string]interface{}{
				"error": "User already exists",
			},
			expectUserCreated: false,
		},
		{
			name: "database error on finding user - quando o DB resolve tirar férias 🏖️",
			requestBody: map[string]string{
				"name":     "João Silva",
				"username": "joaosilva",
				"email":    "joao@example.com",
				"password": "senha123456",
			},
			setupMock: func(mockRepo *mocks.MockUsersRepository) {
				mockRepo.FindByEmailFunc = func(ctx context.Context, email string) (*models.User, error) {
					return nil, errors.New("database connection error")
				}
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": "database connection error",
			},
			expectUserCreated: false,
		},
		{
			name: "database error on creating user - quando o DB decide ser rebelde 😈",
			requestBody: map[string]string{
				"name":     "João Silva",
				"username": "joaosilva",
				"email":    "joao@example.com",
				"password": "senha123456",
			},
			setupMock: func(mockRepo *mocks.MockUsersRepository) {
				mockRepo.FindByEmailFunc = func(ctx context.Context, email string) (*models.User, error) {
					return nil, nil
				}
				mockRepo.CreateFunc = func(ctx context.Context, user *models.User) error {
					return errors.New("failed to insert user")
				}
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": "failed to insert user",
			},
			expectUserCreated: false, // Create is called but fails
		},
		{
			name: "missing required fields - minimalist request, maximalist error 🎨",
			requestBody: map[string]string{
				"name": "João Silva",
				// missing username, email, password
			},
			setupMock: func(mockRepo *mocks.MockUsersRepository) {
				// Nenhuma função será chamada
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "Invalid request data",
			},
			expectUserCreated: false,
		},
		{
			name: "invalid email format - e-mail com mais problemas que novela mexicana 📺",
			requestBody: map[string]string{
				"name":     "João Silva",
				"username": "joaosilva",
				"email":    "invalid-email",
				"password": "senha123456",
			},
			setupMock: func(mockRepo *mocks.MockUsersRepository) {
				// Nenhuma função será chamada devido à validação
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "Invalid request data",
			},
			expectUserCreated: false,
		},
		{
			name: "password too short - senha mais curta que paciência em reunião 🕐",
			requestBody: map[string]string{
				"name":     "João Silva",
				"username": "joaosilva",
				"email":    "joao@example.com",
				"password": "123",
			},
			setupMock: func(mockRepo *mocks.MockUsersRepository) {
				// Nenhuma função será chamada devido à validação
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "Invalid request data",
			},
			expectUserCreated: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := mocks.NewMockUsersRepository()
			tt.setupMock(mockRepo)

			handler := newAuthHandler(mockRepo)

			// Create gin router and request
			router := gin.New()
			router.POST("/register", handler.Register)

			// Convert request body to JSON
			requestBodyBytes, err := json.Marshal(tt.requestBody)
			assert.NoError(t, err)

			req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(requestBodyBytes))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()

			// Act
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expectedStatusCode, w.Code)

			var responseBody map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &responseBody)
			assert.NoError(t, err)

			for key, expectedValue := range tt.expectedBody {
				assert.Equal(t, expectedValue, responseBody[key])
			}

			// Verify repository calls
			if tt.expectUserCreated {
				assert.Equal(t, 1, mockRepo.CreateCallCount, "Create should be called once")
				assert.NotNil(t, mockRepo.LastCreateUser, "User should be passed to Create")
				assert.NotEmpty(t, mockRepo.LastCreateUser.ID, "User ID should be generated")

				// Verify password was hashed (should not be the original password)
				if reqBody, ok := tt.requestBody.(map[string]string); ok {
					assert.NotEqual(t, reqBody["password"], mockRepo.LastCreateUser.Password, "Password should be hashed")
				}
			} else {
				// For database error on creating user, Create is called but fails
				if tt.name == "database error on creating user - quando o DB decide ser rebelde 😈" {
					assert.Equal(t, 1, mockRepo.CreateCallCount, "Create should be called once but fail")
				} else {
					assert.Equal(t, 0, mockRepo.CreateCallCount, "Create should not be called")
				}
			}
		})
	}
}

func TestAuthHandler_Login(t *testing.T) {
	// Login: onde senhas são testadas e tokens nascem! 🔓
	gin.SetMode(gin.TestMode)

	// Helper para criar um usuário válido
	createValidUser := func() *models.User {
		user, _ := models.NewUser("João Silva", "joaosilva", "joao@example.com", "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi") // bcrypt hash of "password"
		return user
	}

	tests := []struct {
		name               string
		requestBody        interface{}
		setupMock          func(*mocks.MockUsersRepository)
		expectedStatusCode int
		expectedBodyKeys   []string
		expectCookies      bool
	}{
		{
			name: "successful login - entrada triunfal como herói de filme 🦸‍♂️",
			requestBody: map[string]string{
				"email":    "joao@example.com",
				"password": "password",
			},
			setupMock: func(mockRepo *mocks.MockUsersRepository) {
				user := createValidUser()
				mockRepo.FindByEmailFunc = func(ctx context.Context, email string) (*models.User, error) {
					return user, nil
				}
				mockRepo.UpdateFunc = func(ctx context.Context, user *models.User) (*models.User, error) {
					return user, nil
				}
			},
			expectedStatusCode: http.StatusOK,
			expectedBodyKeys:   []string{"success"},
			expectCookies:      true,
		},
		{
			name: "invalid request body - JSON mais confuso que manual do IKEA 🛠️",
			requestBody: map[string]interface{}{
				"email":    "",
				"password": "123",
			},
			setupMock: func(mockRepo *mocks.MockUsersRepository) {
				// Nenhuma função será chamada
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedBodyKeys:   []string{"error"},
			expectCookies:      false,
		},
		{
			name: "user not found - procurando usuário que não existe 👻",
			requestBody: map[string]string{
				"email":    "ghost@example.com",
				"password": "password123",
			},
			setupMock: func(mockRepo *mocks.MockUsersRepository) {
				mockRepo.FindByEmailFunc = func(ctx context.Context, email string) (*models.User, error) {
					return nil, nil // User not found
				}
			},
			expectedStatusCode: http.StatusUnauthorized,
			expectedBodyKeys:   []string{"error"},
			expectCookies:      false,
		},
		{
			name: "wrong password - senha mais errada que previsão do tempo ⛈️",
			requestBody: map[string]string{
				"email":    "joao@example.com",
				"password": "wrongpassword",
			},
			setupMock: func(mockRepo *mocks.MockUsersRepository) {
				user := createValidUser()
				mockRepo.FindByEmailFunc = func(ctx context.Context, email string) (*models.User, error) {
					return user, nil
				}
			},
			expectedStatusCode: http.StatusUnauthorized,
			expectedBodyKeys:   []string{"error"},
			expectCookies:      false,
		},
		{
			name: "database error on finding user - DB decidiu sair para almoçar 🍽️",
			requestBody: map[string]string{
				"email":    "joao@example.com",
				"password": "password",
			},
			setupMock: func(mockRepo *mocks.MockUsersRepository) {
				mockRepo.FindByEmailFunc = func(ctx context.Context, email string) (*models.User, error) {
					return nil, errors.New("database connection timeout")
				}
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedBodyKeys:   []string{"error"},
			expectCookies:      false,
		},
		{
			name: "database error on update - DB travou na hora H 🎯",
			requestBody: map[string]string{
				"email":    "joao@example.com",
				"password": "password",
			},
			setupMock: func(mockRepo *mocks.MockUsersRepository) {
				user := createValidUser()
				mockRepo.FindByEmailFunc = func(ctx context.Context, email string) (*models.User, error) {
					return user, nil
				}
				mockRepo.UpdateFunc = func(ctx context.Context, user *models.User) (*models.User, error) {
					return nil, errors.New("failed to update user")
				}
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedBodyKeys:   []string{"error"},
			expectCookies:      false,
		},
		{
			name: "invalid email format - e-mail com mais problemas que aplicativo do governo 🏛️",
			requestBody: map[string]string{
				"email":    "not-an-email",
				"password": "password123",
			},
			setupMock: func(mockRepo *mocks.MockUsersRepository) {
				// Validation will fail before any repo calls
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedBodyKeys:   []string{"error"},
			expectCookies:      false,
		},
		{
			name: "short password - senha mais curta que comercial de TV 📺",
			requestBody: map[string]string{
				"email":    "joao@example.com",
				"password": "123",
			},
			setupMock: func(mockRepo *mocks.MockUsersRepository) {
				// Validation will fail before any repo calls
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedBodyKeys:   []string{"error"},
			expectCookies:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := mocks.NewMockUsersRepository()
			tt.setupMock(mockRepo)

			handler := newAuthHandler(mockRepo)

			// Create gin router and request
			router := gin.New()
			router.POST("/login", handler.Login)

			// Convert request body to JSON
			requestBodyBytes, err := json.Marshal(tt.requestBody)
			assert.NoError(t, err)

			req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(requestBodyBytes))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()

			// Act
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expectedStatusCode, w.Code)

			var responseBody map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &responseBody)
			assert.NoError(t, err)

			// Check that expected keys exist in response
			for _, key := range tt.expectedBodyKeys {
				assert.Contains(t, responseBody, key)
			}

			// Check cookies if expected
			if tt.expectCookies {
				cookies := w.Header()["Set-Cookie"]
				assert.NotEmpty(t, cookies, "Should have cookies set")

				// Look for access_token and refresh_token in any of the cookie headers
				hasAccessToken := false
				hasRefreshToken := false

				for _, cookie := range cookies {
					if strings.Contains(cookie, "access_token") {
						hasAccessToken = true
					}
					if strings.Contains(cookie, "refresh_token") {
						hasRefreshToken = true
					}
				}

				assert.True(t, hasAccessToken, "Should have access_token cookie")
				assert.True(t, hasRefreshToken, "Should have refresh_token cookie")
			}
		})
	}
}

func TestAuthHandler_RefreshToken(t *testing.T) {
	// Refresh Token: porque até tokens precisam de uma pausa para o café ☕
	gin.SetMode(gin.TestMode)

	validRefreshToken := "valid_refresh_token_12345"

	createValidUser := func() *models.User {
		user, _ := models.NewUser("João Silva", "joaosilva", "joao@example.com", "hashedpassword")
		user.RefreshToken = validRefreshToken
		return user
	}

	tests := []struct {
		name               string
		cookieValue        string
		setupMock          func(*mocks.MockUsersRepository)
		expectedStatusCode int
		expectedBodyKeys   []string
		expectNewCookies   bool
	}{
		{
			name:        "successful token refresh - renovação mais suave que jazz 🎷",
			cookieValue: validRefreshToken,
			setupMock: func(mockRepo *mocks.MockUsersRepository) {
				user := createValidUser()
				mockRepo.FindByRefreshTokenFunc = func(ctx context.Context, token string) (*models.User, error) {
					return user, nil
				}
				mockRepo.UpdateFunc = func(ctx context.Context, user *models.User) (*models.User, error) {
					return user, nil
				}
			},
			expectedStatusCode: http.StatusOK,
			expectedBodyKeys:   []string{"message"},
			expectNewCookies:   true,
		},
		{
			name:        "missing refresh token - cookie mais perdido que turista sem GPS 🗺️",
			cookieValue: "",
			setupMock: func(mockRepo *mocks.MockUsersRepository) {
				// No repo calls should happen
			},
			expectedStatusCode: http.StatusUnauthorized,
			expectedBodyKeys:   []string{"error"},
			expectNewCookies:   false,
		},
		{
			name:        "invalid refresh token - token mais falso que nota de 3 reais 💸",
			cookieValue: "invalid_token",
			setupMock: func(mockRepo *mocks.MockUsersRepository) {
				mockRepo.FindByRefreshTokenFunc = func(ctx context.Context, token string) (*models.User, error) {
					return nil, nil // User not found
				}
			},
			expectedStatusCode: http.StatusUnauthorized,
			expectedBodyKeys:   []string{"error"},
			expectNewCookies:   false,
		},
		{
			name:        "database error on finding user - DB em greve 🪧",
			cookieValue: validRefreshToken,
			setupMock: func(mockRepo *mocks.MockUsersRepository) {
				mockRepo.FindByRefreshTokenFunc = func(ctx context.Context, token string) (*models.User, error) {
					return nil, errors.New("database connection failed")
				}
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedBodyKeys:   []string{"error"},
			expectNewCookies:   false,
		},
		{
			name:        "database error on update - DB travou na atualização 🔧",
			cookieValue: validRefreshToken,
			setupMock: func(mockRepo *mocks.MockUsersRepository) {
				user := createValidUser()
				mockRepo.FindByRefreshTokenFunc = func(ctx context.Context, token string) (*models.User, error) {
					return user, nil
				}
				mockRepo.UpdateFunc = func(ctx context.Context, user *models.User) (*models.User, error) {
					return nil, errors.New("failed to update user")
				}
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedBodyKeys:   []string{"error"},
			expectNewCookies:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := mocks.NewMockUsersRepository()
			tt.setupMock(mockRepo)

			handler := newAuthHandler(mockRepo)

			// Create gin router and request
			router := gin.New()
			router.POST("/refresh", handler.RefreshToken)

			req, err := http.NewRequest("POST", "/refresh", nil)
			assert.NoError(t, err)

			// Add refresh token cookie if provided
			if tt.cookieValue != "" {
				cookie := &http.Cookie{
					Name:  "refresh_token",
					Value: tt.cookieValue,
				}
				req.AddCookie(cookie)
			}

			w := httptest.NewRecorder()

			// Act
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expectedStatusCode, w.Code)

			var responseBody map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &responseBody)
			assert.NoError(t, err)

			// Check that expected keys exist in response
			for _, key := range tt.expectedBodyKeys {
				assert.Contains(t, responseBody, key)
			}

			// Check new cookies if expected
			if tt.expectNewCookies {
				cookies := w.Header()["Set-Cookie"]
				assert.NotEmpty(t, cookies, "Should have cookies set")

				hasAccessToken := false
				hasRefreshToken := false

				for _, cookie := range cookies {
					if strings.Contains(cookie, "access_token") {
						hasAccessToken = true
					}
					if strings.Contains(cookie, "refresh_token") {
						hasRefreshToken = true
					}
				}

				assert.True(t, hasAccessToken, "Should have access_token cookie")
				assert.True(t, hasRefreshToken, "Should have refresh_token cookie")
			}
		})
	}
}

func TestAuthHandler_Logout(t *testing.T) {
	// Logout: a despedida mais elegante desde "Até a vista, baby!" 👋
	gin.SetMode(gin.TestMode)

	validRefreshToken := "valid_refresh_token_12345"

	createValidUser := func() *models.User {
		user, _ := models.NewUser("João Silva", "joaosilva", "joao@example.com", "hashedpassword")
		user.RefreshToken = validRefreshToken
		return user
	}

	tests := []struct {
		name                 string
		refreshTokenCookie   string
		setupMock            func(*mocks.MockUsersRepository)
		expectedStatusCode   int
		expectedMessage      string
		expectClearedCookies bool
	}{
		{
			name:               "successful logout with token cleanup - despedida com classe 🎩",
			refreshTokenCookie: validRefreshToken,
			setupMock: func(mockRepo *mocks.MockUsersRepository) {
				user := createValidUser()
				mockRepo.FindByRefreshTokenFunc = func(ctx context.Context, token string) (*models.User, error) {
					return user, nil
				}
				mockRepo.UpdateFunc = func(ctx context.Context, user *models.User) (*models.User, error) {
					return user, nil
				}
			},
			expectedStatusCode:   http.StatusOK,
			expectedMessage:      "Logged out successfully",
			expectClearedCookies: true,
		},
		{
			name:               "logout without refresh token - saída sem cerimônia 🚪",
			refreshTokenCookie: "",
			setupMock: func(mockRepo *mocks.MockUsersRepository) {
				// No repo calls should happen when there's no cookie
			},
			expectedStatusCode:   http.StatusOK,
			expectedMessage:      "Logged out successfully",
			expectClearedCookies: true,
		},
		{
			name:               "logout with invalid token - tentando limpar o que não existe 🧹",
			refreshTokenCookie: "invalid_token",
			setupMock: func(mockRepo *mocks.MockUsersRepository) {
				mockRepo.FindByRefreshTokenFunc = func(ctx context.Context, token string) (*models.User, error) {
					return nil, nil // User not found
				}
			},
			expectedStatusCode:   http.StatusOK,
			expectedMessage:      "Logged out successfully",
			expectClearedCookies: true,
		},
		{
			name:               "database error during cleanup - DB rebelde, mas logout funciona 😤",
			refreshTokenCookie: validRefreshToken,
			setupMock: func(mockRepo *mocks.MockUsersRepository) {
				mockRepo.FindByRefreshTokenFunc = func(ctx context.Context, token string) (*models.User, error) {
					return nil, errors.New("database error")
				}
			},
			expectedStatusCode:   http.StatusOK,
			expectedMessage:      "Logged out successfully",
			expectClearedCookies: true,
		},
		{
			name:               "update error during cleanup - erro na limpeza, mas logout OK 🔧",
			refreshTokenCookie: validRefreshToken,
			setupMock: func(mockRepo *mocks.MockUsersRepository) {
				user := createValidUser()
				mockRepo.FindByRefreshTokenFunc = func(ctx context.Context, token string) (*models.User, error) {
					return user, nil
				}
				mockRepo.UpdateFunc = func(ctx context.Context, user *models.User) (*models.User, error) {
					return nil, errors.New("update failed")
				}
			},
			expectedStatusCode:   http.StatusOK,
			expectedMessage:      "Logged out successfully",
			expectClearedCookies: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := mocks.NewMockUsersRepository()
			tt.setupMock(mockRepo)

			handler := newAuthHandler(mockRepo)

			// Create gin router and request
			router := gin.New()
			router.POST("/logout", handler.Logout)

			req, err := http.NewRequest("POST", "/logout", nil)
			assert.NoError(t, err)

			// Add refresh token cookie if provided
			if tt.refreshTokenCookie != "" {
				cookie := &http.Cookie{
					Name:  "refresh_token",
					Value: tt.refreshTokenCookie,
				}
				req.AddCookie(cookie)
			}

			w := httptest.NewRecorder()

			// Act
			router.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, tt.expectedStatusCode, w.Code)

			var responseBody map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &responseBody)
			assert.NoError(t, err)

			assert.Equal(t, tt.expectedMessage, responseBody["message"])

			// Check that cookies are cleared
			if tt.expectClearedCookies {
				cookies := w.Header()["Set-Cookie"]
				assert.NotEmpty(t, cookies, "Should have cookies set for clearing")

				// Cookies should be cleared (Max-Age=0 or Max-Age=-1)
				hasAccessTokenCleared := false

				for _, cookie := range cookies {
					if strings.Contains(cookie, "access_token") && (strings.Contains(cookie, "Max-Age=0") || strings.Contains(cookie, "Max-Age=-1")) {
						hasAccessTokenCleared = true
					}
				}

				assert.True(t, hasAccessTokenCleared, "access_token should be cleared")
				// Note: refresh_token might not always be present in response if there was no original token
			}
		})
	}
}
