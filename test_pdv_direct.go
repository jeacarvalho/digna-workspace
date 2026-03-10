package main

import (
	"fmt"
	"html/template"
	"os"
)

func main() {
	fmt.Println("🧪 Testando template PDV diretamente...")

	// Testar se o template pdv.html existe e pode ser parseado
	_, err := template.ParseFiles("modules/ui_web/templates/pdv.html")
	if err != nil {
		fmt.Printf("❌ Erro ao parsear pdv.html: %v\n", err)
		os.Exit(1)
	}

	// Testar se o template tem conteúdo
	fmt.Printf("✅ Template pdv.html parseado com sucesso\n")

	// Verificar se referencia componentes
	content, err := os.ReadFile("modules/ui_web/templates/pdv.html")
	if err != nil {
		fmt.Printf("❌ Erro ao ler pdv.html: %v\n", err)
		os.Exit(1)
	}

	contentStr := string(content)

	// Verificar referências a componentes
	checks := []string{
		"components/logo.html",
		"{{template",
		"{{define",
		"{{block",
	}

	fmt.Println("\n🔍 Analisando template pdv.html:")
	for _, check := range checks {
		if contains(contentStr, check) {
			fmt.Printf("   ✅ Contém: %s\n", check)
		}
	}

	// Verificar tamanho
	fmt.Printf("   📏 Tamanho do template: %d bytes\n", len(contentStr))

	// Mostrar primeiras linhas
	fmt.Println("\n📄 Primeiras 10 linhas do template:")
	lines := splitLines(contentStr)
	for i := 0; i < 10 && i < len(lines); i++ {
		fmt.Printf("   %d: %s\n", i+1, lines[i])
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 100 && (s[:100] == substr || contains(s[1:], substr)))
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i, c := range s {
		if c == '\n' {
			line := s[start:i]
			if len(line) > 100 {
				line = line[:100] + "..."
			}
			lines = append(lines, line)
			start = i + 1
		}
	}
	if start < len(s) {
		line := s[start:]
		if len(line) > 100 {
			line = line[:100] + "..."
		}
		lines = append(lines, line)
	}
	return lines
}
