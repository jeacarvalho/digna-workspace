package document

import (
	"bytes"
	"fmt"
	"html/template"
	"time"

	"github.com/providentia/digna/lifecycle/pkg/lifecycle"
)

const identityCardTemplate = `═══════════════════════════════════════════════════════
          CARTÃO DE IDENTIFICAÇÃO - PROVIDENTIA FOUNDATION
═══════════════════════════════════════════════════════

ENTIDADE: {{.EntityName}}
ID: {{.EntityID}}
STATUS: {{.EntityStatus}}

{{if eq .EntityStatus "FORMALIZED"}}
CNPJ: {{.CNPJ}}
NIRE: {{.NIRE}}
{{else}}
[ ENTIDADE EM FASE DE CONSTITUIÇÃO ]
{{end}}

DATA DE EMISSÃO: {{.IssuedAt}}
VALIDADE: {{.ValidUntil}}

───────────────────────────────────────────────────────
Este documento é válido apenas junto à Fundação Providentia
para fins de identificação institucional.
═══════════════════════════════════════════════════════
`

type IdentityCardData struct {
	EntityName   string
	EntityID     string
	EntityStatus string
	CNPJ         string
	NIRE         string
	IssuedAt     string
	ValidUntil   string
}

type IdentityGenerator struct {
	lifecycleManager lifecycle.LifecycleManager
}

func NewIdentityGenerator(lm lifecycle.LifecycleManager) *IdentityGenerator {
	return &IdentityGenerator{
		lifecycleManager: lm,
	}
}

func (ig *IdentityGenerator) GenerateIdentityCard(entityID string, entityName string, status string) (string, error) {
	data := IdentityCardData{
		EntityName:   entityName,
		EntityID:     entityID,
		EntityStatus: status,
		CNPJ:         "00.000.000/0000-00",
		NIRE:         "00.000.000.000",
		IssuedAt:     time.Now().Format("2006-01-02"),
		ValidUntil:   time.Now().AddDate(1, 0, 0).Format("2006-01-02"),
	}

	tmpl, err := template.New("identity").Parse(identityCardTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}
