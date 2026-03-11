#!/bin/bash

# Define a base para as skills
BASE_DIR="docs/skills"

echo "🚀 Iniciando a criação da estrutura de Skills para o Digna..."

# 1. Skill: Backend Go
mkdir -p "$BASE_DIR/developing-digna-backend"
cat <<EOF > "$BASE_DIR/developing-digna-backend/SKILL.md"
yaml
name: developing-digna-backend
description: Especialista em Clean Architecture, DDD e TDD para o Motor Lume. Garante integridade financeira e proíbe pontos flutuantes.

# Instruções de Backend Go
Implementar lógica de negócio robusta no Motor Lume seguindo o rigor contábil.

### Regras Mandatórias
* **Anti-Float:** Proibido o uso de float32/64 em cálculos financeiros ou de tempo. Use exclusivamente int64 para centavos e minutos.
* **Soma Zero:** Todo lançamento deve validar que Débitos + Créditos == 0 via EntryValidator.
* **TDD:** Escrever testes unitários antes da lógica de domínio em internal/domain.

### Protocolo de Shutdown
* Implementar Graceful Shutdown para garantir o fechamento íntegro das conexões SQLite.
EOF

# 2. Skill: Frontend HTMX
mkdir -p "$BASE_DIR/rendering-digna-frontend"
cat <<EOF > "$BASE_DIR/rendering-digna-frontend/SKILL.md"
yaml
name: rendering-digna-frontend
description: Especialista em HTMX e templates Cache-Proof. Resolve problemas de versão garantindo renderização direta do disco.

# Instruções de Frontend
Garantir que a interface reflita sempre a última versão publicada sem interferência de cache.

### Protocolo Cache-Proof
* **Parsing On-Demand:** Proibido o uso de template.ParseGlob ou globais.
* **Direct-from-Disk:** O Handler deve carregar o arquivo .html em cada requisição usando template.ParseFiles.
* **Nomenclatura:** Usar apenas templates com sufixo _simple.html.

### Identidade Visual
* Aplicar paleta: Azul Soberania (#2A5CAA), Verde Suor (#4A7F3E) e Laranja Energia (#F57F17).
EOF

# 3. Skill: Lógica de Solidariedade
mkdir -p "$BASE_DIR/applying-solidarity-logic"
cat <<EOF > "$BASE_DIR/applying-solidarity-logic/SKILL.md"
yaml
name: applying-solidarity-logic
description: Guardião da ITG 2002 e da Tecnologia Social. Traduz jargões contábeis em linguagem popular de autogestão.

# Instruções de Negócio Social
Implementar a primazia do trabalho sobre o capital e a contabilidade invisível.

### Regras Sociotécnicas
* **Linguagem:** Substituir jargões (Ativo/Passivo) por termos de ação (Dinheiro na Gaveta / Contas a Pagar).
* **ITG 2002:** Bloquear mandatoriamente 10% para Reserva Legal e 5% para FATES antes de qualquer rateio.
* **Valor do Suor:** Registrar trabalho estritamente em minutos (int64).
EOF

# 4. Skill: Soberania de Dados
mkdir -p "$BASE_DIR/managing-sovereign-data"
cat <<EOF > "$BASE_DIR/managing-sovereign-data/SKILL.md"
yaml
name: managing-sovereign-data
description: Gerencia o ciclo de vida dos bancos SQLite. Garante isolamento físico por tenant e poder de saída dos dados.

# Instruções de Infraestrutura
Proteger a soberania dos dados de cada empreendimento.

### Regras de Isolamento
* **Physical Isolation:** Cada tenant possui seu próprio arquivo em data/entities/{entity_id}.db.
* **Lifecycle:** Uso obrigatório do LifecycleManager para abertura e migração de bases.
* **Exit Power:** Projetar schemas para exportação SQL simples e soberana.
EOF

# 5. Skill: Auditoria Fiscal
mkdir -p "$BASE_DIR/auditing-fiscal-compliance"
cat <<EOF > "$BASE_DIR/auditing-fiscal-compliance/SKILL.md"
yaml
name: auditing-fiscal-compliance
description: Opera o Accountant Dashboard e exportações SPED. Garante acesso Read-Only e mapeamento referencial.

# Instruções de Auditoria
Traduzir a autogestão popular em conformidade institucional.

### Segurança Fiscal
* **Read-Only:** Auditorias devem usar obrigatoriamente ?mode=ro no SQLite.
* **Tradução:** Mapear contas locais para o Plano de Contas Referencial da Receita Federal.
* **Hash de Integridade:** Registrar SHA256 de cada lote fiscal exportado.
EOF

echo "✅ Estrutura de Skills criada com sucesso em: $BASE_DIR"
