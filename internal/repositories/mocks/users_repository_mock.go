package mocks

import (
	"auth-golang/internal/models"
	"context"
)

// MockUsersRepository é o mock do repository para testes
// Como diria o Uncle Bob: "Mock it 'til you make it!" 🎭
type MockUsersRepository struct {
	// Funções que podem ser definidas para controlar o comportamento
	CreateFunc             func(ctx context.Context, user *models.User) error
	FindByIDFunc           func(ctx context.Context, id string) (*models.User, error)
	FindByEmailFunc        func(ctx context.Context, email string) (*models.User, error)
	FindByRefreshTokenFunc func(ctx context.Context, refreshToken string) (*models.User, error)
	UpdateFunc             func(ctx context.Context, user *models.User) (*models.User, error)

	// Contadores para verificar quantas vezes as funções foram chamadas
	CreateCallCount             int
	FindByIDCallCount           int
	FindByEmailCallCount        int
	FindByRefreshTokenCallCount int
	UpdateCallCount             int

	// Últimos argumentos passados para as funções
	LastCreateUser              *models.User
	LastFindByIDParam           string
	LastFindByEmailParam        string
	LastFindByRefreshTokenParam string
	LastUpdateUser              *models.User
}

func NewMockUsersRepository() *MockUsersRepository {
	return &MockUsersRepository{}
}

func (m *MockUsersRepository) Create(ctx context.Context, user *models.User) error {
	m.CreateCallCount++
	m.LastCreateUser = user

	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, user)
	}

	return nil
}

func (m *MockUsersRepository) FindByID(ctx context.Context, id string) (*models.User, error) {
	m.FindByIDCallCount++
	m.LastFindByIDParam = id

	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, id)
	}

	return nil, nil
}

func (m *MockUsersRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	m.FindByEmailCallCount++
	m.LastFindByEmailParam = email

	if m.FindByEmailFunc != nil {
		return m.FindByEmailFunc(ctx, email)
	}

	return nil, nil
}

func (m *MockUsersRepository) FindByRefreshToken(ctx context.Context, refreshToken string) (*models.User, error) {
	m.FindByRefreshTokenCallCount++
	m.LastFindByRefreshTokenParam = refreshToken

	if m.FindByRefreshTokenFunc != nil {
		return m.FindByRefreshTokenFunc(ctx, refreshToken)
	}

	return nil, nil
}

func (m *MockUsersRepository) Update(ctx context.Context, user *models.User) (*models.User, error) {
	m.UpdateCallCount++
	m.LastUpdateUser = user

	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, user)
	}

	return user, nil
}

// Reset limpa todos os contadores e dados armazenados
// Útil para testes que precisam de um estado limpo
func (m *MockUsersRepository) Reset() {
	m.CreateCallCount = 0
	m.FindByIDCallCount = 0
	m.FindByEmailCallCount = 0
	m.FindByRefreshTokenCallCount = 0
	m.UpdateCallCount = 0

	m.LastCreateUser = nil
	m.LastFindByIDParam = ""
	m.LastFindByEmailParam = ""
	m.LastFindByRefreshTokenParam = ""
	m.LastUpdateUser = nil

	m.CreateFunc = nil
	m.FindByIDFunc = nil
	m.FindByEmailFunc = nil
	m.FindByRefreshTokenFunc = nil
	m.UpdateFunc = nil
}
