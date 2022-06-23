package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	models "github.com/link5401/fiber-chatbot/models"
)

type InputMessage = models.InputMesssage
type ResponseMessage = models.ResponseMessage
type Intent = models.Intent
type Prompt = models.Prompt
type HTTPError = models.HTTPError

var DB *sql.DB

/*
 *getSecretKey(): locate secret.env and load database secrets
 */
func GetSecretKey() {
	err := godotenv.Load("./dev/secret.env")
	if err != nil {
		log.Fatal("Cant load Secret file")
	}
}

/*
 *checkForErr(): checks if the error parameter is nil.
 @param (error): the error passed in to be checked
 @return (bool): A bool value to determine if there was an error or not
*/
func CheckForErr(err error) bool {
	if err != nil {
		fmt.Println(err)
		return true
	}
	return false
}

/*
 *getConnString(): get Connection string used for connecting to db
 @return (string): psqlconn that is used for sql.Open()
*/
func GetConnString() string {
	//keys needed to connect to the database
	GetSecretKey()
	dbpassword := os.Getenv("dbpassword")
	dbname := os.Getenv("dbname")
	host := os.Getenv("host")
	user := os.Getenv("user")

	//generate the string
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, 5432, user, dbpassword, dbname)
	return psqlconn
}

/*
 *getCurrentID(): returns the largest current id in "intent" TABLE
 @return int: ...
*/
func getCurrentID() int {
	var currentLargestID int
	rows, err := DB.Query(currentIDQuery)
	CheckForErr(err)
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&currentLargestID)
		CheckForErr(err)
	}
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
		rows, err := DB.Query(findPromptMessageQuery, "%"+inputMessage.MessageContent+"%")
		CheckForErr(err)
		defer rows.Close()
		for rows.Next() {
			rows.Scan(&promptQueryResult)
			CheckForErr(err)
		}
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
	rows, err := DB.Query(findResponseMessageQuery, "%"+inputMessage.MessageContent+"%")
	CheckForErr(err)
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&messageContent)
		CheckForErr(err)
	}
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
	_, i_err := DB.Exec(insertIntentQuery, newIntent.IntentName, alltrainingPhrases)
	CheckForErr(i_err)

	//check if the bot should insert prompt or message content
	if allPromptQuestion != "" {
		_, p_err := DB.Exec(insertPromptQuery, getCurrentID(), allPromptQuestion)
		CheckForErr(p_err)
	} else if newIntent.Reply.MessageContent != "" {
		_, r_err := DB.Exec(insertResponseQuery, getCurrentID(), newIntent.Reply.MessageContent)
		CheckForErr(r_err)
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
func queryForDeleteIntent(intentName string) ([]byte, error) {
	var intentID int
	rows, err := DB.Query(findIntentIDQuery, intentName)
	CheckForErr(err)
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&intentID)
		CheckForErr(err)
	}

	_, err = DB.Exec(deletePromptQuery, intentID)
	CheckForErr(err)

	_, err = DB.Exec(deleteResponseQuery, intentID)
	CheckForErr(err)

	_, err = DB.Exec(deleteIntentQuery, intentID)
	CheckForErr(err)

	return makeReplyJSON("admin", "deleted "+intentName)
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
func queryForAllIntents() ([]byte, error) {
	var (
		id               int
		intent_name      string
		training_phrases sql.NullString
		message_content  sql.NullString
		prompts          sql.NullString
		intent           []Intent
	)
	rows, err := DB.Query(listAllIntent)
	CheckForErr(err)
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&id, &intent_name, &training_phrases, &message_content, &prompts); err != nil {
			panic(err)
		}
		intent = append(intent, Intent{
			IntentName:      intent_name,
			TrainingPhrases: stringToSlice(training_phrases.String),
			Reply: ResponseMessage{
				MessageContent: message_content.String,
			},
			Prompts: Prompt{
				PromptQuestion: stringToSlice(prompts.String),
			},
		})
		fmt.Println(prompts)

	}

	i, err := json.Marshal(intent)

	return i, err
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
	rows, err := DB.Query(findIntentIDQuery, intentName)
	CheckForErr(err)
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&intentID)
		CheckForErr(err)
	}
	_, err = DB.Query(updateIntent, intent.NewName, intent.GetAllTrainingPhrases(), intentID)
	CheckForErr(err)
	_, err = DB.Query(updatePrompt, intent.GetAllPromptQuestion(), intentID)
	CheckForErr(err)
	_, err = DB.Query(updateResponse, intent.Reply.MessageContent, intentID)
	CheckForErr(err)
	i, err := json.Marshal(intent)
	return i, err
}
