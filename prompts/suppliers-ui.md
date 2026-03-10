# Implementação: Interface Web para Suppliers (Fornecedores)

**Tipo:** Feature  
**Módulo:** ui_web  
**Objetivo:** Implementar interface web para gestão de fornecedores (CRUD HTMX)  
**Decisões:** seguir padrão MemberHandler, usar cards layout, implementar testes unitários/integração e E2E Playwright

---

## 📋 Contexto

Os fornecedores já possuem:
- Domínio: `modules/supply/internal/domain/supplier.go`
- Repository: `modules/supply/internal/repository/` (métodos CRUD)
- API: `modules/supply/pkg/supply/api.go` (público, acessível)

## 🎯 Requisitos

### Funcional
- Listar fornecedores em cards
- Adicionar novo fornecedor (nome, contato)
- Editar fornecedor existente  
- Alternar status (ativo/inativo)
- Validação: nome obrigatório

### Técnico
- Handler: `SupplierHandler` estendendo `BaseHandler`
- Template: `suppliers_simple.html` com design "Soberania e Suor"
- Rotas HTMX: `/suppliers`, `/api/suppliers`, `/api/suppliers/{id}/toggle-status`
- Integração com `SupplyAPI` (já existe em `pkg/supply`)

### Qualidade
- Testes unitários (>90% coverage)
- Testes de integração com SQLite
- Smoke test: `./scripts/dev/smoke_test_new_feature.sh "Suppliers" "/suppliers"`
- Registro automático no `main.go`

## 🔗 Integrações

### Backend existente
```go
// SupplyAPI já tem:
RegisterSupplier(ctx, req) → SupplierResponse
GetSupplier(ctx, entityID, supplierID) → *Supplier
ListSuppliers(ctx, entityID) → []*Supplier
```

### Frontend patterns
- Seguir `MemberHandler` como referência
- Usar `TemplateManager` do `BaseHandler`
- Cards com ações HTMX inline
- Feedback visual com swaps

## 📅 Cronograma estimado
1. Dia 1: SupplierHandler + testes unitários
2. Dia 2: Template suppliers_simple.html  
3. Dia 3: Integração navegação + testes
4. Dia 4: Validação final + smoke test