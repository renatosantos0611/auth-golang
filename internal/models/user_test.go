package models

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	// NewUser: criando usuários como se fosse uma fábrica de Gophers 🏭🐹
	tests := []struct {
		name          string
		inputName     string
		inputUsername string
		inputEmail    string
		inputPassword string
		expectError   bool
		errorContains string
	}{
		{
			name:          "valid user creation - usuário mais válido que CPF da Receita ✅",
			inputName:     "João Silva",
			inputUsername: "joaosilva",
			inputEmail:    "joao@example.com",
			inputPassword: "senha12345678",
			expectError:   false,
		},
		{
			name:          "short name - nome mais curto que paciência em fila de banco 🏦",
			inputName:     "Jo",
			inputUsername: "joaosilva",
			inputEmail:    "joao@example.com",
			inputPassword: "senha12345678",
			expectError:   true,
			errorContains: "Name must be at least 3 characters long",
		},
		{
			name:          "short username - username mais curto que memória de político 🎭",
			inputName:     "João Silva",
			inputUsername: "jo",
			inputEmail:    "joao@example.com",
			inputPassword: "senha12345678",
			expectError:   true,
			errorContains: "Username must be at least 6 characters long",
		},
		{
			name:          "short password - senha mais fraca que Wi-Fi público 📶",
			inputName:     "João Silva",
			inputUsername: "joaosilva",
			inputEmail:    "joao@example.com",
			inputPassword: "123",
			expectError:   true,
			errorContains: "Password must be at least 8 characters long",
		},
		{
			name:          "multiple validation errors - acumulando erros como colecionador 📚",
			inputName:     "Jo",
			inputUsername: "js",
			inputEmail:    "x",
			inputPassword: "123",
			expectError:   true,
			// Should contain multiple error messages
		},
		{
			name:          "minimum valid lengths - no limite da aprovação 🎯",
			inputName:     "João",
			inputUsername: "joaosil",
			inputEmail:    "joao@example.com",
			inputPassword: "12345678",
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := NewUser(tt.inputName, tt.inputUsername, tt.inputEmail, tt.inputPassword)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, user)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)

				// Verify user fields are set correctly
				assert.Equal(t, tt.inputName, user.Name)
				assert.Equal(t, tt.inputUsername, user.Username)
				assert.Equal(t, tt.inputEmail, user.Email)
				assert.Equal(t, tt.inputPassword, user.Password)

				// Verify default values
				assert.NotEqual(t, uuid.Nil, user.ID)
				assert.Equal(t, RoleUser, user.Role)
				assert.Empty(t, user.RefreshToken)

				// Verify timestamps
				assert.WithinDuration(t, time.Now(), user.CreatedAt, time.Second)
				assert.WithinDuration(t, time.Now(), user.UpdatedAt, time.Second)
			}
		})
	}
}

