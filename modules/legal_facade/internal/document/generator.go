package document

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"html/template"
	"strings"
	"time"

	"github.com/providentia/digna/lifecycle/pkg/lifecycle"
)

const assemblyTemplate = `# ATA DE ASSEMBLEIA GERAL EXTRAORDINÁRIA

**Entidade:** {{.EntityName}}  
**Data:** {{.AssemblyDate}}  
**Status:** {{.EntityStatus}}

---

## DECISÕES REGISTRADAS

{{range .Decisions}}
### {{.Index}}. {{.Title}}

- **Status:** {{.Status}}
- **Hash de Auditoria:** ` + "`{{.Hash}}`" + `
- **Data do Registro:** {{.CreatedAt}}

---
{{end}}

## RATIFICAÇÃO

As decisões acima foram registradas no sistema CADSOL (Cadastro de Decisões Soberanas) e são consideradas válidas para fins de auditoria e transparência institucional.

**Assinaturas digitais requeridas:**

- [ ] Presidente da Assembleia
- [ ] Secretário(a)  
- [ ] Tesoureiro(a)

---

*Documento gerado automaticamente pelo Digna - Providentia Foundation*  
*{{.GeneratedAt}}*
`

type AssemblyData struct {
	EntityName   string
	EntityStatus string
	AssemblyDate string
	Decisions    []DecisionData
	GeneratedAt  string
}

type DecisionData struct {
	Index     int
	Title     string
	Status    string
	Hash      string
	CreatedAt string
}

type Generator struct {
	legalRepo LegalRepository
}

func NewGenerator(lm lifecycle.LifecycleManager) *Generator {
	return &Generator{
		legalRepo: NewSQLiteLegalRepository(lm),
	}
}

