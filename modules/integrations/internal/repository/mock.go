// Package repository implementa integrações mockadas em SQLite
// Facilita transição futura para implementações HTTP reais
package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/providentia/digna/integrations/internal/domain"
)

// MockIntegrationRepository implementa domain.IntegrationRepository com mocks
type MockIntegrationRepository struct {
	db *sql.DB
}

// NewMockIntegrationRepository cria repositório mock
func NewMockIntegrationRepository(db *sql.DB) *MockIntegrationRepository {
	return &MockIntegrationRepository{
		db: db,
	}
}

// InitSchema cria tabelas necessárias para logs
func (r *MockIntegrationRepository) InitSchema() error {
	_, err := r.db.Exec(`
		CREATE TABLE IF NOT EXISTS integration_logs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			entity_id TEXT NOT NULL,
			integration TEXT NOT NULL,
			status TEXT NOT NULL,
			request TEXT,
			response TEXT,
			error TEXT,
			timestamp INTEGER DEFAULT (strftime('%s', 'now'))
		);
		CREATE INDEX IF NOT EXISTS idx_integration_logs_entity 
			ON integration_logs(entity_id, integration, timestamp DESC);
	`)
	return err
}

// Log registra uma integração
func (r *MockIntegrationRepository) Log(ctx context.Context, log *domain.IntegrationLog) error {
	requestJSON, _ := json.Marshal(log.Request)
	responseJSON, _ := json.Marshal(log.Response)

	_, err := r.db.Exec(
		"INSERT INTO integration_logs (entity_id, integration, status, request, response, error, timestamp) VALUES (?, ?, ?, ?, ?, ?, ?)",
		log.EntityID, log.Integration, log.Status, string(requestJSON), string(responseJSON), log.Error, log.Timestamp.Unix(),
	)
	return err
}

