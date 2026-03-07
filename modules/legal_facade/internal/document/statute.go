package document

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"html/template"
	"time"

	"github.com/providentia/digna/lifecycle/pkg/lifecycle"
)

const statuteTemplate = `# ESTATUTO SOCIAL

## {{.EntityName}}

**CNPJ:** {{.CNPJ}}  
**NIRE:** {{.NIRE}}  
**Data de Registro:** {{.RegistrationDate}}  
**Sede:** {{.Address}}

---

## CAPÍTULO I - DA DENOMINAÇÃO, SEDE E FINALIDADE

**Art. 1º** {{.EntityName}}, adiante designada simplesmente "Cooperativa", é uma sociedade cooperativa, constituída sob a forma de pessoa jurídica de direito privado, sem fins lucrativos, regida por este Estatuto e pelas disposições legais em vigor.

**Art. 2º** A Cooperativa tem sede e foro na {{.Address}}, podendo abrir filiais ou delegacias em quaisquer localidades.

**Art. 3º** A Cooperativa tem por finalidade:
- I - Promover a organização econômica e social de seus cooperados;
- II - Buscar o desenvolvimento econômico e social sustentável;
- III - Proporcionar beneficios mútuos aos cooperados;
- IV - Contribuir para a edificação de uma economia solidária.

---

## CAPÍTULO II - DOS COOPERADOS

**Art. 4º** Poderão ser cooperados as pessoas físicas ou jurídicas que:
- a) Tenham interesse na realização dos objetivos sociais;
- b) Aceitem as disposições deste Estatuto;
- c) Sejam admitidas pelo Conselho de Administração.

**Art. 5º** São direitos dos cooperados:
- I - Participar das Assembleias Gerais;
- II - Votar e ser votado para os cargos eletivos;
- III - Requisitar a convocação de Assembleia Geral Extraordinária;
- IV - Ter acesso às informações sobre a gestão da Cooperativa.

**Art. 6º** São deveres dos cooperados:
- I - Contribuir para o desenvolvimento da Cooperativa;
- II - Cumprir este Estatuto e as decisões das Assembleias;
- III - Participar das atividades cooperativas;
- IV - Honrar suas responsabilidades financeiras.

---

## CAPÍTULO III - DA ADMINISTRAÇÃO

**Art. 7º** A Cooperativa é administrada por:
- I - Assembleia Geral (órgão soberano);
- II - Conselho de Administração (executivo);
- III - Conselho Fiscal (fiscalizador).

**Art. 8º** A Assembleia Geral é o órgão soberano da Cooperativa, composta por todos os cooperados em pleno gozo de seus direitos.

**Art. 9º** Ao Conselho de Administração compete:
- a) Administrar os negócios da Cooperativa;
- b) Executar as decisões da Assembleia Geral;
- c) Apresentar o balanço anual e relatório de gestão;
- d) Convocar Assembleias Gerais.

---

## CAPÍTULO IV - DO PATRIMÔNIO E DA DISTRIBUIÇÃO DE SOBRAS

**Art. 10º** O patrimônio da Cooperativa é constituído por:
- a) Bens móveis e imóveis;
- b) Reservas estatutárias;
- c) Contribuições dos cooperados;
- d) Outras rendas e benefícios.

**Art. 11º** As sobras líquidas do exercício serão distribuídas proporcionalmente:
- I - 70% aos cooperados, em função das operações realizadas;
- II - 20% ao Fundo de Reserva;
- III - 10% ao Fundo de Assistência Técnica, Educacional e Social.

**Art. 12º** O rateio previsto no inciso I do artigo anterior será calculado com base no trabalho desenvolvido por cada cooperado (ITG 2002), utilizando a fórmula: (Horas do Cooperado / Total de Horas) × Excedente Disponível.

---

## CAPÍTULO V - DA FONTE DE DADOS E AUDITORIA

**Art. 13º** A Cooperativa utiliza o sistema **Digna** da Providentia Foundation para:
- Registro contábil e financeiro em banco de dados SQLite por tenant;
- Contabilização de horas de trabalho conforme ITG 2002;
- Geração de atas de assembleia com hash criptográfico (CADSOL);
- Cálculo automatizado de rateio de sobras;
- Sincronização segura com dados agregados apenas.

**Art. 14º** Todas as decisões da Cooperativa são registradas no sistema CADSOL (Cadastro de Decisões Soberanas) com hash SHA256 para garantir auditoria imutável.

**Art. 15º** O status da entidade no sistema é: **{{.Status}}**

---

## DISPOSIÇÕES FINAIS

**Art. 16º** Os casos omissos neste Estatuto serão resolvidos pelo Conselho de Administração, observadas as disposições legais.

**Art. 17º** Este Estatuto entra em vigor na data de seu registro.

---

*Documento gerado automaticamente pelo Sistema Digna - Providentia Foundation*  
*Hash de Integridade: {{.DocumentHash}}*  
*Data de Geração: {{.GeneratedAt}}*

**Assinaturas:**

- [ ] Presidente da Assembleia Constitutiva
- [ ] Secretário(a)  
- [ ] Tesoureiro(a)
`

type StatuteData struct {
	EntityName       string
	CNPJ             string
	NIRE             string
	RegistrationDate string
	Address          string
	Status           string
	DocumentHash     string
	GeneratedAt      string
}

type StatuteGenerator struct {
	lifecycleManager lifecycle.LifecycleManager
}

func NewStatuteGenerator(lm lifecycle.LifecycleManager) *StatuteGenerator {
	return &StatuteGenerator{
		lifecycleManager: lm,
	}
}

func (sg *StatuteGenerator) GenerateStatute(entityID string, entityName string, status string) (string, error) {
	if status != "FORMALIZED" {
		return "", fmt.Errorf("entity must be FORMALIZED to generate statute")
	}

	data := StatuteData{
		EntityName:       entityName,
		CNPJ:             "00.000.000/0001-91",
		NIRE:             "12.345.678.901",
		RegistrationDate: time.Now().AddDate(0, -1, 0).Format("2006-01-02"),
		Address:          "Rua da Cooperativa, 123 - Centro",
		Status:           status,
		DocumentHash:     sg.generateDocumentHash(entityID, entityName),
		GeneratedAt:      time.Now().Format("2006-01-02 15:04:05"),
	}

	tmpl, err := template.New("statute").Parse(statuteTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse statute template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute statute template: %w", err)
	}

	return buf.String(), nil
}

func (sg *StatuteGenerator) generateDocumentHash(entityID string, entityName string) string {
	data := fmt.Sprintf("%s:%s:%s:STATUTE_v1", entityID, entityName, time.Now().Format("20060102"))
	return generateQuickHash(data)
}

func generateQuickHash(data string) string {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])[:16]
}
