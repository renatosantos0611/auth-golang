package auth

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword(t *testing.T) {
	// HashPassword: transformando senhas como Batman transforma Gotham 🦇
	tests := []struct {
		name        string
		password    string
		expectError bool
	}{
		{
			name:        "valid password - senha forte como Hulk 💪",
			password:    "minhaSenhaSegura123",
			expectError: false,
		},
		{
			name:        "short password - mas ainda funciona 🔐",
			password:    "123",
			expectError: false,
		},
		{
			name:        "empty password - vazio como coração do ex 💔",
			password:    "",
			expectError: false, // bcrypt accepts empty passwords
		},
		{
			name:        "long password - senha épica como saga do Senhor dos Anéis 📚",
			password:    strings.Repeat("a", 70), // Max 72 bytes for bcrypt
			expectError: false,
		},
		{
			name:        "unicode password - senha internacional 🌍",
			password:    "minhaSenha🔒çãoÀÇÜ",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hashedPassword, err := HashPassword(tt.password)

			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, hashedPassword)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, hashedPassword)
				assert.NotEqual(t, tt.password, hashedPassword, "Password should be hashed")

				// Verify it's a valid bcrypt hash
				assert.True(t, strings.HasPrefix(hashedPassword, "$2a$"), "Should be a bcrypt hash")

				// Verify the hash can be used to verify the original password
				err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(tt.password))
				assert.NoError(t, err, "Hash should verify against original password")
			}
		})
	}
}

func TestCheckPasswordHash(t *testing.T) {
	// CheckPasswordHash: verificador mais confiável que detector de mentiras 🕵️‍♀️
	// First, create a known hash for testing
	testPassword := "minhaSenhaSecreta123"
	testHash, err := HashPassword(testPassword)
	assert.NoError(t, err)

	tests := []struct {
		name        string
		password    string
		hash        string
		shouldMatch bool
	}{
		{
			name:        "correct password - match perfeito como puzzle resolvido 🧩",
			password:    testPassword,
			hash:        testHash,
			shouldMatch: true,
		},
		{
			name:        "wrong password - erro mais claro que água 💧",
			password:    "senhaErrada",
			hash:        testHash,
			shouldMatch: false,
		},
		{
			name:        "empty password against hash - vazio tentando se passar por cheio 🕳️",
			password:    "",
			hash:        testHash,
			shouldMatch: false,
		},
		{
			name:        "password against empty hash - comparando com o nada 👻",
			password:    testPassword,
			hash:        "",
			shouldMatch: false,
		},
		{
			name:        "empty password against empty hash - duplo vazio 🔄",
			password:    "",
			hash:        "",
			shouldMatch: false,
		},
		{
			name:        "password against invalid hash - hash mais inválido que nota de 3 reais 💸",
			password:    testPassword,
			hash:        "not_a_valid_bcrypt_hash",
			shouldMatch: false,
		},
		{
			name:        "case sensitive password - maiúscula faz diferença 📝",
			password:    strings.ToUpper(testPassword),
			hash:        testHash,
			shouldMatch: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CheckPasswordHash(tt.password, tt.hash)
			assert.Equal(t, tt.shouldMatch, result)
		})
	}
}

func TestPasswordHashingRoundTrip(t *testing.T) {
	// Teste completo de ida e volta - como viagem redonda bem planejada ✈️
	passwords := []string{
		"senhaSimples",
		"Senha@Complexa123!",
		"🔒senhaComEmoji🔐",
		strings.Repeat("a", 70), // Within bcrypt limit
		"",
	}

	for _, password := range passwords {
		t.Run("round trip for: "+password, func(t *testing.T) {
			// Hash the password
			hash, err := HashPassword(password)
			assert.NoError(t, err)
			assert.NotEmpty(t, hash)

			// Verify the password matches the hash
			assert.True(t, CheckPasswordHash(password, hash))

			// Verify a different password doesn't match
			assert.False(t, CheckPasswordHash(password+"wrong", hash))
		})
	}
}

func TestGenerateAccessToken(t *testing.T) {
	// GenerateAccessToken: criando tokens como padeiro cria pães 🥖
	// Setup environment
	originalSecret := os.Getenv("ACCESS_TOKEN_SECRET")
	defer func() {
		if originalSecret != "" {
			os.Setenv("ACCESS_TOKEN_SECRET", originalSecret)
		} else {
			os.Unsetenv("ACCESS_TOKEN_SECRET")
		}
	}()

	testSecret := "test-access-secret-super-secure-key"
	os.Setenv("ACCESS_TOKEN_SECRET", testSecret)

	tests := []struct {
		name        string
		userID      string
		expectError bool
	}{
		{
			name:        "valid user ID - token fresquinho saindo do forno 🔥",
			userID:      uuid.New().String(),
			expectError: false,
		},
		{
			name:        "empty user ID - ID mais vazio que estádio em pandemia 🏟️",
			userID:      "",
			expectError: false, // JWT allows empty subject
		},
		{
			name:        "long user ID - ID épico como nome de personagem de fantasia 🐉",
			userID:      strings.Repeat("a", 1000),
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GenerateAccessToken(tt.userID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)

				// Verify token structure (should have 3 parts separated by dots)
				parts := strings.Split(token, ".")
				assert.Len(t, parts, 3, "JWT should have 3 parts")

				// Verify token can be validated
				parsedToken, claims, err := ValidateToken(token, testSecret)
				assert.NoError(t, err)
				assert.NotNil(t, parsedToken)
				assert.True(t, parsedToken.Valid)

				// Verify claims
				assert.Equal(t, tt.userID, claims["sub"])

				// Verify expiration is set correctly (should be ~15 minutes from now)
				exp, ok := claims["exp"].(float64)
				assert.True(t, ok, "Expiration should be a number")
				expTime := time.Unix(int64(exp), 0)
				assert.True(t, expTime.After(time.Now().Add(14*time.Minute)))
				assert.True(t, expTime.Before(time.Now().Add(16*time.Minute)))
			}
		})
	}
}