// GetLogs obtém logs
func (r *MockIntegrationRepository) GetLogs(ctx context.Context, entityID string, integration string, limit int) ([]domain.IntegrationLog, error) {
	rows, err := r.db.Query(
		"SELECT id, entity_id, integration, status, request, response, error, timestamp FROM integration_logs WHERE entity_id = ? AND integration = ? ORDER BY timestamp DESC LIMIT ?",
		entityID, integration, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []domain.IntegrationLog
	for rows.Next() {
		var log domain.IntegrationLog
		var timestamp int64
		err := rows.Scan(&log.ID, &log.EntityID, &log.Integration, &log.Status, &log.Request, &log.Response, &log.Error, &timestamp)
		if err != nil {
			return nil, err
		}
		log.Timestamp = time.Unix(timestamp, 0)
		logs = append(logs, log)
	}
	return logs, nil
}

// ReceitaFederal retorna implementação mock
func (r *MockIntegrationRepository) ReceitaFederal() domain.ReceitaFederalRepository {
	return &MockReceitaFederalRepository{db: r.db}
}

// MTE retorna implementação mock
func (r *MockIntegrationRepository) MTE() domain.MTERepository {
	return &MockMTERepository{db: r.db}
}

// MDS retorna implementação mock
func (r *MockIntegrationRepository) MDS() domain.MDSRepository {
	return &MockMDSRepository{db: r.db}
}

// IBGE retorna implementação mock
func (r *MockIntegrationRepository) IBGE() domain.IBGERepository {
	return &MockIBGERepository{db: r.db}
}

// SEFAZ retorna implementação mock
func (r *MockIntegrationRepository) SEFAZ() domain.SEFAZRepository {
	return &MockSEFAZRepository{db: r.db}
}

// BNDES retorna implementação mock
func (r *MockIntegrationRepository) BNDES() domain.BNDESRepository {
	return &MockBNDESRepository{db: r.db}
}

// SEBRAE retorna implementação mock
func (r *MockIntegrationRepository) SEBRAE() domain.SEBRAERepository {
	return &MockSEBRAERepository{db: r.db}
}

// Providentia retorna implementação mock
func (r *MockIntegrationRepository) Providentia() domain.ProvidentiaRepository {
	return &MockProvidentiaRepository{db: r.db}
}

// ==================== MOCK RECEITA FEDERAL ====================

type MockReceitaFederalRepository struct {
	db *sql.DB
}

func (r *MockReceitaFederalRepository) ConsultarCNPJ(ctx context.Context, cnpj string) (*domain.CNPJData, error) {
	// Simula consulta retornando dados mock
	data := &domain.CNPJData{
		CNPJ:             cnpj,
		RazaoSocial:      "Cooperativa de Produção Agropecuária Ltda",
		NomeFantasia:     "CoopAgro",
		Situacao:         "ATIVA",
		DataAbertura:     time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
		NaturezaJuridica: " Cooperativa",
		CapitalSocial:    10000.00,
		Endereco: domain.Endereco{
			Logradouro: "Rua das Flores",
			Numero:     "123",
			Bairro:     "Centro",
			Cidade:     "São Paulo",
			UF:         "SP",
			CEP:        "01000-000",
		},
	}

	r.logIntegration("RECEITA_FEDERAL", "ConsultarCNPJ", "SUCCESS", cnpj, "CNPJ consultado com sucesso")
	return data, nil
}

func (r *MockReceitaFederalRepository) EmitirDARF(ctx context.Context, darf *domain.DARFRequest) (*domain.DARFResponse, error) {
	response := &domain.DARFResponse{
		Numero:      fmt.Sprintf("DARF%09d", rand.Int63n(1000000000)),
		CodigoBarra: fmt.Sprintf("%044d", rand.Int63()),
		LinhaDigitavel: fmt.Sprintf("%d%d.%d%d %d%d.%d%d %d%d.%d%d %d%d",
			rand.Intn(10), rand.Intn(10), rand.Intn(10), rand.Intn(10),
			rand.Intn(10), rand.Intn(10), rand.Intn(10), rand.Intn(10),
			rand.Intn(10), rand.Intn(10), rand.Intn(10), rand.Intn(10),
			rand.Intn(10), rand.Intn(10)),
		Valor:      darf.Valor,
		Vencimento: darf.Vencimento,
		Status:     "EMITIDA",
	}

	r.logIntegration("RECEITA_FEDERAL", "EmitirDARF", "SUCCESS", darf, response)
	return response, nil
}

func (r *MockReceitaFederalRepository) ConsultarDARF(ctx context.Context, numero string) (*domain.DARFResponse, error) {
	return &domain.DARFResponse{
		Numero:     numero,
		Status:     "PAGO",
		Valor:      150000,
		Vencimento: time.Now().AddDate(0, 0, 30),
	}, nil
}

func (r *MockReceitaFederalRepository) logIntegration(integration, method, status string, request, response interface{}) {
	// Log assíncrono não-bloqueante
	go func() {
		reqJSON, _ := json.Marshal(request)
		respJSON, _ := json.Marshal(response)
		r.db.Exec(
			"INSERT INTO integration_logs (entity_id, integration, status, request, response, timestamp) VALUES (?, ?, ?, ?, ?, ?)",
			"global", integration+"/"+method, status, string(reqJSON), string(respJSON), time.Now().Unix(),
		)
	}()
}

// ==================== MOCK MTE ====================

type MockMTERepository struct {
	db *sql.DB
}

func (r *MockMTERepository) RegistrarCAT(ctx context.Context, cat *domain.CATRequest) (*domain.CATResponse, error) {
	response := &domain.CATResponse{
		Numero:    fmt.Sprintf("CAT%012d", rand.Int63n(1000000000000)),
		Protocolo: fmt.Sprintf("PROT%015d", rand.Int63n(1000000000000000)),
		Status:    "REGISTRADO",
		DataEnvio: time.Now(),
	}
	r.logIntegration("MTE", "RegistrarCAT", "SUCCESS", cat, response)
	return response, nil
}

func (r *MockMTERepository) EnviarRAIS(ctx context.Context, rais *domain.RAISRequest) (*domain.RAISResponse, error) {
	response := &domain.RAISResponse{
		Protocolo: fmt.Sprintf("RAIS%013d", rand.Int63n(10000000000000)),
		Status:    "ENVIADO",
		Ano:       rais.Ano,
		CNPJ:      rais.CNPJ,
		DataEnvio: time.Now(),
	}
	r.logIntegration("MTE", "EnviarRAIS", "SUCCESS", rais, response)
	return response, nil
}

func (r *MockMTERepository) ConsultarRAIS(ctx context.Context, ano int, cnpj string) (*domain.RAISResponse, error) {
	return &domain.RAISResponse{
		Protocolo: "RAIS20240000001",
		Status:    "PROCESSADO",
		Ano:       ano,
		CNPJ:      cnpj,
		DataEnvio: time.Now().AddDate(0, -1, 0),
	}, nil
}

func (r *MockMTERepository) EnviarESocial(ctx context.Context, evento *domain.ESocialEvent) (*domain.ESocialResponse, error) {
	response := &domain.ESocialResponse{
		Protocolo: fmt.Sprintf("ES%015d", rand.Int63n(1000000000000000)),
		Status:    "PROCESSADO",
		DataEnvio: time.Now(),
		Mensagem:  "Evento processado com sucesso",
	}
	r.logIntegration("MTE", "EnviarESocial", "SUCCESS", evento, response)
	return response, nil
}

func (r *MockMTERepository) logIntegration(integration, method, status string, request, response interface{}) {
	go func() {
		reqJSON, _ := json.Marshal(request)
		respJSON, _ := json.Marshal(response)
		r.db.Exec(
			"INSERT INTO integration_logs (entity_id, integration, status, request, response, timestamp) VALUES (?, ?, ?, ?, ?, ?)",
			"global", integration+"/"+method, status, string(reqJSON), string(respJSON), time.Now().Unix(),
		)
	}()
}

// ==================== MOCK MDS ====================

type MockMDSRepository struct {
	db *sql.DB
}

func (r *MockMDSRepository) CadastrarFamilia(ctx context.Context, familia *domain.FamiliaCadUnico) (*domain.CadUnicoResponse, error) {
	response := &domain.CadUnicoResponse{
		Codigo:    fmt.Sprintf("CAD%011d", rand.Int63n(100000000000)),
		Protocolo: fmt.Sprintf("PROT%015d", rand.Int63n(1000000000000000)),
		Status:    "ATIVO",
		Score:     rand.Float64() * 100,
	}
	r.logIntegration("MDS", "CadastrarFamilia", "SUCCESS", familia, response)
	return response, nil
}

func (r *MockMDSRepository) ConsultarFamilia(ctx context.Context, codigo string) (*domain.FamiliaCadUnico, error) {
	return &domain.FamiliaCadUnico{
		Codigo:         codigo,
		NomeReferencia: "Maria Silva",
		CPFReferencia:  "123.456.789-00",
		RendaTotal:     250000,
		Integrantes: []domain.PessoaCadUnico{
			{CPF: "123.456.789-00", Nome: "Maria Silva", Parentesco: "Responsavel", Trabalha: true, Renda: 150000},
			{CPF: "987.654.321-00", Nome: "João Silva", Parentesco: "Conjuge", Trabalha: true, Renda: 100000},
		},
	}, nil
}

func (r *MockMDSRepository) AtualizarRenda(ctx context.Context, codigo string, renda int64) error {
	r.logIntegration("MDS", "AtualizarRenda", "SUCCESS", map[string]interface{}{"codigo": codigo, "renda": renda}, "Renda atualizada")
	return nil
}

func (r *MockMDSRepository) EnviarRelatorioSocial(ctx context.Context, relatorio *domain.RelatorioSocial) (*domain.RelatorioResponse, error) {
	response := &domain.RelatorioResponse{
		Protocolo: fmt.Sprintf("RELSOC%012d", rand.Int63n(1000000000000)),
		Status:    "RECEBIDO",
		DataEnvio: time.Now(),
	}
	r.logIntegration("MDS", "EnviarRelatorioSocial", "SUCCESS", relatorio, response)
	return response, nil
}

func (r *MockMDSRepository) logIntegration(integration, method, status string, request, response interface{}) {
	go func() {
		reqJSON, _ := json.Marshal(request)
		respJSON, _ := json.Marshal(response)
		r.db.Exec(
			"INSERT INTO integration_logs (entity_id, integration, status, request, response, timestamp) VALUES (?, ?, ?, ?, ?, ?)",
			"global", integration+"/"+method, status, string(reqJSON), string(respJSON), time.Now().Unix(),
		)
	}()
}

// ==================== MOCK IBGE ====================

type MockIBGERepository struct {
	db *sql.DB
}

func (r *MockIBGERepository) EnviarPesquisa(ctx context.Context, pesquisa *domain.PesquisaIBGE) (*domain.PesquisaResponse, error) {
	response := &domain.PesquisaResponse{
		Protocolo: fmt.Sprintf("PESQ%013d", rand.Int63n(10000000000000)),
		Status:    "RECEBIDO",
		DataEnvio: time.Now(),
	}
	r.logIntegration("IBGE", "EnviarPesquisa", "SUCCESS", pesquisa, response)
	return response, nil
}

func (r *MockIBGERepository) ConsultarDados(ctx context.Context, indicador string, params map[string]string) (*domain.DadosIBGE, error) {
	return &domain.DadosIBGE{
		Indicador: indicador,
		Valor:     rand.Float64() * 10000,
		Periodo:   params["periodo"],
		Unidade:   "R$",
	}, nil
}

func (r *MockIBGERepository) CadastrarProduto(ctx context.Context, produto *domain.ProdutoIBGE) (*domain.ProdutoResponse, error) {
	response := &domain.ProdutoResponse{
		CodigoIBGE: fmt.Sprintf("IBGE%09d", rand.Int63n(1000000000)),
		Status:     "CADASTRADO",
	}
	r.logIntegration("IBGE", "CadastrarProduto", "SUCCESS", produto, response)
	return response, nil
}

func (r *MockIBGERepository) logIntegration(integration, method, status string, request, response interface{}) {
	go func() {
		reqJSON, _ := json.Marshal(request)
		respJSON, _ := json.Marshal(response)
		r.db.Exec(
			"INSERT INTO integration_logs (entity_id, integration, status, request, response, timestamp) VALUES (?, ?, ?, ?, ?, ?)",
			"global", integration+"/"+method, status, string(reqJSON), string(respJSON), time.Now().Unix(),
		)
	}()
}

// ==================== MOCK SEFAZ ====================

type MockSEFAZRepository struct {
	db *sql.DB
}

func (r *MockSEFAZRepository) EmitirNFe(ctx context.Context, nfe *domain.NFeRequest) (*domain.NFeResponse, error) {
	chave := fmt.Sprintf("35%02d%04d%012d%02d%09d%01d%08d%08d",
		time.Now().Year()%100,
		time.Now().Month(),
		rand.Int63n(1000000000000),
		rand.Intn(100),
		nfe.Numero,
		rand.Intn(10),
		rand.Intn(100000000),
		rand.Intn(100000000),
	)

	now := time.Now()
	response := &domain.NFeResponse{
		Chave:           chave,
		Protocolo:       fmt.Sprintf("PROT%015d", rand.Int63n(1000000000000000)),
		Status:          "AUTORIZADA",
		DataAutorizacao: &now,
		XML:             "<NFe>...</NFe>",
	}
	r.logIntegration("SEFAZ", "EmitirNFe", "SUCCESS", nfe, response)
	return response, nil
}

func (r *MockSEFAZRepository) CancelarNFe(ctx context.Context, chave string, justificativa string) (*domain.NFeResponse, error) {
	response := &domain.NFeResponse{
		Chave:     chave,
		Status:    "CANCELADA",
		Protocolo: fmt.Sprintf("CANC%015d", rand.Int63n(1000000000000000)),
	}
	r.logIntegration("SEFAZ", "CancelarNFe", "SUCCESS", chave, response)
	return response, nil
}

func (r *MockSEFAZRepository) ConsultarNFe(ctx context.Context, chave string) (*domain.NFeResponse, error) {
	return &domain.NFeResponse{
		Chave:     chave,
		Status:    "AUTORIZADA",
		Protocolo: "PROT123456789012345",
	}, nil
}

func (r *MockSEFAZRepository) EnviarManifestoDestinatario(ctx context.Context, manifesto *domain.ManifestoRequest) (*domain.ManifestoResponse, error) {
	response := &domain.ManifestoResponse{
		Protocolo: fmt.Sprintf("MANIF%014d", rand.Int63n(100000000000000)),
		Status:    "REGISTRADO",
	}
	r.logIntegration("SEFAZ", "EnviarManifestoDestinatario", "SUCCESS", manifesto, response)
	return response, nil
}

func (r *MockSEFAZRepository) EmitirNFS(ctx context.Context, nfs *domain.NFSRequest) (*domain.NFSResponse, error) {
	response := &domain.NFSResponse{
		Numero:            nfs.Numero,
		CodigoVerificacao: fmt.Sprintf("%08d", rand.Int63n(100000000)),
		Status:            "AUTORIZADA",
		Link:              fmt.Sprintf("https://nfse.prefeitura.sp.gov.br/%d", nfs.Numero),
	}
	r.logIntegration("SEFAZ", "EmitirNFS", "SUCCESS", nfs, response)
	return response, nil
}

func (r *MockSEFAZRepository) logIntegration(integration, method, status string, request, response interface{}) {
	go func() {
		reqJSON, _ := json.Marshal(request)
		respJSON, _ := json.Marshal(response)
		r.db.Exec(
			"INSERT INTO integration_logs (entity_id, integration, status, request, response, timestamp) VALUES (?, ?, ?, ?, ?, ?)",
			"global", integration+"/"+method, status, string(reqJSON), string(respJSON), time.Now().Unix(),
		)
	}()
}

// ==================== MOCK BNDES ====================

type MockBNDESRepository struct {
	db *sql.DB
}

func (r *MockBNDESRepository) ConsultarLinhasCredito(ctx context.Context, cnpj string) ([]domain.LinhaCredito, error) {
	return []domain.LinhaCredito{
		{
			Codigo:      "BNDES_AUTOMATIC",
			Nome:        "BNDES Automático",
			Descricao:   "Crédito automático para ME/EPP",
			ValorMaximo: 100000000,
			PrazoMaximo: 60,
			TaxaJuros:   0.85,
		},
		{
			Codigo:      "PRONAMP",
			Nome:        "PRONAMP Investimento",
			Descricao:   "Programa Nacional de Microcrédito",
			ValorMaximo: 21000000,
			PrazoMaximo: 84,
			TaxaJuros:   0.65,
		},
	}, nil
}

func (r *MockBNDESRepository) SimularCredito(ctx context.Context, simulacao *domain.SimulacaoCredito) (*domain.ResultadoSimulacao, error) {
	parcela := simulacao.Valor / int64(simulacao.Prazo)
	juros := float64(simulacao.Valor) * 0.0085 * float64(simulacao.Prazo)

	response := &domain.ResultadoSimulacao{
		ValorParcela:    parcela + int64(juros/float64(simulacao.Prazo)),
		TotalJuros:      int64(juros),
		CET:             8.5,
		Prazo:           simulacao.Prazo,
		PrimeiraParcela: time.Now().AddDate(0, 1, 0),
	}
	r.logIntegration("BNDES", "SimularCredito", "SUCCESS", simulacao, response)
	return response, nil
}

func (r *MockBNDESRepository) SolicitarCredito(ctx context.Context, solicitacao *domain.SolicitacaoCredito) (*domain.SolicitacaoResponse, error) {
	response := &domain.SolicitacaoResponse{
		Numero:       fmt.Sprintf("BNDES%012d", rand.Int63n(1000000000000)),
		Protocolo:    fmt.Sprintf("PROT%015d", rand.Int63n(1000000000000000)),
		Status:       "EM_ANALISE",
		DataPrevisao: ptrTime(time.Now().AddDate(0, 1, 0)),
	}
	r.logIntegration("BNDES", "SolicitarCredito", "SUCCESS", solicitacao, response)
	return response, nil
}

func (r *MockBNDESRepository) ConsultarSolicitacao(ctx context.Context, numero string) (*domain.SolicitacaoCredito, error) {
	return &domain.SolicitacaoCredito{
		Numero:          numero,
		Linha:           "BNDES_AUTOMATIC",
		Valor:           5000000,
		Prazo:           36,
		Carencia:        6,
		Status:          "EM_ANALISE",
		DataSolicitacao: time.Now().AddDate(0, -1, 0),
	}, nil
}

func (r *MockBNDESRepository) logIntegration(integration, method, status string, request, response interface{}) {
	go func() {
		reqJSON, _ := json.Marshal(request)
		respJSON, _ := json.Marshal(response)
		r.db.Exec(
			"INSERT INTO integration_logs (entity_id, integration, status, request, response, timestamp) VALUES (?, ?, ?, ?, ?, ?)",
			"global", integration+"/"+method, status, string(reqJSON), string(respJSON), time.Now().Unix(),
		)
	}()
}

func ptrTime(t time.Time) *time.Time {
	return &t
}

// ==================== MOCK SEBRAE ====================

type MockSEBRAERepository struct {
	db *sql.DB
}

func (r *MockSEBRAERepository) ConsultarCursos(ctx context.Context, cnpj string) ([]domain.CursoSEBRAE, error) {
	return []domain.CursoSEBRAE{
		{
			Codigo:       "SEBRAE001",
			Nome:         "Gestão Financeira para Cooperativas",
			Descricao:    "Aprenda a gerenciar as finanças da sua cooperativa",
			CargaHoraria: 20,
			Modalidade:   "ONLINE",
			Vagas:        50,
		},
		{
			Codigo:       "SEBRAE002",
			Nome:         "Marketing Digital",
			Descricao:    "Estratégias de marketing para pequenos negócios",
			CargaHoraria: 16,
			Modalidade:   "EAD",
			Vagas:        100,
		},
	}, nil
}

func (r *MockSEBRAERepository) InscreverCurso(ctx context.Context, inscricao *domain.InscricaoCurso) (*domain.InscricaoResponse, error) {
	response := &domain.InscricaoResponse{
		Codigo:    fmt.Sprintf("INSC%010d", rand.Int63n(10000000000)),
		Protocolo: fmt.Sprintf("PROT%015d", rand.Int63n(1000000000000000)),
		Status:    "INSCRITO",
	}
	r.logIntegration("SEBRAE", "InscreverCurso", "SUCCESS", inscricao, response)
	return response, nil
}

func (r *MockSEBRAERepository) ConsultarInscricao(ctx context.Context, codigo string) (*domain.InscricaoCurso, error) {
	return &domain.InscricaoCurso{
		Codigo:        codigo,
		CursoCodigo:   "SEBRAE001",
		Status:        "CONFIRMADO",
		DataInscricao: time.Now().AddDate(0, -1, 0),
	}, nil
}

func (r *MockSEBRAERepository) ConsultarConsultoria(ctx context.Context, cnpj string) ([]domain.ProgramaConsultoria, error) {
	return []domain.ProgramaConsultoria{
		{
			Codigo:    "CONS001",
			Nome:      "Consultoria em Gestão",
			Descricao: "Apoio especializado em gestão empresarial",
			Area:      "GESTAO",
			Duracao:   "3 meses",
		},
	}, nil
}

func (r *MockSEBRAERepository) logIntegration(integration, method, status string, request, response interface{}) {
	go func() {
		reqJSON, _ := json.Marshal(request)
		respJSON, _ := json.Marshal(response)
		r.db.Exec(
			"INSERT INTO integration_logs (entity_id, integration, status, request, response, timestamp) VALUES (?, ?, ?, ?, ?, ?)",
			"global", integration+"/"+method, status, string(reqJSON), string(respJSON), time.Now().Unix(),
		)
	}()
}

// ==================== MOCK PROVIDENTIA ====================

type MockProvidentiaRepository struct {
	db *sql.DB
}

func (r *MockProvidentiaRepository) SyncPackage(ctx context.Context, pkg *domain.SyncPackage) (*domain.SyncResponse, error) {
	response := &domain.SyncResponse{
		Protocolo: fmt.Sprintf("SYNC%014d", rand.Int63n(100000000000000)),
		Status:    "RECEBIDO",
		Timestamp: time.Now().Unix(),
	}
	r.logIntegration("PROVIDENTIA", "SyncPackage", "SUCCESS", pkg, response)
	return response, nil
}

func (r *MockProvidentiaRepository) GetMarketplaceData(ctx context.Context, entityID string) (*domain.MarketplaceData, error) {
	return &domain.MarketplaceData{
		Offers: []domain.MarketplaceOffer{
			{ID: "1", EntityID: entityID, ProductName: "Mel Orgânico", Quantity: 100, Price: 2500, Unit: "kg", Active: true},
			{ID: "2", EntityID: entityID, ProductName: "Mel com Própolis", Quantity: 50, Price: 3500, Unit: "kg", Active: true},
		},
		TotalActive: 2,
	}, nil
}

func (r *MockProvidentiaRepository) RegisterOffer(ctx context.Context, offer *domain.MarketplaceOffer) (*domain.OfferResponse, error) {
	response := &domain.OfferResponse{
		ID:     fmt.Sprintf("OFFER%010d", rand.Int63n(10000000000)),
		Status: "ATIVO",
		Link:   fmt.Sprintf("https://marketplace.providentia.org/offer/%d", rand.Int63()),
	}
	r.logIntegration("PROVIDENTIA", "RegisterOffer", "SUCCESS", offer, response)
	return response, nil
}

func (r *MockProvidentiaRepository) logIntegration(integration, method, status string, request, response interface{}) {
	go func() {
		reqJSON, _ := json.Marshal(request)
		respJSON, _ := json.Marshal(response)
		r.db.Exec(
			"INSERT INTO integration_logs (entity_id, integration, status, request, response, timestamp) VALUES (?, ?, ?, ?, ?, ?)",
			"global", integration+"/"+method, status, string(reqJSON), string(respJSON), time.Now().Unix(),
		)
	}()
}
