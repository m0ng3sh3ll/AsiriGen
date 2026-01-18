package generator

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"asirigen/internal/i18n"
	"asirigen/internal/patterns"
)

// Config contém a configuração do gerador
type Config struct {
	MinLength       int
	MaxLength       int
	MinYear         int
	MaxYear         int
	LeetMode        bool
	Verbose         bool
	MaxPasswords    int
	Templates       []string
	Locale          *i18n.LocaleData
	AllCombinations bool
}

// Generator é a estrutura principal do gerador
type Generator struct {
	config         Config
	leetMap        map[rune][]string
	commonSuffixes []string
	commonPrefixes []string
	years          []int
	patterns       []patterns.PasswordPattern
	corporateCtx   *patterns.CorporatePatterns
}

// NewGenerator cria uma nova instância do gerador
func NewGenerator(config Config) *Generator {
	gen := &Generator{
		config: config,
		leetMap: map[rune][]string{
			'a': {"@", "4", "^", "/-\\", "A"},
			'b': {"8", "|3", "13", "B"},
			'c': {"(", "[", "<", "{", "C"},
			'd': {"|)", "|}", "D"},
			'e': {"3", "&", "€", "E"},
			'f': {"|=", "ph", "F"},
			'g': {"9", "6", "&", "G"},
			'h': {"#", "|-|", "H"},
			'i': {"1", "!", "|", "][", "I"},
			'j': {"_|", "J"},
			'k': {"|<", "1<", "K"},
			'l': {"1", "|", "7", "|_", "L"},
			'm': {"/\\/", "|\\/|", "M"},
			'n': {"/\\/", "|\\|", "N"},
			'o': {"0", "*", "()", "O"},
			'p': {"|*", "P"},
			'q': {"0_", "9", "Q"},
			'r': {"|2", "R"},
			's': {"5", "$", "z", "§", "S"},
			't': {"7", "+", "†", "T"},
			'u': {"|_|", "v", "U"},
			'v': {"\\/", "V"},
			'w': {"\\/\\/", "vv", "W"},
			'x': {"%", "><", "X"},
			'y': {"j", "`/", "Y"},
			'z': {"2", "%", "7_", "Z"},
		},
		commonSuffixes: []string{
			"123", "456", "789", "000", "111", "222", "333",
			"!@#", "!@#$", "!@#$%", "!@#$%^", "!@#$%^&",
			"2020", "2021", "2022", "2023", "2024", "2025",
			"01", "02", "03", "04", "05", "06", "07", "08", "09", "10",
			"11", "12", "13", "14", "15", "16", "17", "18", "19", "20",
		},
		commonPrefixes: []string{
			"", "!", "@", "#", "$", "%", "^", "&", "*",
		},
		patterns: patterns.GetRealisticPatterns(),
	}

	// Aplicar configurações do locale se disponível
	if config.Locale != nil {
		if len(config.Locale.CommonSuffixes) > 0 {
			gen.commonSuffixes = config.Locale.CommonSuffixes
		}
	}

	// Gerar lista de anos
	gen.generateYears()

	// Configurar valores padrão se não especificados
	if config.MaxPasswords == 0 && !config.AllCombinations {
		gen.config.MaxPasswords = 100000
	}

	return gen
}

// generateYears gera a lista de anos baseada na configuração
func (g *Generator) generateYears() {
	for year := g.config.MinYear; year <= g.config.MaxYear; year++ {
		g.years = append(g.years, year)
	}
}

// GenerateFromCompany gera wordlist baseada no nome da corporação
func (g *Generator) GenerateFromCompany(company string) ([]string, error) {
	if g.config.Verbose {
		fmt.Printf("Gerando wordlist para corporação: %s\n", company)
	}

	var wordlist []string
	company = strings.ToLower(company)

	// Variações básicas do nome da corporação
	variations := g.generateCompanyVariations(company)
	wordlist = append(wordlist, variations...)

	// Adicionar sufixos e prefixos
	wordlist = g.addSuffixesAndPrefixes(wordlist)

	// Aplicar leetspeak se habilitado
	if g.config.LeetMode {
		leetVariations := g.generateLeetVariations(wordlist)
		wordlist = append(wordlist, leetVariations...)
	}

	// Remover duplicatas
	wordlist = g.removeDuplicates(wordlist)

	if g.config.Verbose {
		fmt.Printf("Geradas %d variações para %s\n", len(wordlist), company)
	}

	return wordlist, nil
}

