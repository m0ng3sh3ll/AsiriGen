package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"asirigen/internal/generator"
	"asirigen/internal/i18n"
	"asirigen/internal/patterns"
	"asirigen/internal/scraper"

	"github.com/spf13/cobra"
)

// generateCmd representa o comando de geração
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Gerar wordlist personalizada",
	Long: `Gera uma wordlist personalizada baseada nos parâmetros fornecidos.
Pode usar nome de corporação, palavras comuns, ou ambos.`,
	RunE: runGenerate,
}

// Variáveis globais para flags
var (
	templates    []string
	patternsFile string
	targetUrl    string
	lang         string
)

func init() {
	generateCmd.Flags().StringSliceVarP(&templates, "template", "t", []string{}, "Templates personalizados (ex: {company}@{year})")
	generateCmd.Flags().StringVarP(&patternsFile, "patterns-file", "P", "patterns.yaml", "Arquivo YAML com templates personalizados")
	generateCmd.Flags().StringVar(&targetUrl, "url", "", "URL alvo para extrair palavras-chave (OSINT)")
	generateCmd.Flags().StringVar(&lang, "lang", "en", "Idioma para geração (pt-br, en-us)")
}

func ensurePatternsFile() error {
	// Se patterns.yaml não existir, cria um default
	if _, err := os.Stat("patterns.yaml"); os.IsNotExist(err) {
		defaultContent := `patterns:
  - "{company}{year}"
  - "{company}@{year}"
  - "{company}#{year}"
  - "{company}{sep}{word}"
  - "{word}{sep}{company}"
  - "admin{sep}{year}"
  - "{company}2024"
  - "{company}2025"
  - "{word}123"
  - "{company}!"
`
		if err := os.WriteFile("patterns.yaml", []byte(defaultContent), 0644); err != nil {
			return fmt.Errorf("erro ao criar patterns.yaml padrão: %v", err)
		}
		fmt.Fprintf(os.Stderr, "ℹ️  Arquivo 'patterns.yaml' não encontrado. Um arquivo padrão foi criado.\n")
	}
	return nil
}

