title: Modelos de Domínio e Dados
status: implemented
version: 2.1
last_updated: 2026-03-27
---

# Modelos de Domínio e Dados - Ecossistema Digna

> **Nota:** Este documento consolida todos os modelos de domínio do Ecossistema Digna (PDF v1.0), preservando todo o trabalho validado nas Sprints 1-16, incorporando os novos módulos como Fases 3 e 4, e incluindo o Sistema de Ajuda Educativa (RF-30) decidido na sessão de 27/03/2026.

---

## 1. Domain Model (Modelo de Domínio)

O domínio do Digna reflete os princípios da autogestão e da contabilidade invisível, priorizando as relações humanas sobre o capital financeiro, mas agora atuando também como **ponte institucional** para a conformidade legal, acesso a crédito e intercooperação em rede.

### 1.1 Entidades Principais (Módulo 1 - ERP)

#### Enterprise (Empreendimento de Economia Solidária - EES)
Representa o coletivo produtivo. Pode transitar gradualmente por três estados, respeitando o tempo político do grupo:

| Estado | Descrição |
|--------|-----------|
| **DREAM (Sonho)** | Grupo informal, focado na união produtiva inicial |
| **INCUBATED (Incubado)** | Em processo de estruturação, recebendo apoio pedagógico (ITCPs, ONGs) |
| **FORMALIZED (Formalizado)** | Cooperativa ou Associação com CNPJ e estatuto base (Pronto para CADSOL e obrigações fiscais) |

```go
type Enterprise struct {
    ID               string    // UUID
    Name             string    // Nome fantasia ou razão social
    CNPJ             string    // CPF/CNPJ (vazio para DREAM)
    CNAE             string    // Código CNAE principal
    Municipio        string    // Município
    UF               string    // Unidade Federativa
    FaturamentoAnual int64     // int64 - Anti-Float (centavos)
    RegimeTributario string    // MEI, Simples, Isento, etc.
    DataAbertura     int64     // Unix timestamp
    SituacaoFiscal   string    // Ativa, Suspensa, etc.
    Status           string    // DREAM, INCUBATED, FORMALIZED
    CreatedAt        int64     // Unix timestamp
    UpdatedAt        int64     // Unix timestamp
}
```

#### WorkLog (Registro ITG 2002)
Registro de trabalho cooperativo. Converte o suor em Capital Social de Trabalho (mensurado em minutos).

```go
type WorkLog struct {
    ID        string // UUID
    EntityID  string // Vínculo com entidade
    MemberID  string // Vínculo com membro
    Activity  string // Descrição da atividade
    Minutes   int64  // int64 - Anti-Float (minutos trabalhados)
    Date      int64  // Unix timestamp
    CreatedAt int64  // Unix timestamp
}
```

#### Decision (Decisão Democrática)
Decisão coletiva tomada e registrada em Assembleia. Base para a geração das Atas em Markdown (CADSOL).

```go
type Decision struct {
    ID          string // UUID
    EntityID    string // Vínculo com entidade
    Title       string // Título da decisão
    Description string // Descrição detalhada
    VoteType    string // ABERTA, SECRETA
    Result      string // APROVADA, REJEITADA, ADIADA
    VotosSim    int64  // Contagem de votos favoráveis
    VotosNao    int64  // Contagem de votos contrários
    VotosNulos  int64  // Contagem de votos nulos/brancos
    HashSHA256  string // Hash de integridade do registro
    CreatedAt   int64  // Unix timestamp
}
```

#### LegalDocument (Dossiê CADSOL e Atas)
**Camada Cotidiana:** O livro de atas oficial do grupo. Agrupamento das decisões exportado para provar que a gestão é democrática. Com a nova atualização, o presidente assina pelo celular com a senha oficial do Gov.br, sem precisar pisar no cartório, e o voto secreto de cada um é garantido pelo sistema.

**Camada Técnica:** Agrupamento das `Decisions` exportado em formato Markdown/PDF para fins legais, comprovando a gestão democrática perante o Estado. O sistema agora supera o mero uso de Hash SHA256 integrando Assinatura Eletrônica Avançada/Qualificada (ICP-Brasil/Gov.br) para a mesa diretora e assegurando a anonimização em votações (IN DREI nº 79/2020 e Lei nº 14.063/2020).

```go
type LegalDocument struct {
    ID          string // UUID
    EntityID    string // Vínculo com entidade
    Type        string // ATA, ESTATUTO, DOSSIE_CADSOL, MTSE
    Content     string // Conteúdo do documento (Markdown)
    HashSHA256  string // Hash de integridade
    Signature   string // Assinatura eletrônica (Gov.br/ICP-Brasil)
    Status      string // DRAFT, SIGNED, SUBMITTED
    CreatedAt   int64  // Unix timestamp
}
```

