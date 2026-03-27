title: Protocolos Técnicos - Ecossistema Digna
status: implemented
version: 2.0
last_updated: 2026-03-27
---

# Protocolos Técnicos - Ecossistema Digna

> **Nota:** Este documento reflete os protocolos integrados do Ecossistema Digna (PDF v1.0), preservando toda a infraestrutura validada nas Sprints 1-16, incorporando os novos protocolos para os Módulos 2, 3 e 4, e incluindo o Sistema de Ajuda Educativa (RF-30) decidido na sessão de 27/03/2026.

---

## 1. Protocolo de Sincronização [ATUALIZADO]

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
     "decision_count": 0,
     "eligibility_complete": false,
     "public_profile_published": false,
     "credit_matches_count": 0
  },
   "delta_count": 3,
   "module_versions": {
     "erp": "1.0",
     "indicators": "1.0",
     "portal": "1.0",
     "rede": "1.0",
     "help": "1.0"
   }
}
```

### 1.3 Processo de Sincronização

1. **DETECT:** Query deltas desde last_sync_at
   - entries: alterações em lançamentos
   - work_logs: novos registros de trabalho
   - decisions_log: novas decisões
   - fiscal_exports: novos lotes extraídos pelo Contador Social
   - eligibility_profile: atualizações no perfil de elegibilidade [NOVO - PDF v1.0]
   - public_profile: publicações na Rede Digna [NOVO - PDF v1.0]
   - need_posts: necessidades publicadas na Rede [NOVO - PDF v1.0]
   - help_topics: novos tópicos de ajuda [NOVO - RF-30, Sessão 27/03/2026]

2. **AGGREGATE:** Calcular métricas agregadas
   - Soma de vendas (total_sales)
   - Soma de horas (total_work_hours)
   - Contagem de membros (total_members)
   - Status atual (legal_status)
   - Completude do perfil de elegibilidade (eligibility_complete) [NOVO]
   - Perfil público publicado (public_profile_published) [NOVO]
   - Quantidade de matches de crédito (credit_matches_count) [NOVO]

3. **HASH:** Gerar chain digest
   - SHA256 da cadeia contábil atual
   - Inclui todos os hashes de decisões
   - Inclui hash do perfil de elegibilidade [NOVO]

4. **SIGN:** Assinar pacote
   - Usar entity_id como chave
   - Gera signature para autenticidade

5. **TRANSMIT:** Enviar para agregador
   - JSON ~500 bytes (expandido para novos módulos)
   - Apenas dados agregados
   - Dados sensíveis nunca transmitidos

### 1.4 Privacidade - Campos Incluídos vs Protegidos [ATUALIZADO]

| Campo | Incluído | Descrição |
|-------|----------|-----------|
| entity_id | ✅ | ID da entidade |
| total_sales | ✅ | Total vendas (int64) |
| total_work_hours | ✅ | Total horas |
| total_members | ✅ | Quantidade sócios |
| legal_status | ✅ | DREAM ou FORMALIZED |
| chain_digest | ✅ | Hash de integridade |
| signature | ✅ | Assinatura digital |
| fiscal_batch_hash | ✅ | Hash de integridade do último Lote SPED |
| eligibility_complete | ✅ | Perfil de elegibilidade completo (bool) |
| public_profile_published | ✅ | Perfil público publicado na Rede (bool) |
| credit_matches_count | ✅ | Quantidade de matches de crédito encontrados |
| help_topics_viewed | ✅ | Contagem de tópicos de ajuda visualizados |
| **member_id** | ❌ | Dados sensíveis protegidos |
| **entry_details** | ❌ | Transações detalhadas |
| **posting_id** | ❌ | IDs internos |
| **cadunico_status** | ❌ | Status CadÚnico (sensível) |
| **inadimplencia** | ❌ | Status de inadimplência (sensível) |
| **credit_purpose** | ❌ | Finalidade do crédito (sensível) |

---

## 2. Modelo de Segurança [ATUALIZADO]

### 2.1 Isolamento de Dados

Cada entidade possui banco próprio:
- **Path:** `data/entities/{entity_id}.db`
- **Isolamento:** Físico total
- **Acesso Cruzado:** Proibido entre tenants
- **Banco Central:** `data/entities/central.db` para relações inter-tenant (RF-12, indicadores, programas, help_topics)

### 2.2 Acesso do Contador Social (Painel Multi-tenant)

O acesso do contador parceiro aos dados do empreendimento ocorre estritamente em modo de leitura (Read-Only) e mediante delegação de acesso prévia. O painel apenas consulta e compila as transações, sem nunca quebrar o isolamento do arquivo local `.sqlite`.

**Novo:** Contador pode visualizar e consolidar perfil de elegibilidade, mas não pode modificar dados sensíveis sem autorização explícita.

### 2.3 Integridade

- **Hash SHA256 para auditoria:** Cada decisão gera hash do conteúdo
- **Chain digest:** Cada bloco contábil gera hash de integridade
- **Imutabilidade:** Garantida por design
- **Novo:** Hash do perfil de elegibilidade para auditoria de matches de crédito
- **Novo:** Hash de tópicos de ajuda para versionamento de conteúdo pedagógico

### 2.4 Transporte

- **Pacotes assinados digitalmente:** Assinatura com EntityID
- **Verificação de integridade:** Hash validation
- **Non-repudiation:** Timestamp + nonce
- **Novo:** TLS 1.3 obrigatório para transmissão de dados de elegibilidade

### 2.5 Matriz de Ameaças e Mitigações [ATUALIZADO]

| Ameaça | Mitigação |
|--------|-----------|
| Acesso não autorizado | Isolamento por arquivo + autenticação Gov.br |
| Alteração de dados | Hash SHA256 + logging |
| Interceptação | TLS 1.3 em transporte |
| Replay attack | Timestamp + nonce |
| Vazamento de dados sensíveis (CadÚnico, etc.) | Campos sensíveis nunca transmitidos, apenas flags de completude |
| Match de crédito fraudulento | Hash do perfil + validação pelo Contador Social |
| Perfil público expõe dados sensíveis | Validação rigorosa antes de publicar na Rede |
| Conteúdo de ajuda desatualizado | Versionamento de tópicos + hash de integridade |

---

## 3. Protocolo Econômico [MANTIDO]

### 3.1 Introdução

O Economic Protocol do Digna define as regras econômicas fundamentais que governam a operação dos Empreendimentos de Economia Solidária (EES) dentro do sistema.

### 3.2 Princípios Fundamentais

#### 3.2.1 Primazia do Trabalho
O trabalho humano é reconhecido como a principal fonte de valor econômico.

#### 3.2.2 Autogestão
Toda decisão econômica relevante deve ser tomada coletivamente.

#### 3.2.3 Transparência
Todas as operações econômicas são registradas em ledger verificável.

#### 3.2.4 Soberania de Dados
Cada empreendimento mantém controle sobre seus próprios dados.

#### 3.2.5 Ponte Institucional (Aliança Contábil)
A conformidade não deve ser um fardo.

### 3.3 Unidades de Valor

| Unidade | Uso |
|---------|-----|
| **Moeda Nacional (R$)** | Vendas, compras, contabilidade financeira |
| **Trabalho (minutos)** | Cálculo de participação, distribuição de sobras |
| **Bens Substantivos (Futuro)** | Sementes, animais, bens produtivos |

### 3.4 a 3.13 [MANTIDO - Ver versão anterior para detalhes completos]

---

## 4. Protocolo de Integração com APIs Externas [NOVO - PDF v1.0]

### 4.1 Banco Central do Brasil (BCB)

| API | Frequência | Cache TTL | Tratamento de Erro |
|-----|------------|-----------|-------------------|
| SGS (SELIC, IPCA, CDI) | Diária | 24 horas | Circuit breaker + fallback último valor |
| PTAX (Câmbio) | Diária (úteis) | 24 horas | Circuit breaker + fallback último valor |
| Focus (Expectativas) | Semanal | 7 dias | Circuit breaker + fallback último valor |

**Protocolo de Coleta:**
```go
type IndicatorCollector struct {
    cache IndicatorCache
    circuitBreaker *CircuitBreaker
}

