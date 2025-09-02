# Testes E2E para API de Autenticação

Este diretório contém os testes end-to-end (e2e) para a API de autenticação em Go.

## Como Executar

1. **Iniciar o banco de dados:**

   ```bash
   docker-compose up -d
   ```

2. **Executar os testes:**
   ```bash
   go test ./tests/e2e/ -v
   ```

## Endpoints Testados

- `POST /api/auth/register`: Testa o registro de novo usuário
- `POST /api/auth/login`: Testa o login e obtenção de tokens via cookies
- `GET /api/user`: Testa acesso ao perfil do usuário (requer autenticação)
- `POST /api/auth/refresh`: Testa renovação de tokens
- `POST /api/auth/logout`: Testa logout

## Estrutura dos Testes

- `TestMain`: Configuração global, incluindo conexão com banco de teste
- Testes individuais para cada endpoint
- Uso de `httptest.NewServer` para simular o servidor
- Gerenciamento de cookies para autenticação

## Notas

- Os testes usam dados únicos para evitar conflitos
- O banco de teste é configurado via variáveis de ambiente
- Certifique-se de que o PostgreSQL esteja rodando na porta 5435 (via docker-compose)
