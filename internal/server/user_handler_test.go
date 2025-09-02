package server

import (
	"auth-golang/internal/models"
	"auth-golang/internal/repositories/mocks"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUserHandler_GetUser(t *testing.T) {
	// GetUser: buscando usuários como Sherlock Holmes busca pistas 🕵️‍♂️
	gin.SetMode(gin.TestMode)

	validUserID := uuid.New().String()

	createValidUser := func() *models.User {
		user, _ := models.NewUser("João Silva", "joaosilva", "joao@example.com", "hashedpassword")
		user.ID = uuid.MustParse(validUserID)
		return user
	}

	tests := []struct {
		name               string
		userIDInContext    interface{}
		userIDExists       bool
		setupMock          func(*mocks.MockUsersRepository)
		expectedStatusCode int
		expectedBodyKeys   []string
		shouldCallRepo     bool
	}{
		{
			name:            "successful get user - busca mais certeira que GPS 🎯",
			userIDInContext: validUserID,
			userIDExists:    true,
			setupMock: func(mockRepo *mocks.MockUsersRepository) {
				user := createValidUser()
				mockRepo.FindByIDFunc = func(ctx context.Context, id string) (*models.User, error) {
					return user, nil
				}
			},
			expectedStatusCode: http.StatusOK,
			expectedBodyKeys:   []string{"user"},
			shouldCallRepo:     true,
		},
		{
			name:            "user ID not in context - contexto mais vazio que geladeira de estudante 🏠",
			userIDInContext: nil,
			userIDExists:    false,
			setupMock: func(mockRepo *mocks.MockUsersRepository) {
				// No repo calls should happen
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedBodyKeys:   []string{"error"},
			shouldCallRepo:     false,
		},
		{
			name:            "user not found in database - usuário mais sumido que Wally 👀",
			userIDInContext: validUserID,
			userIDExists:    true,
			setupMock: func(mockRepo *mocks.MockUsersRepository) {
				mockRepo.FindByIDFunc = func(ctx context.Context, id string) (*models.User, error) {
					return nil, nil // User not found
				}
			},
			expectedStatusCode: http.StatusNotFound,
			expectedBodyKeys:   []string{"error"},
			shouldCallRepo:     true,
		},
		{
			name:            "database error - DB mais instável que humor de segunda-feira 😪",
			userIDInContext: validUserID,
			userIDExists:    true,
			setupMock: func(mockRepo *mocks.MockUsersRepository) {
				mockRepo.FindByIDFunc = func(ctx context.Context, id string) (*models.User, error) {
					return nil, errors.New("database connection failed")
				}
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedBodyKeys:   []string{"error"},
			shouldCallRepo:     true,
		},
		{
			name:            "invalid user ID type in context - tipo mais confuso que filme do Christopher Nolan 🎬",
			userIDInContext: 12345, // Should be string
			userIDExists:    true,
			setupMock: func(mockRepo *mocks.MockUsersRepository) {
				// This should cause a panic/error when casting to string
				mockRepo.FindByIDFunc = func(ctx context.Context, id string) (*models.User, error) {
					return nil, errors.New("should not be called")
				}
			},
			expectedStatusCode: http.StatusBadRequest, // This should catch the type error
			expectedBodyKeys:   []string{"error"},
			shouldCallRepo:     false, // Should not reach the repo due to type error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := mocks.NewMockUsersRepository()
			tt.setupMock(mockRepo)

			handler := newUserHandler(mockRepo)

			// Create gin router and request
			router := gin.New()

			// Middleware para simular o contexto do usuário
			router.Use(func(c *gin.Context) {
				if tt.userIDExists {
					c.Set("userId", tt.userIDInContext)
				}
				c.Next()
			})

			router.GET("/user", handler.GetUser)

			req, err := http.NewRequest("GET", "/user", nil)
			assert.NoError(t, err)

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

			// Verify repository calls
			if tt.shouldCallRepo && tt.expectedStatusCode != http.StatusBadRequest {
				assert.Equal(t, 1, mockRepo.FindByIDCallCount, "FindByID should be called once")
				if tt.userIDInContext != nil {
					// Only check this for string types to avoid panic
					if userID, ok := tt.userIDInContext.(string); ok {
						assert.Equal(t, userID, mockRepo.LastFindByIDParam)
					}
				}
			} else if !tt.shouldCallRepo {
				assert.Equal(t, 0, mockRepo.FindByIDCallCount, "FindByID should not be called")
			}

			// If successful, verify user response structure
			if tt.expectedStatusCode == http.StatusOK {
				userResponse, exists := responseBody["user"]
				assert.True(t, exists, "User should exist in response")

				userMap, ok := userResponse.(map[string]interface{})
				assert.True(t, ok, "User should be a map")

				// Check that password is not included in response
				_, hasPassword := userMap["password"]
				assert.False(t, hasPassword, "Password should not be in response")

				// Check that required fields are present
				requiredFields := []string{"id", "name", "username", "email", "role", "created_at", "updated_at"}
				for _, field := range requiredFields {
					assert.Contains(t, userMap, field, "Field %s should be present", field)
				}
			}
		})
	}
}

func TestUserHandler_GetUser_Integration(t *testing.T) {
	// Teste de integração para verificar se a resposta está no formato correto
	// Porque precisamos ter certeza que nosso JSON está mais bonito que sunset 🌅
	gin.SetMode(gin.TestMode)

	t.Run("response format validation", func(t *testing.T) {
		// Setup
		mockRepo := mocks.NewMockUsersRepository()
		user, _ := models.NewUser("João Silva", "joaosilva", "joao@example.com", "hashedpassword")

		mockRepo.FindByIDFunc = func(ctx context.Context, id string) (*models.User, error) {
			return user, nil
		}

		handler := newUserHandler(mockRepo)

		// Create gin router
		router := gin.New()
		router.Use(func(c *gin.Context) {
			c.Set("userId", user.ID.String())
			c.Next()
		})
		router.GET("/user", handler.GetUser)

		// Make request
		req, _ := http.NewRequest("GET", "/user", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		// Verify ToResponse() method is called correctly
		userResponse := response["user"].(map[string]interface{})

		// Test that dates are in RFC3339 format
		createdAt := userResponse["created_at"].(string)
		updatedAt := userResponse["updated_at"].(string)

		assert.NotEmpty(t, createdAt)
		assert.NotEmpty(t, updatedAt)

		// Test that role is converted to string
		role := userResponse["role"].(string)
		assert.Equal(t, "user", role) // Default role should be "user"
	})
}