#### Fund (Fundos Obrigatórios)
Reservas estatutárias e legais blindadas pelo sistema (Ex: Reserva Legal e FATES).

```go
type Fund struct {
    ID         string // UUID
    EntityID   string // Vínculo com entidade
    Type       string // RESERVA_LEGAL, FATES, OUTRO
    Amount     int64  // int64 - Anti-Float (centavos)
    Period     string // Período de referência (YYYY-MM)
    CreatedAt  int64  // Unix timestamp
}
```

#### EnterpriseAccountant (Vínculo Contábil Temporal) [NOVO]
**Camada Cotidiana:** A "chave temporária" que a cooperativa empresta para o contador social arrumar a casa, podendo ser recolhida a qualquer hora.

**Camada Técnica:** Relacionamento N:N que estabelece a delegação de acesso temporal entre um `Enterprise` e um Contador parceiro. Atua como um "Filtro de Vigência", liberando acesso ao painel unicamente em modo `Read-Only` para os períodos fiscais previamente autorizados pela autogestão.

```go
type EnterpriseAccountant struct {
    ID           string // UUID
    EnterpriseID string // Vínculo com entidade
    AccountantID string // Vínculo com contador (usuário global)
    Status       string // ACTIVE, INACTIVE
    StartDate    int64  // Unix timestamp - Início do vínculo
    EndDate      int64  // Unix timestamp - Fim do vínculo (0 se ativo)
    DelegatedBy  string // ID do usuário que delegou o acesso
    CreatedAt    int64  // Unix timestamp
    UpdatedAt    int64  // Unix timestamp
}
```

#### TaxCompliance & ReinfEvent [NOVO - Adequação Estatal]
**Camada Cotidiana:** O mensageiro silencioso do aplicativo. Ele avisa os "robôs" da Receita Federal sobre pagamentos de frete e serviços para a cooperativa não ser multada, e garante que o dinheiro do "Ato Cooperativo" não sofra a mesma mordida de imposto das grandes empresas.

**Camada Técnica:** Motor de geração de XMLs de retenção (assinados via certificado A1/A3) e envio síncrono/assíncrono para Web Services da EFD-Reinf (ex: evento de fechamento R-2099) para alimentar a DCTFWeb. Detém a inteligência contábil de expurgar as receitas de Atos Cooperativos na composição do Bloco M da ECF (e-Lalur/e-Lacs), blindando a isenção prevista na Lei 5.764/71 e LC 214/2025.

```go
type ReinfEvent struct {
    ID           string // UUID
    EntityID     string // Vínculo com entidade
    EventType    string // R-2010, R-2020, R-2099, etc.
    XMLContent   string // Conteúdo XML assinado
    HashSHA256   string // Hash de integridade
    Status       string // PENDING, SENT, CONFIRMED, ERROR
    SentAt       int64  // Unix timestamp
    CreatedAt    int64  // Unix timestamp
}
```

#### SanitaryDossier (MTSE) [NOVO - Adequação Estatal]
**Camada Cotidiana:** O gerador automatizado do "caderno da fábrica", ajudando a pequena agroindústria (queijaria, casa de mel, pescados) a desenhar sua planta e procedimentos de higiene para conquistar o selo de venda oficial sem precisar de consultorias caras.

**Camada Técnica:** Estrutura de dados que parametriza fluxogramas de maquinário, capacidade diária e potabilidade de água, exportando automaticamente o Memorial Técnico Sanitário de Estabelecimento (MTSE) em conformidade estrita com o RIISPOA (Portaria MAPA nº 393/2021) para peticionamento no sistema SEI.

```go
type SanitaryDossier struct {
    ID              string // UUID
    EntityID        string // Vínculo com entidade
    TipoProduto     string // QUEIJO, MEL, CARNES, LATICINIOS, etc.
    CapacidadeDiaria int64 // Litros/kg por dia
    FluxoMaquinario string // Descrição do fluxo
    OrigemAgua      string // POÇO, REDE, NASCENTE
    Content         string // Conteúdo do MTSE (Markdown/PDF)
    Status          string // DRAFT, SUBMITTED, APPROVED
    CreatedAt       int64  // Unix timestamp
}
```

#### DASMEI (Cálculo Automático DAS MEI) [NOVO - PDF v1.0, Seção 4.1]
**Camada Cotidiana:** O sistema calcula sozinho o boleto mensal do MEI e avisa quando está perto de vencer, sem o empreendedor precisar fazer conta ou decorar datas.

**Camada Técnica:** Registro mensal do DAS MEI com cálculo baseado em tabela versionada de salário mínimo (5% + ICMS/ISS fixos), alertas de vencimento e histórico de pagamentos.

