---
applyTo: "**/*_e2e_test.go"
---

# Instruções para Testes E2E

## Visão Geral

Testes end-to-end (e2e) são projetados para verificar o fluxo completo de uma aplicação, simulando interações reais do usuário do início ao fim. Em Golang, os testes e2e garantem que todos os componentes (endpoints de API, interações com banco de dados, middleware, etc.) funcionem juntos perfeitamente. Este documento descreve as melhores práticas para escrever testes e2e em Golang, enfatizando princípios de clean code como legibilidade, manutenibilidade e separação de responsabilidades.

## Melhores Práticas

### 1. **Escolha o Framework Adequado**

- Use o pacote padrão `testing` para testes simples, combinado com `httptest` para testes e2e baseados em HTTP.
- Para testes mais complexos, no estilo BDD, considere usar [Ginkgo](https://onsi.github.io/ginkgo/) com [Gomega](https://onsi.github.io/gomega/) para asserções expressivas.
- Evite dependência excessiva de ferramentas externas, a menos que necessário; mantenha-se fiel ao Golang idiomático.

### 2. **Estrutura do Projeto**

- Organize os testes e2e em um diretório dedicado, por exemplo, `tests/e2e/` ou `e2e/`.
- Separe os arquivos de teste por funcionalidade ou módulo, por exemplo, `auth_e2e_test.go`, `user_e2e_test.go`.
- Mantenha auxiliares de teste e fixtures em um subdiretório `testutils/` para promover a reutilização.

### 3. **Nomeação e Organização de Testes**

- Use nomes descritivos para funções de teste que indiquem claramente o cenário, por exemplo, `TestFluxoDeRegistroELoginDoUsuario`.
- Agrupe testes relacionados usando subtestes com `t.Run()` para melhor organização e execução paralela.
- Siga o padrão: `Test[Funcionalidade][Ação][Resultado]`.

### 4. **Configuração e Limpeza**

- Use `TestMain` para configuração/limpeza global, se necessário (por exemplo, iniciar um servidor de teste).
- Para configuração por teste, use `t.Cleanup()` para garantir que os recursos sejam limpos após cada teste.
- Inicialize bancos de dados de teste ou use bancos em memória (por exemplo, SQLite em memória) para evitar efeitos colaterais.
- Exemplo:
  ```go
  func TestMain(m *testing.M) {
      // Configurar banco de dados de teste
      os.Exit(m.Run())
  }
  ```

### 5. **Gerenciamento de Dados**

- Use fixtures ou fábricas de teste para criar dados de teste em vez de valores codificados.
- Limpe os dados de teste após cada teste para evitar interferências entre testes.
- Evite usar dados de produção; gere dados sintéticos que cubram casos extremos.
- Implemente funções de semeadura de dados para ambientes de teste consistentes.

### 6. **Asserções e Verificação**

- Use asserções claras e específicas para verificar o comportamento esperado.
- Para testes HTTP, verifique códigos de status, cabeçalhos de resposta e corpos JSON.
- Utilize bibliotecas como `testify/assert` para asserções legíveis.
- Exemplo:
  ```go
  assert.Equal(t, http.StatusOK, resp.StatusCode)
  var user models.User
  json.Unmarshal(resp.Body, &user)
  assert.NotEmpty(t, user.ID)
  ```

### 7. **Tratamento de Dependências Externas**

- Faça mock ou stub de serviços externos (por exemplo, provedores de e-mail, gateways de pagamento) usando interfaces e injeção de dependência.
- Para testes de banco de dados, use transações para reverter alterações.
- Se usar Docker, inicie contêineres isolados para serviços como bancos de dados.

### 8. **Princípios de Clean Code**

- **Legibilidade**: Escreva testes que leiam como documentação. Use comentários com moderação, mas de forma eficaz.
- **DRY (Não Se Repita)**: Extraia código de configuração comum em funções auxiliares.
- **Responsabilidade Única**: Cada teste deve focar em um cenário.
- **Manutenibilidade**: Refatore testes conforme o código evolui; evite testes frágeis.
- **Tratamento de Erros**: Use `t.Fatal()` ou `t.Errorf()` para falhas, e garanta que os testes falhem rapidamente.
- Mantenha o código de teste tão limpo quanto o código de produção; aplique idiomas Go como nomes de variáveis curtos e retornos antecipados.

### 9. **Performance e Execução**

- Execute testes em paralelo sempre que possível usando `t.Parallel()`.
- Use tags de build (por exemplo, `// +build e2e`) para separar testes e2e de testes unitários.
- Monitore o tempo de execução dos testes e otimize testes lentos.

### 10. **Execução e Depuração de Testes E2E**

- Execute testes e2e com: `go test -tags=e2e ./tests/e2e/...`
- Use saída verbosa: `go test -v`
- Depure com logging ou anexando um depurador ao processo de teste.
- Integre com pipelines de CI/CD para testes e2e automatizados.

### 11. **Armadilhas Comuns a Evitar**

- Escrever testes muito acoplados a detalhes de implementação.
- Ignorar condições de corrida em testes concorrentes.
- Não lidar adequadamente com timeouts ou testes instáveis.
- Excesso de mocking, o que pode levar a falsos positivos.

### 12. **Exemplo de Estrutura de Teste E2E**

```go
package e2e

import (
    "net/http"
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestFluxoDeRegistroELoginDoUsuario(t *testing.T) {
    t.Parallel()

    // Configuração
    client := &http.Client{}
    baseURL := "http://localhost:8080"

    // Testar registro
    resp, err := client.Post(baseURL+"/register", "application/json", strings.NewReader(`{"email":"test@example.com","password":"password"}`))
    assert.NoError(t, err)
    assert.Equal(t, http.StatusCreated, resp.StatusCode)

    // Testar login
    resp, err = client.Post(baseURL+"/login", "application/json", strings.NewReader(`{"email":"test@example.com","password":"password"}`))
    assert.NoError(t, err)
    assert.Equal(t, http.StatusOK, resp.StatusCode)

    // Verificar resposta
    // ... asserções adicionais
}
```

Seguindo essas diretrizes, seus testes e2e serão robustos, manuteníveis e alinhados com as melhores práticas de Golang e princípios de clean code.
