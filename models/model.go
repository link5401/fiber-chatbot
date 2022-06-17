package models

type InputMesssage struct {
	UserID         string
	MessageContent string
}
type ResponseMessage struct {
	UserID         string
	MessageContent string
}
type Prompt struct {
	ParamName      string
	ParamType      string
	PromptQuestion string
}

type Intent struct {
	IntentName      string
	TrainingPhrases []string
	Reply           ResponseMessage
	Prompt          []Prompt
	LAQ             int
}
