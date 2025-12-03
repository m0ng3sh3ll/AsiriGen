package patterns

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// PasswordPattern representa um padrão de senha comum
type PasswordPattern struct {
	Name        string
	Probability float64
	Generator   func(base string, year int) []string
}

// CorporatePatterns contém padrões específicos de corporações
type CorporatePatterns struct {
	CompanyName    string
	Abbreviations  []string
	CommonSuffixes []string
	Industry       string
	FoundedYear    int
}

// GetRealisticPatterns retorna padrões de senha baseados em análise real
func GetRealisticPatterns() []PasswordPattern {
	return []PasswordPattern{
		{
			Name:        "Company + Year",
			Probability: 0.35,
			Generator:   generateCompanyYearPattern,
		},
		{
			Name:        "Company + Common Word",
			Probability: 0.25,
			Generator:   generateCompanyWordPattern,
		},
		{
			Name:        "Company + Numbers",
			Probability: 0.20,
			Generator:   generateCompanyNumbersPattern,
		},
		{
			Name:        "Leet Speak Variations",
			Probability: 0.15,
			Generator:   generateLeetPattern,
		},
		{
			Name:        "Keyboard Patterns",
			Probability: 0.05,
			Generator:   generateKeyboardPattern,
		},
	}
}

// generateCompanyYearPattern gera padrões empresa + ano
func generateCompanyYearPattern(base string, year int) []string {
	var patterns []string
	company := strings.ToLower(base)

	// Anos mais comuns (últimos 5 anos + anos especiais)
	years := []int{year, year - 1, year - 2, 2020, 2021, 2022, 2023, 2024, 2025}

	for _, y := range years {
		patterns = append(patterns,
			company+string(rune(y%100)),     // microsoft24
			company+string(rune(y)),         // microsoft2024
			company+"_"+string(rune(y%100)), // microsoft_24
			company+"@"+string(rune(y%100)), // microsoft@24
		)
	}

	return patterns
}

// generateCompanyWordPattern gera padrões empresa + palavra comum
func generateCompanyWordPattern(base string, year int) []string {
	var patterns []string
	company := strings.ToLower(base)

	// Palavras mais comuns em senhas corporativas
	commonWords := []string{
		"admin", "user", "test", "demo", "guest", "root",
		"password", "pass", "login", "access", "system",
		"server", "database", "network", "security",
		"welcome", "default", "temp", "backup", "mudar", "trocar", "password",
	}

	separators := []string{"", "_", "-", ".", "@", "#", "$"}

	for _, word := range commonWords {
		for _, sep := range separators {
			patterns = append(patterns,
				company+sep+word,
				word+sep+company,
				strings.Title(company)+sep+word,
				word+sep+strings.Title(company),
			)
		}
	}

	return patterns
}

// generateCompanyNumbersPattern gera padrões empresa + números
func generateCompanyNumbersPattern(base string, year int) []string {
	var patterns []string
	company := strings.ToLower(base)

	// Números mais comuns em senhas
	commonNumbers := []string{
		"123", "456", "789", "000", "111", "222", "333",
		"01", "02", "03", "04", "05", "06", "07", "08", "09", "10",
		"11", "12", "13", "14", "15", "16", "17", "18", "19", "20",
		"99", "88", "77", "66", "55", "44", "33", "22", "11", "00",
	}

	for _, num := range commonNumbers {
		patterns = append(patterns,
			company+num,
			num+company,
			company+"_"+num,
			company+"@"+num,
		)
	}

	return patterns
}

// generateLeetPattern gera variações leetspeak mais realistas
func generateLeetPattern(base string, year int) []string {
	var patterns []string

	// Aplicar leetspeak apenas em palavras curtas (mais realista)
	if len(base) <= 8 {
		leetMap := map[rune]string{
			'a': "@", 'e': "3", 'i': "1", 'o': "0", 's': "$", 't': "7",
		}

		company := strings.ToLower(base)
		leet := company
		for old, new := range leetMap {
			leet = strings.ReplaceAll(leet, string(old), new)
		}

		if leet != company {
			patterns = append(patterns, leet)
		}
	}

	return patterns
}

// generateKeyboardPattern gera padrões de teclado
func generateKeyboardPattern(base string, year int) []string {
	var patterns []string

	// Padrões de teclado comuns
	keyboardPatterns := []string{
		"qwerty", "asdf", "zxcv", "123456", "654321",
		"qwertyui", "asdfgh", "zxcvbn", "qazwsx",
		"1qaz2wsx", "qwerty123", "asdf1234",
	}

	for _, pattern := range keyboardPatterns {
		patterns = append(patterns, pattern)
	}

	return patterns
}

// GetCorporateContext retorna contexto corporativo para geração mais inteligente
func GetCorporateContext(companyName string) *CorporatePatterns {
	company := strings.ToLower(companyName)

	// Gerar abreviações comuns
	abbreviations := generateAbbreviations(company)

	// Sufixos baseados no tipo de empresa
	suffixes := getIndustrySuffixes(company)

	return &CorporatePatterns{
		CompanyName:    company,
		Abbreviations:  abbreviations,
		CommonSuffixes: suffixes,
		Industry:       detectIndustry(company),
		FoundedYear:    estimateFoundedYear(company),
	}
}

