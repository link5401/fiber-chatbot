package models

import (
	"time"
)

type InputMesssage struct {
	UserID         string `gorm:"ignoreMigration"`
	MessageContent string `gorm:"ignoreMigration"`
}
type ResponseMessage struct {
	IntentID       int    `json:"intent_id" swaggerignore:"true"`
	UserID         string `json:"user_id"`
	MessageContent string `json:"message_content"`
}
type StringArray []string
type Prompt struct {
	// ParamName      string
	// ParamType      string
	IntentID       int         `json:"intent_id" swaggerignore:"true"`
	PromptQuestion StringArray `gorm:"type:text[]"`
}

type Intent struct {
	ID              int             `json:"id" gorm:"primaryKey,unique" swaggerignore:"true"`
	NewName         string          `gorm:"ignoreMigration"`
	IntentName      string          `json:"IntentName" gorm:"unique"`
	TrainingPhrases StringArray     `gorm:"type:text[]"`
	Reply           ResponseMessage `gorm:"foreignkey:IntentID;references:ID;constraints:OnDelete:CASCADE"`
	Prompts         Prompt          `gorm:"foreignkey:IntentID;references:ID;constraints:OnDelete:CASCADE"`
	CreatedAt       time.Time       `json:"created_at" swaggerignore:"true"`
	UpdatedAt       time.Time       `json:"updated_at" swaggerignore:"true"`
	DeletedFlag     bool            `swaggerignore:"true"`
}
type HTTPError struct {
	status  string
	message string
}

func (i Intent) GetAllPromptQuestion() string {
	var s string = ""
	for k, p := range i.Prompts.PromptQuestion {
		s = s + p
		if k != len(i.Prompts.PromptQuestion)-1 {
			s += ","
		}
	}
	return s
}
func (i Intent) GetAllTrainingPhrases() string {
	var s string = ""
	for k, tp := range i.TrainingPhrases {
		s = s + tp
		if k != len(i.TrainingPhrases)-1 {
			s += ","
		}
	}
	return s
}