func TestGenerateRefreshToken(t *testing.T) {
	// GenerateRefreshToken: tokens de renovação como água da fonte da juventude ⛲
	originalSecret := os.Getenv("REFRESH_TOKEN_SECRET")
	defer func() {
		if originalSecret != "" {
			os.Setenv("REFRESH_TOKEN_SECRET", originalSecret)
		} else {
			os.Unsetenv("REFRESH_TOKEN_SECRET")
		}
	}()

	testSecret := "test-refresh-secret-super-secure-key"
	os.Setenv("REFRESH_TOKEN_SECRET", testSecret)

	tests := []struct {
		name        string
		userID      string
		expectError bool
	}{
		{
			name:        "valid user ID - refresh token fresquinho 🆕",
			userID:      uuid.New().String(),
			expectError: false,
		},
		{
			name:        "empty user ID - vazio mas funcional 🈳",
			userID:      "",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GenerateRefreshToken(tt.userID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)

				// Verify token structure
				parts := strings.Split(token, ".")
				assert.Len(t, parts, 3, "JWT should have 3 parts")

				// Verify token can be validated
				parsedToken, claims, err := ValidateToken(token, testSecret)
				assert.NoError(t, err)
				assert.NotNil(t, parsedToken)
				assert.True(t, parsedToken.Valid)

				// Verify claims
				assert.Equal(t, tt.userID, claims["sub"])

				// Verify expiration is set correctly (should be ~7 days from now)
				exp, ok := claims["exp"].(float64)
				assert.True(t, ok, "Expiration should be a number")
				expTime := time.Unix(int64(exp), 0)
				expected7Days := time.Now().Add(7 * 24 * time.Hour)
				assert.True(t, expTime.After(expected7Days.Add(-1*time.Hour)))
				assert.True(t, expTime.Before(expected7Days.Add(1*time.Hour)))
			}
		})
	}
}

func TestValidateToken(t *testing.T) {
	// ValidateToken: validador mais criterioso que fiscal da Receita 👨‍💼
	testSecret := "test-validation-secret"
	userID := uuid.New().String()

	// Create a valid token for testing
	validToken, err := func() (string, error) {
		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["sub"] = userID
		claims["exp"] = time.Now().Add(1 * time.Hour).Unix()
		return token.SignedString([]byte(testSecret))
	}()
	assert.NoError(t, err)

	// Create an expired token
	expiredToken, err := func() (string, error) {
		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["sub"] = userID
		claims["exp"] = time.Now().Add(-1 * time.Hour).Unix() // Expired 1 hour ago
		return token.SignedString([]byte(testSecret))
	}()
	assert.NoError(t, err)

	tests := []struct {
		name          string
		tokenString   string
		secret        string
		expectError   bool
		errorContains string
	}{
		{
			name:        "valid token - aprovado com distinção 🏆",
			tokenString: validToken,
			secret:      testSecret,
			expectError: false,
		},
		{
			name:          "expired token - mais vencido que leite esquecido 🥛",
			tokenString:   expiredToken,
			secret:        testSecret,
			expectError:   true,
			errorContains: "Token is expired",
		},
		{
			name:        "wrong secret - chave mais errada que GPS quebrado 🧭",
			tokenString: validToken,
			secret:      "wrong-secret",
			expectError: true,
		},
		{
			name:        "malformed token - token mais quebrado que promessa política 🏛️",
			tokenString: "not.a.valid.jwt.token",
			secret:      testSecret,
			expectError: true,
		},
		{
			name:        "empty token - vazio como alma penada 👻",
			tokenString: "",
			secret:      testSecret,
			expectError: true,
		},
		{
			name:        "invalid structure - estrutura mais confusa que código legado 🕸️",
			tokenString: "header.payload", // Missing signature
			secret:      testSecret,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, claims, err := ValidateToken(tt.tokenString, tt.secret)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, token)
				assert.Nil(t, claims)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, token)
				assert.True(t, token.Valid)
				assert.NotNil(t, claims)

				// Verify claims content
				assert.Equal(t, userID, claims["sub"])
				assert.Contains(t, claims, "exp")
			}
		})
	}
}

func TestTokensWithMissingSecrets(t *testing.T) {
	// Testando comportamento quando segredos estão ausentes
	// Porque nem sempre o ambiente coopera 🌪️

	// Save original values
	originalAccessSecret := os.Getenv("ACCESS_TOKEN_SECRET")
	originalRefreshSecret := os.Getenv("REFRESH_TOKEN_SECRET")

	// Clean up after test
	defer func() {
		if originalAccessSecret != "" {
			os.Setenv("ACCESS_TOKEN_SECRET", originalAccessSecret)
		} else {
			os.Unsetenv("ACCESS_TOKEN_SECRET")
		}
		if originalRefreshSecret != "" {
			os.Setenv("REFRESH_TOKEN_SECRET", originalRefreshSecret)
		} else {
			os.Unsetenv("REFRESH_TOKEN_SECRET")
		}
	}()

	userID := uuid.New().String()

	t.Run("access token with empty secret", func(t *testing.T) {
		os.Setenv("ACCESS_TOKEN_SECRET", "")

		token, err := GenerateAccessToken(userID)

		// Should still work with empty secret (though not secure)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("refresh token with empty secret", func(t *testing.T) {
		os.Setenv("REFRESH_TOKEN_SECRET", "")

		token, err := GenerateRefreshToken(userID)

		// Should still work with empty secret (though not secure)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})
}
