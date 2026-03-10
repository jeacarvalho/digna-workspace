package template

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// templateEntry armazena um template e seu timestamp
type templateEntry struct {
	template *template.Template
	modTime  time.Time
}

// TemplateManager gerencia templates com hot-reload em desenvolvimento
type TemplateManager struct {
	mu              sync.RWMutex
	templates       map[string]*templateEntry
	templateDir     string
	developmentMode bool
	funcMap         template.FuncMap
}

// NewTemplateManager cria um novo gerenciador de templates
func NewTemplateManager(templateDir string, development bool) *TemplateManager {
	return &TemplateManager{
		templates:       make(map[string]*templateEntry),
		templateDir:     templateDir,
		developmentMode: development,
		funcMap:         make(template.FuncMap),
	}
}

// AddFunc adiciona uma função ao template
func (tm *TemplateManager) AddFunc(name string, fn interface{}) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	tm.funcMap[name] = fn
}

// loadSingleTemplate carrega um único template do disco
func (tm *TemplateManager) loadSingleTemplate(name string) (*templateEntry, error) {
	templatePath := filepath.Join(tm.templateDir, name)

	// Verificar se arquivo existe
	info, err := os.Stat(templatePath)
	if err != nil {
		return nil, fmt.Errorf("template not found: %s", name)
	}

	// Ler conteúdo
	content, err := os.ReadFile(templatePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read template %s: %w", name, err)
	}

	// Parsear template - usar New com o nome correto e depois Parse
	tmpl := template.New(name).Funcs(tm.funcMap)
	tmpl, err = tmpl.Parse(string(content))
	if err != nil {
		return nil, fmt.Errorf("failed to parse template %s: %w", name, err)
	}

	return &templateEntry{
		template: tmpl,
		modTime:  info.ModTime(),
	}, nil
}

// GetTemplate obtém um template pelo nome
func (tm *TemplateManager) GetTemplate(name string) (*template.Template, error) {
	tm.mu.RLock()
	entry, exists := tm.templates[name]
	tm.mu.RUnlock()

	// Em modo desenvolvimento, verificar se o arquivo foi modificado
	if exists && tm.developmentMode {
		templatePath := filepath.Join(tm.templateDir, name)
		if info, err := os.Stat(templatePath); err == nil {
			if info.ModTime().After(entry.modTime) {
				// Arquivo foi modificado, recarregar
				exists = false
			}
		}
	}

	if exists {
		return entry.template, nil
	}

	// Carregar template do disco
	entry, err := tm.loadSingleTemplate(name)
	if err != nil {
		return nil, err
	}

	// Cachear
	tm.mu.Lock()
	tm.templates[name] = entry
	tm.mu.Unlock()

	return entry.template, nil
}

// ExecuteTemplate executa um template específico
func (tm *TemplateManager) ExecuteTemplate(templateName string, data interface{}) (string, error) {
	tmpl, err := tm.GetTemplate(templateName)
	if err != nil {
		return "", err
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template %s: %w", templateName, err)
	}

	return buf.String(), nil
}

// PreloadTemplates pré-carrega todos os templates
func (tm *TemplateManager) PreloadTemplates() error {
	fmt.Printf("[TEMPLATE] Preloading all templates from %s\n", tm.templateDir)

	// Encontrar todos os arquivos .html
	var templateFiles []string
	err := filepath.WalkDir(tm.templateDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && strings.HasSuffix(strings.ToLower(path), ".html") {
			relPath, err := filepath.Rel(tm.templateDir, path)
			if err == nil {
				templateFiles = append(templateFiles, filepath.ToSlash(relPath))
			}
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to find template files: %w", err)
	}

	// Pré-carregar cada template
	for _, tmplFile := range templateFiles {
		if _, err := tm.GetTemplate(tmplFile); err != nil {
			fmt.Printf("[TEMPLATE WARNING] Failed to preload %s: %v\n", tmplFile, err)
		}
	}

	fmt.Printf("[TEMPLATE] Preloaded %d templates\n", len(templateFiles))
	return nil
}