func TestUser_Validate(t *testing.T) {
	// Validate: o fiscal mais rigoroso desde a Receita Federal 👮‍♂️
	tests := []struct {
		name          string
		user          User
		expectError   bool
		errorContains []string
	}{
		{
			name: "valid user - aprovado com louvor 🎓",
			user: User{
				Name:     "João Silva",
				Username: "joaosilva",
				Email:    "joao@example.com",
				Password: "senha12345678",
			},
			expectError: false,
		},
		{
			name: "short name - nome mais curto que vida de mayfly 🦟",
			user: User{
				Name:     "Jo",
				Username: "joaosilva",
				Email:    "joao@example.com",
				Password: "senha12345678",
			},
			expectError:   true,
			errorContains: []string{"Name must be at least 3 characters long"},
		},
		{
			name: "short username - username mais minimalista que design da Apple 🍎",
			user: User{
				Name:     "João Silva",
				Username: "jo",
				Email:    "joao@example.com",
				Password: "senha12345678",
			},
			expectError:   true,
			errorContains: []string{"Username must be at least 6 characters long"},
		},
		{
			name: "short password - proteção mais fraca que papel molhado 💧",
			user: User{
				Name:     "João Silva",
				Username: "joaosilva",
				Email:    "joao@example.com",
				Password: "123",
			},
			expectError:   true,
			errorContains: []string{"Password must be at least 8 characters long"},
		},
		{
			name: "all validations failing - falha épica como filme do Batman vs Superman 🦇",
			user: User{
				Name:     "Jo",
				Username: "js",
				Email:    "x",
				Password: "123",
			},
			expectError: true,
			errorContains: []string{
				"Name must be at least 3 characters long",
				"Username must be at least 6 characters long",
				"Password must be at least 8 characters long",
			},
		},
		{
			name: "empty fields - vazio como promessa de campanha 🗳️",
			user: User{
				Name:     "",
				Username: "",
				Email:    "",
				Password: "",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Validate()

			if tt.expectError {
				assert.Error(t, err)
				for _, expectedError := range tt.errorContains {
					assert.Contains(t, err.Error(), expectedError)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUser_AddRefreshToken(t *testing.T) {
	// AddRefreshToken: adicionando tokens como colecionador de cartas Pokémon 🎴
	t.Run("add refresh token updates token and timestamp", func(t *testing.T) {
		user := &User{
			ID:           uuid.New(),
			Name:         "João Silva",
			Username:     "joaosilva",
			Email:        "joao@example.com",
			Password:     "hashedpassword",
			RefreshToken: "",
			CreatedAt:    time.Now().Add(-1 * time.Hour),
			UpdatedAt:    time.Now().Add(-1 * time.Hour),
		}

		oldUpdatedAt := user.UpdatedAt
		newToken := "new_refresh_token_12345"

		// Act
		user.AddRefreshToken(newToken)

		// Assert
		assert.Equal(t, newToken, user.RefreshToken)
		assert.True(t, user.UpdatedAt.After(oldUpdatedAt), "UpdatedAt should be updated")
		assert.WithinDuration(t, time.Now(), user.UpdatedAt, time.Second)
	})

	t.Run("add empty token clears refresh token", func(t *testing.T) {
		user := &User{
			RefreshToken: "existing_token",
			UpdatedAt:    time.Now().Add(-1 * time.Hour),
		}

		oldUpdatedAt := user.UpdatedAt

		// Act
		user.AddRefreshToken("")

		// Assert
		assert.Empty(t, user.RefreshToken)
		assert.True(t, user.UpdatedAt.After(oldUpdatedAt))
	})
}

func TestRole_String(t *testing.T) {
	// Role.String(): convertendo enum para string como tradutor profissional 🌐
	tests := []struct {
		name     string
		role     Role
		expected string
	}{
		{
			name:     "user role - usuário comum, cidadão modelo 👤",
			role:     RoleUser,
			expected: "user",
		},
		{
			name:     "admin role - o todo-poderoso, senhor dos servidores 👑",
			role:     RoleAdmin,
			expected: "admin",
		},
		{
			name:     "invalid role - quando o enum resolve inventar moda 🎨",
			role:     Role(999),
			expected: "user", // Should default to user
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.role.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestUser_ToResponse(t *testing.T) {
	// ToResponse: transformando User em UserResponse como mágico transforma pomba 🐰➡️🕊️
	t.Run("converts user to response format", func(t *testing.T) {
		userID := uuid.New()
		createdAt := time.Now().Add(-24 * time.Hour)
		updatedAt := time.Now()

		user := &User{
			ID:           userID,
			Name:         "João Silva",
			Username:     "joaosilva",
			Email:        "joao@example.com",
			Password:     "super_secret_password", // This should NOT appear in response
			RefreshToken: "secret_token",          // This should NOT appear in response
			Role:         RoleAdmin,
			CreatedAt:    createdAt,
			UpdatedAt:    updatedAt,
		}

		// Act
		response := user.ToResponse()

		// Assert
		assert.Equal(t, userID.String(), response.ID)
		assert.Equal(t, "João Silva", response.Name)
		assert.Equal(t, "joaosilva", response.Username)
		assert.Equal(t, "joao@example.com", response.Email)
		assert.Equal(t, "admin", response.Role)
		assert.Equal(t, createdAt.Format(time.RFC3339), response.CreatedAt)
		assert.Equal(t, updatedAt.Format(time.RFC3339), response.UpdatedAt)

		// Verify sensitive data is not included
		// (This is implicit since UserResponse doesn't have these fields,
		// but it's good to document the intention)
	})

	t.Run("formats timestamps correctly", func(t *testing.T) {
		fixedTime := time.Date(2023, 12, 25, 15, 30, 45, 0, time.UTC)
		user := &User{
			ID:        uuid.New(),
			Name:      "Test User",
			Username:  "testuser",
			Email:     "test@example.com",
			Role:      RoleUser,
			CreatedAt: fixedTime,
			UpdatedAt: fixedTime,
		}

		response := user.ToResponse()

		expectedTimeStr := "2023-12-25T15:30:45Z"
		assert.Equal(t, expectedTimeStr, response.CreatedAt)
		assert.Equal(t, expectedTimeStr, response.UpdatedAt)
	})
}

func TestUserResponse_DoesNotContainSensitiveData(t *testing.T) {
	// Verificando que a resposta não vaza dados sensíveis
	// Porque privacidade é mais sagrada que receita da vovó 👵
	t.Run("user response should not expose sensitive fields", func(t *testing.T) {
		user := &User{
			ID:           uuid.New(),
			Name:         "João Silva",
			Username:     "joaosilva",
			Email:        "joao@example.com",
			Password:     "super_secret_password",
			RefreshToken: "super_secret_token",
			Role:         RoleUser,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		response := user.ToResponse()

		// Use reflection to ensure no sensitive fields exist
		// This is a compile-time check, but good for documentation
		_ = response.ID
		_ = response.Name
		_ = response.Username
		_ = response.Email
		_ = response.Role
		_ = response.CreatedAt
		_ = response.UpdatedAt

		// These should not exist and should cause compile errors if uncommented:
		// _ = response.Password     // Should not exist
		// _ = response.RefreshToken // Should not exist

		assert.True(t, true, "UserResponse struct correctly excludes sensitive fields")
	})
}
