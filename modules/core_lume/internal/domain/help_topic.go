package domain

import (
	"errors"
	"fmt"
	"time"
)

// Categorias de tópicos de ajuda
const (
	CategoriaCredito    = "CREDITO"
	CategoriaTributario = "TRIBUTARIO"
	CategoriaGovernanca = "GOVERNANCA"
	CategoriaGeral      = "GERAL"
)

// HelpTopic representa um tópico de ajuda educativa
type HelpTopic struct {
	ID           string // UUID
	Key          string // Chave única (ex: "cadunico", "inadimplencia")
	Title        string // Título em linguagem popular
	Summary      string // Resumo em 1 frase (para tooltips)
	Explanation  string // Explicação completa em linguagem popular
	WhyAsked     string // "Por que perguntamos isso?"
	Legislation  string // Legislação relacionada
	NextSteps    string // Próximos passos acionáveis
	OfficialLink string // Link para fonte oficial (ex: gov.br)
	Category     string // Categoria: CREDITO, TRIBUTARIO, GOVERNANCA, GERAL
	Tags         string // Tags para busca (comma-separated)

	// Metadados
	ViewCount int64 // Quantas vezes foi visualizado
	CreatedAt int64 // Unix timestamp
	UpdatedAt int64 // Unix timestamp
}

var (
	ErrHelpTopicInvalidKey         = errors.New("key is required and must be unique")
	ErrHelpTopicInvalidTitle       = errors.New("title is required")
	ErrHelpTopicInvalidExplanation = errors.New("explanation is required")
	ErrHelpTopicInvalidCategory    = errors.New("category is required")
	ErrHelpTopicNotFound           = errors.New("help topic not found")
)

// Validate verifica se o tópico é válido
func (h *HelpTopic) Validate() error {
	if h.Key == "" {
		return ErrHelpTopicInvalidKey
	}
	if h.Title == "" {
		return ErrHelpTopicInvalidTitle
	}
	if h.Explanation == "" {
		return ErrHelpTopicInvalidExplanation
	}
	if h.Category == "" {
		return ErrHelpTopicInvalidCategory
	}
	if !isValidCategory(h.Category) {
		return ErrHelpTopicInvalidCategory
	}
	return nil
}

// IsComplete verifica se todos os campos obrigatórios estão preenchidos
func (h *HelpTopic) IsComplete() bool {
	return h.Key != "" && h.Title != "" && h.Explanation != "" && h.Category != ""
}

// IncrementView incrementa o contador de visualizações
func (h *HelpTopic) IncrementView() {
	h.ViewCount++
	h.UpdatedAt = time.Now().Unix()
}

// String retorna representação textual do tópico
func (h *HelpTopic) String() string {
	return fmt.Sprintf("HelpTopic{Key: %s, Title: %s, Category: %s}",
		h.Key, h.Title, h.Category)
}

// GetTagsArray retorna as tags como slice
func (h *HelpTopic) GetTagsArray() []string {
	if h.Tags == "" {
		return []string{}
	}
	// Simple split by comma
	var tags []string
	current := ""
	for _, char := range h.Tags {
		if char == ',' {
			if current != "" {
				tags = append(tags, current)
				current = ""
			}
		} else {
			current += string(char)
		}
	}
	if current != "" {
		tags = append(tags, current)
	}
	return tags
}

// MatchesSearch verifica se o tópico corresponde à busca
func (h *HelpTopic) MatchesSearch(query string) bool {
	if query == "" {
		return true
	}

	// Search in key, title, summary, explanation, and tags
	searchable := h.Key + " " + h.Title + " " + h.Summary + " " + h.Explanation + " " + h.Tags
	return containsIgnoreCase(searchable, query)
}

// isValidCategory verifica se a categoria é válida
func isValidCategory(category string) bool {
	switch category {
	case CategoriaCredito, CategoriaTributario, CategoriaGovernanca, CategoriaGeral:
		return true
	}
	return false
}