```go
type DASMEI struct {
    ID              string // UUID
    EntityID        string // Vínculo com entidade
    Competencia     string // YYYY-MM
    ValorDevido     int64  // int64 - Anti-Float (centavos)
    ValorPago       int64  // int64 - Anti-Float (centavos)
    DataVencimento  int64  // Unix timestamp
    DataPagamento   int64  // Unix timestamp (0 se não pago)
    Status          string // PENDENTE, PAGO, VENCIDO
    SalarioMinimo   int64  // Salário mínimo de referência (centavos)
    CreatedAt       int64  // Unix timestamp
    UpdatedAt       int64  // Unix timestamp
}
```

---

### 1.2 Entidades do Motor de Indicadores (Módulo 2) [NOVO - PDF v1.0, Seção 5]

#### EconomicIndicator (Indicador Econômico)
Registro de indicadores coletados de APIs externas (BCB, IBGE).

```go
type EconomicIndicator struct {
    ID            string // UUID
    Codigo        string // Código da série (ex: 433 para IPCA)
    Valor         int64  // int64 - Anti-Float (ex: IPCA em centavos de %)
    DataReferencia int64 // Unix timestamp - Data de referência do indicador
    Fonte         string // BCB_SGS, BCB_PTAX, IBGE_SIDRA, etc.
    CriadoEm      int64  // Unix timestamp - Quando foi coletado
}
```

#### IndicatorCache (Cache Local)
Tabela de cache para evitar chamadas redundantes às APIs externas.

```go
type IndicatorCache struct {
    ID           string // UUID
    IndicatorKey string // Chave única (fonte + codigo)
    Value        int64  // int64 - Valor armazenado
    ExpiresAt    int64  // Unix timestamp - Quando o cache expira
    CreatedAt    int64  // Unix timestamp
}
```

---

### 1.3 Entidades do Perfil de Elegibilidade e Portal (Módulo 3) [NOVO - PDF v1.0, Seção 6 e 9]

#### EligibilityProfile (Perfil de Elegibilidade) [NOVO]
Conjunto de campos que o Portal precisa para executar o match com programas de financiamento. Composto por dados que o ERP já captura automaticamente e dados complementares de preenchimento único.

**Princípio Central:** Nenhum usuário precisa preencher o mesmo dado duas vezes.

```go
type EligibilityProfile struct {
    ID         string // UUID
    EntityID   string // Vínculo com entidade (único por entidade)
    
    // Dados já capturados pelo ERP (referência, não duplicação)
    CNPJ             string // Copiado de Enterprise
    CNAE             string // Copiado de Enterprise
    Municipio        string // Copiado de Enterprise
    UF               string // Copiado de Enterprise
    FaturamentoAnual int64  // int64 - Anti-Float (centavos)
    RegimeTributario string
    DataAbertura     int64 // Unix timestamp
    SituacaoFiscal   string
    
    // CAMPOS COMPLEMENTARES (preenchimento único)
    InscritoCadUnico    bool   // Habilita programas sociais
    SocioMulher         bool   // Prioridade em linhas com foco de gênero
    InadimplenciaAtiva  bool   // Direciona ao Desenrola antes de crédito novo
    FinalidadeCredito   string // Enum: CAPITAL_GIRO, EQUIPAMENTO, REFORMA, OUTRO
    ValorNecessario     int64  // int64 - Anti-Float (centavos)
    TipoEntidade        string // Enum: MEI, ME, EPP, Cooperativa, OSC, OSCIP, PF
    ContabilidadeFormal bool   // Requisito de alguns programas
    
    // Metadados
    PreenchidoEm  int64 // Unix timestamp - primeiro preenchimento
    AtualizadoEm  int64 // Unix timestamp - última atualização
    PreenchidoPor string // ID do usuário que preencheu
    
    CreatedAt int64
    UpdatedAt int64
}

// Enums para campos restritos
const (
    FinalidadeCapitalGiro   = "CAPITAL_GIRO"
    FinalidadeEquipamento   = "EQUIPAMENTO"
    FinalidadeReforma       = "REFORMA"
    FinalidadeOutro         = "OUTRO"
    
    TipoEntidadeMEI         = "MEI"
    TipoEntidadeME          = "ME"
    TipoEntidadeEPP         = "EPP"
    TipoEntidadeCooperativa = "Cooperativa"
    TipoEntidadeOSC         = "OSC"
    TipoEntidadeOSCIP       = "OSCIP"
    TipoEntidadePF          = "PF"
)
```

#### FinancingProgram (Programa de Financiamento) [NOVO]
Catálogo de programas de crédito disponíveis para match.

