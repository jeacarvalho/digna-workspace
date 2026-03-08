# Makefile para Digna Workspace
# Facilita execução de testes em múltiplos módulos Go

.PHONY: help test test-all test-core test-integration test-distribution build clean lint

# Cores para output
BLUE := \033[36m
GREEN := \033[32m
YELLOW := \033[33m
RED := \033[31m
NC := \033[0m

ROOT_DIR := $(CURDIR)
MODULES := modules/core_lume modules/ui_web modules/distribution modules/lifecycle modules/legal_facade modules/reporting modules/sync_engine modules/integrations modules/integration_test

help:
	@echo "$(BLUE)Digna Workspace - Comandos disponíveis:$(NC)"
	@echo ""
	@echo "  $(GREEN)make test$(NC)          - Roda testes em todos os módulos"
	@echo "  $(GREEN)make test-core$(NC)     - Roda testes do core_lume"
	@echo "  $(GREEN)make test-integration$(NC) - Roda testes de integração"
	@echo "  $(GREEN)make test-distribution$(NC) - Roda testes do distribution"
	@echo "  $(GREEN)make test-coverage$(NC) - Roda testes com cobertura"
	@echo "  $(GREEN)make build$(NC)         - Builda todos os módulos"
	@echo "  $(GREEN)make clean$(NC)         - Limpa arquivos de build e caches"
	@echo "  $(GREEN)make lint$(NC)          - Roda linter em todos os módulos"
	@echo ""

test:
	@echo "$(BLUE)🧪 Executando testes em todos os módulos...$(NC)"
	@for module in $(MODULES); do \
		echo "$(YELLOW)▶ Testando $$module...$(NC)"; \
		cd $(ROOT_DIR)/$$module && go test ./... -v 2>&1 | grep -E "(PASS|FAIL|---)" || true; \
		echo ""; \
	done

test-core:
	@echo "$(BLUE)🧪 Executando testes do core_lume...$(NC)"
	@cd $(ROOT_DIR)/modules/core_lume && go test ./... -v

test-integration:
	@echo "$(BLUE)🧪 Executando testes de integração...$(NC)"
	@cd $(ROOT_DIR)/modules/integration_test && go test ./... -v

test-distribution:
	@echo "$(BLUE)🧪 Executando testes do distribution...$(NC)"
	@cd $(ROOT_DIR)/modules/distribution && go test ./... -v

test-coverage:
	@echo "$(BLUE)📊 Executando testes com cobertura...$(NC)"
	@for module in $(MODULES); do \
		echo "$(YELLOW)▶ Cobertura em $$module...$(NC)"; \
		cd $(ROOT_DIR)/$$module && go test ./... -cover 2>&1 | grep -E "(coverage|ok|FAIL)" || true; \
		echo ""; \
	done

build:
	@echo "$(BLUE)🔨 Buildando módulos...$(NC)"
	@cd $(ROOT_DIR)/modules/ui_web && go build -o ../../bin/digna-server .
	@echo "$(GREEN)✅ Build concluído: bin/digna-server$(NC)"

clean:
	@echo "$(BLUE)🧹 Limpando arquivos...$(NC)"
	@rm -rf $(ROOT_DIR)/bin/ $(ROOT_DIR)/dist/
	@find $(ROOT_DIR) -name "*.db" -delete
	@find $(ROOT_DIR) -name "*.db-journal" -delete
	@find $(ROOT_DIR) -name "*.db-wal" -delete
	@go clean -cache
	@echo "$(GREEN)✅ Limpeza concluída$(NC)"

lint:
	@echo "$(BLUE)🔍 Rodando linter...$(NC)"
	@for module in $(MODULES); do \
		echo "$(YELLOW)▶ Lint em $$module...$(NC)"; \
		cd $(ROOT_DIR)/$$module && go vet ./... 2>&1 | head -20 || true; \
	done

# Testes rápidos (sem verbose)
test-quick:
	@echo "$(BLUE)⚡ Testes rápidos...$(NC)"
	@for module in $(MODULES); do \
		cd $(ROOT_DIR)/$$module && go test ./... 2>&1 | grep -E "(ok|FAIL)" || true; \
	done
