package utils

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	models "chatbot_fiber.com/app/models"
	"github.com/gofiber/fiber/v2"
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
 *replyIntent(c *fiber.Ctx): function for handling conversations logic
 @param (c *fiber.Ctx): context of fiber, mainly used for taking request body and parse it.
 ?Handling
 *The body will consists of a InputMessage JSON from the website.
 *This function takes out the MessageContent field.
 *Queries it to find correct id in TABLE intent using LIKE.
 *Continues to query that id that find correct response message.
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
		WHERE lower("phrase") LIKE lower($1)
	)`

func ReplyIntent(c *fiber.Ctx) error {
	//Parses POST request
	inputMessage := new(InputMessage)
	if err := c.BodyParser(inputMessage); err != nil {
		return c.SendStatus(200)
	}
	var messageContent string
	rows, err := DB.Query(findResponseMessageQuery, "%"+inputMessage.MessageContent+"%")
	CheckForErr(err)
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&messageContent)
		CheckForErr(err)
	}
	reply := ResponseMessage{
		UserID:         inputMessage.UserID,
		MessageContent: messageContent,
	}
	r, err := json.Marshal(reply)
	CheckForErr(err)
	return c.SendString(string(r))
}

/*
 ! Only for testing
*/
// func queryForStuff(db sql.DB) string {
// 	intentName := "Hello"
// 	rows, err := db.Query(`SELECT "id" FROM "intent" WHERE "intent_name" = $1`, intentName)
// 	checkForErr(err)
// 	defer rows.Close()

// 	var intentID string
// 	for rows.Next() {
// 		err = rows.Scan(&intentID)
// 		checkForErr(err)
// 	}

// 	var message_content string
// 	Mrows, Merr := db.Query(`SELECT "message_content" FROM "response_message" WHERE "intent_id" = $1`, intentID)
// 	checkForErr(Merr)
// 	defer Mrows.Close()
// 	for Mrows.Next() {
// 		Merr = Mrows.Scan(&message_content)
// 		checkForErr(Merr)
// 	}
// 	return message_content
// }
