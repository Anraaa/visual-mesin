package ai

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/anraaa/visual-mesin/internal/models"
)

type IntentResult struct {
	IntentName string
	Confidence float64
	SchemaMap  *models.AiSchemaMap
}

type IntentDetector struct {
	ollama *OllamaClient
}

func NewIntentDetector(ollama *OllamaClient) *IntentDetector {
	return &IntentDetector{ollama: ollama}
}

func (d *IntentDetector) Detect(query string, schemaMaps []models.AiSchemaMap) (*IntentResult, error) {
	best := d.keywordMatch(query, schemaMaps)
	if best != nil && best.Confidence > 0.7 {
		return best, nil
	}

	result, err := d.llmClassify(query, schemaMaps)
	if err != nil {
		if best != nil {
			return best, nil
		}
		return nil, fmt.Errorf("failed to detect intent: %w", err)
	}

	return result, nil
}

func (d *IntentDetector) keywordMatch(query string, schemaMaps []models.AiSchemaMap) *IntentResult {
	q := strings.ToLower(query)
	var best *IntentResult

	for _, sm := range schemaMaps {
		var keywords []string
		if err := json.Unmarshal([]byte(sm.Keywords), &keywords); err != nil {
			continue
		}

		matchCount := 0
		for _, kw := range keywords {
			if strings.Contains(q, strings.ToLower(kw)) {
				matchCount++
			}
		}

		if matchCount > 0 {
			confidence := float64(matchCount) / float64(len(keywords))
			if best == nil || confidence > best.Confidence {
				best = &IntentResult{
					IntentName: sm.IntentName,
					Confidence: confidence,
					SchemaMap:  &sm,
				}
			}
		}
	}

	return best
}

func (d *IntentDetector) llmClassify(query string, schemaMaps []models.AiSchemaMap) (*IntentResult, error) {
	var intentList []string
	for _, sm := range schemaMaps {
		intentList = append(intentList, fmt.Sprintf("- %s: %s", sm.IntentName, safeStr(sm.Description)))
	}

	system := `Anda adalah asisten klasifikasi intent untuk sistem produksi ban.
Tugas Anda: dari daftar intent di bawah, pilih SATU intent yang paling sesuai dengan pertanyaan user.
Balas hanya dengan nama intent (satu kata), tanpa penjelasan tambahan.`

	prompt := fmt.Sprintf(`Daftar intent yang tersedia:
%s

Pertanyaan user: "%s"

Intent yang paling sesuai:`, strings.Join(intentList, "\n"), query)

	response, err := d.ollama.Generate(system, prompt)
	if err != nil {
		return nil, err
	}

	response = strings.TrimSpace(response)
	for _, sm := range schemaMaps {
		if strings.EqualFold(sm.IntentName, response) {
			return &IntentResult{
				IntentName: sm.IntentName,
				Confidence: 0.85,
				SchemaMap:  &sm,
			}, nil
		}
	}

	// fuzzy match
	responseLower := strings.ToLower(response)
	for _, sm := range schemaMaps {
		if strings.Contains(responseLower, strings.ToLower(sm.IntentName)) {
			return &IntentResult{
				IntentName: sm.IntentName,
				Confidence: 0.8,
				SchemaMap:  &sm,
			}, nil
		}
	}

	return nil, fmt.Errorf("could not classify intent from response: %s", response)
}

func safeStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
