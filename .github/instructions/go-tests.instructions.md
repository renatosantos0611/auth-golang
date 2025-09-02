---
applyTo: "**/*_test.go"
---

# Instruções para Testes Unitários em Go - 100% Coverage

## 🎯 Objetivo

Criar testes unitários robustos, legíveis e com cobertura completa seguindo as melhores práticas Go.

## 📋 Regras Essenciais

### 🏗️ Estrutura dos Testes

- **Table-driven tests** sempre que possível
- Nomes descritivos: `TestFunctionName_Scenario_ExpectedResult`
- Use subtests (`t.Run`) para organizar casos
- Apenas testes unitários (não integração/e2e)

### 📊 Cobertura 100%

- **Teste todos os caminhos**: happy path, edge cases, error paths
- **Cenários obrigatórios**:
  - Entrada válida (caso normal)
  - Entrada inválida (validação)
  - Erros de dependências externas
  - Casos limite (nil, vazio, máximo)
  - Condições de erro específicas

### 🧪 Padrão Table-Driven

```go
func TestFunction(t *testing.T) {
    tests := []struct {
        name    string
        input   Input
        want    Output
        wantErr bool
    }{
        {
            name: "success case",
            input: Input{...},
            want: Output{...},
            wantErr: false,
        },
        {
            name: "error case",
            input: Input{...},
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := Function(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
                t.Errorf("got = %v, want %v", got, tt.want)
            }
        })
    }
}
```

### 🔧 Mocks e Stubs

- Use interfaces para dependências
- Implemente mocks simples ou use `testify/mock`
- Isole todas as dependências externas
- Teste diferentes comportamentos do mock

### 📏 Boas Práticas

- **Uma verificação principal por teste**
- Use `t.Helper()` em funções auxiliares
- Evite goroutines desnecessárias
- Sem prints (`fmt.Println`), apenas `t.Errorf`
- Variáveis com nomes claros
- Setup e cleanup organizados

### 🎯 Cenários para 100% Coverage

1. **Funções públicas**: todos os caminhos
2. **Handlers HTTP**: status codes, validações, erros
3. **Repositories**: CRUD operations, erros de DB
4. **Services**: lógica de negócio, validações
5. **Middleware**: autorização, erros, headers
6. **Models**: validações, métodos

### 🚫 O que NÃO fazer

- Testes dependentes de ordem
- Testes que dependem de estado global
- Mocks excessivamente complexos
- Testes que testam implementação, não comportamento
