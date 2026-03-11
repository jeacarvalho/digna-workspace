# Implementação: Geração de Dossiê Jurídico

**Data:** 2026-03-11  
**Módulo:** `legal_facade`  
**Status:** ✅ IMPLEMENTADO E TESTADO

## O que foi implementado

### 1. Função `GenerateDossier` no módulo `legal_facade`
- **Localização:** `modules/legal_facade/internal/document/generator.go:230-331`
- **API pública:** `modules/legal_facade/pkg/document/document.go:21-24`
- **Assinatura:** `func (g *Generator) GenerateDossier(entityID string, entityName string, status string) (string, string, error)`
- **Retorno:** `(conteúdo_do_dossiê, hash_sha256, erro)`

### 2. Funcionalidades do dossiê
- Gera documento completo de formalização CADSOL
- Inclui histórico de decisões soberanas
- Verifica critérios de formalização (mínimo 3 decisões)
- Calcula hash SHA256 de integridade
- Adapta conteúdo baseado no status da entidade (DREAM, FORMALIZED, etc.)
- Inclui seções: documentos anexos, disposições finais, assinaturas requeridas

### 3. Sistema de hash de integridade
- Hash SHA256 calculado sobre o conteúdo do documento
- Padrão: `conteúdo:entityID:DIGNA_DOSSIER_SALT_v1`
- Incluído no documento para validação posterior
- Segue mesmo padrão do `core_lume` para consistência

## Como usar

```go
import "github.com/providentia/digna/legal_facade/pkg/document"

// 1. Criar gerador
generator := document.NewGenerator(lifecycleManager)

// 2. Gerar dossiê
dossierContent, dossierHash, err := generator.GenerateDossier(
    "entidade_123",
    "Cooperativa Teste",
    "DREAM", // ou "FORMALIZED", etc.
)

if err != nil {
    // tratar erro
}

// 3. Usar conteúdo e hash
fmt.Printf("Dossiê gerado (%d caracteres)\n", len(dossierContent))
fmt.Printf("Hash de integridade: %s\n", dossierHash)
```

## Testes implementados

### Arquivo: `modules/legal_facade/dossier_test.go`
1. **TestGenerateDossier/Gerar_Dossiê_com_Decisões** - Testa geração com 3 decisões
2. **TestGenerateDossier/Gerar_Dossiê_sem_Decisões** - Testa geração sem decisões
3. **TestGenerateDossier/Gerar_Dossiê_Entidade_Formalizada** - Testa para entidade formalizada
4. **TestDossierHashIntegrity** - Valida integridade do hash

### Cobertura:
- ✅ Geração básica do documento
- ✅ Inclusão de hash no conteúdo
- ✅ Formatação correta para diferentes status
- ✅ Validação de critérios de formalização
- ✅ Integridade do hash SHA256

## Integração com sistema existente

### Módulo `legal_facade` já possuía:
- `GenerateAssemblyMinutes()` - Geração de atas de assembleia
- `GenerateIdentityCard()` - Geração de carteira de identidade
- `GenerateStatute()` - Geração de estatuto social
- `FormalizationSimulator` - Simulação de formalização

### Nova função se integra com:
1. **FormalizationSimulator** - Usa `CheckFormalizationCriteria()` para validar elegibilidade
2. **LegalRepository** - Busca decisões via `GetAllDecisions()`
3. **Sistema de hash** - Reutiliza padrão SHA256 do `core_lume`

## Padrões seguidos

1. **Consistência com código existente** - Mesma estrutura de templates e funções auxiliares
2. **Error handling** - Retorna erros descritivos seguindo padrão Go
3. **Testabilidade** - Funções puras, injetáveis, com testes unitários
4. **Documentação** - Comentários em código e exemplos de uso
5. **Segurança** - Hash SHA256 para integridade do documento

## Próximos passos possíveis

1. **Integração com frontend** - Expor via API REST/GraphQL
2. **Exportação para PDF** - Converter markdown para PDF
3. **Assinatura digital** - Integrar com sistema de assinaturas
4. **Cache de documentos** - Otimizar geração frequente
5. **Templates customizáveis** - Permitir personalização por organização

## Aprendizados técnicos

1. **Go templates com funções customizadas** - Uso de `template.FuncMap` para lógica condicional
2. **Cálculo de hash recursivo** - Desafio de calcular hash de documento que inclui o próprio hash
3. **Separação de preocupações** - Template antes/depois da seção de hash
4. **Testes de integridade** - Validação de hash SHA256 em testes
5. **Compatibilidade com sistema existente** - Reutilização de padrões e estruturas

## Status final

✅ **Funcionalidade completa e testada**  
✅ **Integrada com módulo existente**  
✅ **Documentação criada**  
✅ **Testes passando**  
✅ **Pronta para uso em produção**