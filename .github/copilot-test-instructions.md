# Instruções do GitHub Copilot para Projeto Go - Autenticação

## 🎯 Objetivo do Projeto

Este é um projeto de API de autenticação em Go que implementa funcionalidades de registro, login e gerenciamento de usuários usando JWT tokens, Gin framework e MongoDB.

## 📋 Boas Práticas Go e Clean Code

### 🏗️ Estrutura de Projeto

- **Siga a estrutura Standard Go Project Layout**
- **cmd/**: Pontos de entrada da aplicação
- **internal/**: Código privado da aplicação
- **pkg/**: Bibliotecas que podem ser usadas por aplicações externas
- **Mantenha pacotes pequenos e focados em uma responsabilidade**

### 🔤 Convenções de Nomenclatura

```go
// ✅ BOM: Nomes claros e expressivos
type UserRepository interface {
    CreateUser(ctx context.Context, user *User) error
    GetUserByEmail(ctx context.Context, email string) (*User, error)
}

// ❌ RUIM: Nomes genéricos ou abreviados
type UsrRepo interface {
    Create(u *User) error
    Get(email string) *User
}
```

### 📦 Organização de Pacotes

- **Um pacote por diretório**
- **Nomes de pacotes em minúsculas, sem underscores**
- **Evite pacotes "util" ou "common"**
- **Prefira nomes substantivos para pacotes**

### 🔐 Gerenciamento de Erros

```go
// ✅ BOM: Tratamento explícito de erros
func (r *UserRepository) CreateUser(ctx context.Context, user *User) error {
    if user == nil {
        return errors.New("user cannot be nil")
    }

    result, err := r.collection.InsertOne(ctx, user)
    if err != nil {
        return fmt.Errorf("failed to create user: %w", err)
    }

    return nil
}

// ❌ RUIM: Ignorar erros
func (r *UserRepository) CreateUser(ctx context.Context, user *User) {
    r.collection.InsertOne(ctx, user) // Ignora erro
}
```

### 🎭 Interfaces

```go
// ✅ BOM: Interfaces pequenas e focadas
type UserCreator interface {
    CreateUser(ctx context.Context, user *User) error
}

type UserFinder interface {
    GetUserByEmail(ctx context.Context, email string) (*User, error)
}

// ✅ BOM: Composição de interfaces
type UserRepository interface {
    UserCreator
    UserFinder
}
```

### 🏛️ Dependency Injection

```go
// ✅ BOM: Injeção de dependência via construtor
type AuthHandler struct {
    userRepo UserRepository
    tokenSvc TokenService
    logger   Logger
}

func NewAuthHandler(userRepo UserRepository, tokenSvc TokenService, logger Logger) *AuthHandler {
    return &AuthHandler{
        userRepo: userRepo,
        tokenSvc: tokenSvc,
        logger:   logger,
    }
}
```

### 🔄 Context Usage

```go
// ✅ BOM: Sempre passe context como primeiro parâmetro
func (s *UserService) GetUser(ctx context.Context, userID string) (*User, error) {
    return s.repo.GetUserByID(ctx, userID)
}

// ✅ BOM: Use context para timeouts e cancelamento
func (r *UserRepository) GetUserByID(ctx context.Context, userID string) (*User, error) {
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()

    // operação com timeout
}
```

### 🔒 Validação e Sanitização

```go
// ✅ BOM: Validação de entrada
type CreateUserRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=8"`
    Name     string `json:"name" binding:"required,min=2,max=100"`
}

func (r *CreateUserRequest) Validate() error {
    if !isValidEmail(r.Email) {
        return errors.New("invalid email format")
    }

    if len(r.Password) < 8 {
        return errors.New("password must be at least 8 characters")
    }

    return nil
}
```

### 🔐 Segurança

```go
// ✅ BOM: Hash de senhas
func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(bytes), err
}

// ✅ BOM: Verificação de senha
func CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}

// ✅ BOM: JWT com claims estruturados
type JWTClaims struct {
    UserID string `json:"user_id"`
    Email  string `json:"email"`
    jwt.RegisteredClaims
}
```

### 📝 Logging

```go
// ✅ BOM: Logging estruturado
func (h *AuthHandler) Login(c *gin.Context) {
    h.logger.Info("login attempt",
        "ip", c.ClientIP(),
        "user_agent", c.GetHeader("User-Agent"),
    )

    // lógica de login

    h.logger.Info("login successful",
        "user_id", user.ID,
        "email", user.Email,
    )
}
```

### 🧪 Testes

```go
// ✅ BOM: Testes com table-driven tests
func TestUserValidation(t *testing.T) {
    tests := []struct {
        name    string
        user    User
        wantErr bool
    }{
        {
            name: "valid user",
            user: User{Email: "test@example.com", Password: "password123"},
            wantErr: false,
        },
        {
            name: "invalid email",
            user: User{Email: "invalid-email", Password: "password123"},
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.user.Validate()
            if (err != nil) != tt.wantErr {
                t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### 🔧 Configuração

```go
// ✅ BOM: Configuração estruturada
type Config struct {
    Server   ServerConfig   `yaml:"server"`
    Database DatabaseConfig `yaml:"database"`
    JWT      JWTConfig      `yaml:"jwt"`
}

type ServerConfig struct {
    Port         int           `yaml:"port" env:"SERVER_PORT" default:"8080"`
    ReadTimeout  time.Duration `yaml:"read_timeout" env:"SERVER_READ_TIMEOUT" default:"10s"`
    WriteTimeout time.Duration `yaml:"write_timeout" env:"SERVER_WRITE_TIMEOUT" default:"10s"`
}
```

## 🎯 Instruções Específicas para o Copilot

### Para Handlers HTTP:

- Sempre valide a entrada usando binding tags do Gin
- Retorne códigos de status HTTP apropriados
- Implemente logging para auditoria
- Use middleware para autenticação/autorização

### Para Repositories:

- Sempre use context.Context como primeiro parâmetro
- Implemente timeouts apropriados
- Trate erros de forma específica (não encontrado vs erro do banco)
- Use transações quando necessário

### Para Services:

- Mantenha a lógica de negócio separada dos handlers
- Implemente validações de negócio
- Use interfaces para desacoplamento
- Documente comportamentos complexos

### Para Models:

- Use tags apropriadas para JSON, BSON, validation
- Implemente métodos de validação
- Mantenha modelos simples e focados
- Use ponteiros apenas quando necessário

### Para Middleware:

- Mantenha middleware leve e focado
- Propague o context corretamente
- Implemente logging adequado
- Trate erros de forma consistente

## 🚀 Comandos Úteis

```bash
# Executar testes
go test ./...

# Executar com coverage
go test -cover ./...

# Lint
golangci-lint run

# Formatar código
go fmt ./...

# Verificar módulos
go mod tidy

# Executar aplicação
go run cmd/api/main.go
```

## 📚 Recursos Adicionais

- [Effective Go](https://golang.org/doc/effective_go.html)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Standard Go Project Layout](https://github.com/golang-standards/project-layout)
- [Clean Architecture em Go](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