```go
type FinancingProgram struct {
    ID                string // UUID
    Nome              string // Nome do programa
    Fonte             string // FEDERAL, ESTADUAL, MUNICIPAL
    ValorMaximo       int64  // int64 - Anti-Float (centavos)
    TaxaJuros         int64  // int64 - Juros em centavos de % (ex: 0 = juros zero)
    PrazoMaximoMeses  int64  // Prazo máximo em meses
    CarenciaMeses     int64  // Carência em meses
    Requisitos        string // JSON com requisitos de elegibilidade
    PublicoPrioritario string // Descrição do público-alvo
    Ativo             bool   // Programa está ativo?
    LinkEdital        string // Link para edital/documentação
    CreatedAt         int64  // Unix timestamp
    UpdatedAt         int64  // Unix timestamp
}
```

#### ProgramMatch (Match de Elegibilidade) [NOVO]
Registro do cruzamento entre perfil da entidade e programa de financiamento.

```go
type ProgramMatch struct {
    ID          string // UUID
    EntityID    string // Vínculo com entidade
    ProgramID   string // Vínculo com programa
    Elegibilidade string // ELEGIVEL, NAO_ELEGIVEL, PARCIAL
    Motivo      string // Explicação do match (ou motivo de não elegibilidade)
    DocumentosPendentes string // JSON com lista de documentos faltantes
    CreatedAt   int64  // Unix timestamp
}
```

---

### 1.4 Entidades da Rede Digna (Módulo 4) [NOVO - PDF v1.0, Seção 7]

#### PublicProfile (Perfil Público da Entidade) [NOVO]
Informações visíveis para a rede, sem expor dados sensíveis.

```go
type PublicProfile struct {
    EntityID       string   // Hash anonimizado
    NomeFantasia   string
    Missao         string   // Texto livre
    Produtos       []string // Lista de categorias
    Servicos       []string // Lista de capacidades
    Municipio      string
    UF             string
    ContatoPublico string   // Email/telefone para negócios
    FotoLogo       string   // URL do asset
    CreatedAt      int64    // Unix timestamp
    UpdatedAt      int64    // Unix timestamp
}
```

#### NeedPost (Mural de Necessidades) [NOVO]
Demandas de compra publicadas visíveis para a rede.

```go
type NeedPost struct {
    ID            string // UUID
    PublisherID   string // Hash anonimizado
    Categoria     string // INSUMO, EQUIPAMENTO, SERVICO, OUTRO
    Descricao     string
    Quantidade    string
    PrazoDesejado int64  // Unix timestamp
    Municipio     string // Para matching geográfico
    UF            string
    Status        string // ABERTO, EM_NEGOCIACAO, CONCLUIDO
    CreatedAt     int64  // Unix timestamp
    UpdatedAt     int64  // Unix timestamp
}
```

#### SolidarityTransaction (Histórico de Transações Solidárias) [NOVO]
Registro de transações entre entidades da rede para evidência de atividade econômica.

```go
type SolidarityTransaction struct {
    ID            string // UUID
    BuyerID       string // Hash anonimizado do comprador
    SellerID      string // Hash anonimizado do vendedor
    Valor         int64  // int64 - Anti-Float (centavos)
    Descricao     string
    Data          int64  // Unix timestamp
    Status        string // COMPLETADA, CANCELADA
    CreatedAt     int64  // Unix timestamp
}
```

---

### 1.5 Entidades do Sistema de Ajuda Educativa (RF-30) [NOVO - Decisão da Sessão 27/03/2026]

#### HelpTopic (Tópico de Ajuda) [NOVO]
Registro de conceitos técnicos traduzidos em linguagem popular, com linkagem entre elementos de UI e registros de ajuda no banco.

**Princípio Central:** Nenhum usuário deve se sentir humilhado por não entender um termo. O sistema ensina enquanto é operado.

```go
type HelpTopic struct {
    ID           string // UUID
    Key          string // Chave única (ex: "cadunico", "inadimplencia")
    Title        string // Título em linguagem popular
    Summary      string // Resumo em 1 frase (para tooltips)
    Explanation  string // Explicação completa em linguagem popular
    WhyAsked     string // "Por que perguntamos isso?"
    Legislation  string // Legislação relacionada (ex: "Lei nº XXXX/2020")
    NextSteps    string // Próximos passos acionáveis
    OfficialLink string // Link para fonte oficial (ex: gov.br)
    Category     string // Categoria: CREDITO, TRIBUTARIO, GOVERNANCA, GERAL
    Tags         string // Tags para busca (JSON array ou comma-separated)
    
    // Metadados
    ViewCount    int64  // Quantas vezes foi visualizado
    CreatedAt    int64  // Unix timestamp
    UpdatedAt    int64  // Unix timestamp
}

// Categorias de tópicos
const (
    CategoriaCredito    = "CREDITO"
    CategoriaTributario = "TRIBUTARIO"
    CategoriaGovernanca = "GOVERNANCA"
    CategoriaGeral      = "GERAL"
)
```

