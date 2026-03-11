# 🚫 Antipadrões - Projeto Digna

**Última atualização:** 11/03/2026
**Fonte:** Aprendizados de sessões anteriores

---

## ❌ O que NÃO fazer

Esta lista é baseada em erros comuns identificados durante implementações.


## 🔄 Sessão 10/03/2026

1. **ANTIPADRÃO:** Adicionar antipadrão: "Não analisar aprendizados ao final da sessão"

## 🔄 Sessão 11/03/2026

1. **ANTIPADRÃO:** "Implementar funcionalidade sem verificar se já existe" - `legal_facade` já tem 80% da funcionalidade necessária
2. **ANTIPADRÃO:** "Documentar aprendizados em arquivos temporários" - `.agent_context.md` é excluído no `end_session.sh`
3. **ANTIPADRÃO:** "Não consultar `docs/skills/` antes de implementar" - 5 skills críticas já documentadas
4. **ANTIPADRÃO:** "Recriar padrões já implementados" - SHA256, file download, cache-proof templates já existem
5. **ANTIPADRÃO:** "Usar `process_task.sh --file=\"nome.md\"` sem validar parsing" - Problema com caracteres especiais
