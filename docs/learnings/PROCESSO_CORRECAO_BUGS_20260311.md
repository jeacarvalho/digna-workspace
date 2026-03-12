# 📚 APRENDIZADOS: Processo de Correção de Bugs - O Que Faltou e Como Corrigir

**Data:** 11/03/2026  
**Contexto:** Tarefa "Corrigir bugs diversos após implementações até aqui" foi considerada concluída enquanto bugs ainda existiam  
**Problema:** Processo de validação permitiu conclusão sem verificar correção real dos bugs

---

## 🚨 FALHAS IDENTIFICADAS NO PROCESSO

### 1. **Validação Superficial de Testes E2E**
- ✅ Contava arquivos com `*e2e*test*.go`
- ❌ **NÃO executava** os testes Playwright
- ❌ **NÃO validava** se testes cobriam requisitos específicos do prompt

### 2. **Falta de Validação de Requisitos Específicos**
- Prompt pedia: *"Testar com dados reais (cafe_digna, contador_social)"*
- ✅ Testes usavam `entity_id=test_coop`
- ❌ **NÃO usavam** `cafe_digna` como solicitado
- ❌ **NÃO testavam** o erro específico mencionado

### 3. **Ausência de Smoke Test Obrigatório para Bugs**
- Tarefa era de **correção de bug**
- Smoke test era apenas "recomendado"
- ❌ **NÃO obrigatório** para tarefas de correção
- ❌ **NÃO testava** URLs problemáticas mencionadas

### 4. **Validação de "Testes Criados" vs "Bugs Corrigidos"**
- ✅ Verificava que arquivos de teste foram criados
- ❌ **NÃO verificava** se os bugs foram realmente corrigidos
- ❌ **NÃO testava** em ambiente real (servidor ativo)

---

## 🔧 CORREÇÕES IMPLEMENTADAS

### 1. **Novo Script: `validate_task_requirements.sh`**
```bash
# Valida requisitos específicos do prompt
./scripts/validate_task_requirements.sh --task=TASK_ID

# O que faz:
# 1. Extrai requisitos do prompt (Playwright, dados reais, etc.)
# 2. Valida se testes E2E existem e EXECUTAM
# 3. Verifica uso de IDs específicos (cafe_digna, etc.)
# 4. Valida correção de erros específicos mencionados
```

### 2. **Smoke Test Específico para Bugs: `smoke_test_bug_fixes.sh`**
```bash
# Smoke test obrigatório para tarefas de correção
./scripts/dev/smoke_test_bug_fixes.sh --task=TASK_ID

# O que faz:
# 1. Extrai URLs problemáticas do prompt
# 2. Inicia servidor real
# 3. Testa cada URL problemática
# 4. Verifica se erros específicos foram corrigidos
# 5. FAIL HARD se bugs persistem
```

### 3. **Integração no Fluxo de Conclusão**
Atualizado `conclude_task.sh`:
```bash
# Nova validação (passo 1.7)
echo "7. Validação de requisitos específicos do prompt..."
if [ -f "./scripts/validate_task_requirements.sh" ]; then
    ./scripts/validate_task_requirements.sh --task=${TASK_ID}
    
    # Se é tarefa de correção de bug, FAIL HARD
    if echo "$TASK_NAME" | grep -qi "corrigir\|bug\|erro\|fix"; then
        echo "❌❌❌ TAREFA DE CORREÇÃO DE BUG - NÃO PODE SER CONCLUÍDA ❌❌❌"
        VALIDATION_PASSED=false
    fi
fi
```

### 4. **Regras Específicas para Tarefas de Correção**
1. **Smoke test OBRIGATÓRIO** para tarefas com "corrigir", "bug", "erro", "fix" no nome
2. **Validação de URLs problemáticas** OBRIGATÓRIA
3. **Teste com dados reais** OBRIGATÓRIO se mencionado no prompt
4. **FAIL HARD** se bugs não foram corrigidos

---

## 🎯 CHECKLIST REVISADO PARA TAREFAS DE CORREÇÃO

### ✅ ANTES DE CONSIDERAR "CONCLUÍDA":
1. **Testes E2E Playwright EXECUTADOS** (não apenas criados)
2. **IDs reais testados** (cafe_digna, contador_social se mencionados)
3. **Smoke test executado** com servidor real
4. **URLs problemáticas testadas** e validadas
5. **Erros específicos** mencionados no prompt **VERIFICADOS como corrigidos**
6. **Logs de erro** analisados para garantir correção

### 🚨 GATILHOS PARA FAIL HARD:
- Tarefa contém "corrigir", "bug", "erro", "fix" no nome
- Prompt menciona URLs específicas com erros
- Prompt menciona erros específicos (database connection closed, etc.)
- Smoke test falha

---

## 📊 MÉTRICAS DE QUALIDADE ADICIONAIS

### Para tarefas de correção:
- **Taxa de reincidência de bugs**: 0% aceitável
- **Cobertura de testes para o bug específico**: 100% obrigatória
- **Evidências de correção**: logs antes/depois obrigatórios

### Validação cruzada:
1. Prompt → Requisitos identificados
2. Requisitos → Testes criados
3. Testes → Execução bem-sucedida
4. Execução → Smoke test passando
5. Smoke test → Bugs corrigidos

---

## 💡 APRENDIZADOS GERAIS

### Processo > Código
- **Processos robustos previnem entregas com bugs**
- Validação automática NÃO substitui validação do que foi solicitado
- "Testes criados" ≠ "Bugs corrigidos"

### Especificidade é Fundamental
- Testar com `entity_id=test_coop` ≠ testar com `entity_id=cafe_digna`
- Testes unitários ≠ testes de integração com servidor real
- Contar arquivos de teste ≠ executar testes

### Fail Fast, Fail Hard
- Para correções de bug: **FAIL HARD** se não corrigido
- Melhor falhar na validação que entregar bug não corrigido
- Smoke test obrigatório para validação real

---

## 🚀 PRÓXIMOS PASSOS

1. **Implementar** o novo processo em todas as tarefas futuras
2. **Retroativamente validar** tarefas de correção anteriores
3. **Criar dashboard** de métricas de qualidade de correções
4. **Documentar casos de teste** para bugs comuns
5. **Automatizar** mais validações baseadas em padrões de prompt

---

**Status do Processo Corrigido:** ✅ IMPLEMENTADO  
**Próxima Tarefa de Correção:** Será validada com novo processo  
**Expectativa:** 0% de bugs não corrigidos em entregas "concluídas"