package ai

import (
	"fmt"
	"strings"

	"github.com/anraaa/visual-mesin/internal/models"
)

type SQLGenerator struct {
	ollama *OllamaClient
}

func NewSQLGenerator(ollama *OllamaClient) *SQLGenerator {
	return &SQLGenerator{ollama: ollama}
}

func (g *SQLGenerator) Generate(query string, schemaMap *models.AiSchemaMap) (string, error) {
	if schemaMap != nil && schemaMap.FewShotExamples != nil && *schemaMap.FewShotExamples != "" {
		sql, err := g.generateWithFewShot(query, schemaMap)
		if err == nil {
			return sql, nil
		}
	}

	return g.generateDirect(query, schemaMap)
}

func (g *SQLGenerator) generateWithFewShot(query string, schemaMap *models.AiSchemaMap) (string, error) {
	system := `Anda adalah asisten SQL untuk database produksi ban.
Anda HANYA boleh membuat query SELECT (read-only).
Jangan membuat query INSERT, UPDATE, DELETE, DROP, ALTER, TRUNCATE, atau CREATE.
Balas HANYA dengan SQL query, tanpa markdown, tanpa penjelasan, tanpa backtick.`

	prompt := fmt.Sprintf(`Tabel: %s

Konteks skema:
%s

Contoh:
%s

Pertanyaan user: "%s"

SQL query:`,
		schemaMap.TablesInvolved,
		safeStr(schemaMap.SchemaContext),
		*schemaMap.FewShotExamples,
		query)

	response, err := g.ollama.Generate(system, prompt)
	if err != nil {
		return "", err
	}

	return g.cleanSQL(response), nil
}

func (g *SQLGenerator) generateDirect(query string, schemaMap *models.AiSchemaMap) (string, error) {
	tables := "unknown"
	if schemaMap != nil {
		tables = schemaMap.TablesInvolved
	}

	context := ""
	if schemaMap != nil && schemaMap.SchemaContext != nil {
		context = *schemaMap.SchemaContext
	}

	system := `Anda adalah asisten SQL untuk database produksi ban read-only.
Anda HANYA boleh membuat query SELECT.
Balas HANYA dengan SQL query tanpa markdown, backtick, atau penjelasan.`

	prompt := fmt.Sprintf(`Tabel tersedia: %s

Skema tabel:
%s

Pertanyaan: "%s"

SQL query SELECT:`,
		tables, context, query)

	response, err := g.ollama.Generate(system, prompt)
	if err != nil {
		return "", err
	}

	return g.cleanSQL(response), nil
}

func (g *SQLGenerator) cleanSQL(sql string) string {
	sql = strings.TrimSpace(sql)
	sql = strings.TrimPrefix(sql, "```sql")
	sql = strings.TrimPrefix(sql, "```")
	sql = strings.TrimSuffix(sql, "```")
	sql = strings.TrimSpace(sql)
	return sql
}
