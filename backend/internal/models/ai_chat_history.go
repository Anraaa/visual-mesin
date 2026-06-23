package models

import (
	"time"
)

type AiChatHistory struct {
	ID             uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID         uint       `gorm:"not null;index:idx_user_id" json:"user_id"`
	SessionID      string     `gorm:"type:varchar(255);not null" json:"session_id"`
	Question       string     `gorm:"type:text;not null" json:"question"`
	DetectedIntent *string    `gorm:"type:varchar(100)" json:"detected_intent,omitempty"`
	GeneratedSQL   *string    `gorm:"type:text" json:"generated_sql,omitempty"`
	SQLStatus      string     `gorm:"type:enum('pending','valid','invalid','error');not null;default:'pending'" json:"sql_status"`
	AiResponse     *string    `gorm:"type:longtext" json:"ai_response,omitempty"`
	Status         string     `gorm:"type:enum('queued','processing','completed','failed','rejected');not null;default:'queued'" json:"status"`
	RejectionReason *string   `gorm:"type:varchar(255)" json:"rejection_reason,omitempty"`
	StartedAt      *time.Time `json:"started_at,omitempty"`
	CompletedAt    *time.Time `json:"completed_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`

	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (AiChatHistory) TableName() string {
	return "ai_chat_history"
}

type ChatRequest struct {
	Question  string `json:"question" binding:"required"`
	SessionID string `json:"session_id"`
}

type SessionHistory struct {
	SessionID string `json:"session_id"`
	Question  string `json:"question"`
	CreatedAt string `json:"created_at"`
}
