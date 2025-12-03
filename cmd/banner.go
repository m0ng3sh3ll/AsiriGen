package cmd

import (
	"fmt"
)

// Constantes de Versão
const (
	Version = "1.0.0" // Major update com OSINT/Templates
	Author  = "m0ng3Sh3ll"
)

// PrintBanner exibe o banner e informações da versão
func PrintBanner() {
	// Arte ASCII (Espaço reservado para você)
	banner := `
    _         _      _  ____            
   / \    ___(_)____(_)/ ___| ___ _ __  
  / _ \  / __| | '__| | |  _ / _ \ '_ \ 
 / ___ \ \__ \ | |  | | |_| |  __/ | | |
/_/   \_\___/_|_|   |_|\____|\___|_| |_|
                                        
`
	fmt.Println(banner)
	fmt.Printf("   🚀 AsiriGen v%s\n", Version)
	fmt.Printf("   🔧 Autor: %s\n", Author)
	fmt.Println("   ───────────────────────────────────────────")
	fmt.Println()
}
