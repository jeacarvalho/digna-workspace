// Package domain define as interfaces de integração externa seguindo DDD
// Todas as integrações são abstrações (interfaces) para manter o domínio
// desacoplado de detalhes de infraestrutura
package domain

import (
	"context"
	"time"
)

// IntegrationStatus representa o status de uma integração
type IntegrationStatus string

const (
	StatusSuccess IntegrationStatus = "SUCCESS"
	StatusPending IntegrationStatus = "PENDING"
	StatusFailed  IntegrationStatus = "FAILED"
	StatusMock    IntegrationStatus = "MOCK"
)

// IntegrationLog representa um registro de integração
type IntegrationLog struct {
	ID          int64
	EntityID    string
	Integration string
	Status      IntegrationStatus
	Request     string
	Response    string
	Error       string
	Timestamp   time.Time
}

// ==================== RECEITA FEDERAL ====================

// ReceitaFederalRepository interface para integração com Receita Federal
type ReceitaFederalRepository interface {
	// ConsultarCNPJ consulta dados de um CNPJ na Receita Federal
	ConsultarCNPJ(ctx context.Context, cnpj string) (*CNPJData, error)
	// EmitirDARF emite guia DARF para pagamento de impostos
	EmitirDARF(ctx context.Context, darf *DARFRequest) (*DARFResponse, error)
	// ConsultarDARF consulta status de uma guia DARF
	ConsultarDARF(ctx context.Context, numero string) (*DARFResponse, error)
}

// CNPJData representa dados de uma empresa na Receita Federal
type CNPJData struct {
	CNPJ             string
	RazaoSocial      string
	NomeFantasia     string
	Situacao         string
	DataAbertura     time.Time
	Endereco         Endereco
	NaturezaJuridica string
	CapitalSocial    float64
}

// Endereco representa endereço
type Endereco struct {
	Logradouro  string
	Numero      string
	Complemento string
	Bairro      string
	Cidade      string
	UF          string
	CEP         string
}

// DARFRequest representa requisição de emissão de DARF
type DARFRequest struct {
	EntityID      string
	CPF           string
	CodigoReceita int
	Periodo       string
	Valor         int64 // centavos
	Vencimento    time.Time
}

// DARFResponse representa resposta de DARF
type DARFResponse struct {
	Numero         string
	CodigoBarra    string
	LinhaDigitavel string
	Valor          int64
	Vencimento     time.Time
	Status         string
}

// ==================== MTE (MINISTÉRIO DO TRABALHO) ====================

// MTERepository interface para integração com Ministério do Trabalho
type MTERepository interface {
	// RegistrarCAT registra Comunicação de Acidente de Trabalho
	RegistrarCAT(ctx context.Context, cat *CATRequest) (*CATResponse, error)
	// EnviarRAIS envia declaração RAIS
	EnviarRAIS(ctx context.Context, rais *RAISRequest) (*RAISResponse, error)
	// ConsultarRAIS consulta declaração RAIS
	ConsultarRAIS(ctx context.Context, ano int, cnpj string) (*RAISResponse, error)
	// EnviarESocial envia eventos do eSocial
	EnviarESocial(ctx context.Context, evento *ESocialEvent) (*ESocialResponse, error)
}

// CATRequest representa requisição de CAT
type CATRequest struct {
	EntityID      string
	CNPJ          string
	DataAcidente  time.Time
	TipoAcidente  string
	HoraAcidente  string
	LocalAcidente string
	Descricao     string
	Comunicante   string
}

// CATResponse representa resposta de CAT
type CATResponse struct {
	Numero    string
	Protocolo string
	Status    string
	DataEnvio time.Time
}

// RAISRequest representa requisição de RAIS
type RAISRequest struct {
	EntityID      string
	Ano           int
	CNPJ          string
	Trabalhadores []TrabalhadorRAIS
}

// TrabalhadorRAIS representa dados de trabalhador na RAIS
type TrabalhadorRAIS struct {
	CPF           string
	Nome          string
	DataAdmissao  time.Time
	CBO           string // Código Brasileiro de Ocupação
	Salario       int64  // centavos
	HorasSemanais int
}

// RAISResponse representa resposta de RAIS
type RAISResponse struct {
	Protocolo string
	Status    string
	Ano       int
	CNPJ      string
	DataEnvio time.Time
}

// ESocialEvent representa evento do eSocial
type ESocialEvent struct {
	EntityID string
	Tipo     string // S-1000, S-2200, S-2299, etc.
	XML      string
}