// generateAbbreviations gera abreviações realistas
func generateAbbreviations(company string) []string {
	var abbrs []string

	// Primeiras letras de cada palavra
	words := strings.Fields(company)
	if len(words) > 1 {
		var abbr strings.Builder
		for _, w := range words {
			if len(w) > 0 {
				abbr.WriteRune(rune(w[0]))
			}
		}
		abbrs = append(abbrs, abbr.String())
		abbrs = append(abbrs, strings.ToUpper(abbr.String()))
	}

	// Primeiras 3-4 letras
	if len(company) >= 3 {
		abbrs = append(abbrs, company[:3])
		abbrs = append(abbrs, strings.ToUpper(company[:3]))
	}
	if len(company) >= 4 {
		abbrs = append(abbrs, company[:4])
		abbrs = append(abbrs, strings.ToUpper(company[:4]))
	}

	return abbrs
}

// getIndustrySuffixes retorna sufixos baseados na indústria
func getIndustrySuffixes(company string) []string {
	// Detectar indústria baseada no nome
	if strings.Contains(company, "tech") || strings.Contains(company, "soft") {
		return []string{"dev", "admin", "user", "test", "demo", "api", "app"}
	}
	if strings.Contains(company, "bank") || strings.Contains(company, "finance") {
		return []string{"bank", "fin", "money", "cash", "account", "client"}
	}
	if strings.Contains(company, "health") || strings.Contains(company, "med") {
		return []string{"health", "med", "patient", "doctor", "nurse", "care"}
	}

	// Sufixos genéricos
	return []string{"admin", "user", "test", "demo", "guest", "root"}
}

// detectIndustry detecta a indústria da empresa
func detectIndustry(company string) string {
	company = strings.ToLower(company)

	if strings.Contains(company, "tech") || strings.Contains(company, "soft") || strings.Contains(company, "it") {
		return "technology"
	}
	if strings.Contains(company, "bank") || strings.Contains(company, "finance") || strings.Contains(company, "credit") {
		return "finance"
	}
	if strings.Contains(company, "health") || strings.Contains(company, "med") || strings.Contains(company, "hospital") {
		return "healthcare"
	}
	if strings.Contains(company, "edu") || strings.Contains(company, "school") || strings.Contains(company, "university") {
		return "education"
	}
	if strings.Contains(company, "gov") || strings.Contains(company, "government") {
		return "government"
	}

	return "general"
}

// estimateFoundedYear estima o ano de fundação (para contexto)
func estimateFoundedYear(company string) int {
	// Empresas conhecidas
	knownYears := map[string]int{
		"microsoft": 1975,
		"apple":     1976,
		"google":    1998,
		"facebook":  2004,
		"amazon":    1994,
		"tesla":     2003,
	}

	if year, exists := knownYears[strings.ToLower(company)]; exists {
		return year
	}

	// Estimativa baseada no nome (empresas mais antigas tendem a ter nomes mais simples)
	if len(company) <= 6 {
		return 1980 + rand.Intn(20) // 1980-2000
	}
	return 2000 + rand.Intn(25) // 2000-2025
}

// GetSeasonalPatterns retorna padrões sazonais
func GetSeasonalPatterns(seasons []string, months []string) []string {
	now := time.Now()
	year := now.Year()
	month := int(now.Month()) // 1-12

	var patterns []string

	// Se listas estiverem vazias, usar fallback inglês
	if len(seasons) < 4 {
		seasons = []string{"spring", "summer", "fall", "winter"}
	}
	if len(months) < 12 {
		months = []string{
			"january", "february", "march", "april", "may", "june",
			"july", "august", "september", "october", "november", "december",
		}
	}

	// Mapeamento simples de estação (Hemisfério Norte padrão, ajustável futuramente)
	// Primavera: Mar-Mai, Verão: Jun-Ago, Outono: Set-Nov, Inverno: Dez-Fev
	var currentSeason string
	var currentMonths []string

	if month >= 3 && month <= 5 {
		currentSeason = seasons[0]
		currentMonths = []string{months[2], months[3], months[4]}
	} else if month >= 6 && month <= 8 {
		currentSeason = seasons[1]
		currentMonths = []string{months[5], months[6], months[7]}
	} else if month >= 9 && month <= 11 {
		currentSeason = seasons[2]
		currentMonths = []string{months[8], months[9], months[10]}
	} else {
		currentSeason = seasons[3]
		// Dezembro, Janeiro, Fevereiro
		currentMonths = []string{months[11], months[0], months[1]}
	}

	patterns = append(patterns, currentSeason)
	patterns = append(patterns, currentMonths...)

	// Anos relevantes
	patterns = append(patterns,
		fmt.Sprintf("%d", year),
		fmt.Sprintf("%d", year%100),
		fmt.Sprintf("%d", year-1),
		fmt.Sprintf("%d", (year-1)%100),
	)

	return patterns
}
