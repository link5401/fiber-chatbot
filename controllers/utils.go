package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	models "github.com/link5401/fiber-chatbot/models"
	"gorm.io/gorm"
)

type InputMessage = models.InputMesssage
type ResponseMessage = models.ResponseMessage
type Intent = models.Intent
type Prompt = models.Prompt
type HTTPError = models.HTTPError

const TimeFormat = "2 Jan 2006 15:04:05"

var DB *gorm.DB

/*
 *checkForErr(): checks if the error parameter is nil.
 @param (error): the error passed in to be checked
 @return (bool): A bool value to determine if there was an error or not
*/
func CheckForErr(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}

/*
 *getCurrentID(): returns the largest current id in "intent" TABLE
 @return int: ...
*/
func getCurrentID() int {
	var currentLargestID int
	DB.Raw(currentIDQuery).Scan(&currentLargestID)
	return currentLargestID
}

/*
 *stringToSlice(s string) []string: gets a string of format {"a","b","c"}, converts it to an array of ["a","b","c"]
 @param  (string) : the string of the format mentioned above
 @return ([]string) an array of strings
*/
func stringToSlice(s string) []string {
	s = strings.Replace(s, "{", "", -1)
	s = strings.Replace(s, "}", "", -1)
	s = strings.Replace(s, "\"", "", -1)
	slice := strings.Split(s, ",")
	return slice
}

/*
 *makeReplyJSON(): Conjures up a JSON Reply of Type ReponseMessage to return
 @param inputMessage: the message that the user sends is passed here
 @param messageContent: the Response to that message.
 @return ([]byte, error): The JSON mentioned above
*/
func makeReplyJSON(UserID string, messageContent string) ([]byte, error) {
	reply := ResponseMessage{
		UserID:         UserID,
		MessageContent: messageContent,
	}
	r, err := json.Marshal(reply)
	CheckForErr(err)
	return r, err
}

/*
 *queryForPrompt: function to query correct prompt from DB
 @param (InputMessage): Contains message content, userID
 @return ([]byte, error): a JSON that contains the current prompt, userID
 ?Handling
 *This function takes out the MessageContent field.
 *Queries it to find correct id in TABLE intent using LIKE.
 *Continues to query that id that find correct prompt.
 *Prompt flow is based on @indexLastAsked.
*/
/*
 @indexLastAsked: variable to track prompt progress.
 @isPrompt: boolean to tell other functions if the bot is in prompt state or not
 @promptQueryResult: Save prompt query result
*/
var (
	indexLastAsked    int  = -1
	isPrompt          bool = false
	promptQueryResult string
)

func queryForPrompt(inputMessage InputMessage) ([]byte, error) {
	//Query for results
	if !isPrompt {
		DB.Raw(findPromptMessageQuery, "%"+inputMessage.MessageContent+"%").Scan(&promptQueryResult)
	}
	//Check if there exists a prompt for this training_phrase
	if promptQueryResult != "" {
		isPrompt = true
		prompts := stringToSlice(promptQueryResult)
		if indexLastAsked < len(prompts)-1 {
			indexLastAsked++
			return makeReplyJSON(inputMessage.UserID, prompts[indexLastAsked])
		} else {
			//If run out of prompts, returns a response message
			promptQueryResult = ""
			indexLastAsked = -1
			isPrompt = false
			return makeReplyJSON(inputMessage.UserID, "All done")
		}
	}
	return makeReplyJSON(inputMessage.UserID, "")
}

/*
 *queryForResponse: function to query correct response message from DB
 @param (InputMessage): Contains message content, userID
 @return ([]byte, error): a JSON that contains the repsonse message, userID
 ?Handling
 *This function takes out the MessageContent field.
 *Queries it to find correct id in TABLE intent using LIKE.
 *Continues to query that id that find correct response message.
*/
func queryForResponse(inputMessage InputMessage) ([]byte, error) {
	//Query for results
	var messageContent string
	DB.Raw(findResponseMessageQuery, "%"+inputMessage.MessageContent+"%").Scan(&messageContent)

	//Marshal results into JSON
	if messageContent == "" {
		return makeReplyJSON(inputMessage.UserID, "I dont get what you are saying")
	}
	return makeReplyJSON(inputMessage.UserID, messageContent)
}

