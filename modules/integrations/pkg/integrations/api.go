// Package integrations expõe a API pública do módulo de integrações
package integrations

import (
	"database/sql"

	"github.com/providentia/digna/integrations/internal/domain"
	"github.com/providentia/digna/integrations/internal/repository"
	"github.com/providentia/digna/integrations/internal/service"
)

// IntegrationService é o serviço principal de integrações
type IntegrationService = service.IntegrationService

// Interfaces exportadas do domínio
type (
	// ReceitaFederalRepository interface para Receita Federal
	ReceitaFederalRepository = domain.ReceitaFederalRepository
	// MTERepository interface para MTE
	MTERepository = domain.MTERepository
	// MDSRepository interface para MDS
	MDSRepository = domain.MDSRepository
	// IBGERepository interface para IBGE
	IBGERepository = domain.IBGERepository
	// SEFAZRepository interface para SEFAZ
	SEFAZRepository = domain.SEFAZRepository
	// BNDESRepository interface para BNDES
	BNDESRepository = domain.BNDESRepository
	// SEBRAERepository interface para SEBRAE
	SEBRAERepository = domain.SEBRAERepository
	// ProvidentiaRepository interface para Providentia
	ProvidentiaRepository = domain.ProvidentiaRepository
)

// Estruturas de dados exportadas
type (
	CNPJData            = domain.CNPJData
	Endereco            = domain.Endereco
	DARFRequest         = domain.DARFRequest
	DARFResponse        = domain.DARFResponse
	CATRequest          = domain.CATRequest
	CATResponse         = domain.CATResponse
	RAISRequest         = domain.RAISRequest
	TrabalhadorRAIS     = domain.TrabalhadorRAIS
	RAISResponse        = domain.RAISResponse
	ESocialEvent        = domain.ESocialEvent
	ESocialResponse     = domain.ESocialResponse
	FamiliaCadUnico     = domain.FamiliaCadUnico
	PessoaCadUnico      = domain.PessoaCadUnico
	CadUnicoResponse    = domain.CadUnicoResponse
	RelatorioSocial     = domain.RelatorioSocial
	DistribuicaoSocial  = domain.DistribuicaoSocial
	RelatorioResponse   = domain.RelatorioResponse
	PesquisaIBGE        = domain.PesquisaIBGE
	PesquisaResponse    = domain.PesquisaResponse
	DadosIBGE           = domain.DadosIBGE
	ProdutoIBGE         = domain.ProdutoIBGE
	ProdutoResponse     = domain.ProdutoResponse
	DestinatarioNFe     = domain.DestinatarioNFe
	ItemNFe             = domain.ItemNFe
	NFeRequest          = domain.NFeRequest
	NFeResponse         = domain.NFeResponse
	ManifestoRequest    = domain.ManifestoRequest
	ManifestoResponse   = domain.ManifestoResponse
	NFSRequest          = domain.NFSRequest
	NFSResponse         = domain.NFSResponse
	LinhaCredito        = domain.LinhaCredito
	SimulacaoCredito    = domain.SimulacaoCredito
	ResultadoSimulacao  = domain.ResultadoSimulacao
	SolicitacaoCredito  = domain.SolicitacaoCredito
	SolicitacaoResponse = domain.SolicitacaoResponse
	CursoSEBRAE         = domain.CursoSEBRAE
	InscricaoCurso      = domain.InscricaoCurso
	InscricaoResponse   = domain.InscricaoResponse
	ProgramaConsultoria = domain.ProgramaConsultoria
	SyncPackage         = domain.SyncPackage
	AggregatedMetrics   = domain.AggregatedMetrics
	SyncResponse        = domain.SyncResponse
	MarketplaceData     = domain.MarketplaceData
	MarketplaceOffer    = domain.MarketplaceOffer
	OfferResponse       = domain.OfferResponse
)

// IntegrationStatus exportado
type IntegrationStatus = domain.IntegrationStatus

const (
	StatusSuccess = domain.StatusSuccess
	StatusPending = domain.StatusPending
	StatusFailed  = domain.StatusFailed
	StatusMock    = domain.StatusMock
)

// NewMockIntegrationService cria serviço com implementações mock
func NewMockIntegrationService(db *sql.DB) (*IntegrationService, error) {
	mockRepo := repository.NewMockIntegrationRepository(db)
	if err := mockRepo.InitSchema(); err != nil {
		return nil, err
	}
	return service.NewIntegrationService(mockRepo), nil
}

// NewIntegrationService cria serviço com repositório customizado
func NewIntegrationService(repo domain.IntegrationRepository) *IntegrationService {
	return service.NewIntegrationService(repo)
}
