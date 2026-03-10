# Testes End-to-End (E2E) com Playwright

## Visão Geral

Os testes E2E com Playwright validam o fluxo completo do sistema Digna, simulando ações de usuário real através da interface web. Seguem a sequência de 7 passos definida para validação completa do negócio.

## Instalação e Configuração

### Pré-requisitos
- Node.js 16+
- Servidor Digna rodando na porta 8090

### Instalação
```bash
npm install --save-dev @playwright/test
npx playwright install
```

### Configuração
- `playwright.config.js` - Configuração principal
- `tests/digna-e2e.spec.js` - Testes E2E
- Scripts npm disponíveis no `package.json`

## Sequência de Testes E2E

### 1. Login no Sistema
- Acessa `/login`
- Seleciona empresa `cafe_digna`
- Insere senha `cd0123`
- Valida redirecionamento para `/dashboard`

### 2. Criar Item de Estoque (se não existir)
- Navega para `/supply/stock`
- Verifica se existem itens
- Cria novo item se necessário
- Valida criação bem-sucedida

### 3. Criar Membro (se não existir)
- Navega para `/members`
- Verifica se existem membros
- Cria novo membro se necessário
- Valida criação bem-sucedida

### 4. Criar Fornecedor (se não existir)
- Navega para `/supply/suppliers`
- Verifica se existem fornecedores
- Cria novo fornecedor se necessário
- Valida criação bem-sucedida

### 5. Registrar Compra do Item
- Navega para `/supply/purchases`
- Cria nova compra
- Seleciona fornecedor e item
- Define quantidade e preço
- Valida registro da compra

### 6. Registrar Venda no PDV
- Navega para `/pdv`
- Seleciona item do estoque
- Define quantidade
- Finaliza venda com pagamento
- Valida recibo gerado

### 7. Validar Saldo e Horas Trabalhadas
- Navega para `/cash`
- Verifica saldo disponível
- Navega para `/social/clock`
- Registra horas trabalhadas
- Retorna ao dashboard para validação final

## Execução dos Testes

### Comandos Disponíveis
```bash
# Executar todos os testes E2E
npm run test:e2e

# Executar com interface gráfica
npm run test:e2e:ui

# Executar em modo debug
npm run test:e2e:debug

# Executar com navegador visível
npm run test:e2e:headed

# Executar apenas no Chrome
npm run test:e2e:chrome

# Abrir relatório de testes
npm run report
```

### Execução Manual
```bash
# Iniciar servidor Digna
cd modules/ui_web && go run main.go

# Em outro terminal, executar testes
npx playwright test

# Executar teste específico
npx playwright test tests/digna-e2e.spec.js --grep "Fluxo completo"
```

## Estrutura dos Testes

### Arquitetura
```
tests/
├── digna-e2e.spec.js     # Testes principais
└── (futuros testes)

playwright.config.js      # Configuração
package.json              # Scripts npm
```

### Funções Auxiliares
- `login(page)` - Realiza login
- `createStockItemIfNeeded(page)` - Cria item de estoque
- `createMemberIfNeeded(page)` - Cria membro
- `createSupplierIfNeeded(page)` - Cria fornecedor
- `registerPurchase(page)` - Registra compra
- `registerSaleInPDV(page)` - Registra venda
- `verifyCashBalanceAndWorkHours(page)` - Valida saldo e horas

## Dados de Teste

### Credenciais
```javascript
const TEST_ENTITY = 'cafe_digna';
const TEST_PASSWORD = 'cd0123';
```

### Dados de Exemplo
- Item: "Produto Teste E2E", unidade: "un", preço: R$ 25,00
- Membro: "Membro E2E Test", CPF: "11122233344"
- Fornecedor: "Fornecedor E2E Test", CNPJ: "99887766000155"
- Compra: 5 unidades a R$ 20,00 cada
- Venda: 1 unidade por R$ 30,00
- Horas trabalhadas: 4 horas

## Validações

### Validações de Interface
- URLs corretas após navegação
- Elementos visíveis na página
- Mensagens de sucesso/erro
- Dados exibidos corretamente em tabelas

### Validações de Negócio
- Fluxo completo de estoque (criação → compra → venda)
- Integridade financeira (saldo em caixa)
- Registro de horas trabalhadas
- Dashboard com estatísticas atualizadas

## Solução de Problemas

### Servidor Não Disponível
```bash
# Verificar se o servidor está rodando
curl http://localhost:8090/health

# Iniciar servidor
cd modules/ui_web && go run main.go
```

### Timeout nos Testes
- Aumentar timeout no `playwright.config.js`
- Verificar se elementos estão carregando corretamente
- Usar `page.waitForSelector()` para elementos dinâmicos

### Elementos Não Encontrados
- Verificar seletor correto com DevTools
- Usar `page.pause()` para debug
- Verificar se página carregou completamente

### Relatórios
```bash
# Gerar relatório HTML
npx playwright show-report

# Verificar screenshots em caso de falha
ls test-results/
```

## Integração com Workflow

### Pré-commit
Adicionar hook para executar testes E2E antes de marcar tarefas como completas.

### CI/CD
Configurar execução automática em pipeline de CI.

### Monitoramento
Usar relatórios do Playwright para identificar problemas recorrentes.

## Próximos Passos

1. **Cobertura Completa**: Adicionar testes para todos os módulos
2. **Dados de Teste**: Criar fixtures para dados de teste reutilizáveis
3. **Paralelização**: Executar testes em paralelo para velocidade
4. **Relatórios**: Integrar com ferramentas de relatório (Allure, etc.)
5. **CI/CD**: Configurar execução automática no GitHub Actions

## Referências

- [Documentação Playwright](https://playwright.dev/docs/intro)
- [Best Practices E2E](https://playwright.dev/docs/best-practices)
- [Testes com Auth](https://playwright.dev/docs/auth)