/*
 *queryForInserIntent(): Executes INSERT queries that relates to Intent
 @param (Intent) an intent passed in by context
 @return (string) a resulting string showing the results added
 ?Handling
 *First, Insert the intent name as well as all the training phrases, this is NOT NULL
 *Second, if theres prompt, insert prompts. Otherwise, insert response message
 *Return the resulting string
*/
func queryForInsertIntent(newIntent Intent) ([]byte, error) {
	alltrainingPhrases := newIntent.GetAllTrainingPhrases()
	allPromptQuestion := newIntent.GetAllPromptQuestion()

	//insert intent name, training phrases
	err := DB.Exec(insertIntentQuery, newIntent.IntentName, alltrainingPhrases, time.Now().Format(TimeFormat)).Error
	if err != nil {
		return makeReplyJSON("admin", err.Error())
	}
	//check if the bot should insert prompt or message content
	fmt.Println(allPromptQuestion)
	if allPromptQuestion != "" {

		DB.Exec(insertPromptQuery, getCurrentID(), allPromptQuestion)

	} else if newIntent.Reply.MessageContent != "" {
		DB.Exec(insertResponseQuery, getCurrentID(), newIntent.Reply.MessageContent)

	}

	return makeReplyJSON("admin", "added "+newIntent.IntentName)
}

/*
 *queryForDeleteIntent() takes in an intentName to find its intentID. Then proceeds to delete all entries with that ID
 @param (string) an intent name that is served in context body
 ?Handling
 *Queries for ID
 *Delete all rows with that ID
*/
func queryForDeleteIntent(intent Intent) ([]byte, error) {
	var (
		intentID int
	)
	DB.Raw(findIntentIDQuery, intent.IntentName).Scan(&intentID)
	DB.Exec(updateDeletedQuery, intentID)
	DB.Exec(modifyLogQuery, time.Now().Format(TimeFormat), intentID)

	return makeReplyJSON("admin", "deleted")
}

/*
 *queryForAllIntents(): listing all intents existing in DB
 @return ([]byte, error) JSON of all intents
 ?Handling
 *Queries For All intents and their foreign keys
 *Put them in the format of JSON
 *Appends them.
 *Returns
*/
func queryForAllIntents() []Intent {
	var (
		id               int
		intent_name      string
		training_phrases sql.NullString
		message_content  sql.NullString
		prompt_question  sql.NullString
		created_at       sql.NullTime
		updated_at       sql.NullTime
		DeletedFlag      sql.NullBool
		intent           []Intent
	)
	rows, err := DB.Raw(listAllIntent).Rows()
	fmt.Println(CheckForErr(err))
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&id, &intent_name, &training_phrases, &created_at, &updated_at, &DeletedFlag, &message_content, &prompt_question); err != nil {
			panic(err)
		}
		if !DeletedFlag.Bool {
			intent = append(intent, Intent{
				ID:              id,
				IntentName:      intent_name,
				TrainingPhrases: stringToSlice(training_phrases.String),
				Reply: models.ResponseMessage{
					IntentID:       id,
					MessageContent: message_content.String,
				},
				Prompts: models.Prompt{
					IntentID:       id,
					PromptQuestion: stringToSlice(prompt_question.String),
				},
				CreatedAt: created_at.Time,
				UpdatedAt: updated_at.Time,
			})
		}

	}

	return intent
}

/*
 *queryForModifyIntent(): query call to modify an intent
 @return ([]byte, error) JSON of new intent
 ?Handling
 *Queries for intent ID
 *Update based on new info
 *Returns
*/

func queryForModifyIntent(intent Intent) ([]byte, error) {
	var (
		intentID   int
		intentName string = intent.IntentName
	)
	DB.Raw(findIntentIDQuery, intentName).Scan(&intentID)

	DB.Exec(updateIntent, intent.NewName, intent.GetAllTrainingPhrases(), intentID)

	DB.Exec(updatePrompt, intent.GetAllPromptQuestion(), intentID)

	DB.Exec(updateResponse, intent.Reply.MessageContent, intentID)

	DB.Exec(modifyLogQuery, time.Now().Format(TimeFormat), intentID)
	fmt.Println(intent.ID)
	i, err := json.Marshal(intent)
	return i, err
}