// ESocialResponse representa resposta do eSocial
type ESocialResponse struct {
	Protocolo   string
	Status      string
	DataEnvio   time.Time
	DataRetorno *time.Time
	Mensagem    string
}

// ==================== MDS (MINISTÉRIO DO DESENVOLVIMENTO SOCIAL) ====================

// MDSRepository interface para integração com MDS
type MDSRepository interface {
	// CadastrarFamilia cadastra família no CadÚnico
	CadastrarFamilia(ctx context.Context, familia *FamiliaCadUnico) (*CadUnicoResponse, error)
	// ConsultarFamilia consulta dados de família no CadÚnico
	ConsultarFamilia(ctx context.Context, codigo string) (*FamiliaCadUnico, error)
	// AtualizarRenda atualiza renda familiar
	AtualizarRenda(ctx context.Context, codigo string, renda int64) error
	// EnviarRelatorioSocial envia relatório de impacto social
	EnviarRelatorioSocial(ctx context.Context, relatorio *RelatorioSocial) (*RelatorioResponse, error)
}

// FamiliaCadUnico representa família no CadÚnico
type FamiliaCadUnico struct {
	Codigo         string
	NomeReferencia string
	CPFReferencia  string
	Endereco       Endereco
	RendaTotal     int64
	Integrantes    []PessoaCadUnico
}

// PessoaCadUnico representa pessoa no CadÚnico
type PessoaCadUnico struct {
	CPF        string
	Nome       string
	DataNasc   time.Time
	Parentesco string
	Trabalha   bool
	Renda      int64
}

// CadUnicoResponse representa resposta do CadÚnico
type CadUnicoResponse struct {
	Codigo    string
	Protocolo string
	Status    string
	Score     float64
}

// RelatorioSocial representa relatório de impacto social
type RelatorioSocial struct {
	EntityID     string
	Periodo      string
	TotalHoras   int64
	Membros      int
	Distribuicao []DistribuicaoSocial
	Observacoes  string
}

// DistribuicaoSocial representa distribuição de recursos
type DistribuicaoSocial struct {
	MembroID string
	Horas    int64
	Valor    int64
	Tipo     string // TRABALHO, INVESTIMENTO, etc.
}

// RelatorioResponse representa resposta de relatório
type RelatorioResponse struct {
	Protocolo string
	Status    string
	DataEnvio time.Time
}

// ==================== IBGE ====================

// IBGERepository interface para integração com IBGE
type IBGERepository interface {
	// EnviarPesquisa envia dados de pesquisa (PNAD, PAM, etc.)
	EnviarPesquisa(ctx context.Context, pesquisa *PesquisaIBGE) (*PesquisaResponse, error)
	// ConsultarDados consulta dados estatísticos
	ConsultarDados(ctx context.Context, indicador string, params map[string]string) (*DadosIBGE, error)
	// CadastrarProduto cadastra produto na classificação CNAE/IBGE
	CadastrarProduto(ctx context.Context, produto *ProdutoIBGE) (*ProdutoResponse, error)
}

// PesquisaIBGE representa pesquisa do IBGE
type PesquisaIBGE struct {
	Tipo      string // PNAD, PAM, PIA, etc.
	EntityID  string
	Ano       int
	Trimestre int
	Dados     map[string]interface{}
}

// PesquisaResponse representa resposta de pesquisa
type PesquisaResponse struct {
	Protocolo string
	Status    string
	DataEnvio time.Time
}

// DadosIBGE representa dados do IBGE
type DadosIBGE struct {
	Indicador string
	Valor     float64
	Periodo   string
	Unidade   string
}

// ProdutoIBGE representa produto para classificação
type ProdutoIBGE struct {
	Codigo     string
	Nome       string
	Descricao  string
	CNAE       string
	Atividade  string
	Quantidade int64
	Valor      int64
}

// ProdutoResponse representa resposta de produto
type ProdutoResponse struct {
	CodigoIBGE string
	Status     string
}

// ==================== SEFAZ / NOTAS FISCAIS ====================

// SEFAZRepository interface para integração com SEFAZ
type SEFAZRepository interface {
	// EmitirNFe emite Nota Fiscal Eletrônica
	EmitirNFe(ctx context.Context, nfe *NFeRequest) (*NFeResponse, error)
	// CancelarNFe cancela NFe
	CancelarNFe(ctx context.Context, chave string, justificativa string) (*NFeResponse, error)
	// ConsultarNFe consulta status de NFe
	ConsultarNFe(ctx context.Context, chave string) (*NFeResponse, error)
	// EnviarManifestoDestinatario envia manifesto do destinatário
	EnviarManifestoDestinatario(ctx context.Context, manifesto *ManifestoRequest) (*ManifestoResponse, error)
	// EmitirNFS e emite Nota Fiscal de Serviços (quando aplicável)
	EmitirNFS(ctx context.Context, nfs *NFSRequest) (*NFSResponse, error)
}

