***

```markdown
---
title: Product Backlog
status: proposed
version: 1.3
last_updated: 2026-03-09
---

# Product Backlog - Digna

---

## Funcionalidades Concluídas

### Core Contábil
- [x] PDV Operacional (RF-01)
- [x] Registro de Trabalho ITG 2002 (RF-02)
- [x] Motor de Reservas Obrigatórias (RF-03)
- [x] Dossiê de Formalização (RF-04)
- [x] Sincronização Offline-First (RF-05)
- [x] Intercooperação B2B (RF-06)
- [x] Gestão de Compras (RF-07) ✅ **SPRINT 13**
- [x] Controle de Estoque (RF-08) ✅ **SPRINT 13**
- [x] Gestão de Caixa (RF-09)

### Gestão de Membros (Sprint 10)
- [x] Cadastro de membros
- [x] Roles (COORDINATOR, MEMBER, ADVISOR)
- [x] Status (ACTIVE, INACTIVE)
- [x] Skills/habilidades

### Formalização (Sprint 11)
- [x] Transição automática DREAM → FORMALIZED
- [x] Simulador de formalização
- [x] CheckFormalizationCriteria()
- [x] AutoTransitionIfReady()

### Integrações Governamentais
- [x] Receita Federal (CNPJ, DARF) - Mock
- [x] MTE (RAIS, CAT, eSocial) - Mock
- [x] MDS (CadÚnico, Relatório Social) - Mock
- [x] IBGE (Pesquisas, PAM, CNAE) - Mock
- [x] SEFAZ (NFe, NFS-e, Manifesto) - Mock
- [x] BNDES (Linhas de Crédito, Simulação) - Mock
- [x] SEBRAE (Cursos, Consultoria) - Mock
- [x] Providentia (Sync, Marketplace) - Mock

### Surplus Calculator
- [x] CalculateSocialSurplus()
- [x] CalculateWithDeductions() - 10% Reserva Legal + 5% FATES
- [x] Rateio proporcional por minutos trabalhados
- [x] Tratamento de resíduos

### Testes
- [x] Testes unitários por módulo
- [x] Testes E2E Journey (jornada anual)
- [x] Testes de integrações governamentais

---

## Funcionalidades Futuras

### Alta Prioridade (Fase 1: Integração e Aliança Contábil)
- [x] **[NOVO] Painel do Contador Social (Accountant Dashboard):** Interface Web Multi-tenant em modo Read-Only para auditores voluntários/parceiros. ✅ **SPRINT 12 COMPLETE**
- [x] **[NOVO] Motor de Exportação Fiscal (SPED):** Mapeamento do `Core Lume` para formatos padrões exigidos pela Receita Federal e sistemas contábeis comerciais. ✅ **SPRINT 12 COMPLETE**
- [ ] **[NOVO] Plataforma de Mutirão:** Módulo de gerenciamento para o "Imposto de Renda Solidário" em parceria com Universidades/CRCs.
- [ ] Autenticação Gov.br (Realização do fluxo OAuth2)
- [ ] Exportação CADSOL completa
- [ ] Auditoria pública (Validação visual ITG 2002)
- [ ] Integrações reais (Substituição de Mocks por APIs do governo)

### Média Prioridade
- [ ] **[NOVO] Módulos educativos embutidos:** UI com auxílio visual na formação de preço no PDV (Custo vs. Hora trabalhada).
- [ ] Gestão de membros (UI)
- [x] Gestão de compras (RF-07) ✅ **SPRINT 13**
- [x] Gestão de fornecedores ✅ **SPRINT 13**
- [x] Controle de estoque (RF-08) ✅ **SPRINT 13**
- [ ] Gestão orçamentária (RF-10)
- [ ] Múltiplas moedas sociais (Preparação para Fase 2)

### Baixa Prioridade
- [ ] API GraphQL
- [ ] App mobile nativo (Substituindo o atual PWA, se necessário)
- [ ] Relatórios em PDF
```

***
