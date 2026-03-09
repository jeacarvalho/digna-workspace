# Módulo Educativo do PDV - Exemplo de Integração

## 📋 Visão Geral

O módulo educativo do PDV fornece uma calculadora de preço justo que ajuda empreendedores solidários a valorar corretamente seu trabalho, seguindo os princípios da Economia Solidária.

## 🎯 Funcionalidades

1. **Calculadora Pedagógica:** Cálculo automático baseado em custo material + valor do trabalho
2. **Valoração do Trabalho (ITG 2002):** Conversão de minutos trabalhados em valor monetário
3. **Gráfico Visual:** Representação colorida da composição do preço
4. **Linguagem Coloquial:** Termos acessíveis, sem jargões contábeis

## 🏗️ Estrutura do Módulo

```
modules/pdv_ui/
├── internal/
│   ├── service/
│   │   ├── pricing_service.go        # Lógica matemática (100% int64)
│   │   └── pricing_service_test.go   # Testes unitários
│   └── handler/
│       ├── pricing_handler.go        # HTTP Handler HTMX
│       └── pricing_handler_test.go   # Testes do handler
└── templates/
    └── components/
        └── pricing_calculator.html   # Template Tailwind + HTMX
```

## 🔧 Como Integrar

### 1. Registrar o Handler no Servidor

```go
// No seu main.go ou router principal
pricingHandler, err := handler.NewPricingHandler()
if err != nil {
    log.Fatalf("Failed to create pricing handler: %v", err)
}

mux.HandleFunc("/pdv/pricing/calculate", pricingHandler.HandleCalculatePrice)
```

### 2. Incluir o Componente no PDV

```html
<!-- No template do PDV -->
<div class="container mx-auto px-4 py-8">
    <h1 class="text-2xl font-bold mb-6">Ponto de Venda - Economia Solidária</h1>
    
    <!-- Calculadora de Preço Justo -->
    <div hx-get="/pdv/pricing/calculate" hx-trigger="load">
        <!-- O componente será carregado aqui via HTMX -->
    </div>
    
    <!-- Resto do PDV... -->
</div>
```

### 3. Configurar Dependências HTMX

```html
<!-- No layout principal -->
<head>
    <script src="https://unpkg.com/htmx.org@1.9.10"></script>
    <script src="https://cdn.tailwindcss.com"></script>
</head>
```

## 📊 Exemplo de Cálculo

### Entrada:
- **Custo do Material:** R$ 10,00 (1000 centavos)
- **Tempo de Trabalho:** 60 minutos (1 hora)
- **Valor da Hora:** R$ 20,00 (2000 centavos/hora)

### Cálculo:
1. **Valor do Trabalho:** (60 × 2000) / 60 = 2000 centavos (R$ 20,00)
2. **Total antes do fundo:** 1000 + 2000 = 3000 centavos (R$ 30,00)
3. **Fundo da Cooperativa (5%):** 3000 × 5 / 100 = 150 centavos (R$ 1,50)
4. **Preço Sugerido:** 3000 + 150 = 3150 centavos (R$ 31,50)

### Gráfico Visual:
- 🟦 **Material:** 31.7% (R$ 10,00)
- 🟩 **Trabalho:** 63.5% (R$ 20,00)  
- 🟨 **Fundo:** 4.8% (R$ 1,50)

## ✅ Validações Implementadas

### Anti-Float (Regra de Ouro)
- ✅ Nenhum `float32` ou `float64` usado nos cálculos
- ✅ Todos os valores monetários em `int64` (centavos)
- ✅ Porcentagens calculadas com `int64` (×100 para precisão)

### Linguagem Pedagógica
- ✅ **Termos Usados:** "Custo do Material", "Seu Tempo", "Preço Justo", "Fundo da Cooperativa"
- ✅ **Termos Proibidos:** "Markup", "Lucro", "CPV", "Débito/Crédito", "ROI"
- ✅ **Tom:** Acolhedor, explicativo, focado na valorização do trabalho

