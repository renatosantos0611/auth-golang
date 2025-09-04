# 🔐 Auth Golang API

> **Uma API de autenticação em Go que é mais segura que o cofre do Fort Knox e mais rápida que um Gopher correndo atrás de uma cenoura!** 🐹💨

[![Go Version](https://img.shields.io/badge/Go-1.24.3-00ADD8?style=for-the-badge&logo=go)](https://golang.org/)
[![Gin Framework](https://img.shields.io/badge/Gin-Framework-00ADD8?style=for-the-badge)](https://gin-gonic.com/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-336791?style=for-the-badge&logo=postgresql)](https://www.postgresql.org/)
[![JWT](https://img.shields.io/badge/JWT-Tokens-000000?style=for-the-badge&logo=jsonwebtokens)](https://jwt.io/)
[![Docker](https://img.shields.io/badge/Docker-Compose-2496ED?style=for-the-badge&logo=docker)](https://www.docker.com/)

## 🎯 O Que É Este Projeto?

Esta é uma API RESTful de autenticação robusta construída em Go, que implementa as melhores práticas de segurança. Como um guarda-costas digital, ela protege seus endpoints com autenticação JWT, hash de senhas com bcrypt, e toda a segurança que você esperaria de uma aplicação profissional.

### 🌟 Por Que Este Projeto é Especial?

- **🔒 Segurança de Fort Knox**: Senhas hasheadas com bcrypt e tokens JWT seguros
- **⚡ Performance de Fórmula 1**: Gin framework para velocidade máxima
- **🧪 Testado Como Remédio**: Cobertura de testes robusta (83%+ nos componentes críticos)
- **🐳 Deploy Fácil**: Docker Compose para subir tudo com um comando
- **📚 Documentação Completa**: Testes HTTP prontos para usar
- **🏗️ Arquitetura Limpa**: Clean Architecture que deixaria o Uncle Bob orgulhoso

## 🚀 Features Principais

### 🔐 Autenticação Completa

- ✅ **Registro de usuários** com validação de dados
- ✅ **Login seguro** com hash de senhas
- ✅ **Tokens JWT** (Access + Refresh Token)
- ✅ **Logout** com invalidação de tokens
- ✅ **Refresh de tokens** automático
- ✅ **Middleware de autenticação** para rotas protegidas

### 🛡️ Segurança

- ✅ **Bcrypt** para hash de senhas
- ✅ **JWT tokens** com expiração
- ✅ **CORS** configurado
- ✅ **Validação de dados** robusta
- ✅ **UUIDs** para identificação única

### 🗄️ Banco de Dados

- ✅ **PostgreSQL** como database principal
- ✅ **Migrations** automatizadas
- ✅ **Repository pattern** para abstração

## 📋 Pré-requisitos

Antes de começar esta jornada épica, certifique-se de ter:

- **Go 1.24.3+** (o mais novo é sempre melhor!)
- **Docker & Docker Compose** (para subir o PostgreSQL sem dor de cabeça)
- **Git** (obviamente! 😄)

## 🎬 Quick Start

### 1. Clone o Repositório

```bash
git clone https://github.com/renatosantos0611/auth-golang.git
cd auth-golang
```

### 2. Suba o Database

```bash
# Vai subir um PostgreSQL fresquinho na porta 5435
docker-compose up -d postgres
```

### 3. Configure as Variáveis de Ambiente

```bash
# Crie um arquivo .env na raiz do projeto
cp .env.example .env

# Edite conforme necessário (as configurações padrão já funcionam!)
```

### 4. Execute as Migrations

```bash
# Executa o script que prepara seu database
chmod +x migrate.sh
./migrate.sh
```

### 5. Instale as Dependências

```bash
go mod download
```

### 6. Execute a Aplicação

```bash
go run cmd/api/main.go
```

🎉 **Pronto!** Sua API está rodando em `http://localhost:8080`

## 🧪 Testando a API

### Testes Automáticos (Recomendado!)

```bash
# Execute todos os testes
go test ./...

# Com cobertura detalhada
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Testes E2E (requer a API rodando)
go test ./tests/e2e/...
```

### Testes Manuais (Para os Aventureiros!)

Temos uma coleção completa de arquivos `.http` na pasta `http/` que é mais organizada que a biblioteca de Alexandria:

```bash
# 1. Registre um usuário
# Use: http/auth-register.http

# 2. Faça login
# Use: http/auth-login.http

# 3. Acesse rotas protegidas
# Use: http/user-profile.http

# 4. Teste o refresh token
# Use: http/auth-refresh.http

# 5. Faça logout
# Use: http/auth-logout.http

# 6. Ou execute o fluxo completo
# Use: http/complete-flow.http
```

## 📐 Arquitetura

Esta aplicação segue os princípios da **Clean Architecture**, organizando o código como uma cidade bem planejada:

```
auth-golang/
├── 🎯 cmd/api/              # Ponto de entrada da aplicação
├── 🏗️ internal/             # Código interno da aplicação
│   ├── 🔐 auth/             # Lógica de autenticação e tokens
│   ├── 🗄️ database/         # Configuração do banco de dados
│   ├── 🛡️ middleware/       # Middlewares (autenticação, etc.)
│   ├── 📦 models/           # Estruturas de dados
│   ├── 🏪 repositories/     # Acesso aos dados
│   └── 🌐 server/           # Handlers HTTP e configuração do servidor
├── 🧪 tests/e2e/           # Testes end-to-end
├── 📡 http/                # Arquivos de teste da API
└── 🐳 docker-compose.yml   # Configuração do PostgreSQL
```

## 🔌 Endpoints da API

### 🚪 Autenticação (Públicos)

- `POST /auth/register` - Cadastro de novos usuários
- `POST /auth/login` - Login e obtenção de tokens
- `POST /auth/refresh` - Renovação do access token
- `POST /auth/logout` - Logout e invalidação de tokens

### 👤 Usuário (Protegidos)

- `GET /user/profile` - Dados do usuário autenticado

## 🧪 Cobertura de Testes

Nossa aplicação é mais testada que um carro de Fórmula 1 antes da corrida:

- **📦 Models**: **100%** de cobertura (perfeição absoluta!)
- **🔐 Auth**: **83.3%** de cobertura (quase perfeito!)
- **🌐 Handlers**: **75%** de cobertura (muito bem testado!)
- **🛡️ Middleware**: **100%** de cobertura (blindado!)

## 🛠️ Tecnologias Utilizadas

### Core

- **[Go 1.24.3](https://golang.org/)** - A linguagem que faz os desenvolvedores sorrirem
- **[Gin](https://gin-gonic.com/)** - Framework web mais rápido que a luz
- **[PostgreSQL 16](https://www.postgresql.org/)** - Database mais confiável que um amigo de infância

### Autenticação & Segurança

- **[JWT-Go](https://github.com/golang-jwt/jwt)** - Para tokens que funcionam
- **[bcrypt](https://golang.org/x/crypto/bcrypt)** - Hash de senhas à prova de hackers
- **[UUID](https://github.com/google/uuid)** - IDs únicos como impressões digitais

### Infraestrutura & DevOps

- **[Docker](https://www.docker.com/)** - Containerização sem dor de cabeça
- **[GoDotEnv](https://github.com/joho/godotenv)** - Variáveis de ambiente organizadas

### Testes & Qualidade

- **[Testify](https://github.com/stretchr/testify)** - Testes que realmente testam
- **Mocks** - Para testes isolados e confiáveis

## 🤝 Como Contribuir

Quer fazer parte desta jornada épica? Contribuições são mais bem-vindas que pizza numa sexta-feira!

1. **Fork** o projeto
2. **Crie** uma branch para sua feature (`git checkout -b feature/nova-feature-incrivel`)
3. **Commit** suas mudanças (`git commit -m 'feat: adiciona feature incrível'`)
4. **Push** para a branch (`git push origin feature/nova-feature-incrivel`)
5. **Abra** um Pull Request

### 📝 Padrões de Commit

Usamos **Conventional Commits** em português:

- `feat:` para novas funcionalidades
- `fix:` para correções de bugs
- `docs:` para documentação
- `test:` para testes
- `refactor:` para refatorações

## 📄 Licença

Este projeto está sob a licença MIT. Veja o arquivo [LICENSE](LICENSE) para mais detalhes.

## 👨‍💻 Autor

**Renato Santos** - [@renatosantos0611](https://github.com/renatosantos0611)

---

<div align="center">

**⭐ Se este projeto te ajudou, considere dar uma estrela! ⭐**

_Feito com ❤️ e muito ☕ em Go_

</div>
