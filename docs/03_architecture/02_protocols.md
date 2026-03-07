---
title: Protocolos Técnicos
status: implemented
version: 1.0
last_updated: 2026-03-07
---

# Protocolos Técnicos - Digna

---

## 1. Protocolo de Sincronização

### 1.1 Modelo

**Estratégia:** Delta-based synchronization

O sistema detecta alterações desde a última sincronização e transmite apenas os deltas, não os dados completos.

### 1.2 Estrutura do Pacote de Sync

```json
{
  "entity_id": "cooperativa_mel",
  "timestamp": 1772856840,
  "chain_digest": "d51e6eb402a6984e",
  "signature": "f802343da66e8396",
  "aggregated_data": {
    "total_sales": 7500,
    "total_work_hours": 12,
    "total_members": 2,
    "legal_status": "DREAM",
    "decision_count": 0
  },
  "delta_count": 3
}
```

### 1.3 Processo de Sincronização

```
1. DETECT: Query deltas desde last_sync_at
   - entries: alterações em lançamentos
   - work_logs: novos registros de trabalho
   - decisions_log: novas decisões

2. AGGREGATE: Calcular métricas agregadas
   - Soma de vendas (total_sales)
   - Soma de horas (total_work_hours)
   - Contagem de membros (total_members)
   - Status atual (legal_status)

3. HASH: Gerar chain digest
   - SHA256 da cadeia contábil atual
   - Inclui todos os hashes de decisões

4. SIGN: Assinar pacote
   - Usar entity_id como chave
   - Gera signature para autenticidade

5. TRANSMIT: Enviar para agregador
   - JSON ~400 bytes
   - Apenas dados agregados
```

### 1.4 Privacidade - Campos Incluídos vs Protegidos

| Campo | Incluído | Descrição |
|-------|----------|-----------|
| entity_id | ✅ | ID da entidade |
| total_sales | ✅ | Total vendas (int64) |
| total_work_hours | ✅ | Total horas |
| total_members | ✅ | Quantidade sócios |
| legal_status | ✅ | DREAM ou FORMALIZED |
| chain_digest | ✅ | Hash de integridade |
| signature | ✅ | Assinatura digital |
| member_id | ❌ | Dados sensíveis protegidos |
| entry_details | ❌ | Transações detalhadas |
| posting_id | ❌ | IDs internos |

---

## 2. Modelo de Segurança

### 2.1 Isolamento de Dados

Cada entidade possui banco próprio:
- Path: `data/entities/{entity_id}.db`
- Isolamento físico total
- Sem acesso cruzado entre tenants

### 2.2 Integridade

**Hash SHA256 para auditoria:**
- Cada decisão gera hash do conteúdo
- Cada bloco contábil gera chain digest
- Imutabilidade garantida por design

### 2.3 Transporte

**Pacotes assinados digitalmente:**
- Assinatura com EntityID
- Verificação de integridade
- Non-repudiation

### 2.4 Matriz de Ameaças e Mitigações

| Ameaça | Mitigação |
|--------|-----------|
| Acesso não autorizado | Isolamento por arquivo |
| Alteração de dados | Hash SHA256 + logging |
| Interceptação | TLS em transporte |
| Replay attack | Timestamp + nonce |

---

## 3. Protocolo Econômico

### 3.1 Introdução

O **Economic Protocol do Digna** define as regras econômicas fundamentais que governam a operação dos Empreendimentos de Economia Solidária (EES) dentro do sistema.

Este protocolo não descreve a implementação técnica, mas os **princípios econômicos e contábeis** que orientam seu funcionamento.

### 3.2 Princípios Fundamentais

#### 3.2.1 Primazia do Trabalho

O trabalho humano é reconhecido como a principal fonte de valor econômico.

O protocolo reconhece:
- Horas de trabalho
- Trabalho voluntário
- Contribuição coletiva

como formas legítimas de capital social.

#### 3.2.2 Autogestão