// GenerateFromWords gera wordlist baseada em palavras comuns
func (g *Generator) GenerateFromWords(words []string) ([]string, error) {
	if g.config.Verbose {
		fmt.Printf("Gerando wordlist para %d palavras comuns\n", len(words))
	}

	var wordlist []string

	for _, word := range words {
		word = strings.ToLower(strings.TrimSpace(word))
		if word == "" {
			continue
		}

		// Variações básicas
		variations := g.generateWordVariations(word)
		wordlist = append(wordlist, variations...)

		// Adicionar sufixos e prefixos
		variations = g.addSuffixesAndPrefixes(variations)
		wordlist = append(wordlist, variations...)

		// Aplicar leetspeak se habilitado
		if g.config.LeetMode {
			leetVariations := g.generateLeetVariations(variations)
			wordlist = append(wordlist, leetVariations...)
		}
	}

	// Remover duplicatas
	wordlist = g.removeDuplicates(wordlist)

	if g.config.Verbose {
		fmt.Printf("Geradas %d variações para palavras comuns\n", len(wordlist))
	}

	return wordlist, nil
}

// GenerateCombined gera wordlist combinando corporação e palavras comuns
func (g *Generator) GenerateCombined(company string, words []string) ([]string, error) {
	if g.config.Verbose {
		fmt.Printf("Gerando wordlist combinada para %s e %d palavras\n", company, len(words))
	}

	var wordlist []string

	// Gerar wordlist da corporação
	companyWords, err := g.GenerateFromCompany(company)
	if err != nil {
		return nil, err
	}
	wordlist = append(wordlist, companyWords...)

	// Gerar wordlist das palavras comuns
	commonWords, err := g.GenerateFromWords(words)
	if err != nil {
		return nil, err
	}
	wordlist = append(wordlist, commonWords...)

	// Combinar corporação com palavras comuns
	combinations := g.generateCombinations(company, words)
	wordlist = append(wordlist, combinations...)

	// Remover duplicatas
	wordlist = g.removeDuplicates(wordlist)

	if g.config.Verbose {
		fmt.Printf("Geradas %d variações combinadas\n", len(wordlist))
	}

	return wordlist, nil
}

// generateCompanyVariations gera variações do nome da corporação
func (g *Generator) generateCompanyVariations(company string) []string {
	var variations []string

	// Original
	variations = append(variations, company)

	// Capitalizar primeira letra
	variations = append(variations, strings.Title(company))

	// Tudo maiúsculo
	variations = append(variations, strings.ToUpper(company))

	// Variações para nomes com espaços
	if strings.Contains(company, " ") {
		spaceVariations := g.generateSpaceVariations(company)
		variations = append(variations, spaceVariations...)
	}

	// Abreviações comuns
	abbreviations := g.generateAbbreviations(company)
	variations = append(variations, abbreviations...)

	// Variações com números
	numberVariations := g.generateNumberVariations(company)
	variations = append(variations, numberVariations...)

	return variations
}

// generateWordVariations gera variações de uma palavra
func (g *Generator) generateWordVariations(word string) []string {
	var variations []string

	// Original
	variations = append(variations, word)

	// Capitalizar primeira letra
	variations = append(variations, strings.Title(word))

	// Tudo maiúsculo
	variations = append(variations, strings.ToUpper(word))

	// Variações com números
	numberVariations := g.generateNumberVariations(word)
	variations = append(variations, numberVariations...)

	return variations
}

// generateAbbreviations gera abreviações comuns
func (g *Generator) generateAbbreviations(word string) []string {
	var abbreviations []string

	// Primeiras letras de cada palavra
	words := strings.Fields(word)
	if len(words) > 1 {
		var abbr strings.Builder
		for _, w := range words {
			if len(w) > 0 {
				abbr.WriteRune(rune(w[0]))
			}
		}
		abbreviations = append(abbreviations, abbr.String())
		abbreviations = append(abbreviations, strings.ToUpper(abbr.String()))
	}

	// Primeiras 3-4 letras
	if len(word) >= 3 {
		abbreviations = append(abbreviations, word[:3])
		abbreviations = append(abbreviations, strings.ToUpper(word[:3]))
	}
	if len(word) >= 4 {
		abbreviations = append(abbreviations, word[:4])
		abbreviations = append(abbreviations, strings.ToUpper(word[:4]))
	}

	return abbreviations
}

