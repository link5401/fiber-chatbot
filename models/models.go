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
}
type HTTPError struct {
	status  string
	message string
}

func (i Intent) GetAllPromptQuestion() string {
	var s string = ""
	for k, p := range i.Prompt {
		s = s + p.PromptQuestion
		if k != len(i.Prompt)-1 {
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