// NFeRequest representa requisição de NFe
type NFeRequest struct {
	EntityID     string
	CNPJ         string
	Serie        string
	Numero       int64
	DataEmissao  time.Time
	Destinatario DestinatarioNFe
	Itens        []ItemNFe
	Total        int64
	NaturezaOp   string
}

// DestinatarioNFe representa destinatário de NFe
type DestinatarioNFe struct {
	CNPJCPF   string
	Nome      string
	Endereco  Endereco
	IndIEDest string // 1=Contribuinte, 2=Isento, 9=Não contribuinte
}

// ItemNFe representa item de NFe
type ItemNFe struct {
	Numero        int
	Codigo        string
	Descricao     string
	NCM           string // Nomenclatura Comum do Mercosul
	CFOP          string // Código Fiscal de Operações
	Unidade       string
	Quantidade    float64
	ValorUnitario int64
	ValorTotal    int64
}

// NFeResponse representa resposta de NFe
type NFeResponse struct {
	Chave           string
	Protocolo       string
	Status          string // AUTORIZADA, REJEITADA, DENEGADA
	DataAutorizacao *time.Time
	XML             string
	PDF             []byte
}

// ManifestoRequest representa requisição de manifesto
type ManifestoRequest struct {
	EntityID      string
	CNPJ          string
	ChaveNFe      string
	TipoManifesto string // 210200=Confirmação, 210210=Ciência, 210220=Desconhecimento, 210240=Operação não realizada
}

// ManifestoResponse representa resposta de manifesto
type ManifestoResponse struct {
	Protocolo string
	Status    string
}

// NFSRequest representa requisição de NFS-e
type NFSRequest struct {
	EntityID           string
	CNPJ               string
	InscricaoMunicipal string
	Numero             int64
	DataEmissao        time.Time
	Destinatario       DestinatarioNFe
	Discriminacao      string
	ValorServicos      int64
	ValorDeducoes      int64
	ValorIss           int64
	AliquotaIss        float64
}

// NFSResponse representa resposta de NFS-e
type NFSResponse struct {
	Numero            int64
	CodigoVerificacao string
	Status            string
	Link              string
}

// ==================== BNDES ====================

// BNDESRepository interface para integração com BNDES
type BNDESRepository interface {
	// ConsultarLinhasCredito consulta linhas de crédito disponíveis
	ConsultarLinhasCredito(ctx context.Context, cnpj string) ([]LinhaCredito, error)
	// SimularCredito simula operação de crédito
	SimularCredito(ctx context.Context, simulacao *SimulacaoCredito) (*ResultadoSimulacao, error)
	// SolicitarCredito solicita operação de crédito
	SolicitarCredito(ctx context.Context, solicitacao *SolicitacaoCredito) (*SolicitacaoResponse, error)
	// ConsultarSolicitacao consulta status de solicitação
	ConsultarSolicitacao(ctx context.Context, numero string) (*SolicitacaoCredito, error)
}

// LinhaCredito representa linha de crédito
type LinhaCredito struct {
	Codigo      string
	Nome        string
	Descricao   string
	ValorMaximo int64
	PrazoMaximo int // meses
	TaxaJuros   float64
}

// SimulacaoCredito representa simulação
type SimulacaoCredito struct {
	EntityID string
	CNPJ     string
	Linha    string
	Valor    int64
	Prazo    int
	Carencia int
}

// ResultadoSimulacao representa resultado
type ResultadoSimulacao struct {
	ValorParcela    int64
	TotalJuros      int64
	CET             float64 // Custo Efetivo Total
	Prazo           int
	PrimeiraParcela time.Time
}

// SolicitacaoCredito representa solicitação
type SolicitacaoCredito struct {
	Numero          string
	EntityID        string
	CNPJ            string
	Linha           string
	Valor           int64
	Prazo           int
	Carencia        int
	Finalidade      string
	Status          string
	DataSolicitacao time.Time
}

// SolicitacaoResponse representa resposta
type SolicitacaoResponse struct {
	Numero       string
	Protocolo    string
	Status       string
	DataPrevisao *time.Time
}

// ==================== SEBRAE ====================