// generateNumberVariations gera variações com números
func (g *Generator) generateNumberVariations(word string) []string {
	var variations []string

	// Adicionar anos
	for _, year := range g.years {
		variations = append(variations, word+fmt.Sprintf("%d", year))
		variations = append(variations, word+fmt.Sprintf("%02d", year%100))
	}

	// Adicionar números sequenciais
	for i := 1; i <= 20; i++ {
		variations = append(variations, word+fmt.Sprintf("%d", i))
		variations = append(variations, word+fmt.Sprintf("%02d", i))
		variations = append(variations, word+fmt.Sprintf("%03d", i))
	}

	return variations
}

// addSuffixesAndPrefixes adiciona sufixos e prefixos comuns
func (g *Generator) addSuffixesAndPrefixes(words []string) []string {
	var result []string

	for _, word := range words {
		// Adicionar sufixos
		for _, suffix := range g.commonSuffixes {
			result = append(result, word+suffix)
		}

		// Adicionar prefixos
		for _, prefix := range g.commonPrefixes {
			if prefix != "" {
				result = append(result, prefix+word)
			}
		}

		// Adicionar prefixos e sufixos juntos
		for _, prefix := range g.commonPrefixes {
			for _, suffix := range g.commonSuffixes {
				if prefix != "" {
					result = append(result, prefix+word+suffix)
				}
			}
		}
	}

	return result
}

// generateLeetVariations gera variações leetspeak
func (g *Generator) generateLeetVariations(words []string) []string {
	var result []string

	for _, word := range words {
		leetWords := g.generateLeetWord(word)
		result = append(result, leetWords...)
	}

	return result
}

// generateLeetWord gera variações leetspeak de uma palavra
func (g *Generator) generateLeetWord(word string) []string {
	var result []string
	word = strings.ToLower(word)

	// Gerar todas as combinações possíveis de substituições
	combinations := g.generateLeetCombinations(word)
	result = append(result, combinations...)

	return result
}

// generateLeetCombinations gera combinações leetspeak recursivamente
func (g *Generator) generateLeetCombinations(word string) []string {
	var result []string
	result = append(result, word) // Adicionar original

	// Encontrar posições que podem ser substituídas
	var positions []int
	for i, char := range word {
		if _, exists := g.leetMap[char]; exists {
			positions = append(positions, i)
		}
	}

	// Gerar combinações (limitado para evitar explosão)
	maxCombinations := 1000
	if len(positions) > 10 {
		positions = positions[:10] // Limitar a 10 posições
	}

	g.generateLeetRecursive(word, positions, 0, &result, maxCombinations)

	return result
}

// generateLeetRecursive gera combinações leetspeak recursivamente
func (g *Generator) generateLeetRecursive(word string, positions []int, posIndex int, result *[]string, maxCombinations int) {
	if len(*result) >= maxCombinations || posIndex >= len(positions) {
		return
	}

	pos := positions[posIndex]
	char := rune(word[pos])
	replacements := g.leetMap[char]

	for _, replacement := range replacements {
		newWord := word[:pos] + replacement + word[pos+1:]
		*result = append(*result, newWord)

		// Recursão para próxima posição
		g.generateLeetRecursive(newWord, positions, posIndex+1, result, maxCombinations)
	}
}

// generateCombinations gera combinações entre corporação e palavras
func (g *Generator) generateCombinations(company string, words []string) []string {
	var result []string
	company = strings.ToLower(company)

	for _, word := range words {
		word = strings.ToLower(strings.TrimSpace(word))
		if word == "" {
			continue
		}

		// Combinações com separadores
		separators := []string{"", "_", "-", ".", "@", "#", "$", "%", "&", "*"}
		for _, sep := range separators {
			result = append(result, company+sep+word)
			result = append(result, word+sep+company)
		}

		// Combinações com números
		for _, year := range g.years {
			result = append(result, company+fmt.Sprintf("%d", year)+word)
			result = append(result, word+fmt.Sprintf("%d", year)+company)
		}
	}

	return result
}

// FilterByLength filtra palavras por tamanho
func (g *Generator) FilterByLength(words []string) []string {
	var result []string

	for _, word := range words {
		if len(word) >= g.config.MinLength && len(word) <= g.config.MaxLength {
			result = append(result, word)
		}
	}

	return result
}

