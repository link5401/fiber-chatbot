package utils

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	models "chatbot_fiber.com/app/models"
	"github.com/joho/godotenv"
)

type InputMessage = models.InputMesssage
type ResponseMessage = models.ResponseMessage
type Intent = models.Intent
type Prompt = models.Prompt

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

//db info
const (
	host = "localhost"
	port = 5432
	user = "postgres"
)

/*
 *getConnString(): get Connection string used for connecting to db
 @return (string): psqlconn that is used for sql.Open()
*/
func GetConnString() string {
	//keys needed to connect to the database
	GetSecretKey()
	var dbpassword string = os.Getenv("dbpassword")
	var dbname string = os.Getenv("dbname")
	//generate the string
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, dbpassword, dbname)
	return psqlconn
}

/*
 *stringToSlice(s string) []string: gets a string of format {"a","b","c"}, converts it to an array of ["a","b","c"]
 @param  (string) : the string of the format mentioned above
 @return ([]string) an array of strings
*/
func stringToSlice(s string) []string {
	s1 := strings.Replace(s, "{", "", -1)
	s2 := strings.Replace(s1, "}", "", -1)
	s3 := strings.Replace(s2, "\"", "", -1)
	slice := strings.Split(s3, ",")
	return slice
}

/*
 *makeReplyJSON(): Conjures up a JSON Reply of Type ReponseMessage to return
 @param inputMessage: the message that the user sends is passed here
 @param messageContent: the Response to that message.
 @return ([]byte, error): The JSON mentioned above
*/
func makeReplyJSON(inputMessage InputMessage, messageContent string) ([]byte, error) {
	reply := ResponseMessage{
		UserID:         inputMessage.UserID,
		MessageContent: messageContent,
	}
	r, err := json.Marshal(reply)
	CheckForErr(err)
	return r, err
}

/*
 * This is a string for handling replyIntent()
 */
var findResponseMessageQuery string = `
	SELECT "message_content" 
	FROM "response_message"
	WHERE "intent_id" = (
		SELECT "id"
		FROM (
			SELECT "id", unnest("training_phrases") as "phrase"
			FROM "intent"
		) search_training_phrase 
		WHERE lower("phrase") LIKE lower($1) LIMIT 1
	)`

/*
 * String for finding prompts
 */
var findPromptMessageQuery string = `
	SELECT prompt_question
		FROM "prompt"
		WHERE "intent_id" = (
			SELECT "id"
			FROM (
				SELECT "id", unnest("training_phrases") as "phrase"
				FROM "intent" 
			) search_training_phrase 
			WHERE lower("phrase") LIKE lower($1) LIMIT 1
	)`

// var addIntentQuery string = `
// 	WINSERT INTO intent ("intent_name","training_phrases")
// 	VALUES($1,ARRAY [$2]);`
// var addResponseID string = `INSERT INTO response_message("intent_id")
// 	SELECT MAX(id) FROM intent;`
// var addResponseContent string = `INSERT INTO response_message("message_content") VALUES($1);`

/*
 @indexLastAsked: variable to track prompt progress.
*/
var indexLastAsked int = -1

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
func queryForPrompt(inputMessage InputMessage) ([]byte, error) {
	//Query for results
	var promptQueryResult string
	rows, err := DB.Query(findPromptMessageQuery, "%"+inputMessage.MessageContent+"%")
	CheckForErr(err)
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&promptQueryResult)
		CheckForErr(err)
	}

	//Check if there exists a prompt for this training_phrase
	if promptQueryResult != "" {
		prompts := stringToSlice(promptQueryResult)
		if indexLastAsked < len(prompts)-1 {
			indexLastAsked++
			return makeReplyJSON(inputMessage, prompts[indexLastAsked])
		}
		indexLastAsked = -1
	}
	//If run out of prompts, returns a response message
	return queryForResponse(inputMessage)

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
	return makeReplyJSON(inputMessage, messageContent)
}
