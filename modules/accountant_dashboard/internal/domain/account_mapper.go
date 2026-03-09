package domain

var defaultAccountMappings = []AccountMapping{
	{
		LocalCode:    "1.1.01",
		LocalName:    "Gaveta / Caixa",
		StandardCode: "1.1.01.00.00",
		StandardName: "Disponibilidades - Caixa",
	},
	{
		LocalCode:    "1.1.02",
		LocalName:    "Banco / Conta",
		StandardCode: "1.1.02.00.00",
		StandardName: "Disponibilidades - Bancos Conta Movimento",
	},
	{
		LocalCode:    "1.2.01",
		LocalName:    "Estoque de Insumos",
		StandardCode: "1.2.01.00.00",
		StandardName: "Estoques - Insumos",
	},
	{
		LocalCode:    "2.1.01",
		LocalName:    "Quem Fornece",
		StandardCode: "2.1.01.00.00",
		StandardName: "Fornecedores a Pagar",
	},
	{
		LocalCode:    "2.2.01",
		LocalName:    "Capital Social",
		StandardCode: "2.2.01.00.00",
		StandardName: "Capital Social",
	},
	{
		LocalCode:    "2.2.02",
		LocalName:    "Fundo FATES",
		StandardCode: "2.2.02.00.00",
		StandardName: "Reservas Estatutárias - FATES",
	},
	{
		LocalCode:    "2.2.03",
		LocalName:    "Reserva Legal",
		StandardCode: "2.2.03.00.00",
		StandardName: "Reserva Legal",
	},
	{
		LocalCode:    "3.1.01",
		LocalName:    "Nossas Vendas",
		StandardCode: "3.1.01.00.00",
		StandardName: "Receita Bruta de Vendas",
	},
	{
		LocalCode:    "3.2.01",
		LocalName:    "Despesa com Insumos",
		StandardCode: "3.2.01.00.00",
		StandardName: "Custo dos Insumos Vendidos",
	},
	{
		LocalCode:    "3.2.02",
		LocalName:    "Despesa com Trabalho",
		StandardCode: "3.2.02.00.00",
		StandardName: "Despesas com Mão de Obra",
	},
}

type DefaultAccountMapper struct {
	mappings map[string]AccountMapping
}

func NewDefaultAccountMapper() *DefaultAccountMapper {
	mappings := make(map[string]AccountMapping)
	for _, m := range defaultAccountMappings {
		mappings[m.LocalCode] = m
	}
	return &DefaultAccountMapper{mappings: mappings}
}

func (m *DefaultAccountMapper) GetMapping(localCode string) (AccountMapping, bool) {
	val, ok := m.mappings[localCode]
	return val, ok
}

func (m *DefaultAccountMapper) GetAllMappings() []AccountMapping {
	return defaultAccountMappings
}
