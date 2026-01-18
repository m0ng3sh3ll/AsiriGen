package generator

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// TemplateData contém os dados disponíveis para expansão do template
type TemplateData struct {
	Company     string
	CommonWords []string
	Years       []int
}

// Lista de tokens válidos para validação
var validTokens = map[string]bool{
	"{company}": true,
	"{word}":    true,
	"{year}":    true,
	"{sep}":     true,
	"{num}":     true,
	"{special}": true,
	"{leet}":    true,
	// Novos tokens de data
	"{day}":         true,
	"{month}":       true,
	"{month_num}":   true,
	"{month_short}": true,
	"{season}":      true,
	"{weekday}":     true,
	"{date}":        true,
	// Novos tokens de padrões
	"{reverse}":  true,
	"{keyboard}": true,
}

// ExpandTemplate expande um template em uma lista de senhas
// Suporta tokens: {company}, {word}, {year}, {sep}, {num}, {special}
func (g *Generator) ExpandTemplate(template string, data TemplateData) []string {
	// Validação de tokens desconhecidos
	g.validateTemplateTokens(template)

	// Começamos com o template original
	results := []string{template}

	// 1. Expandir {company}
	// Se o template tem {company}, substituímos por todas as variações da empresa
	if strings.Contains(template, "{company}") {
		var nextResults []string
		companyVariations := g.generateCompanyVariations(data.Company)
		for _, res := range results {
			for _, variation := range companyVariations {
				nextResults = append(nextResults, strings.ReplaceAll(res, "{company}", variation))
			}
		}
		results = nextResults
	}

	// 2. Expandir {word}
	// Se o template tem {word}, substituímos por todas as palavras comuns
	if strings.Contains(template, "{word}") && len(data.CommonWords) > 0 {
		var nextResults []string
		for _, res := range results {
			for _, word := range data.CommonWords {
				// Também podemos aplicar variações na palavra (Title, Upper)
				wordVars := []string{word, strings.Title(word), strings.ToUpper(word)}
				for _, wv := range wordVars {
					nextResults = append(nextResults, strings.ReplaceAll(res, "{word}", wv))
				}
			}
		}
		results = nextResults
	}

	// 3. Expandir {year}
	if strings.Contains(template, "{year}") {
		var nextResults []string
		for _, res := range results {
			for _, year := range data.Years {
				// Formatos de ano: 2024, 24
				yearFull := fmt.Sprintf("%d", year)
				yearShort := fmt.Sprintf("%02d", year%100)
				nextResults = append(nextResults, strings.ReplaceAll(res, "{year}", yearFull))
				nextResults = append(nextResults, strings.ReplaceAll(res, "{year}", yearShort))
			}
		}
		results = nextResults
	}

	// 4. Expandir {sep}
	if strings.Contains(template, "{sep}") {
		separators := []string{"", "_", "-", ".", "@", "#", "$"}
		var nextResults []string
		for _, res := range results {
			for _, sep := range separators {
				nextResults = append(nextResults, strings.ReplaceAll(res, "{sep}", sep))
			}
		}
		results = nextResults
	}

	// 5. Expandir {num}
	if strings.Contains(template, "{num}") {
		numbers := []string{"1", "12", "123", "1234", "12345", "123456", "0", "01", "007"}
		var nextResults []string
		for _, res := range results {
			for _, num := range numbers {
				nextResults = append(nextResults, strings.ReplaceAll(res, "{num}", num))
			}
		}
		results = nextResults
	}

	// 6. Expandir {special}
	if strings.Contains(template, "{special}") {
		specials := []string{"!", "@", "#", "$", "%", "&", "*", "?"}
		var nextResults []string
		for _, res := range results {
			for _, special := range specials {
				nextResults = append(nextResults, strings.ReplaceAll(res, "{special}", special))
			}
		}
		results = nextResults
	}

	// 7. Expandir {leet}
	if strings.Contains(template, "{leet}") {
		var nextResults []string

		// Coletar todas as palavras base para aplicar leet
		var baseWords []string
		if data.Company != "" {
			baseWords = append(baseWords, data.Company)
		}
		baseWords = append(baseWords, data.CommonWords...)

		// Gerar variações leet para cada palavra base
		var leetVariations []string
		for _, word := range baseWords {
			leetVariations = append(leetVariations, g.generateLeetWord(word)...)
		}

		// Substituir no template
		for _, res := range results {
			for _, leet := range leetVariations {
				nextResults = append(nextResults, strings.ReplaceAll(res, "{leet}", leet))
			}
		}
		results = nextResults
	}

	// 8. Expandir {day} - Dias do mês (01-31)
	if strings.Contains(template, "{day}") {
		var nextResults []string
		for _, res := range results {
			for day := 1; day <= 31; day++ {
				nextResults = append(nextResults, strings.ReplaceAll(res, "{day}", fmt.Sprintf("%02d", day)))
			}
		}
		results = nextResults
	}

	// 9. Expandir {month} - Meses por extenso
	if strings.Contains(template, "{month}") && g.config.Locale != nil && len(g.config.Locale.Months) > 0 {
		var nextResults []string
		for _, res := range results {
			for _, month := range g.config.Locale.Months {
				nextResults = append(nextResults, strings.ReplaceAll(res, "{month}", month))
				nextResults = append(nextResults, strings.ReplaceAll(res, "{month}", strings.Title(month)))
			}
		}
		results = nextResults
	}

	// 10. Expandir {month_num} - Meses numéricos (01-12)
	if strings.Contains(template, "{month_num}") {
		var nextResults []string
		for _, res := range results {
			for month := 1; month <= 12; month++ {
				nextResults = append(nextResults, strings.ReplaceAll(res, "{month_num}", fmt.Sprintf("%02d", month)))
			}
		}
		results = nextResults
	}

	// 11. Expandir {month_short} - Abreviações de meses
	if strings.Contains(template, "{month_short}") && g.config.Locale != nil && len(g.config.Locale.MonthsShort) > 0 {
		var nextResults []string
		for _, res := range results {
			for _, month := range g.config.Locale.MonthsShort {
				nextResults = append(nextResults, strings.ReplaceAll(res, "{month_short}", month))
				nextResults = append(nextResults, strings.ReplaceAll(res, "{month_short}", strings.Title(month)))
			}
		}
		results = nextResults
	}

	// 12. Expandir {season} - Estações do ano
	if strings.Contains(template, "{season}") && g.config.Locale != nil && len(g.config.Locale.Seasons) > 0 {
		var nextResults []string
		for _, res := range results {
			for _, season := range g.config.Locale.Seasons {
				nextResults = append(nextResults, strings.ReplaceAll(res, "{season}", season))
				nextResults = append(nextResults, strings.ReplaceAll(res, "{season}", strings.Title(season)))
			}
		}
		results = nextResults
	}

	// 13. Expandir {weekday} - Dias da semana
	if strings.Contains(template, "{weekday}") && g.config.Locale != nil && len(g.config.Locale.Weekdays) > 0 {
		var nextResults []string
		for _, res := range results {
			for _, weekday := range g.config.Locale.Weekdays {
				nextResults = append(nextResults, strings.ReplaceAll(res, "{weekday}", weekday))
				nextResults = append(nextResults, strings.ReplaceAll(res, "{weekday}", strings.Title(weekday)))
			}
		}
		results = nextResults
	}

	// 14. Expandir {date} - Formatos de data REALISTAS (sem barras/hífens)
	if strings.Contains(template, "{date}") {
		var nextResults []string
		// Apenas formatos compactos que são realmente usados em senhas
		dateFormats := []string{
			"01012024", // ddmmyyyy
			"010124",   // ddmmyy
			"20240101", // yyyymmdd
			"240101",   // yymmdd
		}
		for _, res := range results {
			for _, dateFormat := range dateFormats {
				nextResults = append(nextResults, strings.ReplaceAll(res, "{date}", dateFormat))
			}
		}
		results = nextResults
	}

	// 15. Expandir {reverse} - Palavra invertida
	if strings.Contains(template, "{reverse}") {
		var nextResults []string
		var baseWords []string
		if data.Company != "" {
			baseWords = append(baseWords, data.Company)
		}
		baseWords = append(baseWords, data.CommonWords...)

		for _, res := range results {
			for _, word := range baseWords {
				reversed := reverseString(word)
				nextResults = append(nextResults, strings.ReplaceAll(res, "{reverse}", reversed))
				nextResults = append(nextResults, strings.ReplaceAll(res, "{reverse}", strings.Title(reversed)))
			}
		}
		results = nextResults
	}

	// 16. Expandir {keyboard} - Padrões de teclado contextualizados
	if strings.Contains(template, "{keyboard}") {
		var nextResults []string
		keyboardPatterns := []string{"qwe", "asd", "zxc", "qaz", "wsx", "edc"}
		for _, res := range results {
			for _, pattern := range keyboardPatterns {
				nextResults = append(nextResults, strings.ReplaceAll(res, "{keyboard}", pattern))
			}
		}
		results = nextResults
	}

	return results
}

// reverseString inverte uma string
func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// validateTemplateTokens verifica se existem tokens desconhecidos no template
func (g *Generator) validateTemplateTokens(template string) {
	// Regex para encontrar qualquer coisa entre chaves: {algumacoisa}
	re := regexp.MustCompile(`\{[a-zA-Z0-9_]+\}`)
	matches := re.FindAllString(template, -1)

	for _, match := range matches {
		if !validTokens[match] {
			// Se encontrou um token que não está no mapa de válidos, avisa o usuário
			fmt.Fprintf(os.Stderr, "⚠️  AVISO: Token desconhecido encontrado no template '%s': %s\n", template, match)
			fmt.Fprintf(os.Stderr, "   Tokens válidos são:\n")
			fmt.Fprintf(os.Stderr, "   - Básicos: {company}, {word}, {year}, {sep}, {num}, {special}, {leet}\n")
			fmt.Fprintf(os.Stderr, "   - Data: {day}, {month}, {month_num}, {month_short}, {season}, {weekday}, {date}\n")
			fmt.Fprintf(os.Stderr, "   - Padrões: {reverse}, {keyboard}\n")
		}
	}
}
