package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

// versionCmd representa o comando de versão
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Mostrar versão do AsiriGen",
	Long:  `Mostra a versão atual do AsiriGen e informações de build.`,
	Run: func(cmd *cobra.Command, args []string) {
		PrintBanner()
		fmt.Printf("Build Info:\n")
		fmt.Printf("  Version:    %s\n", Version)
		fmt.Printf("  Go Version: %s\n", runtime.Version())
		fmt.Printf("  OS/Arch:    %s/%s\n", runtime.GOOS, runtime.GOARCH)
		fmt.Printf("  Author:     %s\n", Author)
	},
}