func (c *IndicatorCollector) CollectDaily() error {
    // 1. Verificar circuit breaker
    if c.circuitBreaker.IsOpen() {
        return ErrCircuitBreakerOpen
    }
    
    // 2. Coletar da API
    data, err := c.fetchFromBCB()
    if err != nil {
        c.circuitBreaker.RecordFailure()
        return err
    }
    
    // 3. Persistir em cache local (central.db)
    c.cache.Save(data, 24*time.Hour)
    
    // 4. Reset circuit breaker
    c.circuitBreaker.RecordSuccess()
    return nil
}
```

### 4.2 IBGE (SIDRA)

| API | Frequência | Cache TTL | Tratamento de Erro |
|-----|------------|-----------|-------------------|
| IPCA/INPC | Mensal | 30 dias | Fallback IPCA do BCB |
| PNAD (Emprego) | Trimestral | 90 dias | Fallback último valor |

### 4.3 Gov.br (OAuth2 + Assinatura)

**Fluxo de Autenticação:**
1. Redirect para Gov.br OAuth2
2. Usuário autentica com CPF + senha
3. Callback com access_token
4. Token armazenado criptografado no central.db
5. Token refresh automático antes da expiração

**Fluxo de Assinatura Qualificada:**
1. Usuário solicita assinatura de documento (ata, dossiê)
2. Sistema gera hash SHA256 do documento
3. Redirect para Gov.br assinatura
4. Assinatura ICP-Brasil aplicada
5. Hash da assinatura armazenado no legal_documents

### 4.4 APIs de Certidões (PGFN, Estadual, Municipal)

**Protocolo de Consulta:**
- Frequência: Semanal para entidades FORMALIZED
- Cache: 7 dias
- Fallback: Status "não verificado" se API indisponível
- Notificação: Alerta se certidão vencer em < 30 dias

---

## 5. Protocolo de Match de Elegibilidade [NOVO - PDF v1.0]

### 5.1 Fluxo de Match Automático

```
┌─────────────────┐     ┌─────────────────┐     ┌─────────────────┐
│   ERP (Perfil)  │────▶│   Motor Match   │────▶│   Portal (Lista)│
│   + Elegibilidade│     │   (Regras)      │     │   de Oportunidades│
└─────────────────┘     └─────────────────┘     └─────────────────┘
         │                       │                       │
         │                       │                       │
         ▼                       ▼                       ▼
   Dados já existentes    Regras de negócio      Match + Checklist
   (CNPJ, CNAE, etc.)     (programa × perfil)    (documentos faltantes)
