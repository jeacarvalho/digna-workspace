package legal_facade_test

import (
	"os"
	"strings"
	"testing"

	"github.com/providentia/digna/legal_facade/pkg/document"
	"github.com/providentia/digna/lifecycle/pkg/lifecycle"
	"github.com/providentia/digna/pdv_ui/usecase"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func TestGenerateDossier(t *testing.T) {
	dataDir := "../../data/entities"
	defer os.RemoveAll(dataDir)

	lifecycleMgr := lifecycle.NewSQLiteManager()
	defer lifecycleMgr.CloseAll()

	entityID := "Entidade_Dossie_Teste"
	entityName := "Cooperativa Dossiê Teste"

	generator := document.NewGenerator(lifecycleMgr)
	sgHandler := usecase.NewSocialGovernanceHandler(lifecycleMgr)

	t.Run("Gerar_Dossiê_com_Decisões", func(t *testing.T) {
		decisions := []struct {
			title   string
			content string
		}{
			{
				title:   "Aprovação do Estatuto Social",
				content: "Assembleia aprova o estatuto da cooperativa por unanimidade",
			},
			{
				title:   "Eleição do Conselho Administrativo",
				content: "Eleitos membros do conselho para gestão 2026-2028",
			},
			{
				title:   "Aprovação do Plano de Negócios",
				content: "Aprovado plano de negócios para o exercício fiscal",
			},
		}

		for _, d := range decisions {
			req := usecase.DecisionRequest{
				EntityID: entityID,
				Title:    d.title,
				Content:  d.content,
			}
			_, err := sgHandler.RecordDecision(req)
			if err != nil {
				t.Fatalf("failed to record decision '%s': %v", d.title, err)
			}
		}

		dossierContent, dossierHash, err := generator.GenerateDossier(entityID, entityName, "DREAM")
		if err != nil {
			t.Fatalf("failed to generate dossier: %v", err)
		}

		if dossierContent == "" {
			t.Fatal("generated dossier content is empty")
		}

		if dossierHash == "" {
			t.Fatal("dossier hash is empty")
		}

		requiredSections := []string{
			"DOSSIÊ DE FORMALIZAÇÃO CADSOL",
			"Cooperativa Dossiê Teste",
			"HISTÓRICO DE DECISÕES SOBERANAS",
			"CRITÉRIOS DE FORMALIZAÇÃO",
			"DOCUMENTOS ANEXOS",
			"HASH DE INTEGRIDADE DO DOSSIÊ",
			"DISPOSIÇÕES FINAIS",
			"Aprovação do Estatuto Social",
			"Eleição do Conselho Administrativo",
			"Aprovação do Plano de Negócios",
			"Hash do Dossiê:",
			"Sistema Digna",
			"Providentia Foundation",
		}

		for _, section := range requiredSections {
			if !strings.Contains(dossierContent, section) {
				t.Errorf("dossier missing section: '%s'", section)
			}
		}

		if !strings.Contains(dossierContent, dossierHash) {
			t.Errorf("dossier hash not found in content. Hash: %s", dossierHash)
		}

		// Verificar se mostra 3 decisões
		if !strings.Contains(dossierContent, "3 decisões") {
			t.Error("dossier should mention 3 decisions")
		}

		// Verificar se mostra que atingiu o critério (pode ser "✅ ATINGIDO" ou similar)
		if !strings.Contains(dossierContent, "3/3") {
			t.Error("dossier should show 3/3 decisions")
		}

		t.Log("✅ Dossiê gerado com sucesso!")
		t.Logf("Hash do dossiê: %s", dossierHash)
		t.Logf("Tamanho do documento: %d caracteres", len(dossierContent))

		// Verificar se o hash está no conteúdo
		if strings.Contains(dossierContent, dossierHash) {
			t.Log("✅ Hash encontrado no conteúdo do documento")
		} else {
			t.Log("⚠️  Hash NÃO encontrado no conteúdo do documento")
			// Procurar por "Hash do Dossiê:" no conteúdo
			hashIndex := strings.Index(dossierContent, "Hash do Dossiê:")
			if hashIndex > 0 {
				t.Logf("Seção de hash encontrada em: %d", hashIndex)
				t.Logf("Conteúdo ao redor:\n%s", dossierContent[max(0, hashIndex-50):min(len(dossierContent), hashIndex+100)])
			}
		}

		t.Logf("Document preview (primeiros 1000 caracteres):\n---\n%s\n---", dossierContent[:min(1000, len(dossierContent))])
	})

	t.Run("Gerar_Dossiê_sem_Decisões", func(t *testing.T) {
		entityID2 := "Entidade_Sem_Decisoes"
		dossierContent, _, err := generator.GenerateDossier(entityID2, "Entidade Sem Decisões", "DREAM")
		if err != nil {
			t.Fatalf("failed to generate dossier without decisions: %v", err)
		}

		// Verificar se avisa sobre falta de decisões
		if !strings.Contains(dossierContent, "não possui decisões registradas") {
			t.Error("dossier should warn about no decisions")
		}

		// Verificar se mostra 0 decisões
		if !strings.Contains(dossierContent, "0 decisões") && !strings.Contains(dossierContent, "0/3") {
			t.Error("dossier should show 0 decisions")
		}

		t.Log("✅ Dossiê sem decisões gerado corretamente!")
	})

	t.Run("Gerar_Dossiê_Entidade_Formalizada", func(t *testing.T) {
		dossierContent, dossierHash, err := generator.GenerateDossier(entityID, entityName, "FORMALIZED")
		if err != nil {
			t.Fatalf("failed to generate dossier for formalized entity: %v", err)
		}

		if !strings.Contains(dossierContent, "FORMALIZED") {
			t.Error("dossier should show FORMALIZED status")
		}

		t.Log("✅ Dossiê para entidade formalizada gerado com sucesso!")
		t.Logf("Hash: %s", dossierHash)
	})
}

func TestDossierHashIntegrity(t *testing.T) {
	dataDir := "../../data/entities"
	defer os.RemoveAll(dataDir)

	lifecycleMgr := lifecycle.NewSQLiteManager()
	defer lifecycleMgr.CloseAll()

	entityID := "Entidade_Hash_Teste"
	generator := document.NewGenerator(lifecycleMgr)
	sgHandler := usecase.NewSocialGovernanceHandler(lifecycleMgr)

	req := usecase.DecisionRequest{
		EntityID: entityID,
		Title:    "Decisão Teste Hash",
		Content:  "Conteúdo da decisão para teste de hash",
	}
	_, err := sgHandler.RecordDecision(req)
	if err != nil {
		t.Fatalf("failed to record decision: %v", err)
	}

	dossierContent, dossierHash, err := generator.GenerateDossier(entityID, "Entidade Hash Teste", "DREAM")
	if err != nil {
		t.Fatalf("failed to generate dossier: %v", err)
	}

	// Verificar se o hash está no conteúdo
	if !strings.Contains(dossierContent, dossierHash) {
		t.Errorf("dossier hash not found in content. Hash: %s", dossierHash)
	}

	// Verificar formato do hash (SHA256 tem 64 caracteres hex)
	if len(dossierHash) != 64 {
		t.Errorf("invalid hash length: %d (expected 64)", len(dossierHash))
	}

	// Verificar se é hexadecimal válido
	for _, c := range dossierHash {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			t.Errorf("invalid character in hash: %c", c)
			break
		}
	}

	t.Log("✅ Hash integrity validated successfully!")
	t.Logf("Hash: %s", dossierHash)
}
