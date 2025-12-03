# AsiriGen - Makefile
# Facilita o desenvolvimento e build do AsiriGen

.PHONY: build clean test install run help deps

# Variáveis
BINARY_NAME=asirigen
VERSION=2.0.0
BUILD_DIR=builds
GO_VERSION=1.21

# Cores para output
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[1;33m
BLUE=\033[0;34m
NC=\033[0m # No Color

help: ## Mostra esta ajuda
	@echo "$(BLUE)AsiriGen - Gerador de wordlists inteligente$(NC)"
	@echo
	@echo "$(YELLOW)Comandos disponíveis:$(NC)"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  $(GREEN)%-15s$(NC) %s\n", $$1, $$2}' $(MAKEFILE_LIST)

deps: ## Instala dependências
	@echo "$(BLUE)Instalando dependências...$(NC)"
	go mod tidy
	go mod download

build: deps ## Compila o AsiriGen
	@echo "$(BLUE)Compilando AsiriGen...$(NC)"
	go build -ldflags="-s -w -X main.version=$(VERSION)" -o $(BINARY_NAME) .
	@echo "$(GREEN)✓ Compilação concluída: $(BINARY_NAME)$(NC)"

build-all: deps ## Compila para todas as plataformas
	@echo "$(BLUE)Compilando para todas as plataformas...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@echo "Linux AMD64..."
	@GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)_linux_amd64 .
	@echo "Linux ARM64..."
	@GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)_linux_arm64 .
	@echo "Windows AMD64..."
	@GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)_windows_amd64.exe .
	@echo "Windows ARM64..."
	@GOOS=windows GOARCH=arm64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)_windows_arm64.exe .
	@echo "macOS AMD64..."
	@GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)_darwin_amd64 .
	@echo "macOS ARM64..."
	@GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)_darwin_arm64 .
	@echo "$(GREEN)✓ Compilação para todas as plataformas concluída$(NC)"

test: ## Executa os testes
	@echo "$(BLUE)Executando testes...$(NC)"
	go test -v ./...

test-coverage: ## Executa testes com cobertura
	@echo "$(BLUE)Executando testes com cobertura...$(NC)"
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)✓ Relatório de cobertura gerado: coverage.html$(NC)"

run: build ## Compila e executa o AsiriGen
	@echo "$(BLUE)Executando AsiriGen...$(NC)"
	./$(BINARY_NAME) --help

run-example: build ## Executa exemplo com Microsoft
	@echo "$(BLUE)Executando exemplo com Microsoft...$(NC)"
	./$(BINARY_NAME) generate --company Microsoft --min-length 6 --max-length 12 | head -20

run-leet: build ## Executa exemplo com leetspeak
	@echo "$(BLUE)Executando exemplo com leetspeak...$(NC)"
	./$(BINARY_NAME) generate --words "admin,test" --leet --min-length 4 --max-length 12 | head -20

install: build ## Instala o AsiriGen no sistema
	@echo "$(BLUE)Instalando AsiriGen...$(NC)"
	sudo cp $(BINARY_NAME) /usr/local/bin/
	@echo "$(GREEN)✓ AsiriGen instalado em /usr/local/bin/$(NC)"

uninstall: ## Remove o AsiriGen do sistema
	@echo "$(BLUE)Removendo AsiriGen...$(NC)"
	sudo rm -f /usr/local/bin/$(BINARY_NAME)
	@echo "$(GREEN)✓ AsiriGen removido$(NC)"

clean: ## Limpa arquivos de build
	@echo "$(BLUE)Limpando arquivos de build...$(NC)"
	rm -f $(BINARY_NAME)
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html
	@echo "$(GREEN)✓ Limpeza concluída$(NC)"

lint: ## Executa linter
	@echo "$(BLUE)Executando linter...$(NC)"
	golangci-lint run

format: ## Formata o código
	@echo "$(BLUE)Formatando código...$(NC)"
	go fmt ./...

check: format lint test ## Executa todas as verificações

demo: build ## Executa demonstração completa
	@echo "$(BLUE)=== Demonstração do AsiriGen ===$(NC)"
	@echo
	@echo "$(YELLOW)1. Ajuda:$(NC)"
	./$(BINARY_NAME) --help
	@echo
	@echo "$(YELLOW)2. Versão:$(NC)"
	./$(BINARY_NAME) version
	@echo
	@echo "$(YELLOW)3. Exemplo com Microsoft:$(NC)"
	./$(BINARY_NAME) generate --company Microsoft --min-length 6 --max-length 12 | head -10
	@echo
	@echo "$(YELLOW)4. Exemplo com leetspeak:$(NC)"
	./$(BINARY_NAME) generate --words "admin,test" --leet --min-length 4 --max-length 12 | head -10
	@echo
	@echo "$(YELLOW)5. Exemplo combinado:$(NC)"
	./$(BINARY_NAME) generate --company "TechCorp" --words "admin,user" --min-length 6 --max-length 14 | head -10
	@echo
	@echo "$(GREEN)=== Demonstração concluída ===$(NC)"

release: clean build-all ## Cria release com todos os binários
	@echo "$(BLUE)Criando release...$(NC)"
	@mkdir -p release
	@cp -r $(BUILD_DIR)/* release/
	@cp README.md LICENSE release/ 2>/dev/null || true
	@cd release && sha256sum * > checksums.txt
	@echo "$(GREEN)✓ Release criada na pasta 'release/'$(NC)"

# Comando padrão
.DEFAULT_GOAL := help