**Tópicos Seed Obrigatórios (Inicialização do Sistema):**
| Key | Title | Category |
|-----|-------|----------|
| `cadunico` | "O que é o CadÚnico?" | CREDITO |
| `inadimplencia` | "O que é inadimplência?" | CREDITO |
| `cnae` | "O que é CNAE?" | TRIBUTARIO |
| `das_mei` | "O que é o DAS MEI?" | TRIBUTARIO |
| `reserva_legal` | "O que é Reserva Legal?" | GOVERNANCA |
| `fates` | "O que é o FATES?" | GOVERNANCA |

---

## 2. Data Model (Schema v2.0)

O banco de dados é instanciado fisicamente de forma isolada por Enterprise (Soberania do Dado).

**Nota Arquitetural:** O Painel do Contador Social **não possui** banco de transações próprio; ele atua apenas consumindo dados em modo de leitura (Read-Only) dos micro-databases autorizados. E como um contador social pode estar cuidando de várias entidades, as funcionalidades a ele relacionadas devem estar preparadas para pesquisar em diversos db sqlite utilizando o filtro do `EnterpriseAccountant`.

### 2.1 Banco Central (central.db)

Armazena relações inter-tenant e dados globais:

```sql
-- Tabela de vínculos contábeis (RF-12)
CREATE TABLE IF NOT EXISTS enterprise_accountants (
    id TEXT PRIMARY KEY,
    enterprise_id TEXT NOT NULL,
    accountant_id TEXT NOT NULL,
    status TEXT NOT NULL, -- ACTIVE, INACTIVE
    start_date INTEGER NOT NULL,
    end_date INTEGER DEFAULT 0,
    delegated_by TEXT NOT NULL,
    created_at INTEGER,
    updated_at INTEGER,
    UNIQUE(enterprise_id, accountant_id)
);

-- Tabela de indicadores econômicos (RF-18)
CREATE TABLE IF NOT EXISTS economic_indicators (
    id TEXT PRIMARY KEY,
    codigo TEXT NOT NULL,
    valor INTEGER NOT NULL, -- int64, Anti-Float
    data_referencia INTEGER NOT NULL,
    fonte TEXT NOT NULL,
    criado_em INTEGER,
    UNIQUE(codigo, data_referencia, fonte)
);

-- Tabela de cache de indicadores
CREATE TABLE IF NOT EXISTS indicator_cache (
    id TEXT PRIMARY KEY,
    indicator_key TEXT NOT NULL UNIQUE,
    value INTEGER NOT NULL,
    expires_at INTEGER NOT NULL,
    created_at INTEGER
);

-- Tabela de programas de financiamento (RF-20)
CREATE TABLE IF NOT EXISTS financing_programs (
    id TEXT PRIMARY KEY,
    nome TEXT NOT NULL,
    fonte TEXT NOT NULL, -- FEDERAL, ESTADUAL, MUNICIPAL
    valor_maximo INTEGER NOT NULL, -- int64, Anti-Float
    taxa_juros INTEGER NOT NULL, -- int64, centavos de %
    prazo_maximo_meses INTEGER,
    carencia_meses INTEGER,
    requisitos TEXT, -- JSON
    publico_prioritario TEXT,
    ativo INTEGER DEFAULT 1,
    link_edital TEXT,
    created_at INTEGER,
    updated_at INTEGER
);

-- Tabela de tópicos de ajuda (RF-30) [NOVO]
CREATE TABLE IF NOT EXISTS help_topics (
    id TEXT PRIMARY KEY,
    key TEXT NOT NULL UNIQUE,
    title TEXT NOT NULL,
    summary TEXT,
    explanation TEXT NOT NULL,
    why_asked TEXT,
    legislation TEXT,
    next_steps TEXT,
    official_link TEXT,
    category TEXT NOT NULL,
    tags TEXT,
    view_count INTEGER DEFAULT 0,
    created_at INTEGER,
    updated_at INTEGER
);
```

### 2.2 Banco por Entidade (data/entities/{entity_id}.db)

Cada entidade possui seu próprio banco com as seguintes tabelas principais:

