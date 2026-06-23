package ai

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/anraaa/visual-mesin/internal/models"
)

type ChatPipeline struct {
	detector      *IntentDetector
	sqlGenerator  *SQLGenerator
	firewall      *SQLFirewall
	executor      *SQLExecutor
	ollama        *OllamaClient
	schemaRepo    SchemaMapProvider
	chatHistory   ChatHistoryRepo
}

type SchemaMapProvider interface {
	ListActive() ([]models.AiSchemaMap, error)
}

type ChatHistoryRepo interface {
	Create(h *models.AiChatHistory) error
	Update(h *models.AiChatHistory) error
}

func NewChatPipeline(
	detector *IntentDetector,
	sqlGenerator *SQLGenerator,
	firewall *SQLFirewall,
	executor *SQLExecutor,
	ollama *OllamaClient,
	schemaRepo SchemaMapProvider,
	chatHistory ChatHistoryRepo,
) *ChatPipeline {
	return &ChatPipeline{
		detector:     detector,
		sqlGenerator: sqlGenerator,
		firewall:     firewall,
		executor:     executor,
		ollama:       ollama,
		schemaRepo:   schemaRepo,
		chatHistory:  chatHistory,
	}
}

type PipelineResult struct {
	Answer         string           `json:"answer"`
	DetectedIntent string           `json:"detected_intent,omitempty"`
	GeneratedSQL   string           `json:"generated_sql,omitempty"`
	SQLStatus      string           `json:"sql_status"`
	QueryResult    *ExecutorResult  `json:"query_result,omitempty"`
	Latency        string           `json:"latency"`
}

func (p *ChatPipeline) Process(question string, userID uint, sessionID string) (*PipelineResult, error) {
	start := time.Now()
	history := &models.AiChatHistory{
		UserID:    userID,
		SessionID: sessionID,
		Question:  question,
		Status:    "processing",
		SQLStatus: "pending",
	}
	p.chatHistory.Create(history)

	result := &PipelineResult{SQLStatus: "pending"}

	schemaMaps, err := p.schemaRepo.ListActive()
	if err != nil {
		result.SQLStatus = "error"
		result.Answer = fmt.Sprintf("Gagal memuat konfigurasi: %s", err.Error())
		p.failHistory(history, err.Error())
		return result, nil
	}

	intent, err := p.detector.Detect(question, schemaMaps)
	if err != nil {
		result.SQLStatus = "error"
		result.Answer = "Maaf, saya tidak mengerti pertanyaan Anda. Silakan coba dengan pertanyaan yang lebih spesifik tentang data produksi."
		p.failHistory(history, err.Error())
		return result, nil
	}

	result.DetectedIntent = intent.IntentName
	history.DetectedIntent = &intent.IntentName

	if intent.SchemaMap == nil {
		result.Answer = "Saya tidak dapat menentukan konteks data untuk pertanyaan Anda."
		result.SQLStatus = "rejected"
		p.completeHistory(history, result.Answer)
		return result, nil
	}

	sqlQuery, err := p.sqlGenerator.Generate(question, intent.SchemaMap)
	if err != nil {
		result.SQLStatus = "error"
		result.Answer = fmt.Sprintf("Gagal menghasilkan query: %s", err.Error())
		p.failHistory(history, err.Error())
		return result, nil
	}

	result.GeneratedSQL = sqlQuery
	history.GeneratedSQL = &sqlQuery

	fwResult := p.firewall.Validate(sqlQuery)
	if !fwResult.Valid {
		result.SQLStatus = "rejected"
		result.Answer = fmt.Sprintf("Query ditolak oleh firewall: %s", fwResult.Reason)
		reason := fwResult.Reason
		history.SQLStatus = "invalid"
		history.Status = "rejected"
		history.RejectionReason = &reason
		p.chatHistory.Update(history)
		return result, nil
	}

	history.SQLStatus = "valid"

	dbConfig, err := p.resolveDBConfig(intent.SchemaMap)
	if err != nil {
		result.SQLStatus = "error"
		result.Answer = fmt.Sprintf("Gagal menentukan database: %s", err.Error())
		p.failHistory(history, err.Error())
		return result, nil
	}

	execResult, err := p.executor.Execute(sqlQuery, dbConfig)
	if err != nil {
		result.SQLStatus = "error"
		result.Answer = fmt.Sprintf("Gagal menjalankan query: %s", err.Error())
		p.failHistory(history, err.Error())
		return result, nil
	}

	result.QueryResult = execResult
	latency := time.Since(start).Round(time.Millisecond).String()
	result.Latency = latency

	finalAnswer, err := p.formatResponse(question, execResult)
	if err != nil {
		finalAnswer = fmt.Sprintf("Data ditemukan (%d baris). Latensi: %s", execResult.Total, latency)
	}
	result.Answer = finalAnswer

	p.completeHistory(history, finalAnswer)
	return result, nil
}

func (p *ChatPipeline) resolveDBConfig(schemaMap *models.AiSchemaMap) (*DBConfig, error) {
	return &DBConfig{
		Driver:   "mysql",
		Host:     "localhost",
		Port:     3306,
		Username: "app",
		Password: "apppass",
		Database: "visual_mesin",
	}, nil
}

func (p *ChatPipeline) formatResponse(question string, execResult *ExecutorResult) (string, error) {
	if execResult.Total == 0 {
		return "Tidak ada data yang ditemukan untuk pertanyaan Anda.", nil
	}

	if execResult.Total == 1 && len(execResult.Rows) > 0 {
		row := execResult.Rows[0]
		parts := ""
		for _, col := range execResult.Columns {
			if val, ok := row[col]; ok {
				parts += fmt.Sprintf("- %s: %v\n", col, val)
			}
		}
		return fmt.Sprintf("Data ditemukan:\n%s", parts), nil
	}

	system := `Anda adalah asisten yang menjelaskan data produksi ban dalam bahasa Indonesia.
Jelaskan data berikut dengan singkat dan jelas.`

	prompt := fmt.Sprintf(`Pertanyaan: %s

Data (%d baris, %d total):
Kolom: %v
Baris pertama: %v

Jelaskan data ini dengan bahasa yang mudah dipahami:`,
		question,
		len(execResult.Rows),
		execResult.Total,
		execResult.Columns,
		execResult.Rows[0],
	)

	response, err := p.ollama.Generate(system, prompt)
	if err != nil {
		return "", err
	}

	response = fmt.Sprintf("%s\n\n*Menampilkan %d dari %d baris. Latensi: %s*",
		response, len(execResult.Rows), execResult.Total, execResult.Latency)

	return response, nil
}

func (p *ChatPipeline) completeHistory(history *models.AiChatHistory, answer string) {
	now := time.Now()
	history.AiResponse = &answer
	history.Status = "completed"
	history.CompletedAt = &now
	history.StartedAt = &now
	p.chatHistory.Update(history)
}

func (p *ChatPipeline) failHistory(history *models.AiChatHistory, reason string) {
	now := time.Now()
	history.Status = "failed"
	history.SQLStatus = "error"
	history.CompletedAt = &now
	p.chatHistory.Update(history)
}

func NewDBConnectionFn(getDB func(driver, host string, port int, username, password, dbName string) (*sql.DB, error)) func(driver, host string, port int, username, password, dbName string) (*sql.DB, error) {
	return getDB
}