// removeDuplicates remove duplicatas da lista
func (g *Generator) removeDuplicates(words []string) []string {
	seen := make(map[string]bool)
	var result []string

	for _, word := range words {
		if !seen[word] {
			seen[word] = true
			result = append(result, word)
		}
	}

	return result
}

// GenerateIntelligentWordlist gera wordlist inteligente baseada em análise realística (Streaming)
func (g *Generator) GenerateIntelligentWordlist(company string, commonWords []string, out chan<- string) error {
	defer close(out)

	if g.config.Verbose {
		fmt.Printf("🔍 Gerando wordlist inteligente para %s\n", company)
	}

	fmt.Print("🔍 Configurando contexto corporativo... ")
	g.corporateCtx = patterns.GetCorporateContext(company)
	fmt.Println("✅")

	// Buffer temporário para randomização
	var tempBuffer []string

	// Helper para processar lista de palavras
	processList := func(words []string) {
		tempBuffer = append(tempBuffer, words...)
	}

	// REMOVIDO: getHighPriorityPasswords() - senhas genéricas já cobertas por rockyou.txt
	// fmt.Print("🎯 Gerando senhas de alta prioridade... ")
	// processList(g.getHighPriorityPasswords())
	// fmt.Println("✅")

	if company != "" {
		fmt.Print("🏢 Gerando senhas baseadas na empresa... ")
		processList(g.generateCompanyBasedPasswords(company))
		fmt.Println("✅")
	}

	if len(commonWords) > 0 {
		fmt.Print("📝 Gerando senhas baseadas em palavras... ")
		processList(g.generateWordBasedPasswords(commonWords))
		fmt.Println("✅")
	}

	fmt.Print("📅 Gerando padrões sazonais... ")
	processList(g.generateSeasonalPasswords())
	fmt.Println("✅")

	// REMOVIDO: generateKeyboardPasswords() - padrões genéricos já cobertos por rockyou.txt
	// fmt.Print("⌨️  Gerando padrões de teclado... ")
	// processList(g.generateKeyboardPasswords())
	// fmt.Println("✅")

	// Templates Personalizados
	if len(g.config.Templates) > 0 {
		fmt.Print("🎨 Gerando templates personalizados... ")
		data := TemplateData{
			Company:     company,
			CommonWords: commonWords,
			Years:       g.years,
		}

		for _, tmpl := range g.config.Templates {
			processList(g.ExpandTemplate(tmpl, data))
		}
		fmt.Println("✅")
	}

	// Leetspeak - Aplicar em palavras base E gerar variações completas
	if g.config.LeetMode {
		fmt.Print("🔤 Aplicando leetspeak... ")
		var baseWords []string
		if company != "" {
			baseWords = append(baseWords, g.generateCompanyVariations(company)...)
		}
		baseWords = append(baseWords, commonWords...)

		// Gerar variações leet
		var leetWords []string
		for _, word := range baseWords {
			leetWords = append(leetWords, g.generateLeetWord(word)...)
		}

		// Adicionar variações leet básicas
		processList(leetWords)

		// NOVO: Gerar variações completas das palavras leet (números + sufixos)
		processList(g.generateWordBasedPasswords(leetWords))

		fmt.Println("✅")
	}

	// Embaralhar se não for modo "All Combinations" (que deve ser determinístico)
	if !g.config.AllCombinations {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(tempBuffer), func(i, j int) {
			tempBuffer[i], tempBuffer[j] = tempBuffer[j], tempBuffer[i]
		})
	}

	// Emitir para o canal
	seen := make(map[string]bool)
	count := 0
	for _, word := range tempBuffer {
		// Filtrar por tamanho
		if len(word) < g.config.MinLength || len(word) > g.config.MaxLength {
			continue
		}

		// Remover duplicatas
		if seen[word] {
			continue
		}
		seen[word] = true

		out <- word
		count++

		// Verificar limite (apenas se não for AllCombinations)
		if !g.config.AllCombinations && g.config.MaxPasswords > 0 && count >= g.config.MaxPasswords {
			break
		}
	}

	if g.config.Verbose {
		fmt.Printf("🎉 Geradas %d senhas inteligentes\n", count)
	}

	return nil
}

