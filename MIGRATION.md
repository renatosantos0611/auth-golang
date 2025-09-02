# 🔄 Migração MongoDB → PostgreSQL

Este documento descreve a migração do banco de dados de MongoDB para PostgreSQL.

## 📋 Mudanças Realizadas

### 1. Dependências

- ❌ Removido: `go.mongodb.org/mongo-driver/v2`
- ✅ Adicionado: `github.com/lib/pq` (driver PostgreSQL)
- ✅ Adicionado: `github.com/google/uuid` (para UUIDs)

### 2. Modelo de Dados

- **ID**: `bson.ObjectID` → `uuid.UUID`
- **Tags**: `bson` → `db` (para SQL)
- **Campos**: Mantidos todos os campos originais

### 3. Banco de Dados

- **Arquivo**: `internal/database/database.go`
- **Conexão**: MongoDB → PostgreSQL com `database/sql`
- **Configuração**: Variáveis de ambiente atualizadas

### 4. Repositório

- **Arquivo**: `internal/repositories/users-repositories.go`
- **Queries**: BSON → SQL queries
- **Erros**: `mongo.ErrNoDocuments` → `sql.ErrNoRows`

### 5. Docker

- **Container**: MongoDB → PostgreSQL 16
- **Script de inicialização**: `init.sql` com schema da tabela

## 🚀 Como Migrar

### 1. Execute o script de migração:

```bash
./migrate.sh
```

### 2. Ou faça manualmente:

#### Pare o MongoDB:

```bash
docker-compose down
```

#### Configure as variáveis de ambiente:

```bash
cp .env.example .env
# Edite o .env com suas configurações
```

#### Inicie o PostgreSQL:

```bash
docker-compose up -d postgres
```

#### Instale dependências:

```bash
go mod tidy
```

#### Execute a aplicação:

```bash
go run cmd/api/main.go
```

## 📊 Schema da Tabela

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    username VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    refresh_token TEXT,
    role INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
```

## 🔧 Variáveis de Ambiente

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=admin
DB_PASSWORD=123
DB_NAME=auth_golang
```

## 📝 Migração de Dados

Se você tem dados no MongoDB que precisa migrar, você pode:

1. **Exportar do MongoDB:**

```bash
mongoexport --db auth-golang --collection users --out users.json
```

2. **Criar script de importação para PostgreSQL**
3. **Converter os ObjectIDs para UUIDs**

## ✅ Verificação

Para verificar se tudo está funcionando:

```bash
# Teste a conexão com o banco
docker exec postgres psql -U admin -d auth_golang -c "SELECT 1;"

# Veja as tabelas criadas
docker exec postgres psql -U admin -d auth_golang -c "\dt"

# Execute os testes
go test ./...
```

## 🔍 Principais Diferenças

| Aspecto         | MongoDB                 | PostgreSQL      |
| --------------- | ----------------------- | --------------- |
| ID              | ObjectID (24 chars hex) | UUID (36 chars) |
| Schema          | Schemaless              | Schema rígido   |
| Queries         | BSON/Aggregation        | SQL             |
| Relacionamentos | Embedded/References     | Foreign Keys    |
| Transações      | Limited                 | ACID completo   |

## 🐛 Solução de Problemas

### Erro de conexão:

- Verifique se o PostgreSQL está rodando
- Confirme as variáveis de ambiente
- Verifique se a porta 5432 não está em uso

### Erro de dependências:

```bash
go clean -modcache
go mod download
go mod tidy
```

### Erro de permissão no banco:

```bash
docker exec postgres psql -U admin -d auth_golang -c "GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO admin;"
```
