package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	// Testar parsing de templates
	absPath, _ := filepath.Abs("modules/ui_web/templates")
	fmt.Printf("Templates path: %s\n", absPath)

	// Tentar parsear todos os templates
	tmpl, err := template.ParseGlob("modules/ui_web/templates/*.html")
	if err != nil {
		fmt.Printf("Erro ParseGlob(*.html): %v\n", err)
	}

	// Parsear componentes
	_, err = tmpl.ParseGlob("modules/ui_web/templates/components/*.html")
	if err != nil {
		fmt.Printf("Erro ParseGlob(components/*.html): %v\n", err)
	}

	fmt.Println("✅ Templates parseados")

	// Servidor de teste
	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		data := map[string]interface{}{
			"Title": "Test",
			"Companies": []map[string]string{
				{"id": "test", "name": "Test"},
			},
		}

		err := tmpl.ExecuteTemplate(w, "login.html", data)
		if err != nil {
			fmt.Printf("Erro ao executar template: %v\n", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	fmt.Println("Servidor de teste na porta 9999")
	http.ListenAndServe(":9999", nil)
}