// containsIgnoreCase verifica se a string contém a substring (case insensitive)
func containsIgnoreCase(s, substr string) bool {
	if len(substr) > len(s) {
		return false
	}

	sLower := toLower(s)
	substrLower := toLower(substr)

	return contains(sLower, substrLower)
}

// toLower converte string para minúsculas simples
func toLower(s string) string {
	result := ""
	for _, char := range s {
		if char >= 'A' && char <= 'Z' {
			result += string(char + 32) // 'a' - 'A' = 32
		} else {
			result += string(char)
		}
	}
	return result
}

// contains verifica se a string contém a substring
func contains(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	if len(substr) > len(s) {
		return false
	}

	for i := 0; i <= len(s)-len(substr); i++ {
		found := true
		for j := 0; j < len(substr); j++ {
			if s[i+j] != substr[j] {
				found = false
				break
			}
		}
		if found {
			return true
		}
	}
	return false
}

// GetCategoryLabel retorna o label da categoria
func GetCategoryLabel(category string) string {
	labels := map[string]string{
		CategoriaCredito:    "Crédito e Financiamento",
		CategoriaTributario: "Tributos e Impostos",
		CategoriaGovernanca: "Governança e Gestão",
		CategoriaGeral:      "Geral",
	}
	if label, ok := labels[category]; ok {
		return label
	}
	return category
}