// getHighPriorityPasswords retorna senhas de alta prioridade
func (g *Generator) getHighPriorityPasswords() []string {
	return []string{
		"password", "123456", "admin", "user", "test", "guest", "root",
		"welcome", "login", "access", "system", "server", "database",
		"password123", "admin123", "user123", "test123", "guest123",
		"123456789", "qwerty", "abc123", "letmein", "monkey", "dragon",
		"master", "hello", "welcome123", "login123", "access123",
		"default", "temp", "backup", "demo", "sample", "example",
	}
}

// generateCompanyBasedPasswords gera senhas baseadas na empresa
func (g *Generator) generateCompanyBasedPasswords(company string) []string {
	var passwords []string
	company = strings.ToLower(company)

	// Gerar permutações para nomes com espaços
	companyVariations := g.generateCompanyVariations(company)

	// Padrões mais comuns para empresas
	patterns := []string{
		company,                  // microsoft
		strings.Title(company),   // Microsoft
		strings.ToUpper(company), // MICROSOFT
	}

	// Adicionar variações de espaços
	patterns = append(patterns, companyVariations...)

	// Adicionar abreviações
	if g.corporateCtx != nil {
		patterns = append(patterns, g.corporateCtx.Abbreviations...)
	}

	// Gerar variações com anos
	for _, pattern := range patterns {
		for _, year := range g.years {
			passwords = append(passwords,
				pattern+fmt.Sprintf("%d", year),
				pattern+fmt.Sprintf("%02d", year%100),
				pattern+"_"+fmt.Sprintf("%02d", year%100),
				pattern+"@"+fmt.Sprintf("%02d", year%100),
			)
		}
	}

	// Gerar variações com números comuns
	commonNumbers := []string{"123", "456", "789", "000", "111", "222", "333"}
	for _, pattern := range patterns {
		for _, num := range commonNumbers {
			passwords = append(passwords,
				pattern+num,
				num+pattern,
				pattern+"_"+num,
				pattern+"@"+num,
			)
		}
	}

	return passwords
}

// generateSpaceVariations gera variações de nomes com espaços
func (g *Generator) generateSpaceVariations(company string) []string {
	var variations []string

	// Variações comuns para nomes com espaços
	variations = append(variations,
		strings.ReplaceAll(company, " ", ""),  // ferreiracosta
		strings.ReplaceAll(company, " ", "_"), // ferreira_costa
		strings.ReplaceAll(company, " ", "-"), // ferreira-costa
		strings.ReplaceAll(company, " ", "."), // ferreira.costa
		strings.ReplaceAll(company, " ", "@"), // ferreira@costa
	)

	// Variações com Title case
	titleCompany := strings.Title(company)
	variations = append(variations,
		strings.ReplaceAll(titleCompany, " ", ""),  // FerreiraCosta
		strings.ReplaceAll(titleCompany, " ", "_"), // Ferreira_Costa
		strings.ReplaceAll(titleCompany, " ", "-"), // Ferreira-Costa
		strings.ReplaceAll(titleCompany, " ", "."), // Ferreira.Costa
		strings.ReplaceAll(titleCompany, " ", "@"), // Ferreira@Costa
	)

	// Variações com UPPER case
	upperCompany := strings.ToUpper(company)
	variations = append(variations,
		strings.ReplaceAll(upperCompany, " ", ""),  // FERREIRACOSTA
		strings.ReplaceAll(upperCompany, " ", "_"), // FERREIRA_COSTA
		strings.ReplaceAll(upperCompany, " ", "-"), // FERREIRA-COSTA
		strings.ReplaceAll(upperCompany, " ", "."), // FERREIRA.COSTA
		strings.ReplaceAll(upperCompany, " ", "@"), // FERREIRA@COSTA
	)

	// Variações com apenas primeira letra maiúscula
	words := strings.Fields(company)
	if len(words) > 1 {
		// Primeira palavra maiúscula, resto minúscula
		firstWord := strings.Title(words[0])
		restWords := strings.Join(words[1:], "")
		variations = append(variations, firstWord+restWords) // FerreiraCosta

		// Primeira letra de cada palavra maiúscula, sem espaços
		var initials []string
		for _, word := range words {
			if len(word) > 0 {
				initials = append(initials, strings.ToUpper(string(word[0])))
			}
		}
		variations = append(variations, strings.Join(initials, ""))  // FC
		variations = append(variations, strings.Join(initials, ".")) // F.C
		variations = append(variations, strings.Join(initials, "_")) // F_C
		variations = append(variations, strings.Join(initials, "-")) // F-C
	}

	return variations
}

