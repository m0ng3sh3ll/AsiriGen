# 🔐 AsiriGen 2.0 - Gerador de Wordlists Inteligente para Pentesting

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Platform](https://img.shields.io/badge/Platform-Windows%20%7C%20Linux%20%7C%20macOS-lightgrey.svg)]()

**AsiriGen** é uma ferramenta avançada de geração de wordlists desenvolvida especificamente para operações de **red team** e **pentest interno**. Diferente de geradores comuns, o AsiriGen combina **OSINT (Open Source Intelligence)**, **Contexto Corporativo** e **Templates Dinâmicos** para criar listas de senhas altamente prováveis e focadas no alvo.

---

## ✨ Novidades Da Versão

### 🌐 **Coleta de Inteligência (OSINT)**
O AsiriGen agora possui um **scraper inteligente** embutido que utiliza um navegador real (Headless Chrome) para navegar no site do alvo, ignorar certificados SSL e extrair palavras-chave relevantes do conteúdo da página.

### 🎨 **Templates Dinâmicos**
Crie seus próprios padrões de senha usando arquivos YAML simples. Defina formatos como `{company}@{year}` ou `Senh4_{leet}` sem precisar recompilar o código.

### 🌍 **Internacionalização (i18n)**
Suporte nativo para **Português (pt-BR)** e Inglês (en-US). Gera variações com meses, estações e termos corporativos no idioma local do alvo (ex: "Janeiro2024" vs "January2024").

---

## 🚀 Funcionalidades Principais

### 🎯 **Foco em Red Team**
- **Máxima cobertura** de variações de senhas corporativas
- **Geração inteligente** baseada em padrões reais de usuários
- **Contexto corporativo** específico (abreviações, setores, etc.)

### 🏢 **Análise Corporativa**
- Suporte a nomes compostos (ex: "New Corp" → "new_corp", "NewCorp")
- Detecção automática de abreviações e iniciais
- Padrões específicos da indústria (Tech, Finance, Health, etc.)

### 🔤 **Geração Avançada**
- **Modo Leetspeak** contextual (ex: "admin" → "@dm1n")
- Variações de case (minúsculo, maiúsculo, título)
- Combinações com anos (2020-2025)
- Números e símbolos comuns

---

## 📦 Instalação

### Pré-requisitos
- **Go 1.21+**
- **Google Chrome** (para funcionalidade OSINT)

### Instalação via Go
```bash
# Clonar o repositório
git clone github.com/m0ng3sh3ll/AsiriGen.git
cd asirigen

# Instalar dependências
go mod tidy

# Compilar
go build -o asirigen.exe .
```

---

## 📖 Guia de Uso

### Sintaxe Básica
```bash
./asirigen.exe generate [flags]
```

### Flags Principais

| Flag | Descrição | Exemplo |
|------|-----------|---------|
| `--company` | Nome da empresa/alvo | `--company "Microsoft"` |
| `--url` | **[NOVO]** URL para extração OSINT | `--url "https://alvo.com"` |
| `--lang` | **[NOVO]** Idioma (pt-br, en-us) | `--lang pt-br` |
| `--patterns-file` | **[NOVO]** Arquivo de templates YAML | `--patterns-file "pat.yaml"` |
| `--leet` | Ativar modo leetspeak | `--leet` |
| `--min-length` | Tamanho mínimo | `--min-length 8` |

### Comandos Úteis

- `init`: Cria um arquivo `patterns.yaml` padrão no diretório atual.
  ```bash
  ./asirigen.exe init
  ```

### Exemplos Práticos
| `--output, -o` | Arquivo de saída | `-o wordlist.txt` |

### Exemplos Práticos

#### 1. O "Combo Supremo" (Recomendado)
Gera uma wordlist completa usando o nome da empresa, extraindo dados do site oficial, aplicando padrões em português e formatando com templates personalizados.

```bash
./asirigen.exe generate \
  --company "Alvo" \
  --url "https://www.alvo.com.br" \
  --lang pt-br \
  --leet \
  --patterns-file "examples/patterns.yaml" \
  --verbose
```

#### 2. Geração Básica Rápida
```bash
./asirigen.exe generate --company "EmpresaX" --lang pt-br
```

#### 3. Apenas OSINT (Extração de Site)
```bash
./asirigen.exe generate --url "https://site-alvo.com" --verbose
```

---

## 🛠️ Templates Personalizados (YAML)

O sistema de templates permite que você defina exatamente como as senhas devem ser formadas. Edite o arquivo `patterns.yaml` para adicionar seus próprios padrões.

### Variáveis Disponíveis

| Variável | Descrição | Exemplo (Company="Microsoft", Year=2024) |
|----------|-----------|------------------------------------------|
| `{company}` | Variações do nome da empresa (Original, Title, Upper, Iniciais) | `microsoft`, `Microsoft`, `MICROSOFT`, `M.S` |
| `{word}` | Palavras comuns (fornecidas ou extraídas via OSINT) | `admin`, `Admin`, `ADMIN` |
| `{year}` | Variações do ano (4 dígitos e 2 dígitos) | `2024`, `24` |
| `{sep}` | Separadores comuns | `.`, `_`, `-`, `@`, `#` |
| `{num}` | Sequências numéricas comuns | `1`, `123`, `123456`, `01` |
| `{special}` | Caracteres especiais | `!`, `@`, `#`, `$` |
| `{leet}` | Versão leetspeak da palavra base | `m1cr0s0ft`, `@dm1n` |

### Comportamento do Sistema com Templates

1. **Substituição Inteligente**: O sistema varre cada template e substitui os tokens por todas as suas variações possíveis (produto cartesiano). Por exemplo, um único template `{company}_{year}` pode gerar centenas de senhas (todas as variações do nome da empresa combinadas com todas as variações de anos).

2. **Texto Estático**: Qualquer texto fora das chaves `{}` é mantido como está.
   - Exemplo: `Super_{company}!` gerará `Super_Microsoft!`, `Super_microsoft!`, etc.

3. **Validação de Tokens**: O sistema valida automaticamente os templates. Se você digitar um token errado (ex: `{usuar}` em vez de `{word}`), ele emitirá um **aviso no terminal** (`⚠️ AVISO: Token desconhecido...`), mas continuará gerando o restante das senhas normalmente, tratando o token errado como texto literal.

### Exemplos de Templates

```yaml
patterns:
  # Padrões corporativos clássicos
  - "{company}{year}"          # microsoft2024, Microsoft24
  - "{company}{sep}{year}"     # Microsoft.2024, Microsoft_24

  # Padrões com palavras comuns
  - "{word}{sep}{company}"     # admin.microsoft, Admin_Microsoft
  - "{company}{sep}{word}"     # Microsoft.admin

  # Padrões complexos
  - "{company}_{word}_{year}"  # Microsoft_Admin_2024
  - "{leet}#{year}"            # m1cr0s0ft#2024
  - "{word}{num}"              # Admin123, Senha123456
  - "{company}{special}"       # Microsoft!
```

---

## 📁 Estrutura do Projeto

```
asirigen/
├── cmd/                 # Comandos CLI (root, generate, version, banner)
├── internal/
│   ├── generator/       # Motor de geração de senhas
│   ├── scraper/         # Módulo OSINT (Chromedp)
│   ├── patterns/        # Gerenciador de templates e padrões
│   └── i18n/            # Internacionalização (pt-BR/en-US)
├── examples/            # Arquivos de exemplo (YAML, JSON)
└── main.go              # Entrypoint
```

---

## ⚠️ Aviso Legal

Esta ferramenta é destinada **apenas** para:
- Testes de penetração autorizados
- Red team operations
- Auditorias de segurança e pesquisa

**NÃO** use esta ferramenta para atividades maliciosas ou não autorizadas. O uso inadequado é de responsabilidade exclusiva do usuário.

---

**Desenvolvido com 💜 e Go.**