Toda decisão econômica relevante deve ser tomada coletivamente.

O sistema registra:
- Decisões de assembleia
- Distribuição de sobras
- Regras internas do grupo

#### 3.2.3 Transparência

Todas as operações econômicas são registradas em ledger verificável.

Cada operação gera:
- Registro contábil
- Hash de auditoria
- Histórico imutável

#### 3.2.4 Soberania de Dados

Cada empreendimento mantém controle sobre seus próprios dados.

O protocolo exige:
- Banco de dados próprio
- Exportação a qualquer momento
- Sem dependência centralizada

### 3.3 Unidades de Valor

#### 3.3.1 Moeda Nacional

**Real (R$)**

Utilizado para:
- Vendas
- Compras
- Contabilidade financeira

#### 3.3.2 Trabalho

O tempo de trabalho é tratado como capital social.

**Unidade:** minutos de trabalho

Utilizado para:
- Cálculo de participação
- Distribuição de sobras

#### 3.3.3 Bens Substantivos (Futuro)

O protocolo poderá suportar:
- Sementes
- Animais
- Bens produtivos

### 3.4 Registro de Operações

Toda operação econômica é registrada no ledger contábil.

**Tipos principais:**
- Venda
- Compra
- Trabalho
- Decisão coletiva

**Regra obrigatória:**

```
soma dos débitos = soma dos créditos (partidas dobradas)
```

### 3.5 Trabalho Cooperativo (ITG 2002)

O sistema registra trabalho coletivo conforme a ITG 2002.

Cada registro contém:
- Membro
- Atividade
- Minutos trabalhados
- Data

Este registro constitui **capital social de trabalho**.

### 3.6 Distribuição de Sobras

Ao final de um período econômico, o sistema calcula o excedente financeiro.

**Antes da distribuição:**

```
10% Reserva Legal
5% FATES
```

**Após segregação, o excedente pode ser distribuído.**

#### Fórmula básica de distribuição

Distribuição proporcional baseada em trabalho:

```
participação = horas_do_membro / horas_totais
valor_distribuído = participação × excedente
```

### 3.7 Registro de Decisões

Decisões coletivas são registradas como eventos institucionais.

Exemplos:
- Aprovação de estatuto
- Eleição de coordenação
- Definição de regras de rateio

Cada decisão gera:
- Registro
- Hash criptográfico
- Documento institucional

### 3.8 Formalização

O protocolo permite transição progressiva de grupo informal para entidade formal.

**Estados:**

```
DREAM → FORMALIZED
```

A formalização depende de:
- Registro de decisões (mínimo 3)
- Existência de membros
- Histórico econômico mínimo

### 3.9 Intercooperação

O protocolo permite interação econômica entre empreendimentos.

Possibilidades:
- Troca de produtos
- Cooperação produtiva
- Formação de cadeias solidárias

### 3.10 Privacidade Econômica

O protocolo estabelece que apenas dados agregados podem ser compartilhados.

**Dados protegidos:**
- Identidade de membros
- Detalhes de transações
- Dados sensíveis

**Dados compartilháveis:**
- Volume econômico
- Número de membros
- Status institucional

### 3.11 Evolução do Protocolo

O Economic Protocol pode evoluir através de decisões comunitárias.

Mudanças devem ser aprovadas pelo PMC da Fundação Providentia.

Cada alteração deve:
- Manter compatibilidade contábil
- Preservar soberania dos dados
- Respeitar princípios da economia solidária

### 3.12 Escopo

Este protocolo define apenas:
- Regras econômicas
- Princípios institucionais
- Lógica distributiva

A implementação técnica está descrita em `03_architecture/01_system.md`.

---

## 4. Referências Externas

- **ITG 2002** - Norma ITG do Conselho Federal de Contabilidade
- **Lei nº 15.068/2024** - Lei Paul Singer
- **CADSOL** - Cadastro Nacional de Economia Solidária
