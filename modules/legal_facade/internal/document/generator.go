package document

import (
	"bytes"
	"fmt"
	"html/template"
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
