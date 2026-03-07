---
title: Escopo do Produto
status: implemented
version: 1.0
last_updated: 2026-03-07
---

# Escopo do Produto - Digna

---

## 1. Tipo de Produto

Infraestrutura contábil e institucional para economia solidária.

---

## 2. Capacidades Principais

### 2.1 PDV Operacional
Registro simplificado de vendas e compras.

**Funcionalidades:**
- Keyboard numérico para entrada de valores
- Seleção de produtos
- Geração automática de partidas dobradas

### 2.2 Ledger Contábil
Motor de partidas dobradas automático.

**Funcionalidades:**
- Validação soma zero
- Histórico de lançamentos
- Balancete e demonstrativos

### 2.3 Registro de Trabalho
Captura de horas de trabalho cooperativo (ITG 2002).

**Funcionalidades:**
- Cronômetro em tempo real
- Registro manual de minutos
- Conversão em capital social

### 2.4 Documentação Institucional
Geração automática de documentos.

**Funcionalidades:**
- Atas de Assembleia (Markdown)
- Relatórios de impacto social
- Dossiês de formalização

### 2.5 Sincronização Offline
Operação mesmo sem internet.

**Funcionalidades:**
- Delta tracking local
- Sync apenas com dados agregados
- Marketplace B2B

---

## 3. Capacidades Excluídas

O Digna **não** é:

- ERP completo
- Folha de pagamento empresarial
- Contabilidade fiscal complexa
- Gestão financeira bancária
- Sistema de NF-e/DF-e
- Folha de pagamento

---

## 4. Limites do Sistema

| Função | Status | Justificativa |
|--------|--------|---------------|
| Contabilidade Gerencial | ✅ Escopo | Core do sistema |
| ITG 2002 | ✅ Escopo | Requisito legal |
| CADSOL | ✅ Escopo | Integração pública |
| Fiscal (IRPJ/CSLL) | ❌ Fora | Depende de contador |
| RH/FP | ❌ Fora | Não é cooperativa |
| Bancário | ❌ Fora | Integração futura |