### Acessibilidade
- ✅ Contraste visual adequado
- ✅ Gráfico colorido com legenda clara
- ✅ Textos explicativos passo a passo
- ✅ Ícones e dicas visuais

## 🧪 Testes

### Testes Unitários do Serviço
```bash
cd modules/pdv_ui
go test ./internal/service/... -v
```

### Testes do Handler
```bash
cd modules/pdv_ui  
go test ./internal/handler/... -v
```

### Validação Anti-Float
```bash
grep -r "float32\|float64" modules/pdv_ui/internal/ --include="*.go" | grep -v "test.go"
# Apenas para formatação visual no handler (permitido)
```

## 🚀 Integração Completa com PDV

### ✅ Implementado na Sprint Atual
1. **API Pública:** `pkg/pricing/api.go` - Wrapper público para integração entre módulos
2. **Integração no PDV Handler:** `ui_web/internal/handler/pdv_handler.go` - Adição do PricingCalculator
3. **Template PDV Atualizado:** `ui_web/templates/pdv.html` - Container HTMX para calculadora
4. **Testes de Integração:** `ui_web/integration_test.go` - Validação completa do fluxo

### Como Funciona a Integração:
1. **PDV Handler** cria uma instância de `PricingCalculator` via `pricing.NewPricingCalculator()`
2. **Rotas Registradas:** `/pdv/pricing/calculate` adicionada ao mux do PDV
3. **Template PDV** carrega calculadora dinamicamente via HTMX `hx-get`
4. **Cálculos em Tempo Real:** Usuário digita valores e vê resultados instantaneamente

### Estrutura da Integração:
```
modules/
├── pdv_ui/                          # Módulo educativo
│   ├── pkg/pricing/api.go           # API pública para integração
│   ├── internal/service/            # Lógica de cálculo (100% int64)
│   └── internal/handler/            # Handler interno (testes)
└── ui_web/                          # Interface web principal
    ├── internal/handler/pdv_handler.go  # PDV com calculadora integrada
    └── templates/pdv.html           # Template com container HTMX
```

### Próximas Melhorias
1. **Hook no Formulário de Produto:** Preencher automaticamente campos de preço
2. **Histórico de Cálculos:** Salvar cálculos para referência futura
3. **Valores Padrão:** Sugerir valor da hora baseado no histórico da cooperativa
4. **Exemplos Práticos:** Casos reais de diferentes tipos de produtos

## 📚 Princípios Pedagógicos

### 1. Primazia do Trabalho (ITG 2002)
> "O trabalho tem valor intrínseco e deve ser remunerado dignamente."

### 2. Transparência Algorítmica  
> "O usuário deve entender como o preço é calculado, não apenas aceitá-lo."

### 3. Educação Financeira Solidária
> "Ensinar a valorar, não a lucrar. Fortalecer, não explorar."

### 4. Tecnologia Social
> "Ferramentas simples que empoderam, não complexas que excluem."

## 🧪 Status da Integração

### ✅ Concluído
- [x] Módulo educativo implementado e testado (97.1% coverage)
- [x] API pública criada para integração entre módulos
- [x] Integração completa com PDV principal
- [x] Testes de integração validados
- [x] Servidor inicia sem erros
- [x] Linguagem pedagógica validada (termos proibidos ausentes)

### 📊 Métricas de Qualidade
- **Cobertura de Testes:** 97.1% (pdv_ui), 100% (integração)
- **Anti-Float Compliance:** 100% int64 nos cálculos
- **Performance:** Cálculos em tempo real via HTMX
- **Acessibilidade:** Interface responsiva com Tailwind CSS

### 🔍 Validações Realizadas
1. **Testes Unitários:** Todos os cenários de cálculo
2. **Testes de Integração:** Fluxo completo PDV + Calculadora
3. **Validação de Linguagem:** Termos pedagógicos vs jargões contábeis
4. **Build e Deploy:** Servidor compila e inicia sem erros

---
**Status:** ✅ Módulo implementado, integrado e testado  
**Próxima Sprint:** Melhorias de UX e casos de uso avançados