// InitialHelpTopics contém os tópicos iniciais de ajuda
var InitialHelpTopics = []HelpTopic{
	{
		Key:          "cadunico",
		Title:        "O que é o CadÚnico?",
		Summary:      "É o cadastro do governo para programas sociais.",
		Explanation:  "O Cadastro Único (CadÚnico) reúne informações sobre famílias de baixa renda. Estar inscrito permite acesso a programas como Bolsa Família, Tarifa Social de Energia e linhas de crédito especiais.",
		WhyAsked:     "No Digna, informamos isso para encontrar programas de crédito que só atendem quem está no CadÚnico, como o 'Acredita no Primeiro Passo'.",
		Legislation:  "Decreto nº 6.135/2007",
		NextSteps:    "Se não está inscrito, procure o CRAS (Centro de Referência de Assistência Social) do seu município com documentos pessoais e comprovante de residência.",
		OfficialLink: "https://www.gov.br/cadunico",
		Category:     CategoriaCredito,
		Tags:         "cadastro,programa social,crédito",
	},
	{
		Key:          "inadimplencia",
		Title:        "O que é inadimplência?",
		Summary:      "É quando há dívidas não pagas registradas.",
		Explanation:  "Inadimplência significa que você tem contas atrasadas com bancos, lojas ou com o governo. Isso pode aparecer em sistemas como Serasa, SPC ou Dívida Ativa da União.",
		WhyAsked:     "Alguns programas de crédito exigem que você regularize essas dívidas antes de aplicar. Outros podem ajudar você a renegociar.",
		Legislation:  "Lei nº 10.820/2003 (Descontos em Folha)",
		NextSteps:    "Se tem dívidas, o Digna pode ajudar a identificar programas de renegociação como o 'Desenrola Pequenos Negócios' antes de buscar crédito novo.",
		OfficialLink: "https://www.gov.br/economia",
		Category:     CategoriaCredito,
		Tags:         "dívida,renegociação,crédito",
	},
	{
		Key:          "cnae",
		Title:        "O que é CNAE?",
		Summary:      "É o código que diz qual é a atividade do seu negócio.",
		Explanation:  "CNAE significa Classificação Nacional de Atividades Econômicas. É um número que o governo usa para saber se você vende comida, faz costura, presta serviço, etc.",
		WhyAsked:     "Programas de crédito usam o CNAE para saber se seu negócio se enquadra nas regras deles.",
		Legislation:  "Resolução CONCLA nº 1/2006",
		NextSteps:    "Se não sabe seu CNAE, consulte no cartão do CNPJ ou no site da Receita Federal.",
		OfficialLink: "https://www.gov.br/receitafederal",
		Category:     CategoriaTributario,
		Tags:         "atividade,cadastro,receita",
	},
	{
		Key:          "das_mei",
		Title:        "O que é o DAS MEI?",
		Summary:      "É o boleto mensal que o MEI paga.",
		Explanation:  "O DAS (Documento de Arrecadação do Simples Nacional) é o imposto que o Microempreendedor Individual paga todo mês. O valor é 5% do salário mínimo + valores fixos de ICMS e ISS.",
		WhyAsked:     "O Digna calcula automaticamente o valor do DAS e avisa quando está perto de vencer, para você não pagar multa.",
		Legislation:  "Lei Complementar nº 123/2006",
		NextSteps:    "O Digna gera o cálculo automaticamente. Você só precisa pagar até o dia 20 de cada mês.",
		OfficialLink: "https://www.gov.br/empresas-e-negocios",
		Category:     CategoriaTributario,
		Tags:         "MEI,imposto,boleto",
	},
	{
		Key:          "reserva_legal",
		Title:        "O que é Reserva Legal?",
		Summary:      "É uma parte do lucro que a lei manda guardar.",
		Explanation:  "A Reserva Legal é 10% do lucro da cooperativa que deve ser guardado por lei. Esse dinheiro não pode ser distribuído aos sócios — fica guardado para proteger a cooperativa em tempos difíceis.",
		WhyAsked:     "O Digna aplica automaticamente esse bloqueio antes de distribuir as sobras, para cumprir a lei e proteger o grupo.",
		Legislation:  "Lei nº 5.764/71 (Lei Geral das Cooperativas)",
		NextSteps:    "Não precisa fazer nada — o Digna calcula e guarda automaticamente.",
		OfficialLink: "https://www.planalto.gov.br/ccivil_03/leis/l5764.htm",
		Category:     CategoriaGovernanca,
		Tags:         "lucro,reserva,lei",
	},
	{
		Key:          "fates",
		Title:        "O que é o FATES?",
		Summary:      "É um fundo para ajudar outros grupos a se organizarem.",
		Explanation:  "O FATES (Fundo de Assistência Técnica, Educacional e Social) é 5% do lucro da cooperativa que é separado para ajudar outras cooperativas e grupos a se organizarem. É uma forma de solidariedade entre grupos.",
		WhyAsked:     "O Digna aplica automaticamente esse bloqueio antes de distribuir as sobras, para cumprir a lei da Economia Solidária.",
		Legislation:  "Lei nº 15.068/2024 (Lei Paul Singer)",
		NextSteps:    "Não precisa fazer nada — o Digna calcula e guarda automaticamente.",
		OfficialLink: "https://www.gov.br/trabalho-e-emprego",
		Category:     CategoriaGovernanca,
		Tags:         "fundo,solidariedade,lei",
	},
	{
		Key:          "icms_iss",
		Title:        "O que é ICMS e ISS?",
		Summary:      "São impostos sobre vendas e serviços.",
		Explanation:  "ICMS é o imposto que você paga quando vende produtos (comércio). ISS é o imposto que você paga quando presta serviços. Como MEI, você paga valores fixos desses impostos no DAS mensal.",
		WhyAsked:     "O Digna calcula automaticamente esses impostos quando você registra vendas no PDV.",
		Legislation:  "Lei Complementar nº 123/2006 (Simples Nacional)",
		NextSteps:    "O Digna calcula automaticamente. Você só precisa registrar vendas corretamente no PDV.",
		OfficialLink: "https://www.gov.br/fazenda",
		Category:     CategoriaTributario,
		Tags:         "imposto,venda,serviço,ICMS,ISS",
	},
	{
		Key:          "sped",
		Title:        "O que é SPED?",
		Summary:      "É o sistema digital de envio de informações fiscais.",
		Explanation:  "SPED (Sistema Público de Escrituração Digital) é o jeito que o governo exige que empresas enviem informações sobre vendas, compras e impostos pela internet. Substitui as declarações em papel.",
		WhyAsked:     "Seu contador usa o SPED para declarar seus impostos digitalmente para a Receita Federal.",
		Legislation:  "Decreto nº 6.022/2007",
		NextSteps:    "Deixe seu contador cuidar disso. O Digna já organiza os dados que ele precisa.",
		OfficialLink: "https://www.gov.br/receitafederal",
		Category:     CategoriaTributario,
		Tags:         "declaração,digital,fiscal",
	},
	{
		Key:          "efd_reinf",
		Title:        "O que é EFD-Reinf?",
		Summary:      "É uma declaração sobre pagamentos e retenções de impostos.",
		Explanation:  "EFD-Reinf é um sistema que informa ao governo quando sua empresa paga fornecedores, contratos ou faz retenções de impostos. É obrigatório para empresas que pagam mais de determinado valor por ano.",
		WhyAsked:     "Se sua cooperativa tem pagamentos altos a fornecedores, o contador precisa enviar essa declaração.",
		Legislation:  "Lei nº 13.137/2015",
		NextSteps:    "Seu contador cuida disso. O Digna registra os pagamentos que ele precisa declarar.",
		OfficialLink: "https://www.gov.br/receitafederal",
		Category:     CategoriaTributario,
		Tags:         "declaração,pagamento,retencao",
	},
	{
		Key:          "ecf",
		Title:        "O que é ECF?",
		Summary:      "É a declaração de imposto de renda da empresa.",
		Explanation:  "ECF (Escrituração Contábil Fiscal) é a declaração anual que empresas entregam para a Receita Federal informando lucros, prejuízos e pagamento de impostos. É como a declaração de imposto de renda, mas da empresa.",
		WhyAsked:     "Todo ano seu contador precisa enviar a ECF. O Digna ajuda organizando os dados contábeis durante o ano.",
		Legislation:  "Instrução Normativa RFB nº 1.422/2013",
		NextSteps:    "Entregue ao contador os relatórios do Digna no início de cada ano.",
		OfficialLink: "https://www.gov.br/receitafederal",
		Category:     CategoriaTributario,
		Tags:         "declaração,imposto de renda,anual",
	},
	{
		Key:          "itg_2002",
		Title:        "O que é ITG 2002?",
		Summary:      "São regras de como organizar as contas da cooperativa.",
		Explanation:  "ITG 2002 é um documento do governo que ensina como cooperativas devem organizar sua contabilidade. Define como registrar entradas, saídas, sobras e rateios de forma correta.",
		WhyAsked:     "O Digna segue as regras do ITG 2002 para fazer a contabilidade da sua cooperativa corretamente.",
		Legislation:  "ITG 2002 - CFC/CNC",
		NextSteps:    "O Digna já aplica essas regras automaticamente. Seu contador pode verificar isso nos relatórios.",
		OfficialLink: "https://www.cfc.org.br",
		Category:     CategoriaGovernanca,
		Tags:         "contabilidade,regras,cooperativa",
	},
	{
		Key:          "cadsol",
		Title:        "O que é CADSOL?",
		Summary:      "É o cadastro de entidades da Economia Solidária.",
		Explanation:  "CADSOL (Cadastro Nacional de Entidades da Economia Solidária) é um registro que identifica cooperativas, associações e outros grupos de economia solidária no Brasil. Ajuda o governo a fazer políticas públicas.",
		WhyAsked:     "Estar no CADSOL pode dar acesso a programas específicos para economia solidária e linhas de crédito diferenciadas.",
		Legislation:  "Lei nº 13.675/2018",
		NextSteps:    "Verifique se sua cooperativa está cadastrada no site do SIES (Sistema de Informações da Economia Solidária).",
		OfficialLink: "https://www.gov.br/trabalho-e-emprego",
		Category:     CategoriaGeral,
		Tags:         "cadastro,economia solidária,SIES",
	},
	{
		Key:          "exit_power",
		Title:        "O que é Exit Power?",
		Summary:      "É o poder de encerrar o vínculo com o contador.",
		Explanation:  "Exit Power significa que apenas a cooperativa pode desativar o vínculo com o contador. O contador não pode sair sozinho e deixar a cooperativa sem atendimento. Isso protege a cooperativa.",
		WhyAsked:     "O Digna implementa esse controle para garantir que apenas você decide quando trocar de contador.",
		Legislation:  "Autonomia da cooperativa (Lei nº 5.764/71)",
		NextSteps:    "Se precisar trocar de contador, você pode desativar o vínculo atual e criar um novo quando quiser.",
		OfficialLink: "",
		Category:     CategoriaGovernanca,
		Tags:         "contador,vínculo,proteção",
	},
	{
		Key:          "cardinalidade_temporal",
		Title:        "O que é Cardinalidade Temporal?",
		Summary:      "É a regra de ter apenas um contador ativo por vez.",
		Explanation:  "Cardinalidade Temporal significa que sua cooperativa só pode ter um contador ativo em cada momento. Se você contratar um contador novo, o vínculo antigo é desativado automaticamente. Isso evita confusão de dados.",
		WhyAsked:     "O Digna gerencia isso automaticamente para manter a organização e evitar que dois contadores trabalhem ao mesmo tempo.",
		Legislation:  "Organização contábil",
		NextSteps:    "Não precisa se preocupar. O Digna cuida de ativar e desativar vínculos automaticamente.",
		OfficialLink: "",
		Category:     CategoriaGovernanca,
		Tags:         "contador,vínculo,regra",
	},
	{
		Key:          "filtragem_temporal",
		Title:        "O que é Filtragem Temporal?",
		Summary:      "É o controle do que o contador pode ver por período.",
		Explanation:  "Filtragem Temporal significa que contadores inativos só podem ver os dados do período em que estavam ativos. Eles não podem acessar dados de antes ou depois do vínculo. Isso protege a privacidade.",
		WhyAsked:     "O Digna aplica essa regra automaticamente para proteger as informações da sua cooperativa.",
		Legislation:  "LGPD (Lei Geral de Proteção de Dados)",
		NextSteps:    "O Digna controla isso automaticamente. Contadores inativos têm acesso limitado.",
		OfficialLink: "https://www.gov.br/lgpd",
		Category:     CategoriaGovernanca,
		Tags:         "privacidade,proteção,contador",
	},
	{
		Key:          "tipo_entidade",
		Title:        "Qual o Tipo da sua Entidade?",
		Summary:      "É a forma como seu negócio está organizado legalmente.",
		Explanation:  "O Tipo de Entidade indica se você é uma Cooperativa, MEI, Associação, ou outro tipo de organização. Cada tipo tem regras diferentes de funcionamento, impostos e acesso a crédito.",
		WhyAsked:     "Programas de crédito têm regras diferentes para cada tipo de entidade. O Digna usa isso para encontrar as melhores oportunidades para você.",
		Legislation:  "Varia conforme o tipo",
		NextSteps:    "Verifique seu CNPJ ou documento de registro para confirmar seu tipo de entidade.",
		OfficialLink: "",
		Category:     CategoriaGeral,
		Tags:         "cooperativa,MEI,associação,tipo",
	},
	{
		Key:          "finalidade_credito",
		Title:        "Qual a Finalidade do Crédito?",
		Summary:      "É o motivo pelo qual você precisa do dinheiro.",
		Explanation:  "A Finalidade do Crédito explica para que você vai usar o dinheiro: comprar equipamentos, reformar, capital de giro, pagar dívidas, etc. Cada programa de crédito financia finalidades específicas.",
		WhyAsked:     "O Digna usa essa informação para encontrar programas de crédito que financiem exatamente o que você precisa.",
		Legislation:  "",
		NextSteps:    "Seja específico. Programas de crédito têm regras claras sobre o que podem financiar.",
		OfficialLink: "",
		Category:     CategoriaCredito,
		Tags:         "crédito,financiamento,motivo",
	},
	{
		Key:          "capital_social",
		Title:        "O que é Capital Social?",
		Summary:      "É o dinheiro que os sócios colocaram na cooperativa.",
		Explanation:  "Capital Social é o valor total que todos os sócios investiram na cooperativa quando ela foi criada. Esse dinheiro pertence à cooperativa e ajuda a mostrar que o negócio tem estrutura financeira.",
		WhyAsked:     "Alguns programas de crédito exigem um capital social mínimo para aprovar o financiamento.",
		Legislation:  "Lei nº 5.764/71 (Lei Geral das Cooperativas)",
		NextSteps:    "Consulte seu contrato social ou ata de fundação para saber o valor do capital social.",
		OfficialLink: "",
		Category:     CategoriaGovernanca,
		Tags:         "cooperativa,sócios,investimento",
	},
	{
		Key:          "insumo_produto",
		Title:        "Qual a diferença entre Insumo, Produto e Mercadoria?",
		Summary:      "São formas diferentes de classificar itens no estoque.",
		Explanation:  "Insumo é a matéria-prima que você usa para fazer algo (ex: farinha para fazer pão). Produto é o que você fabrica (ex: pão pronto). Mercadoria é o que você compra pronto para revender (ex: refrigerante). Cada um tem impostos diferentes.",
		WhyAsked:     "O Digna precisa saber a diferença para calcular corretamente os impostos e o valor do seu estoque.",
		Legislation:  "Código Tributário Nacional",
		NextSteps:    "Classifique corretamente cada item quando cadastrar no estoque. Se tiver dúvida, consulte seu contador.",
		OfficialLink: "",
		Category:     CategoriaTributario,
		Tags:         "estoque,insumo,produto,mercadoria",
	},
	{
		Key:          "fluxo_caixa",
		Title:        "O que é Fluxo de Caixa?",
		Summary:      "É o registro de entradas e saídas de dinheiro.",
		Explanation:  "Fluxo de Caixa é o acompanhamento de todo dinheiro que entra (vendas) e sai (despesas) da sua cooperativa. Mostra se você está ganhando ou perdendo dinheiro no dia a dia.",
		WhyAsked:     "O Digna calcula automaticamente seu fluxo de caixa para você saber se está no azul ou no vermelho.",
		Legislation:  "",
		NextSteps:    "Registre todas as entradas e saídas no módulo Caixa. O Digna faz os cálculos automaticamente.",
		OfficialLink: "",
		Category:     CategoriaGeral,
		Tags:         "caixa,entrada,saída,dinheiro",
	},
	{
		Key:          "contabilidade_formal",
		Title:        "O que é Contabilidade Formal?",
		Summary:      "É ter um contador registrado e fazer declarações corretamente.",
		Explanation:  "Contabilidade Formal significa que sua cooperativa tem um contador registrado no CRC que cuida dos livros contábeis e faz todas as declarações obrigatórias para o governo (impostos, balanços, etc.).",
		WhyAsked:     "Alguns programas de crédito só financiam cooperativas que têm contabilidade formal. Isso mostra que o negócio é sério e organizado.",
		Legislation:  "Lei nº 12.441/2011",
		NextSteps:    "Se ainda não tem contador, use o módulo 'Contador Social' do Digna para encontrar um.",
		OfficialLink: "",
		Category:     CategoriaGovernanca,
		Tags:         "contador,registro,formal",
	},
}