```sql
-- Enterprise (dados locais da entidade)
CREATE TABLE IF NOT EXISTS enterprises (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    cnpj TEXT,
    cnae TEXT,
    municipio TEXT,
    uf TEXT,
    faturamento_anual INTEGER, -- int64, Anti-Float
    regime_tributario TEXT,
    data_abertura INTEGER,
    situacao_fiscal TEXT,
    status TEXT NOT NULL, -- DREAM, INCUBATED, FORMALIZED
    created_at INTEGER,
    updated_at INTEGER
);

-- WorkLog (Registro ITG 2002)
CREATE TABLE IF NOT EXISTS work_logs (
    id TEXT PRIMARY KEY,
    entity_id TEXT NOT NULL,
    member_id TEXT NOT NULL,
    activity TEXT,
    minutes INTEGER NOT NULL, -- int64, Anti-Float
    date INTEGER NOT NULL,
    created_at INTEGER
);

-- Decision (Decisões de Assembleia)
CREATE TABLE IF NOT EXISTS decisions (
    id TEXT PRIMARY KEY,
    entity_id TEXT NOT NULL,
    title TEXT NOT NULL,
    description TEXT,
    vote_type TEXT, -- ABERTA, SECRETA
    result TEXT, -- APROVADA, REJEITADA, ADIADA
    votos_sim INTEGER,
    votos_nao INTEGER,
    votos_nulos INTEGER,
    hash_sha256 TEXT,
    created_at INTEGER
);

-- LegalDocument (Dossiês e Atas)
CREATE TABLE IF NOT EXISTS legal_documents (
    id TEXT PRIMARY KEY,
    entity_id TEXT NOT NULL,
    type TEXT NOT NULL, -- ATA, ESTATUTO, DOSSIE_CADSOL, MTSE
    content TEXT,
    hash_sha256 TEXT,
    signature TEXT,
    status TEXT, -- DRAFT, SIGNED, SUBMITTED
    created_at INTEGER
);

-- EligibilityProfile (RF-19)
CREATE TABLE IF NOT EXISTS eligibility_profiles (
    id TEXT PRIMARY KEY,
    entity_id TEXT NOT NULL UNIQUE,
    
    -- Dados do ERP (cópia para consulta rápida)
    cnpj TEXT,
    cnae TEXT,
    municipio TEXT,
    uf TEXT,
    faturamento_anual INTEGER,  -- int64, Anti-Float
    regime_tributario TEXT,
    data_abertura INTEGER,
    situacao_fiscal TEXT,
    
    -- Campos complementares
    inscrito_cad_unico INTEGER,      -- 0/1 (bool)
    socio_mulher INTEGER,            -- 0/1 (bool)
    inadimplencia_ativa INTEGER,     -- 0/1 (bool)
    finalidade_credito TEXT,
    valor_necessario INTEGER,        -- int64, Anti-Float
    tipo_entidade TEXT,
    contabilidade_formal INTEGER,    -- 0/1 (bool)
    
    -- Metadados
    preenchido_em INTEGER,
    atualizado_em INTEGER,
    preenchido_por TEXT,
    created_at INTEGER,
    updated_at INTEGER
);

-- ProgramMatch (RF-20)
CREATE TABLE IF NOT EXISTS program_matches (
    id TEXT PRIMARY KEY,
    entity_id TEXT NOT NULL,
    program_id TEXT NOT NULL,
    elegibilidade TEXT NOT NULL, -- ELEGIVEL, NAO_ELEGIVEL, PARCIAL
    motivo TEXT,
    documentos_pendentes TEXT, -- JSON
    created_at INTEGER
);

-- PublicProfile (RF-24)
CREATE TABLE IF NOT EXISTS public_profiles (
    entity_id TEXT PRIMARY KEY,
    nome_fantasia TEXT,
    missao TEXT,
    produtos TEXT, -- JSON array
    servicos TEXT, -- JSON array
    municipio TEXT,
    uf TEXT,
    contato_publico TEXT,
    foto_logo TEXT,
    created_at INTEGER,
    updated_at INTEGER
);

-- NeedPost (RF-25)
CREATE TABLE IF NOT EXISTS need_posts (
    id TEXT PRIMARY KEY,
    publisher_id TEXT NOT NULL,
    categoria TEXT NOT NULL,
    descricao TEXT,
    quantidade TEXT,
    prazo_desejado INTEGER,
    municipio TEXT,
    uf TEXT,
    status TEXT NOT NULL, -- ABERTO, EM_NEGOCIACAO, CONCLUIDO
    created_at INTEGER,
    updated_at INTEGER
);

-- ReinfEvent (RF-14)
CREATE TABLE IF NOT EXISTS reinf_events (
    id TEXT PRIMARY KEY,
    entity_id TEXT NOT NULL,
    event_type TEXT NOT NULL,
    xml_content TEXT,
    hash_sha256 TEXT,
    status TEXT NOT NULL, -- PENDING, SENT, CONFIRMED, ERROR
    sent_at INTEGER,
    created_at INTEGER
);

-- SanitaryDossier (RF-16)
CREATE TABLE IF NOT EXISTS sanitary_dossiers (
    id TEXT PRIMARY KEY,
    entity_id TEXT NOT NULL,
    tipo_produto TEXT,
    capacidade_diaria INTEGER,
    fluxo_maquinario TEXT,
    origem_agua TEXT,
    content TEXT,
    status TEXT NOT NULL, -- DRAFT, SUBMITTED, APPROVED
    created_at INTEGER
);

-- DASMEI (RF-27)
CREATE TABLE IF NOT EXISTS das_mei (
    id TEXT PRIMARY KEY,
    entity_id TEXT NOT NULL,
    competencia TEXT NOT NULL, -- YYYY-MM
    valor_devido INTEGER NOT NULL, -- int64, Anti-Float
    valor_pago INTEGER DEFAULT 0,
    data_vencimento INTEGER NOT NULL,
    data_pagamento INTEGER DEFAULT 0,
    status TEXT NOT NULL, -- PENDENTE, PAGO, VENCIDO
    salario_minimo INTEGER NOT NULL,
    created_at INTEGER,
    updated_at INTEGER,
    UNIQUE(entity_id, competencia)
);
```