// SEBRAERepository interface para integração com SEBRAE
type SEBRAERepository interface {
	// ConsultarCursos consulta cursos de capacitação disponíveis
	ConsultarCursos(ctx context.Context, cnpj string) ([]CursoSEBRAE, error)
	// InscreverCurso inscreve cooperado em curso
	InscreverCurso(ctx context.Context, inscricao *InscricaoCurso) (*InscricaoResponse, error)
	// ConsultarInscricao consulta status de inscrição
	ConsultarInscricao(ctx context.Context, codigo string) (*InscricaoCurso, error)
	// ConsultarConsultoria consulta programas de consultoria
	ConsultarConsultoria(ctx context.Context, cnpj string) ([]ProgramaConsultoria, error)
}

// CursoSEBRAE representa curso
type CursoSEBRAE struct {
	Codigo       string
	Nome         string
	Descricao    string
	CargaHoraria int
	Modalidade   string // PRESENCIAL, ONLINE, EAD
	DataInicio   *time.Time
	Vagas        int
}

// InscricaoCurso representa inscrição
type InscricaoCurso struct {
	Codigo           string
	EntityID         string
	CNPJ             string
	CursoCodigo      string
	CPFParticipante  string
	NomeParticipante string
	Status           string // INSCRITO, CONFIRMADO, CANCELADO, CONCLUIDO
	DataInscricao    time.Time
}

// InscricaoResponse representa resposta
type InscricaoResponse struct {
	Codigo    string
	Protocolo string
	Status    string
}

// ProgramaConsultoria representa programa
type ProgramaConsultoria struct {
	Codigo    string
	Nome      string
	Descricao string
	Area      string // GESTAO, MARKETING, FINANCAS, PRODUCAO
	Duracao   string
}

// ==================== PROVIDENTIA / SYNC ====================

// ProvidentiaRepository interface para integração com nuvem Providentia
type ProvidentiaRepository interface {
	// SyncPackage envia pacote de sincronização
	SyncPackage(ctx context.Context, pkg *SyncPackage) (*SyncResponse, error)
	// GetMarketplaceData obtém dados do marketplace solidário
	GetMarketplaceData(ctx context.Context, entityID string) (*MarketplaceData, error)
	// RegisterOffer registra oferta no marketplace
	RegisterOffer(ctx context.Context, offer *MarketplaceOffer) (*OfferResponse, error)
}

// SyncPackage representa pacote de sincronização
type SyncPackage struct {
	EntityID       string
	Timestamp      int64
	PeriodStart    int64
	PeriodEnd      int64
	ChainDigest    string
	Signature      string
	AggregatedData AggregatedMetrics
	DeltaCount     int64
}

// AggregatedMetrics representa métricas agregadas
type AggregatedMetrics struct {
	TotalSales     int64
	TotalWorkHours int64
	TotalMembers   int64
	LegalStatus    string
	ActiveOffers   int64
	DecisionCount  int64
}

// SyncResponse representa resposta
type SyncResponse struct {
	Protocolo string
	Status    string
	Timestamp int64
}

// MarketplaceData representa dados do marketplace
type MarketplaceData struct {
	Offers      []MarketplaceOffer
	TotalActive int
}

// MarketplaceOffer representa oferta
type MarketplaceOffer struct {
	ID          string
	EntityID    string
	ProductName string
	Description string
	Quantity    int64
	Price       int64
	Unit        string
	Active      bool
}

// OfferResponse representa resposta
type OfferResponse struct {
	ID     string
	Status string
	Link   string
}

// ==================== INTEGRATION MANAGER ====================

// IntegrationRepository interface unificada para acesso a todas as integrações
type IntegrationRepository interface {
	// Log registra uma integração
	Log(ctx context.Context, log *IntegrationLog) error
	// GetLogs obtém logs de integração
	GetLogs(ctx context.Context, entityID string, integration string, limit int) ([]IntegrationLog, error)

	// ReceitaFederal acesso a Receita Federal
	ReceitaFederal() ReceitaFederalRepository
	// MTE acesso a MTE
	MTE() MTERepository
	// MDS acesso a MDS
	MDS() MDSRepository
	// IBGE acesso a IBGE
	IBGE() IBGERepository
	// SEFAZ acesso a SEFAZ
	SEFAZ() SEFAZRepository
	// BNDES acesso a BNDES
	BNDES() BNDESRepository
	// SEBRAE acesso a SEBRAE
	SEBRAE() SEBRAERepository
	// Providentia acesso a Providentia
	Providentia() ProvidentiaRepository
}
