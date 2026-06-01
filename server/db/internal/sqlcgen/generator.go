package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Templates []TemplateConfig `yaml:"templates"`
}

type TemplateConfig struct {
	Template string     `yaml:"template"`
	Output   string     `yaml:"output"`
	Queries  []QueryDef `yaml:"queries"`
}

type QueryDef struct {
	Name    string `yaml:"name"`
	OrderBy string `yaml:"order_by"`
}

type TemplateData struct {
	Queries []QueryDef
}

func Generate() {
	cfg := mustLoadConfig("internal/sqlcgen/config.yaml")

	for _, tplCfg := range cfg.Templates {
		generateTemplate(tplCfg)
	}
}

func generateTemplate(tplCfg TemplateConfig) {
	tplPath := filepath.Join("queries/templates", tplCfg.Template)

	tpl := mustLoadTemplate(tplPath)

	data := TemplateData{
		Queries: tplCfg.Queries,
	}

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		panic(err)
	}

	outPath := filepath.Join("queries/generated", tplCfg.Output)

	if err := os.WriteFile(outPath, buf.Bytes(), 0644); err != nil {
		panic(err)
	}

	fmt.Println("generated:", outPath)
}

func mustLoadConfig(path string) Config {
	b, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var cfg Config
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		panic(err)
	}

	return cfg
}

func mustLoadTemplate(path string) *template.Template {
	b, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return template.Must(template.New("sql").Parse(string(b)))
}