```

### 5.2 Regras de Match

| Critério | Fonte | Validação |
|----------|-------|-----------|
| Tipo de Entidade (MEI, ME, etc.) | ERP | Exato |
| Faturamento Anual | ERP | Faixa (mín/máx) |
| Município/UF | ERP | Exato ou lista |
| Tempo de CNPJ | ERP | >= X meses |
| CadÚnico | Elegibilidade (campo sensível) | Booleano |
| Gênero (mulher) | Elegibilidade (campo sensível) | Booleano |
| Inadimplência | Elegibilidade (campo sensível) | Booleano (direciona para Desenrola) |
| Certidões Negativas | Integração PGFN | Válido/Não válido |

### 5.3 Estrutura do Resultado de Match

```go
type MatchResult struct {
    ProgramID       string   // ID do programa
    ProgramName     string   // Nome amigável
    Eligible        bool     // Elegível ou não
    MatchScore      int64    // Score de match (0-100)
    MissingDocs     []string // Documentos faltantes
    Deadline        int64    // Prazo de inscrição (Unix timestamp)
    EffectiveRate   int64    // Taxa efetiva de juros (int64, centavos de %)
    MaxAmount       int64    // Valor máximo (int64, centavos)
    Recommendation  string   // Recomendação em linguagem popular
}
```

### 5.4 Privacidade no Match

**Princípio:** Dados sensíveis (CadÚnico, inadimplência, gênero) NUNCA são transmitidos para APIs externas.

**Implementação:**
- Match executado LOCALMENTE no módulo `portal_opportunities`
- Apenas resultado do match (elegível/não elegível) é armazenado
- Contador Social pode visualizar resultado, mas não campos sensíveis brutos
- Para submissão de candidatura, usuário autoriza envio específico por programa

---

## 6. Protocolo da Rede Digna [NOVO - PDF v1.0]

### 6.1 Publicação de Perfil Público

**Fluxo:**
1. Usuário preenche perfil público (campos não sensíveis)
2. Sistema valida que nenhum dado sensível está incluído
3. Usuário confirma publicação
4. Perfil marcado como `public_profile_published = true`
5. Perfil sincronizado para agregador (apenas campos públicos)

**Campos Públicos Permitidos:**
- ✅ Nome fantasia
- ✅ Missão (texto livre)
- ✅ Produtos (categorias)
- ✅ Serviços (categorias)
- ✅ Município, UF
- ✅ Contato público (email/telefone)
- ✅ Foto/logo
- ❌ CNPJ/CPF (hash anonimizado apenas)
- ❌ Faturamento
- ❌ Membros
- ❌ Dados financeiros

### 6.2 Mural de Necessidades

**Estrutura de Postagem:**
```go
type NeedPost struct {
    ID            string // UUID
    PublisherID   string // Hash anonimizado (não entity_id direto)
    Categoria     string // INSUMO, EQUIPAMENTO, SERVICO, OUTRO
    Descricao     string // Texto livre (moderado)
    Quantidade    string // Ex: "100kg", "5 unidades"
    PrazoDesejado int64  // Unix timestamp
    Municipio     string // Para matching geográfico
    UF            string
    Status        string // ABERTO, EM_NEGOCIACAO, CONCLUIDO
    CreatedAt     int64  // Unix timestamp
}
```

### 6.3 Protocolo de Match B2B

**Critérios de Matching:**
1. **Geográfico:** Priorizar mesma cidade → mesma UF → região
2. **Setorial:** Afinidade de CNAE/categoria
3. **Temporal:** Prazos compatíveis
4. **Reputação:** Histórico de transações (futuro)

**Algoritmo:**
```go
func MatchBuyerToSeller(buyer NeedPost, sellers []PublicProfile) []MatchResult {
    var results []MatchResult
    
    for _, seller := range sellers {
        score := 0
        
        // Geográfico (40 pontos)
        if buyer.Municipio == seller.Municipio {
            score += 40
        } else if buyer.UF == seller.UF {
            score += 20
        }
        
        // Setorial (40 pontos)
        if hasCategoryOverlap(buyer.Categoria, seller.Produtos) {
            score += 40
        }
        
        // Temporal (20 pontos)
        if isPrazoCompatible(buyer.PrazoDesejado) {
            score += 20
        }
        
        if score >= 60 { // Threshold mínimo
            results = append(results, MatchResult{
                SellerID: seller.EntityID,
                Score: score,
                // ...
            })
        }
    }
    
    // Ordenar por score descendente
    sort.Slice(results, func(i, j int) bool {
        return results[i].Score > results[j].Score
    })
    
    return results
}
```

### 6.4 Privacidade na Rede

**Regras:**
- EntityID sempre anonimizado (hash) no mural
- Negociação inicial via sistema (não expõe contato direto)
- Contato direto apenas após ambas partes concordarem
- Histórico de transações agregado, não detalhado

---

## 7. Protocolo de Cache e Performance [NOVO]

### 7.1 Estratégias de Cache por Módulo

| Módulo | Dado | TTL | Estratégia |
|--------|------|-----|------------|
| **indicators_engine** | SELIC, IPCA | 24h | Cache local + refresh diário |
| **indicators_engine** | Câmbio | 24h (úteis) | Cache local + refresh diário |
| **portal_opportunities** | Programas | 7 dias | Cache local + refresh semanal |
| **portal_opportunities** | Matches | 24h | Recalcular diário |
| **rede_digna** | Perfis públicos | 1h | Cache local + invalidação por evento |
| **rede_digna** | Mural | 15min | Cache local + polling |
| **help_engine** | Tópicos de ajuda | 1h | Cache local + invalidação por atualização |

### 7.2 Estrutura de Cache (central.db)

```sql
CREATE TABLE IF NOT EXISTS cache_indicators (
    id TEXT PRIMARY KEY,
    indicator_key TEXT NOT NULL UNIQUE,
    value INTEGER NOT NULL,
    source TEXT NOT NULL,
    created_at INTEGER NOT NULL,
    expires_at INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS cache_programs (
    id TEXT PRIMARY KEY,
    program_id TEXT NOT NULL UNIQUE,
    data TEXT NOT NULL, -- JSON
    created_at INTEGER NOT NULL,
    expires_at INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS cache_public_profiles (
    entity_id TEXT PRIMARY KEY,
    profile_data TEXT NOT NULL, -- JSON
    created_at INTEGER NOT NULL,
    expires_at INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS cache_help_topics (
    topic_key TEXT PRIMARY KEY,
    topic_data TEXT NOT NULL, -- JSON
    created_at INTEGER NOT NULL,
    expires_at INTEGER NOT NULL
);
```

### 7.3 Invalidação de Cache

**Eventos que invalidam cache:**
- Atualização de perfil de elegibilidade → Invalida matches de crédito
- Publicação de perfil público → Invalida cache da Rede
- Nova decisão de assembleia → Invalida cache de formalização
- Exportação fiscal → Invalida cache de certidões
- Atualização de tópico de ajuda → Invalida cache de ajuda [NOVO - RF-30]

---

## 8. Protocolo de Notificações [NOVO]

### 8.1 Tipos de Notificação

| Tipo | Gatilho | Canal | Prioridade |
|------|---------|-------|------------|
| **Vencimento DAS** | Dia 15, 19, 20 de cada mês | In-app + Email | Alta |
| **Match de Crédito** | Novo programa elegível encontrado | In-app + Email | Alta |
| **Prazo de Edital** | 7 dias antes do vencimento | In-app + Email + WhatsApp | Crítica |
| **Certidão Vencendo** | 30 dias antes do vencimento | In-app + Email | Média |
| **Match B2B** | Nova oportunidade na Rede | In-app | Baixa |
| **Nova Decisão** | Assembleia registrada | In-app | Média |
| **Ajuda Contextual** | Campo técnico acessado pela primeira vez | In-app (tooltip) | Baixa |

### 8.2 Estrutura de Notificação

```go
type Notification struct {
    ID        string // UUID
    EntityID  string // Destinatário
    Type      string // DAS, CREDIT_MATCH, DEADLINE, etc.
    Title     string // Título em linguagem popular
    Message   string // Mensagem completa
    ActionURL string // Link para ação relacionada
    Priority  string // LOW, MEDIUM, HIGH, CRITICAL
    Read      bool   // Lida ou não
    CreatedAt int64  // Unix timestamp
}
```

### 8.3 Preferências do Usuário

- Usuário pode configurar canais preferidos por tipo de notificação
- Silenciar notificações por período (férias, etc.)
- Opt-out de notificações não críticas

---

## 9. Protocolo de Auditoria e Logging [ATUALIZADO]

### 9.1 Eventos Auditáveis

| Evento | Dados Auditados | Retenção |
|--------|-----------------|----------|
| Login/Logout | User ID, timestamp, IP | 2 anos |
| Exportação Fiscal | Entity ID, período, hash do lote | 5 anos |
| Alteração de Perfil | Campo antigo, campo novo, user ID | 2 anos |
| Match de Crédito | Programa, resultado, timestamp | 2 anos |
| Publicação na Rede | Conteúdo publicado, timestamp | 2 anos |
| Assinatura Gov.br | Documento, hash, timestamp | 5 anos |
| Acesso a Ajuda | Tópico acessado, timestamp, entity_id | 1 ano |

### 9.2 Estrutura de Log

```json
{
  "timestamp": 1772856840,
  "level": "INFO",
  "module": "portal_opportunities",
  "entity_id": "hash_anonimizado",
  "action": "credit_match_executed",
  "details": {
    "programs_evaluated": 15,
    "matches_found": 3,
    "duration_ms": 45
  },
  "user_id": "hash_anonimizado"
}
```

### 9.3 Privacidade em Logs

- EntityID e UserID sempre anonimizados em logs externos
- Dados sensíveis (CadÚnico, inadimplência) nunca logados em claro
- Logs de produção separados de logs de desenvolvimento

---

## 10. Protocolo de Backup e Recovery [ATUALIZADO]

### 10.1 Backup por Entidade

**Frequência:** Diária (automática)
**Formato:** SQLite dump + metadata JSON
**Retenção:** 30 dias (configurável)
**Local:** Volume externo + cloud (opcional)

### 10.2 Backup do Banco Central

**Frequência:** Diária (automática)
**Conteúdo:**
- Vínculos contábeis (EnterpriseAccountant)
- Indicadores econômicos (cache)
- Programas de financiamento (catálogo)
- Perfis públicos (Rede Digna)
- Tópicos de ajuda (Help System) [NOVO - RF-30]

### 10.3 Recovery

**Tempo Alvo:** < 1 hora para entidade crítica
**Processo:**
1. Identificar backup mais recente válido
2. Validar integridade (hash)
3. Restaurar em arquivo temporário
4. Validar schema e dados
5. Substituir arquivo original
6. Validar acesso

---

## 11. Protocolo do Sistema de Ajuda Educativa [NOVO - RF-30, Sessão 27/03/2026]

### 11.1 Linkagem UI → Banco de Ajuda

**Fluxo:**
1. UI detecta campo técnico (ex: "Inscrito no CadÚnico?")
2. Botão "?" linka para `/help/topic/{key}` (ex: `/help/topic/cadunico`)
3. Sistema busca `HelpTopic` no `central.db` por `key`
4. Renderiza explicação em linguagem popular + legislação + próximo passo
5. Incrementa `ViewCount` para métricas de uso

### 11.2 Estrutura de Tópico de Ajuda

```go
type HelpTopic struct {
    Key          string // Chave única (ex: "cadunico", "inadimplencia")
    Title        string // Título em linguagem popular
    Summary      string // Resumo em 1 frase (para tooltips)
    Explanation  string // Explicação completa em linguagem popular
    WhyAsked     string // "Por que perguntamos isso?"
    Legislation  string // Legislação relacionada
    NextSteps    string // Próximos passos acionáveis
    OfficialLink string // Link para fonte oficial (ex: gov.br)
    Category     string // Categoria: CREDITO, TRIBUTARIO, GOVERNANCA, GERAL
    Tags         string // Tags para busca (JSON array ou comma-separated)
    ViewCount    int64  // Quantas vezes foi visualizado
    CreatedAt    int64  // Unix timestamp
    UpdatedAt    int64  // Unix timestamp
}
```

### 11.3 Categorias de Tópicos

| Categoria | Exemplos de Tópicos |
|-----------|---------------------|
| **CRÉDITO** | CadÚnico, Inadimplência, Pronampe, PNMPO |
| **TRIBUTÁRIO** | CNAE, DAS MEI, Simples Nacional, Certidões |
| **GOVERNANÇA** | Reserva Legal, FATES, Assembleias, CADSOL |
| **GERAL** | Soberania de Dados, Contabilidade Invisível, ITG 2002 |

### 11.4 Critérios de Qualidade de Conteúdo

- **Linguagem:** Usuário com 5ª série consegue entender sem ajuda externa
- **Performance:** Tooltip carrega em < 500ms via HTMX
- **Jargão:** Zero termos técnicos sem explicação ("cadastramento", "regularização fiscal", etc.)
- **Ação:** Sempre inclui "próximo passo" acionável (ex: "procure o CRAS")

### 11.5 Versionamento de Conteúdo

- Cada atualização de tópico gera novo hash de integridade
- Histórico de versões mantido para auditoria
- Invalidação de cache automática upon update

---

## 12. Referências Externas [ATUALIZADO]

| Referência | Aplicação no Digna |
|------------|-------------------|
| ITG 2002 (CFC) | Norma de contabilidade para EES |
| Lei nº 5.764/71 | Lei Geral das Cooperativas |
| LC nº 214/2025 | Atos Cooperativos e tributação |
| Lei nº 15.068/2024 | Lei Paul Singer (Economia Solidária) |
| Lei nº 14.063/2020 | Assinaturas Eletrônicas Gov.br |
| IN DREI nº 79/2020 | Reuniões e Assembleias Digitais |
| Decreto nº 12.784/2025 | SINAES/CADSOL |
| Portaria MAPA nº 393/2021 | Memorial Técnico Sanitário (MTSE) |
| BCB API Documentation | Motor de Indicadores |
| IBGE SIDRA API | Indicadores sociais |

---

**Status:** ✅ ATUALIZADO COM PROTOCOLOS DO ECOSSISTEMA (PDF v1.0) + RF-30 (Decisão de Design 27/03/2026)  
**Próxima Ação:** Atualizar `03_architecture/03_improvements.md` com melhorias identificadas para os novos módulos  
**Versão Anterior:** 1.0 (2026-03-09)  
**Versão Atual:** 2.0 (2026-03-27)