func runGenerate(cmd *cobra.Command, args []string) error {
	// Validar flags
	if err := validateFlags(); err != nil {
		return err
	}

	// Garantir que patterns.yaml existe se for o arquivo padrão
	if patternsFile == "patterns.yaml" {
		if err := ensurePatternsFile(); err != nil {
			return err
		}
	}

	// Coleta de Inteligência (OSINT)
	if targetUrl != "" {
		if verbose {
			fmt.Fprintf(os.Stderr, "🔍 Iniciando coleta de inteligência em %s...\n", targetUrl)
		}
		scrapedWords, err := scraper.ScrapeUrl(targetUrl, scraper.DefaultConfig())
		if err != nil {
			fmt.Fprintf(os.Stderr, "⚠️ Aviso: Falha ao extrair dados da URL: %v\n", err)
		} else {
			commonWords = append(commonWords, scrapedWords...)
			if verbose {
				fmt.Fprintf(os.Stderr, "✅ Encontradas %d palavras relevantes na URL\n", len(scrapedWords))
			}
		}
	}

	// Carregar configurações de internacionalização
	locale, err := i18n.LoadLocale(lang)
	if err != nil {
		return fmt.Errorf("erro ao carregar idioma: %v", err)
	}
	
	// Adicionar palavras comuns do idioma selecionado
	commonWords = append(commonWords, locale.CommonWords...)

	// Verificar se pelo menos uma fonte foi fornecida
	if companyName == "" && len(commonWords) == 0 {
		return fmt.Errorf("deve fornecer pelo menos um nome de corporação (--company) ou palavras comuns (--words)")
	}

	// Carregar templates do arquivo se fornecido
	if patternsFile != "" {
		filePatterns, err := patterns.LoadPatternsFromFile(patternsFile)
		if err != nil {
			// Se o arquivo não for encontrado e não for o padrão (já tratado), erro
			if !os.IsNotExist(err) || patternsFile != "patterns.yaml" {
				return fmt.Errorf("erro ao carregar patterns do arquivo %s: %v", patternsFile, err)
			}
		}
		if len(filePatterns) > 0 {
			templates = append(templates, filePatterns...)
		}
	}

	// --- VERIFICAÇÃO DE ARQUIVO DE SAÍDA (ANTES DE GERAR) ---

	// Definir nome padrão do arquivo de saída se não especificado
	if outputFile == "" {
		if companyName != "" {
			outputFile = fmt.Sprintf("wordlist_%s.txt", strings.ReplaceAll(strings.ToLower(companyName), " ", "_"))
		} else {
			outputFile = "asirigen_wordlist.txt"
		}
		if verbose {
			fmt.Fprintf(os.Stderr, "ℹ️  Nenhum arquivo de saída especificado. Usando padrão: %s\n", outputFile)
		}
	}

	// Carregar senhas existentes se o arquivo já existir
	existingPasswords := make(map[string]bool)
	if _, err := os.Stat(outputFile); err == nil {
		fmt.Fprintf(os.Stderr, "⚠️  O arquivo '%s' já existe. Lendo senhas existentes...\n", outputFile)
		file, err := os.Open(outputFile)
		if err != nil {
			return fmt.Errorf("erro ao abrir arquivo existente: %v", err)
		}
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			existingPasswords[scanner.Text()] = true
		}
		file.Close()
		fmt.Fprintf(os.Stderr, "📊 %d senhas carregadas do arquivo existente.\n", len(existingPasswords))
	}

	// Se --all for usado, fazer estimativa e pedir confirmação
	if allCombinations {
		// Estimativa grosseira: (variações empresa + (palavras * variações)) * (anos + numeros + leet)
		// Esta é uma simplificação, a lógica real é complexa. Vamos fazer uma contagem simulada rápida ou apenas alertar.
		// Como simular é custoso, vamos alertar o usuário sobre o potencial tamanho.
		
		fmt.Fprintf(os.Stderr, "\n⚠️  ATENÇÃO: O modo --all irá gerar TODAS as combinações possíveis.\n")
		fmt.Fprintf(os.Stderr, "Isso pode resultar em arquivos muito grandes (GBs) e levar muito tempo.\n")
		fmt.Fprintf(os.Stderr, "Deseja continuar? [S/n]: ")
		
		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToLower(response))

		if response != "" && response != "s" && response != "sim" && response != "y" && response != "yes" {
			fmt.Println("❌ Operação cancelada pelo usuário.")
			return nil
		}
		
		startTime := time.Now()
		defer func() {
			duration := time.Since(startTime)
			fmt.Fprintf(os.Stderr, "⏱️  Tempo total de geração: %s\n", duration)
		}()
	}

	// Output writer (com diff para novas senhas)
	var writer func(string) error
	var newPasswordsFile *os.File
	var newPasswordsCount int

	// Se houver senhas existentes, criar arquivo separado para as novas
	if len(existingPasswords) > 0 {
		newFilename := strings.TrimSuffix(outputFile, ".txt") + "_new.txt"
		fmt.Fprintf(os.Stderr, "🆕 Novas senhas serão salvas em: %s\n", newFilename)
		
		f, err := os.Create(newFilename)
		if err != nil {
			return fmt.Errorf("erro ao criar arquivo de novas senhas: %v", err)
		}
		newPasswordsFile = f
		defer newPasswordsFile.Close()
	}

	// Abrir arquivo principal em modo append
	file, err := os.OpenFile(outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("erro ao abrir arquivo %s: %v", outputFile, err)
	}
	defer file.Close()
	
	writer = func(word string) error {
		// Se a senha já existe, ignorar (não duplicar no arquivo principal)
		if existingPasswords[word] {
			return nil
		}
		
		// Se é nova, adicionar ao arquivo principal
		if _, err := file.WriteString(word + "\n"); err != nil {
			return err
		}
		
		// E também ao arquivo de novas senhas, se existir
		if newPasswordsFile != nil {
			if _, err := newPasswordsFile.WriteString(word + "\n"); err != nil {
				return err
			}
			newPasswordsCount++
		}
		
		// Atualizar mapa para evitar duplicatas na mesma execução
		existingPasswords[word] = true
		return nil
	}

	// --- INICIAR GERAÇÃO ---

	// Criar gerador
	gen := generator.NewGenerator(generator.Config{
		MinLength:       minLength,
		MaxLength:       maxLength,
		MinYear:         minYear,
		MaxYear:         maxYear,
		LeetMode:        leetMode,
		Verbose:         verbose,
		MaxPasswords:    maxPasswords,
		Templates:       templates,
		Locale:          locale,
		AllCombinations: allCombinations,
	})

	// Canal para receber senhas
	out := make(chan string, 10000) // Buffer maior para --all
	errChan := make(chan error, 1)

	// Iniciar geração em goroutine
	go func() {
		if err := gen.GenerateIntelligentWordlist(companyName, commonWords, out); err != nil {
			errChan <- err
		}
		close(errChan)
	}()

	// Consumir canal
	count := 0
	totalNew := 0
	
	// Spinner simples
	spinChars := []rune{'|', '/', '-', '\\'}
	spinIdx := 0
	lastUpdate := time.Now()

	for word := range out {
		isNew := !existingPasswords[word]
		
		if err := writer(word); err != nil {
			return fmt.Errorf("erro ao escrever saída: %v", err)
		}
		
		count++
		if isNew {
			totalNew++
		}

		// Atualizar progresso a cada 100ms
		if verbose || time.Since(lastUpdate) > 100*time.Millisecond {
			fmt.Fprintf(os.Stderr, "\r%c Gerando: %d senhas processadas", spinChars[spinIdx], count)
			spinIdx = (spinIdx + 1) % len(spinChars)
			lastUpdate = time.Now()
		}
	}
	fmt.Fprintf(os.Stderr, "\r✅ Geração concluída!                  \n")

	// Verificar erros do gerador
	if err := <-errChan; err != nil {
		return err
	}

	if verbose || outputFile != "" {
		fmt.Fprintf(os.Stderr, "📄 Total processado: %d\n", count)
		if newPasswordsFile != nil {
			fmt.Fprintf(os.Stderr, "🆕 Novas senhas encontradas: %d (salvas em %s)\n", newPasswordsCount, newPasswordsFile.Name())
		}
		fmt.Fprintf(os.Stderr, "💾 Wordlist final atualizada em: %s\n", outputFile)
	}

	return nil
}
