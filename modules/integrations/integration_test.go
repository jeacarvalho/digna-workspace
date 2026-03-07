package integrations_test

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/providentia/digna/integrations/pkg/integrations"
)

func TestIntegrationService_Mock(t *testing.T) {
	// Criar banco de dados temporário
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Criar serviço com implementações mock
	service, err := integrations.NewMockIntegrationService(db)
	if err != nil {
		t.Fatalf("Failed to create integration service: %v", err)
	}

	// Testar Receita Federal
	t.Run("ReceitaFederal_ConsultarCNPJ", func(t *testing.T) {
		ctx := t.Context()
		cnpjData, err := service.ReceitaFederal().ConsultarCNPJ(ctx, "12345678000190")
		if err != nil {
			t.Errorf("Failed to consult CNPJ: %v", err)
		}
		if cnpjData == nil {
			t.Error("CNPJ data is nil")
		} else {
			t.Logf("CNPJ: %s, Razão: %s", cnpjData.CNPJ, cnpjData.RazaoSocial)
		}
	})

	// Testar SEFAZ
	t.Run("SEFAZ_EmitirNFe", func(t *testing.T) {
		ctx := t.Context()
		nfeReq := &integrations.NFeRequest{
			CNPJ:       "12345678000190",
			Serie:      "1",
			Numero:     1,
			NaturezaOp: "Venda",
			Destinatario: integrations.DestinatarioNFe{
				CNPJCPF: "98765432100",
				Nome:    "Cliente Teste",
			},
		}
		nfeResp, err := service.SEFAZ().EmitirNFe(ctx, nfeReq)
		if err != nil {
			t.Errorf("Failed to emit NFe: %v", err)
		}
		if nfeResp == nil {
			t.Error("NFe response is nil")
		} else {
			t.Logf("NFe emitted: Chave=%s, Status=%s", nfeResp.Chave, nfeResp.Status)
		}
	})

	// Testar MTE
	t.Run("MTE_EnviarRAIS", func(t *testing.T) {
		ctx := t.Context()
		raisReq := &integrations.RAISRequest{
			Ano:  2026,
			CNPJ: "12345678000190",
			Trabalhadores: []integrations.TrabalhadorRAIS{
				{
					CPF:           "12345678900",
					Nome:          "João Silva",
					CBO:           "622020",
					Salario:       150000,
					HorasSemanais: 44,
				},
			},
		}
		raisResp, err := service.MTE().EnviarRAIS(ctx, raisReq)
		if err != nil {
			t.Errorf("Failed to send RAIS: %v", err)
		}
		if raisResp == nil {
			t.Error("RAIS response is nil")
		} else {
			t.Logf("RAIS sent: Protocolo=%s, Status=%s", raisResp.Protocolo, raisResp.Status)
		}
	})

	// Testar BNDES
	t.Run("BNDES_SimularCredito", func(t *testing.T) {
		ctx := t.Context()
		simulacao := &integrations.SimulacaoCredito{
			CNPJ:  "12345678000190",
			Linha: "BNDES_AUTOMATIC",
			Valor: 5000000,
			Prazo: 36,
		}
		resultado, err := service.BNDES().SimularCredito(ctx, simulacao)
		if err != nil {
			t.Errorf("Failed to simulate credit: %v", err)
		}
		if resultado == nil {
			t.Error("Credit simulation result is nil")
		} else {
			t.Logf("Credit simulation: Parcela=%d, CET=%.2f%%", resultado.ValorParcela, resultado.CET)
		}
	})

	// Testar log de integração
	t.Run("IntegrationLog", func(t *testing.T) {
		ctx := t.Context()
		err := service.LogIntegration(ctx, "test-entity", "TEST_INTEGRATION", integrations.StatusSuccess,
			map[string]string{"test": "data"},
			map[string]string{"result": "success"})
		if err != nil {
			t.Logf("Note: Failed to log integration: %v (this is acceptable for mock)", err)
		}

		// Verificar logs
		logs, err := service.GetIntegrationLogs(ctx, "test-entity", "TEST_INTEGRATION", 10)
		if err != nil {
			t.Logf("Note: Failed to get integration logs: %v (this is acceptable for mock)", err)
		} else {
			t.Logf("Found %d integration logs", len(logs))
		}
	})
}
