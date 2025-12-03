package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// Flags globais
	outputFile      string
	minLength       int
	maxLength       int
	minYear         int
	maxYear         int
	companyName     string
	commonWords     []string
	leetMode        bool
	verbose         bool
	maxPasswords    int
	allCombinations bool
)

// rootCmd representa o comando base
var rootCmd = &cobra.Command{
	Use:   "asirigen",
	Short: "AsiriGen - Gerador de wordlists",
	Long: `AsiriGen é uma ferramenta avançada para geração de wordlists baseadas em nomes de corporações,
palavras comuns e padrões de senha.

Funcionalidades:
- Geração inteligente baseada em nomes de corporações
- Coleta de Inteligência (OSINT) via URL
- Sistema de Templates Dinâmicos (YAML)
- Internacionalização (i18n) para PT-BR e EN-US
- Suporte a leetspeak e substituições de caracteres
- Controle granular de tamanho e anos`,
	Example: `  # Gerar wordlist básica para Microsoft
  asirigen generate --company Microsoft

  # Usar OSINT para extrair palavras de um site
  asirigen generate --company "TechCorp" --url "https://techcorp.com" --verbose

  # Gerar com templates personalizados e em Português
  asirigen generate --company "Alvo" --lang pt-br --patterns-file "patterns.yaml"

  # Combo completo: OSINT + i18n + Templates + Leet
  asirigen generate --company "Target" --url "https://target.com" --lang pt-br --leet --patterns-file "patterns.yaml" -o wordlist.txt`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Não exibir banner se estiver apenas completando shell ou pedindo versão
		if cmd.Name() != "version" && cmd.Name() != "completion" {
			PrintBanner()
		}
	},
}

func init() {
	// Flags globais
	rootCmd.PersistentFlags().StringVarP(&outputFile, "output", "o", "", "Arquivo de saída (padrão: stdout)")
	rootCmd.PersistentFlags().IntVar(&minLength, "min-length", 4, "Tamanho mínimo da senha")
	rootCmd.PersistentFlags().IntVar(&maxLength, "max-length", 16, "Tamanho máximo da senha")
	rootCmd.PersistentFlags().IntVar(&minYear, "min-year", 2020, "Ano mínimo")
	rootCmd.PersistentFlags().IntVar(&maxYear, "max-year", 2025, "Ano máximo")
	rootCmd.PersistentFlags().StringVar(&companyName, "company", "", "Nome da corporação/alvo")
	rootCmd.PersistentFlags().StringSliceVar(&commonWords, "words", []string{}, "Palavras comuns para usar (separadas por vírgula)")
	rootCmd.PersistentFlags().BoolVar(&leetMode, "leet", false, "Ativar modo leetspeak")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Modo verboso")
	rootCmd.PersistentFlags().IntVar(&maxPasswords, "max-passwords", 10000, "Número máximo de senhas a gerar")
	rootCmd.PersistentFlags().BoolVar(&allCombinations, "all", false, "Gerar TODAS as combinações possíveis (ignora limite de senhas)")

	// Comandos
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(versionCmd)
}

// Execute executa o comando raiz
func Execute() error {
	return rootCmd.Execute()
}

// validateFlags valida as flags fornecidas
func validateFlags() error {
	if minLength < 1 {
		return fmt.Errorf("tamanho mínimo deve ser maior que 0")
	}
	if maxLength < minLength {
		return fmt.Errorf("tamanho máximo deve ser maior ou igual ao tamanho mínimo")
	}
	if minYear < 1900 || minYear > 2100 {
		return fmt.Errorf("ano mínimo deve estar entre 1900 e 2100")
	}
	if maxYear < minYear {
		return fmt.Errorf("ano máximo deve ser maior ou igual ao ano mínimo")
	}
	return nil
}
