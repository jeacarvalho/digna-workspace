# 🎯 Próximos Passos - Projeto Digna

**Última atualização:** 10/03/2026

---

## 🚀 Próxima Tarefa

Escolha uma tarefa do backlog ou crie uma nova:

1. Use `./process_task.sh "sua descrição de tarefa"`
2. Siga o checklist pré-implementação
3. Documente aprendizados com `./conclude_task.sh`

## Correção de Bug HTMX no Estoque (20260310_180750)
**Status:** parcial  
**Concluído em:** 10/03/2026  
**Duração:** ~45 minutos  

### Próximas Ações:
1. **Corrigir problema original**: Tipo `TotalValue` no template ainda precisa conversão `int64` → `float64`
2. **Resolver mismatch de campos**: Adicionar `quantity` ao formulário ou torná-lo opcional no handler
3. **Completar testes E2E**: Executar com autenticação funcionando
4. **Validar bug fix**: Testar criação de item de estoque end-to-end

### Decisões Pendentes:
- Como lidar com campo `quantity` faltante no formulário?
- Implementar sistema de seeding para entidades de teste?
- Refatorar servidor para aceitar porta configurável?

### Links:
- Aprendizados: `docs/learnings/20260310_183000_stock_htmx_bug_fix_learnings.md`
- Test runner: `scripts/dev/e2e_test_runner.sh`
- Testes E2E: `tests/stock-crud.spec.js`

