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

	return results
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
			fmt.Fprintf(os.Stderr, "   Tokens válidos são: {company}, {word}, {year}, {sep}, {num}, {special}, {leet}\n")
		}
	}
}