// generateWordBasedPasswords gera senhas baseadas em palavras comuns
func (g *Generator) generateWordBasedPasswords(words []string) []string {
	var passwords []string

	for _, word := range words {
		word = strings.TrimSpace(word)
		if word == "" {
			continue
		}

		// Processar palavra e adicionar ao resultado
		passwords = append(passwords, g.processWord(word)...)
	}

	return passwords
}

// processWord processa uma única palavra e gera todas as variações
func (g *Generator) processWord(word string) []string {
	var passwords []string

	// Variações de case (preservando a palavra original também)
	variations := []string{
		word,                  // Original (pode ser leet: m@raul, ou normal: localbw)
		strings.ToLower(word), // minúsculo
		strings.Title(word),   // Title Case
		strings.ToUpper(word), // MAIÚSCULO
	}

	// Remover duplicatas (caso word já seja minúsculo, por exemplo)
	uniqueVariations := make(map[string]bool)
	var finalVariations []string
	for _, v := range variations {
		if !uniqueVariations[v] {
			uniqueVariations[v] = true
			finalVariations = append(finalVariations, v)
		}
	}
	variations = finalVariations

	// Adicionar variações simples (sem números)
	passwords = append(passwords, variations...)

	// Adicionar números
	for _, year := range g.years {
		for _, variation := range variations {
			passwords = append(passwords,
				variation+fmt.Sprintf("%d", year),
				variation+fmt.Sprintf("%02d", year%100),
				variation+"_"+fmt.Sprintf("%02d", year%100),
				variation+"@"+fmt.Sprintf("%02d", year%100),
			)
		}
	}

	// Adicionar números comuns
	commonNumbers := []string{"123", "456", "789", "000", "111", "222", "333", "01", "02", "03", "04", "05", "06", "07", "08", "09", "10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "20"}
	for _, variation := range variations {
		for _, num := range commonNumbers {
			passwords = append(passwords,
				variation+num,
				num+variation,
				variation+"_"+num,
				variation+"@"+num,
			)
		}
	}

	// Adicionar sufixos comuns
	commonSuffixes := []string{"!", "@", "#", "$", "%", "^", "&", "*"}
	for _, variation := range variations {
		for _, suffix := range commonSuffixes {
			passwords = append(passwords, variation+suffix)
		}
	}

	// NOVO: Adicionar combinações número + sufixo (ex: Localbw10*, M@raul17*)
	for _, variation := range variations {
		for _, num := range commonNumbers {
			for _, suffix := range commonSuffixes {
				passwords = append(passwords,
					variation+num+suffix,
					variation+"_"+num+suffix,
					variation+"@"+num+suffix,
				)
			}
		}
	}

	// NOVO: Adicionar prefixos (ex: !brasil123, !bw2017)
	commonPrefixes := []string{"!", "@", "#"}
	for _, prefix := range commonPrefixes {
		// Prefixo + palavra
		for _, variation := range variations {
			passwords = append(passwords, prefix+variation)
		}

		// Prefixo + palavra + número
		for _, variation := range variations {
			for _, num := range commonNumbers {
				passwords = append(passwords,
					prefix+variation+num,
					prefix+variation+"_"+num,
					prefix+variation+"@"+num,
				)
			}
		}

		// Prefixo + palavra + ano
		for _, variation := range variations {
			for _, year := range g.years {
				passwords = append(passwords,
					prefix+variation+fmt.Sprintf("%d", year),
					prefix+variation+fmt.Sprintf("%02d", year%100),
				)
			}
		}
	}

	return passwords
}

// generateSeasonalPasswords gera senhas baseadas em padrões sazonais
func (g *Generator) generateSeasonalPasswords() []string {
	if g.config.Locale != nil {
		return patterns.GetSeasonalPatterns(g.config.Locale.Seasons, g.config.Locale.Months)
	}
	return patterns.GetSeasonalPatterns(nil, nil)
}

// generateKeyboardPasswords gera senhas baseadas em padrões de teclado
func (g *Generator) generateKeyboardPasswords() []string {
	return []string{
		"qwerty", "asdf", "zxcv", "123456", "654321",
		"qwertyui", "asdfgh", "zxcvbn", "qazwsx",
		"1qaz2wsx", "qwerty123", "asdf1234", "zxcv1234",
		"qwertyuiop", "asdfghjkl", "zxcvbnm",
		"1q2w3e4r", "qwe123", "asd123", "zxc123",
	}
}
