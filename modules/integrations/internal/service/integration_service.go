// Package service implementa a camada de aplicação para integrações
// Coordena o uso dos repositórios de integração
package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/providentia/digna/integrations/internal/domain"
)

// IntegrationService coordena todas as integrações externas
type IntegrationService struct {
	repo domain.IntegrationRepository
}

// NewIntegrationService cria serviço de integração
func NewIntegrationService(repo domain.IntegrationRepository) *IntegrationService {
	return &IntegrationService{repo: repo}
}

// ReceitaFederal acesso a Receita Federal
func (s *IntegrationService) ReceitaFederal() domain.ReceitaFederalRepository {
	return s.repo.ReceitaFederal()
}

// MTE acesso a MTE
func (s *IntegrationService) MTE() domain.MTERepository {
	return s.repo.MTE()
}

// MDS acesso a MDS
func (s *IntegrationService) MDS() domain.MDSRepository {
	return s.repo.MDS()
}

// IBGE acesso a IBGE
func (s *IntegrationService) IBGE() domain.IBGERepository {
	return s.repo.IBGE()
}

// SEFAZ acesso a SEFAZ
func (s *IntegrationService) SEFAZ() domain.SEFAZRepository {
	return s.repo.SEFAZ()
}

// BNDES acesso a BNDES
func (s *IntegrationService) BNDES() domain.BNDESRepository {
	return s.repo.BNDES()
}

// SEBRAE acesso a SEBRAE
func (s *IntegrationService) SEBRAE() domain.SEBRAERepository {
	return s.repo.SEBRAE()
}

// Providentia acesso a Providentia
func (s *IntegrationService) Providentia() domain.ProvidentiaRepository {
	return s.repo.Providentia()
}

// LogIntegration registra uma integração no log
func (s *IntegrationService) LogIntegration(ctx context.Context, entityID string, integration string, status domain.IntegrationStatus, request, response interface{}) error {
	reqJSON, _ := json.Marshal(request)
	respJSON, _ := json.Marshal(response)

	return s.repo.Log(ctx, &domain.IntegrationLog{
		EntityID:    entityID,
		Integration: integration,
		Status:      status,
		Request:     string(reqJSON),
		Response:    string(respJSON),
		Timestamp:   time.Now(),
	})
}

// GetIntegrationLogs obtém logs de integração
func (s *IntegrationService) GetIntegrationLogs(ctx context.Context, entityID string, integration string, limit int) ([]domain.IntegrationLog, error) {
	return s.repo.GetLogs(ctx, entityID, integration, limit)
}

// Helper methods para casos de uso comuns

// EmitirNFeCompleto emite NFe e registra log
func (s *IntegrationService) EmitirNFeCompleto(ctx context.Context, entityID string, nfe *domain.NFeRequest) (*domain.NFeResponse, error) {
	resp, err := s.SEFAZ().EmitirNFe(ctx, nfe)
	if err != nil {
		s.LogIntegration(ctx, entityID, "SEFAZ.EmitirNFe", domain.StatusFailed, nfe, err.Error())
		return nil, err
	}
	s.LogIntegration(ctx, entityID, "SEFAZ.EmitirNFe", domain.StatusSuccess, nfe, resp)
	return resp, nil
}

// RegistrarRAISCompleto envia RAIS e registra log
func (s *IntegrationService) RegistrarRAISCompleto(ctx context.Context, entityID string, rais *domain.RAISRequest) (*domain.RAISResponse, error) {
	resp, err := s.MTE().EnviarRAIS(ctx, rais)
	if err != nil {
		s.LogIntegration(ctx, entityID, "MTE.EnviarRAIS", domain.StatusFailed, rais, err.Error())
		return nil, err
	}
	s.LogIntegration(ctx, entityID, "MTE.EnviarRAIS", domain.StatusSuccess, rais, resp)
	return resp, nil
}

// EnviarRelatorioSocialCompleto envia relatório ao MDS
func (s *IntegrationService) EnviarRelatorioSocialCompleto(ctx context.Context, entityID string, relatorio *domain.RelatorioSocial) (*domain.RelatorioResponse, error) {
	resp, err := s.MDS().EnviarRelatorioSocial(ctx, relatorio)
	if err != nil {
		s.LogIntegration(ctx, entityID, "MDS.EnviarRelatorioSocial", domain.StatusFailed, relatorio, err.Error())
		return nil, err
	}
	s.LogIntegration(ctx, entityID, "MDS.EnviarRelatorioSocial", domain.StatusSuccess, relatorio, resp)
	return resp, nil
}

// SimularCreditoBNDES simula crédito e retorna resultado
func (s *IntegrationService) SimularCreditoBNDES(ctx context.Context, entityID string, simulacao *domain.SimulacaoCredito) (*domain.ResultadoSimulacao, error) {
	return s.BNDES().SimularCredito(ctx, simulacao)
}
