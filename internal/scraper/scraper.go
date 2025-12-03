package scraper

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

// ScrapeConfig configura o scraper
type ScrapeConfig struct {
	URL        string
	MaxWords   int
	MinWordLen int
	Timeout    time.Duration
	UserAgent  string
	IgnoreCert bool
}

// DefaultConfig retorna a configuração padrão
func DefaultConfig() ScrapeConfig {
	return ScrapeConfig{
		MaxWords:   100, // Aumentei para pegar mais contexto
		MinWordLen: 4,
		Timeout:    60 * time.Second, // Timeout maior para browser
		// UserAgent moderno para passar despercebido
		UserAgent:  "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		IgnoreCert: true,
	}
}

// ScrapeUrl extrai palavras relevantes de uma URL usando Headless Browser
func ScrapeUrl(targetUrl string, config ScrapeConfig) ([]string, error) {
	if !strings.HasPrefix(targetUrl, "http") {
		targetUrl = "https://" + targetUrl
	}

	// Opções do Chrome
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true), // Mude para false se quiser ver o browser abrindo
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("ignore-certificate-errors", config.IgnoreCert),
		chromedp.UserAgent(config.UserAgent),
		// Flags adicionais para evitar detecção de automação básica
		chromedp.Flag("enable-automation", false),
		chromedp.Flag("disable-blink-features", "AutomationControlled"),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// Criar contexto do browser com log silenciado para evitar poluição visual
	ctx, cancel := chromedp.NewContext(allocCtx,
		chromedp.WithLogf(func(string, ...interface{}) {}),
		chromedp.WithErrorf(func(string, ...interface{}) {}),
	)
	defer cancel()

	// Timeout total
	ctx, cancel = context.WithTimeout(ctx, config.Timeout)
	defer cancel()

	var resBody string
	var title string

	fmt.Printf("🌐 Navegando para %s (Modo: Headless Chrome)...\n", targetUrl)

	// Tarefas do Chromedp
	err := chromedp.Run(ctx,
		// Configurar headers extras para parecer humano
		network.SetExtraHTTPHeaders(network.Headers{
			"Accept-Language": "pt-BR,pt;q=0.9,en-US;q=0.8,en;q=0.7",
			"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
		}),
		chromedp.Navigate(targetUrl),
		// Delay humano para carregamento inicial
		chromedp.Sleep(3*time.Second),
		// Obter título para verificação rápida
		chromedp.Title(&title),
		// Scroll para baixo para ativar lazy loading (Human-like)
		chromedp.Evaluate(`window.scrollTo(0, document.body.scrollHeight)`, nil),
		chromedp.Sleep(2*time.Second), // Esperar conteúdo carregar
		// Extrair texto visível do corpo
		chromedp.Text("body", &resBody, chromedp.ByQuery),
	)

	if err != nil {
		return nil, fmt.Errorf("erro no navegador: %v", err)
	}

	// Verificar CAPTCHA ou Bloqueios
	if isBlocked(title, resBody) {
		fmt.Println("⚠️  ALERTA: Possível Captcha ou Bloqueio detectado!")
		fmt.Println("    O site pode estar usando Cloudflare/Incapsula ou requer interação humana.")
		fmt.Println("    Dica: Tente acessar o site manualmente para verificar.")
		return nil, fmt.Errorf("acesso bloqueado ou captcha detectado")
	}

	// Adicionar o título ao texto para análise
	fullText := title + " " + resBody

	return processText(fullText, config)
}

// isBlocked verifica sinais comuns de bloqueio
func isBlocked(title, body string) bool {
	lowerTitle := strings.ToLower(title)
	lowerBody := strings.ToLower(body)

	blockKeywords := []string{
		"attention required", "access denied", "security check",
		"cloudflare", "captcha", "robot", "blocked", "forbidden",
		"just a moment", "wait a moment", "challenge",
		"pardon our interruption",
	}

	for _, keyword := range blockKeywords {
		if strings.Contains(lowerTitle, keyword) {
			return true
		}
		// Verificar corpo apenas se for muito curto (páginas de bloqueio costumam ter pouco texto real)
		if len(body) < 500 && strings.Contains(lowerBody, keyword) {
			return true
		}
	}
	return false
}

// processText limpa e conta as palavras
func processText(text string, config ScrapeConfig) ([]string, error) {
	wordCounts := make(map[string]int)

	// Limpeza e split
	words := cleanAndSplit(text, config.MinWordLen)

	for _, w := range words {
		wordCounts[w]++
	}

	// Converter mapa para slice ordenado por frequência
	type wordPair struct {
		Word  string
		Count int
	}
	var pairs []wordPair
	for k, v := range wordCounts {
		// Filtrar palavras muito comuns irrelevantes (stop words simples)
		if !isStopWord(k) {
			pairs = append(pairs, wordPair{k, v})
		}
	}

	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].Count > pairs[j].Count
	})

	// Selecionar top N palavras
	var result []string
	for i, p := range pairs {
		if i >= config.MaxWords {
			break
		}
		result = append(result, p.Word)
	}

	return result, nil
}

var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9áéíóúâêîôûãõçÁÉÍÓÚÂÊÎÔÛÃÕÇ]`)

func cleanAndSplit(text string, minLen int) []string {
	// Substituir quebras de linha por espaço
	text = strings.ReplaceAll(text, "\n", " ")
	text = strings.ReplaceAll(text, "\r", " ")
	text = strings.ReplaceAll(text, "\t", " ")

	// Substituir não alfanuméricos por espaço
	cleanText := nonAlphanumericRegex.ReplaceAllString(text, " ")
	parts := strings.Fields(cleanText)

	var words []string
	for _, p := range parts {
		w := strings.ToLower(strings.TrimSpace(p))
		if len(w) >= minLen {
			words = append(words, w)
		}
	}
	return words
}

// isStopWord filtra palavras comuns irrelevantes
func isStopWord(word string) bool {
	stopWords := map[string]bool{
		"para": true, "com": true, "não": true, "que": true, "dos": true,
		"são": true, "das": true, "uma": true, "mas": true, "por": true,
		"sobre": true, "entre": true, "seus": true, "muito": true,
		"this": true, "that": true, "with": true, "from": true, "your": true,
		"contact": true, "policy": true, "rights": true, "reserved": true,
		"privacy": true, "terms": true, "menu": true, "home": true,
		"cookie": true, "cookies": true, "site": true, "website": true,
		"copyright": true, "todos": true, "direitos": true, "reservados": true,
	}
	return stopWords[word]
}