func (g *Generator) GenerateAssemblyMinutes(entityID string, entityName string, status string) (string, error) {
	decisionsInfo, err := g.legalRepo.GetAllDecisions(entityID)
	if err != nil {
		return "", fmt.Errorf("failed to query decisions: %w", err)
	}

	var decisions []DecisionData
	index := 1

	for _, d := range decisionsInfo {
		decisions = append(decisions, DecisionData{
			Index:     index,
			Title:     d.Title,
			Status:    d.Status,
			Hash:      d.Hash,
			CreatedAt: time.Unix(d.CreatedAt, 0).Format("2006-01-02 15:04"),
		})
		index++
	}

	if len(decisions) == 0 {
		return "", fmt.Errorf("no decisions found for entity %s", entityID)
	}

	data := AssemblyData{
		EntityName:   entityName,
		EntityStatus: status,
		AssemblyDate: time.Now().Format("2006-01-02"),
		Decisions:    decisions,
		GeneratedAt:  time.Now().Format("2006-01-02 15:04:05"),
	}

	tmpl, err := template.New("assembly").Parse(assemblyTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

const dossierTemplate = `# DOSSIÊ DE FORMALIZAÇÃO CADSOL

**Entidade:** {{.EntityName}}  
**Status:** {{.EntityStatus}}  
**Data de Geração:** {{.GeneratedAt}}  
**CNPJ:** {{.CNPJ}}  
**NIRE:** {{.NIRE}}  
**Endereço:** {{.Address}}

---

## HISTÓRICO DE DECISÕES SOBERANAS

{{if .Decisions}}
A entidade possui **{{.DecisionCount}} decisões** registradas no sistema CADSOL (Cadastro de Decisões Soberanas), demonstrando processo democrático e autogestão.

{{range .Decisions}}
### {{.Index}}. {{.Title}}

- **Status:** {{.Status}}
- **Hash de Auditoria:** ` + "`{{.Hash}}`" + `
- **Data do Registro:** {{.CreatedAt}}
- **Conteúdo Resumido:** {{.ContentSummary}}

{{end}}
{{else}}
**ATENÇÃO:** A entidade não possui decisões registradas no sistema CADSOL.
{{end}}

---

## CRITÉRIOS DE FORMALIZAÇÃO

### Requisitos Mínimos para Formalização CADSOL:
1. ✅ **Mínimo de 3 decisões registradas** (ITG 2002, Art. 14)
2. ✅ **Estatuto social aprovado em assembleia**
3. ✅ **CNPJ e NIRE registrados** (opcional para grupos informais)
4. ✅ **Processo democrático documentado**

### Status Atual:
- **Decisões Registradas:** {{.DecisionCount}}/3 {{if ge .DecisionCount 3}}✅ ATINGIDO{{else}}❌ PENDENTE{{end}}
- **Status da Entidade:** {{.EntityStatus}}
- **Elegível para Formalização:** {{if ge .DecisionCount 3}}✅ SIM{{else}}❌ NÃO (necessário {{sub 3 .DecisionCount}} decisão(ões) adicional(is)){{end}}

---

## DOCUMENTOS ANEXOS

1. **Estatuto Social** - {{if eq .EntityStatus "FORMALIZED"}}✅ GERADO{{else}}⏳ AGUARDANDO FORMALIZAÇÃO{{end}}
2. **Atas de Assembleia** - {{if gt .DecisionCount 0}}✅ {{.DecisionCount}} ATA(S) DISPONÍVEL(IS){{else}}❌ NENHUMA ATA REGISTRADA{{end}}
3. **Registro de Trabalho** - ⏳ EM IMPLEMENTAÇÃO (ITG 2002)
4. **Balanço Patrimonial** - ⏳ EM IMPLEMENTAÇÃO

---

## HASH DE INTEGRIDADE DO DOSSIÊ

Este documento possui um hash criptográfico SHA256 que garante sua integridade e imutabilidade para fins de auditoria pública perante o Ministério do Trabalho e Emprego (SINAES/CADSOL).

**Hash do Dossiê:** ` + "`{{.DocumentHash}}`" + `

**Como validar:**
1. Copie todo o conteúdo deste documento (exceto esta seção)
2. Calcule o hash SHA256 do conteúdo
3. Compare com o hash acima
4. Se coincidirem, o documento é íntegro e não foi alterado

---

## DISPOSIÇÕES FINAIS

1. Este dossiê é gerado automaticamente pelo **Sistema Digna** da Providentia Foundation.
2. As decisões registradas no CADSOL são auditáveis publicamente via hash criptográfico.
3. A formalização é um processo **gradual e pedagógico** - não force a burocracia.
4. Grupos com menos de 3 decisões devem focar em **construir autogestão** antes de formalização.

*Documento gerado automaticamente pelo Sistema Digna - Providentia Foundation*  
*{{.GeneratedAt}}*

**Assinaturas requeridas para formalização:**

- [ ] Presidente da Assembleia
- [ ] Secretário(a)
- [ ] Tesoureiro(a)
- [ ] 2/3 dos membros presentes
`

type DossierData struct {
	EntityName    string
	EntityStatus  string
	GeneratedAt   string
	CNPJ          string
	NIRE          string
	Address       string
	Decisions     []DossierDecisionData
	DecisionCount int
	DocumentHash  string
}

type DossierDecisionData struct {
	Index          int
	Title          string
	Status         string
	Hash           string
	CreatedAt      string
	ContentSummary string
}

func (g *Generator) GenerateDossier(entityID string, entityName string, status string) (string, string, error) {
	// Verificar se tem decisões suficientes usando FormalizationSimulator
	fs := NewFormalizationSimulator(g.legalRepo.(*SQLiteLegalRepository).lifecycleManager)
	_, err := fs.CheckFormalizationCriteria(entityID)
	if err != nil {
		return "", "", fmt.Errorf("failed to check formalization criteria: %w", err)
	}

	// Buscar todas as decisões
	decisionsInfo, err := g.legalRepo.GetAllDecisions(entityID)
	if err != nil {
		return "", "", fmt.Errorf("failed to query decisions: %w", err)
	}

	// Preparar dados das decisões para o dossiê
	var dossierDecisions []DossierDecisionData
	index := 1
	for _, d := range decisionsInfo {
		// Resumir conteúdo (primeiros 100 caracteres)
		contentSummary := "Decisão registrada no sistema CADSOL"
		if len(d.Title) > 50 {
			contentSummary = d.Title[:50] + "..."
		} else if d.Title != "" {
			contentSummary = d.Title
		}

		dossierDecisions = append(dossierDecisions, DossierDecisionData{
			Index:          index,
			Title:          d.Title,
			Status:         d.Status,
			Hash:           d.Hash,
			CreatedAt:      time.Unix(d.CreatedAt, 0).Format("2006-01-02"),
			ContentSummary: contentSummary,
		})
		index++
	}

	// Dados fictícios para exemplo (em produção viriam do banco)
	cnpj := "00.000.000/0001-91"
	nire := "12.345.678.901"
	address := "Rua da Autogestão, 123 - Centro"

	if status == "DREAM" || status == "" {
		cnpj = "A DEFINIR APÓS FORMALIZAÇÃO"
		nire = "A DEFINIR APÓS FORMALIZAÇÃO"
		address = "ENDEREÇO A DEFINIR"
	}

	// Dados para o dossiê (hash será calculado depois)

	// Primeiro, gerar conteúdo SEM a seção de hash para calcular o hash
	// Dividir o template em partes
	templateParts := strings.Split(dossierTemplate, "## HASH DE INTEGRIDADE DO DOSSIÊ")
	if len(templateParts) < 2 {
		return "", "", fmt.Errorf("invalid dossier template format")
	}

	contentBeforeHash := templateParts[0]

	// Gerar conteúdo antes do hash
	tmplBeforeHash, err := template.New("dossier_before_hash").Funcs(template.FuncMap{
		"sub": func(a, b int) int { return a - b },
		"ge":  func(a, b int) bool { return a >= b },
		"gt":  func(a, b int) bool { return a > b },
		"eq":  func(a, b string) bool { return a == b },
	}).Parse(contentBeforeHash)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse dossier template (before hash): %w", err)
	}

	var bufBeforeHash bytes.Buffer
	dataForHash := DossierData{
		EntityName:    entityName,
		EntityStatus:  status,
		GeneratedAt:   time.Now().Format("2006-01-02 15:04:05"),
		CNPJ:          cnpj,
		NIRE:          nire,
		Address:       address,
		Decisions:     dossierDecisions,
		DecisionCount: len(decisionsInfo),
		DocumentHash:  "", // Vazio para esta parte
	}

	if err := tmplBeforeHash.Execute(&bufBeforeHash, dataForHash); err != nil {
		return "", "", fmt.Errorf("failed to execute dossier template (before hash): %w", err)
	}

	contentBeforeHashSection := bufBeforeHash.String()

	// Calcular hash do conteúdo antes da seção de hash
	documentHash := generateDossierHash(contentBeforeHashSection, entityID)

	// Agora gerar o conteúdo completo COM o hash real
	dataWithHash := DossierData{
		EntityName:    entityName,
		EntityStatus:  status,
		GeneratedAt:   time.Now().Format("2006-01-02 15:04:05"),
		CNPJ:          cnpj,
		NIRE:          nire,
		Address:       address,
		Decisions:     dossierDecisions,
		DecisionCount: len(decisionsInfo),
		DocumentHash:  documentHash,
	}

	// Gerar conteúdo completo
	tmpl, err := template.New("dossier").Funcs(template.FuncMap{
		"sub": func(a, b int) int { return a - b },
		"ge":  func(a, b int) bool { return a >= b },
		"gt":  func(a, b int) bool { return a > b },
		"eq":  func(a, b string) bool { return a == b },
	}).Parse(dossierTemplate)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse dossier template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, dataWithHash); err != nil {
		return "", "", fmt.Errorf("failed to execute dossier template: %w", err)
	}

	finalContent := buf.String()

	return finalContent, documentHash, nil
}

func generateDossierHash(content string, entityID string) string {
	// Seguir mesmo padrão do core_lume: content:entityID:salt
	salted := fmt.Sprintf("%s:%s:DIGNA_DOSSIER_SALT_v1", content, entityID)
	hash := sha256.Sum256([]byte(salted))
	return hex.EncodeToString(hash[:])
}