---

## 3. Algoritmos de Negócio e Compliance

### 3.1 Algoritmo de Formalização Gradual e Integração CADSOL
**Objetivo:** Avaliar a maturidade institucional para permitir a transição DREAM -> FORMALIZED.

**Critérios Automatizados:**
1. Mínimo de 3 registros de Decision (Assembleias realizadas provando autogestão)
2. Mínimo de 1 membro ativo com histórico de WorkLog
3. Criação automática do Dossiê Hash SHA256 com validação de Assinaturas ICP-Brasil/Gov.br
4. Integração Automática (MTE): Consumo via integração (módulo `integrations`) das futuras APIs do SINAES/CADSOL (Decreto nº 12.784/2025), substituindo envios em papel pela interoperabilidade digital assim que o grupo atingir o estágio FORMALIZED

### 3.2 Algoritmo de Tradução Fiscal (Ponte do Contador) [NOVO]
**Objetivo:** Converter a contabilidade social gerada invisivelmente pelo produtor em linguagem de conformidade estatal (SPED/Lotes Fiscais/Reinf).

**Processo:**
1. Painel do Contador solicita dados de um Period fechado via filtro temporal `EnterpriseAccountant`
2. O algoritmo compila todas as entries de soma zero (respeitando a ITG 2002)
3. Mapeia as contas locais amigáveis ("Gaveta") para o Plano de Contas Referencial da Receita Federal
4. **Isolamento do Ato Cooperativo:** O sistema segrega as transações com terceiros não cooperados (destinando sobras ao FATES) e estrutura a base de cálculo tributável. As receitas do trabalho intercooperativo são expurgadas automaticamente para preenchimento na Parte A do e-Lalur/e-Lacs (Bloco M da ECF)
5. Gera o pacote CSV/SPED/XMLs da EFD-Reinf e salva o evento de exportação na tabela `fiscal_exports`

### 3.3 Algoritmo de Match de Elegibilidade [NOVO - PDF v1.0]
**Objetivo:** Cruzar automaticamente o perfil da entidade com requisitos de programas de financiamento.

**Processo:**
1. Coleta dados do `EligibilityProfile` (ERP + campos complementares)
2. Coleta taxas atualizadas do `EconomicIndicator` (Motor)
3. Para cada `FinancingProgram` ativo:
   - Verifica requisitos básicos (tipo entidade, faturamento, município)
   - Verifica requisitos específicos (CadÚnico, gênero, inadimplência)
   - Calcula custo efetivo de capital usando taxas do Motor
4. Gera `ProgramMatch` com status e documentos pendentes
5. Ordena por vantagem (menor custo efetivo primeiro)

### 3.4 Algoritmo de Matching da Rede Digna [NOVO - PDF v1.0]
**Objetivo:** Conectar entidades da rede para compras e vendas solidárias.

**Critérios de Matching:**
1. **Geográfico:** Priorizar proximidade (mesmo município/UF)
2. **Setorial:** Afinidade de CNAE/categoria
3. **Temporal:** Prazos compatíveis
4. **Reputação:** Histórico de transações na rede (futuro)

### 3.5 Algoritmo de Help System (RF-30) [NOVO - Decisão 27/03/2026]
**Objetivo:** Traduzir conceitos técnicos em linguagem popular com linkagem UI → banco de ajuda.

**Processo:**
1. UI detecta campo técnico (ex: "Inscrito no CadÚnico?")
2. Botão "?" linka para `/help/topic/{key}` (ex: `/help/topic/cadunico`)
3. Sistema busca `HelpTopic` no `central.db` por `key`
4. Renderiza explicação em linguagem popular + legislação + próximo passo
5. Incrementa `ViewCount` para métricas de uso

---

## 4. Seed Data (Carga Inicial Padrão)

Toda nova base SQLite de um EES nasce com este plano de contas enxuto e adaptado:

