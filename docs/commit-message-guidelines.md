# 📝 Diretrizes de Commit Messages - Conventional Commits

## 🎯 Padrão de Commit Messages

Este documento estabelece as diretrizes para escrita de commit messages seguindo o padrão **Conventional Commits** em português brasileiro. Como diria um gopher experiente: "Commit message ruim é como bug silencioso - só dá problema depois!" 🐛

## 📋 Formato Padrão

### 🏗️ Estrutura Base

```
<tipo>(<escopo>): <descrição>

[corpo opcional]

[rodapé opcional]
```

### 🎨 Exemplos Práticos

```bash
feat(auth): adiciona validação de JWT tokens
fix(database): corrige conexão com postgres em produção
docs(readme): atualiza instruções de instalação
test(user): adiciona testes unitários para criação de usuário
refactor(middleware): simplifica lógica de autenticação
style(handlers): corrige formatação e imports
perf(database): otimiza queries de usuários
build(docker): atualiza imagem base para go 1.21
ci(github): adiciona workflow de testes automatizados
```

## 🏷️ Tipos de Commit

### 📦 Tipos Principais

- **feat**: Nova funcionalidade para o usuário
- **fix**: Correção de bug
- **docs**: Mudanças na documentação
- **style**: Formatação, missing semi colons, etc (sem mudança de código)
- **refactor**: Refatoração de código que não adiciona feature nem corrige bug
- **test**: Adição ou correção de testes
- **chore**: Mudanças no processo de build ou ferramentas auxiliares

### 🔧 Tipos Secundários

- **perf**: Melhoria de performance
- **build**: Mudanças que afetam o sistema de build (gulp, webpack, etc)
- **ci**: Mudanças em arquivos de configuração de CI (Travis, Circle, etc)
- **revert**: Reverte um commit anterior

## 🎯 Escopos Sugeridos

### 🏗️ Componentes Principais

- **auth**: Sistema de autenticação
- **user**: Gerenciamento de usuários
- **database**: Configurações e operações de banco
- **middleware**: Middlewares HTTP
- **handlers**: Handlers HTTP
- **models**: Modelos de dados
- **tests**: Testes automatizados
- **docker**: Configurações Docker
- **config**: Configurações da aplicação

## ✅ Boas Práticas

### 📝 Descrição

- **Máximo 100 caracteres** na primeira linha
- **Use imperativo**: "adiciona" ao invés de "adicionado"
- **Seja específico**: "corrige validação de email" ao invés de "corrige bug"
- **Em português brasileiro**: Sempre! 🇧🇷
- **Sem ponto final**: Na descrição principal
- **Minúscula**: Primeira palavra após os dois pontos

### 🎨 Exemplos Bons vs Ruins

#### ✅ Bons Exemplos

```bash
feat(auth): implementa refresh token automático
fix(user): corrige validação de email inválido
docs(api): adiciona documentação dos endpoints
test(middleware): adiciona testes para middleware de CORS
```

#### ❌ Exemplos Ruins

```bash
fix: bug
update stuff
Fixed the authentication problem.
feat: new feature
```

## 🚀 Casos Especiais

### 🔄 Breaking Changes

Para mudanças que quebram compatibilidade, adicione `!` após o escopo:

```bash
feat(auth)!: remove suporte para JWT v1
refactor(database)!: altera schema de usuários
```

### 📋 Commits com Múltiplas Mudanças

Prefira commits atômicos, mas quando necessário:

```bash
feat(auth): implementa login e logout
- adiciona endpoint de login
- implementa middleware de logout
- atualiza documentação da API
```

## 🎯 Dicas do Mestre Gopher

1. **Commit pequeno é commit feliz** - Mantenha commits atômicos
2. **Teste antes de commitar** - Go vet, go test, go fmt
3. **Rebase quando necessário** - Mantenha o histórico limpo
4. **Use co-authored quando em pair** - Dê crédito ao parceiro

## 📚 Referências

- [Conventional Commits](https://www.conventionalcommits.org/)
- [Keep a Changelog](https://keepachangelog.com/)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)

---

_"Um commit bem escrito é como um gopher bem alimentado - sempre produtivo!"_ 🐹✨
