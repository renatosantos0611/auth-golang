#!/bin/bash

echo "🔄 Migrando de MongoDB para PostgreSQL..."

# Para o container MongoDB se estiver rodando
echo "⏹️  Parando container MongoDB..."
docker compose down

# Remove o volume do MongoDB (opcional - descomente se quiser limpar)
# echo "🗑️  Removendo dados do MongoDB..."
# docker volume rm auth-golang_mongodb_data

# Inicia o container PostgreSQL
echo "🚀 Iniciando PostgreSQL..."
docker compose up -d postgres

# Aguarda o PostgreSQL inicializar
echo "⏳ Aguardando PostgreSQL inicializar..."
sleep 10

# Verifica se o banco está funcionando
echo "🔍 Verificando conexão com PostgreSQL..."
docker exec postgres psql -U admin -d auth_golang -c "SELECT 1;"

if [ $? -eq 0 ]; then
    echo "✅ PostgreSQL está funcionando!"
    echo "📊 Verificando tabelas criadas..."
    docker exec postgres psql -U admin -d auth_golang -c "\dt"
else
    echo "❌ Erro ao conectar com PostgreSQL"
    exit 1
fi

echo "🎉 Migração concluída!"
echo "💡 Lembre-se de:"
echo "   1. Atualizar seu arquivo .env com as variáveis do PostgreSQL"
echo "   2. Executar 'go mod tidy' para baixar as dependências"
echo "   3. Migrar os dados do MongoDB manualmente se necessário"
