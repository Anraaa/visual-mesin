package models

import (
	"time"
)

type AiSchemaMap struct {
	ID               uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	IntentName       string    `gorm:"type:varchar(100);not null" json:"intent_name"`
	Keywords         string    `gorm:"type:json;not null" json:"keywords"`
	TablesInvolved   string    `gorm:"type:json;not null" json:"tables_involved"`
	SchemaContext    *string   `gorm:"type:text" json:"schema_context,omitempty"`
	FewShotExamples  *string   `gorm:"type:json" json:"few_shot_examples,omitempty"`
	Description      *string   `gorm:"type:varchar(255)" json:"description,omitempty"`
	IsActive         bool      `gorm:"not null;default:true" json:"is_active"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func (AiSchemaMap) TableName() string {
	return "ai_schema_map"
}

type AiSchemaMapRequest struct {
	IntentName      string  `json:"intent_name" binding:"required"`
	Keywords        string  `json:"keywords" binding:"required"`
	TablesInvolved  string  `json:"tables_involved" binding:"required"`
	SchemaContext   *string `json:"schema_context"`
	FewShotExamples *string `json:"few_shot_examples"`
	Description     *string `json:"description"`
	IsActive        *bool   `json:"is_active"`
}

type AiSchemaMapUpdateRequest struct {
	IntentName      *string `json:"intent_name"`
	Keywords        *string `json:"keywords"`
	TablesInvolved  *string `json:"tables_involved"`
	SchemaContext   *string `json:"schema_context"`
	FewShotExamples *string `json:"few_shot_examples"`
	Description     *string `json:"description"`
	IsActive        *bool   `json:"is_active"`
}
