# 📁 Pasta HTTP - Testes de Endpoints

Bem-vindo à coleção de testes HTTP da sua API de autenticação! 🎯

Como um conjunto de ferramentas bem organizadas, cada arquivo tem sua função específica para testar diferentes aspectos da aplicação.

## 📋 Arquivos Disponíveis

### 🔧 Configuração

- **`_variables.http`** - Configurações e variáveis de referência

### 🛠️ Endpoints de Autenticação

- **`auth-register.http`** - Registro de novos usuários
- **`auth-login.http`** - Login e obtenção de tokens
- **`auth-refresh.http`** - Renovação de tokens de acesso
- **`auth-logout.http`** - Logout e limpeza de sessão

### 👤 Endpoints de Usuário

- **`user-profile.http`** - Obtenção de dados do usuário (protegido)

### 🎬 Testes Completos

- **`complete-flow.http`** - Fluxo completo da aplicação

## 🚀 Como Usar

### 1. **Configuração Inicial**

```bash
# Certifique-se de que sua aplicação está rodando
go run cmd/api/main.go

# A aplicação deve estar acessível em http://localhost:4000
```

### 2. **Testando Individualmente**

Abra qualquer arquivo `.http` e execute as requisições:

- Clique no ícone "Send Request" acima de cada requisição
- Ou use `Ctrl+Shift+P` → "Rest Client: Send Request"

### 3. **Testando o Fluxo Completo**

Para testar todo o sistema de uma vez:

1. Abra `complete-flow.http`
2. Execute as requisições **na ordem**, de cima para baixo
3. Aguarde cada resposta antes de executar a próxima

## 📝 Detalhes dos Testes

### 🔐 **auth-register.http**

- ✅ Registro com dados válidos
- ❌ Tentativa de registro duplicado
- ❌ Validações (senha curta, email inválido, campos vazios)

### 🔑 **auth-login.http**

- ✅ Login com credenciais válidas
- ❌ Email inexistente
- ❌ Senha incorreta
- ❌ Validações de formato

### 🔄 **auth-refresh.http**

- ✅ Renovação de token válido
- ❌ Token inexistente ou inválido
- ❌ Token expirado

### 👋 **auth-logout.http**

- ✅ Logout padrão
- ✅ Logout múltiplo (sempre funciona)

### 👤 **user-profile.http**

- ✅ Acesso com autenticação válida
- ❌ Acesso sem token
- ❌ Token inválido ou expirado

## 🎯 Resultados Esperados

### ✅ **Sucesso (Status 200/201)**

```json
{
  "message": "User registered successfully"
}
```

```json
{
  "success": true
}
```

```json
{
  "user": {
    "id": "uuid-do-usuario",
    "name": "Nome do Usuário",
    "email": "email@example.com",
    "role": "user"
  }
}
```

### ❌ **Erro (Status 400/401/409/500)**

```json
{
  "error": "Invalid request data"
}
```

```json
{
  "error": "Unauthorized"
}
```

```json
{
  "error": "User already exists"
}
```

## 🔧 Configuração Personalizada

Se sua aplicação roda em porta diferente, edite a variável `@baseUrl` em cada arquivo:

```http
@baseUrl = http://localhost:SUA_PORTA
```

## 🐛 Troubleshooting

### Problema: "Connection refused"

- ✅ Verifique se a aplicação está rodando
- ✅ Confirme a porta correta (padrão: 4000)
- ✅ Verifique se o PostgreSQL está ativo

### Problema: "Unauthorized" no endpoint protegido

- ✅ Execute o login primeiro
- ✅ Use a mesma sessão/cookies
- ✅ Verifique se o token não expirou (15 minutos)

### Problema: Erro de banco de dados

- ✅ Execute `docker compose up -d postgres`
- ✅ Verifique as variáveis de ambiente
- ✅ Confirme a conexão com o banco

## 💡 Dicas Profissionais

1. **Use o fluxo completo** para testar cenários reais
2. **Execute testes na ordem** para manter o estado correto
3. **Monitore os cookies** - eles são essenciais para autenticação
4. **Teste casos de erro** - eles são tão importantes quanto os sucessos
5. **Use dados diferentes** para cada teste para evitar conflitos

---

**Divirta-se testando! 🎉**

_Como diria um gopher: "Testing is caring!"_ 🐹
