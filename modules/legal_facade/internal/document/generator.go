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
	lifecycleManager lifecycle.LifecycleManager
}

func NewGenerator(lm lifecycle.LifecycleManager) *Generator {
	return &Generator{
		lifecycleManager: lm,
	}
}

func (g *Generator) GenerateAssemblyMinutes(entityID string, entityName string, status string) (string, error) {
	db, err := g.lifecycleManager.GetConnection(entityID)
	if err != nil {
		return "", fmt.Errorf("failed to get connection: %w", err)
	}

	rows, err := db.Query(
		"SELECT id, title, content_hash, status, created_at FROM decisions_log ORDER BY created_at DESC",
	)
	if err != nil {
		return "", fmt.Errorf("failed to query decisions: %w", err)
	}
	defer rows.Close()

	var decisions []DecisionData
	index := 1

	for rows.Next() {
		var id int64
		var title, hash, decisionStatus string
		var createdAt int64

		if err := rows.Scan(&id, &title, &hash, &decisionStatus, &createdAt); err != nil {
			return "", fmt.Errorf("failed to scan decision: %w", err)
		}

		decisions = append(decisions, DecisionData{
			Index:     index,
			Title:     title,
			Status:    decisionStatus,
			Hash:      hash,
			CreatedAt: time.Unix(createdAt, 0).Format("2006-01-02 15:04"),
		})
		index++
	}

	if err := rows.Err(); err != nil {
		return "", fmt.Errorf("error iterating decisions: %w", err)
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
