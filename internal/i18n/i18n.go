package i18n

import (
	"embed"
	"encoding/json"
	"fmt"
	"strings"
)

//go:embed *.json
var localeFS embed.FS

// LocaleData contém os dados específicos de um idioma
type LocaleData struct {
	CommonWords    []string `json:"common_words"`
	CommonSuffixes []string `json:"common_suffixes"`
	Seasons        []string `json:"seasons"`
	Months         []string `json:"months"`
	MonthsShort    []string `json:"months_short"`   // Abreviações de meses (Jan, Fev, etc)
	Weekdays       []string `json:"weekdays"`       // Dias da semana
	WeekdaysShort  []string `json:"weekdays_short"` // Abreviações de dias (Seg, Ter, etc)
}

// LoadLocale carrega os dados de um idioma específico
func LoadLocale(lang string) (*LocaleData, error) {
	lang = strings.ToLower(strings.ReplaceAll(lang, "-", "_"))

	// Normalização simples
	if lang == "pt" {
		lang = "pt_br"
	}
	if lang == "en" {
		lang = "en_us"
	}

	filename := lang + ".json"

	data, err := localeFS.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("idioma não encontrado: %s (tente pt-br ou en-us)", lang)
	}

	var locale LocaleData
	if err := json.Unmarshal(data, &locale); err != nil {
		return nil, fmt.Errorf("erro ao processar arquivo de idioma: %v", err)
	}

	return &locale, nil
}