| ID | Código | Nome Amigável | Natureza Contábil (Invisível) | Mapeamento Fiscal [NOVO] |
|------|------|------|------|------|
| 1 | 1.1.01 | Gaveta / Caixa | ASSET (Ativo) | Disponibilidades (Ativo) |
| 2 | 3.1.01 | Nossas Vendas | REVENUE (Receita) | Receita Bruta |
| 3 | 1.1.02 | Banco / Conta | ASSET (Ativo) | Contas Bancárias |
| 4 | 2.1.01 | Quem Fornece | LIABILITY (Passivo) | Fornecedores a Pagar |
| 5 | 3.2.01 | Fundo FATES | EQUITY (Patrimônio Líquido) | Reservas Estatutárias |
| 6 | 3.2.02 | Reserva Legal | EQUITY (Patrimônio Líquido) | Reservas de Lucros |
| 7 | 3.2.03 | Capital Social de Trabalho | EQUITY (Patrimônio Líquido) | Capital Social (ITG 2002) |

**Seed de Tópicos de Ajuda (central.db) [NOVO - RF-30]:**

| Key | Title | Category | Explanation (Resumo) |
|-----|-------|----------|---------------------|
| `cadunico` | "O que é o CadÚnico?" | CREDITO | "É o cadastro do governo para programas sociais" |
| `inadimplencia` | "O que é inadimplência?" | CREDITO | "É quando há dívidas não pagas registradas" |
| `cnae` | "O que é CNAE?" | TRIBUTARIO | "É o código que diz qual é a atividade do seu negócio" |
| `das_mei` | "O que é o DAS MEI?" | TRIBUTARIO | "É o boleto mensal que o MEI paga" |
| `reserva_legal` | "O que é Reserva Legal?" | GOVERNANCA | "É uma parte do lucro que a lei manda guardar" |
| `fates` | "O que é o FATES?" | GOVERNANCA | "É um fundo para ajudar outros grupos a se organizarem" |

---

## 5. Regras de Integridade e Validação

### 5.1 Anti-Float (Regra Sagrada)
Todos os campos monetários e de tempo devem usar `int64`:
- `FaturamentoAnual int64` — Centavos
- `Minutes int64` — Minutos trabalhados
- `ValorNecessario int64` — Centavos
- `ValorMaximo int64` — Centavos (programas de crédito)
- `ViewCount int64` — Contagem de visualizações (Help System)

**Validação:** `grep -r "float[0-9]*" modules/` deve retornar apenas logs/comentários.

### 5.2 Soberania de Dados
- Cada entidade tem seu próprio arquivo `.db` em `data/entities/{entity_id}.db`
- Proibido JOIN entre bancos de entidades diferentes
- `central.db` armazena apenas relações inter-tenant (vínculos, indicadores, programas, help topics)

### 5.3 Cache-Proof para Templates
- Templates `*_simple.html` são documentos HTML completos
- Carregados via `ParseFiles()` no handler (não variáveis globais)

### 5.4 Hash de Integridade
- Todo documento legal (`LegalDocument`) deve ter `HashSHA256`
- Todo evento fiscal (`ReinfEvent`) deve ter `HashSHA256`
- Todo registro de decisão (`Decision`) deve ter `HashSHA256`

---

## 6. Matriz de Rastreabilidade (RF → Modelo)

| Requisito | Modelo(s) de Domínio | Status |
|-----------|---------------------|--------|
| RF-01 (Identidade) | `Enterprise` | ✅ Implementado |
| RF-02 (PDV) | `Enterprise`, `Ledger` (core_lume) | ✅ Implementado |
| RF-03 (Trabalho ITG 2002) | `WorkLog` | ✅ Implementado |
| RF-04 (Dossiê CADSOL) | `LegalDocument`, `Decision` | ✅ Implementado |
| RF-11 (Exportação SPED) | `EnterpriseAccountant`, `LegalDocument` | ✅ Implementado |
| RF-12 (Vínculo Contábil) | `EnterpriseAccountant` | ✅ 95% Implementado |
| RF-14 (EFD-Reinf) | `ReinfEvent` | 📋 Backlog |
| RF-16 (MTSE/MAPA) | `SanitaryDossier` | 📋 Backlog |
| RF-18 (Motor Indicadores) | `EconomicIndicator`, `IndicatorCache` | 📋 Backlog |
| RF-19 (Perfil Elegibilidade) | `EligibilityProfile` | 📋 Backlog |
| RF-20 (Portal Oportunidades) | `FinancingProgram`, `ProgramMatch` | 📋 Backlog |
| RF-24 (Perfil Público) | `PublicProfile` | 📋 Backlog |
| RF-25 (Mural Necessidades) | `NeedPost` | 📋 Backlog |
| RF-27 (DAS MEI) | `DASMEI` | 📋 Backlog |
| RF-30 (Sistema de Ajuda) | `HelpTopic` | 📋 Backlog |

---

**Status:** ✅ ATUALIZADO COM VISÃO DE ECOSSISTEMA (PDF v1.0) + RF-30 (Decisão de Design 27/03/2026)  
**Próxima Ação:** Atualizar `03_architecture/01_system.md` com novos módulos e help_engine  
**Versão Anterior:** 1.4 (2026-03-13)  
**Versão Atual:** 2.1 (2026-03-27)
