package eligibility

import (
	"context"

	"github.com/providentia/digna/core_lume/internal/domain"
	"github.com/providentia/digna/core_lume/internal/repository"
	"github.com/providentia/digna/core_lume/internal/service"
	"github.com/providentia/digna/lifecycle/pkg/lifecycle"
)

// EligibilityProfile represents the public eligibility profile model
type EligibilityProfile struct {
	ID       string
	EntityID string

	// Dados do ERP
	CNPJ             string
	CNAE             string
	Municipio        string
	UF               string
	FaturamentoAnual int64
	RegimeTributario string
	DataAbertura     int64
	SituacaoFiscal   string

	// Campos complementares
	InscritoCadUnico    bool
	SocioMulher         bool
	InadimplenciaAtiva  bool
	FinalidadeCredito   string
	ValorNecessario     int64
	TipoEntidade        string
	ContabilidadeFormal bool

	// Metadados
	PreenchidoEm  int64
	AtualizadoEm  int64
	PreenchidoPor string

	CreatedAt int64
	UpdatedAt int64
}

// EligibilityInput represents input for creating/updating profile
type EligibilityInput struct {
	InscritoCadUnico    *bool   `json:"inscrito_cad_unico,omitempty"`
	SocioMulher         *bool   `json:"socio_mulher,omitempty"`
	InadimplenciaAtiva  *bool   `json:"inadimplencia_ativa,omitempty"`
	FinalidadeCredito   *string `json:"finalidade_credito,omitempty"`
	ValorNecessario     *int64  `json:"valor_necessario,omitempty"`
	TipoEntidade        *string `json:"tipo_entidade,omitempty"`
	ContabilidadeFormal *bool   `json:"contabilidade_formal,omitempty"`
}

// Service provides eligibility profile operations
type Service struct {
	eligService *service.EligibilityService
}

// NewService creates a new eligibility service
func NewService(lm lifecycle.LifecycleManager) *Service {
	eligRepo := repository.NewSQLiteEligibilityRepository(lm)
	memberRepo := repository.NewSQLiteMemberRepository(lm)
	return &Service{
		eligService: service.NewEligibilityService(eligRepo, memberRepo),
	}
}

// CreateOrUpdate creates or updates an eligibility profile
func (s *Service) CreateOrUpdate(ctx context.Context, entityID string, userID string, input EligibilityInput) (*EligibilityProfile, error) {
	// Use system user if empty
	if userID == "" {
		userID = "system"
	}

	domainInput := domain.EligibilityInput{
		InscritoCadUnico:    input.InscritoCadUnico,
		SocioMulher:         input.SocioMulher,
		InadimplenciaAtiva:  input.InadimplenciaAtiva,
		FinalidadeCredito:   input.FinalidadeCredito,
		ValorNecessario:     input.ValorNecessario,
		TipoEntidade:        input.TipoEntidade,
		ContabilidadeFormal: input.ContabilidadeFormal,
	}

	profile, err := s.eligService.CreateOrUpdate(ctx, entityID, userID, domainInput)
	if err != nil {
		return nil, err
	}

	return convertToPublic(profile), nil
}

// GetProfile retrieves the eligibility profile
func (s *Service) GetProfile(ctx context.Context, entityID string) (*EligibilityProfile, error) {
	profile, err := s.eligService.GetProfile(ctx, entityID)
	if err != nil {
		return nil, err
	}
	return convertToPublic(profile), nil
}

// GetOrCreateProfile gets existing profile or creates empty one
func (s *Service) GetOrCreateProfile(ctx context.Context, entityID string, userID string) (*EligibilityProfile, error) {
	profile, err := s.eligService.GetOrCreateProfile(ctx, entityID, userID)
	if err != nil {
		return nil, err
	}
	return convertToPublic(profile), nil
}

// GetCompletionStatus returns the completion percentage
func (s *Service) GetCompletionStatus(ctx context.Context, entityID string) (float64, error) {
	return s.eligService.GetCompletionStatus(ctx, entityID)
}

// EnsureTableExists ensures the table exists for the entity
func (s *Service) EnsureTableExists(entityID string) error {
	return s.eligService.InitTableForEntity(entityID)
}

// CanUserEditProfile checks if user can edit
func (s *Service) CanUserEditProfile(ctx context.Context, entityID string, userID string) (bool, error) {
	return s.eligService.CanUserEditProfile(ctx, entityID, userID)
}

// Helper function to convert internal to public
func convertToPublic(profile *domain.EligibilityProfile) *EligibilityProfile {
	return &EligibilityProfile{
		ID:                  profile.ID,
		EntityID:            profile.EntityID,
		CNPJ:                profile.CNPJ,
		CNAE:                profile.CNAE,
		Municipio:           profile.Municipio,
		UF:                  profile.UF,
		FaturamentoAnual:    profile.FaturamentoAnual,
		RegimeTributario:    profile.RegimeTributario,
		DataAbertura:        profile.DataAbertura,
		SituacaoFiscal:      profile.SituacaoFiscal,
		InscritoCadUnico:    profile.InscritoCadUnico,
		SocioMulher:         profile.SocioMulher,
		InadimplenciaAtiva:  profile.InadimplenciaAtiva,
		FinalidadeCredito:   string(profile.FinalidadeCredito),
		ValorNecessario:     profile.ValorNecessario,
		TipoEntidade:        string(profile.TipoEntidade),
		ContabilidadeFormal: profile.ContabilidadeFormal,
		PreenchidoEm:        profile.PreenchidoEm,
		AtualizadoEm:        profile.AtualizadoEm,
		PreenchidoPor:       profile.PreenchidoPor,
		CreatedAt:           profile.CreatedAt,
		UpdatedAt:           profile.UpdatedAt,
	}
}
