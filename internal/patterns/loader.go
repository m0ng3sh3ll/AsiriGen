package patterns

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// PatternConfig define a estrutura do arquivo patterns.yaml
type PatternConfig struct {
	Patterns []string `yaml:"patterns"`
}

// LoadPatternsFromFile carrega padrões de um arquivo YAML
func LoadPatternsFromFile(filename string) ([]string, error) {
	if filename == "" {
		return nil, nil
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler arquivo de padrões: %v", err)
	}

	var config PatternConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("erro ao fazer parse do arquivo de padrões: %v", err)
	}

	return config.Patterns, nil
}
