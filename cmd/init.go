package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// initCmd representa o comando de inicialização
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Inicializa arquivos de configuração padrão",
	Long: `Cria um arquivo patterns.yaml padrão no diretório atual se ele não existir.
Isso ajuda a começar rapidamente com padrões personalizados.`,
	Run: runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) {
	filename := "patterns.yaml"

	// Verificar se já existe
	if _, err := os.Stat(filename); err == nil {
		fmt.Printf("ℹ️  O arquivo '%s' já existe. Nenhuma ação realizada.\n", filename)
		return
	}

	// Conteúdo padrão
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

	// Criar arquivo
	if err := os.WriteFile(filename, []byte(defaultContent), 0644); err != nil {
		fmt.Printf("❌ Erro ao criar arquivo '%s': %v\n", filename, err)
		os.Exit(1)
	}

	fmt.Printf("✅ Arquivo '%s' criado com sucesso!\n", filename)
	fmt.Println("Agora você pode editá-lo para adicionar seus próprios padrões.")
}

