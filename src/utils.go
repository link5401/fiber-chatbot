package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

/*
*@getSecretKey(): locate secret.env and load database secrets
 */
func getSecretKey() {
	err := godotenv.Load("./dev/secret.env")
	if err != nil {
		log.Fatal("Cant load Secret file")
	}
}

/*
*@utils checkForErr(): checks if the error parameter is nil.
*@param (error): the error passed in to be checked
*@return (bool): A bool value to determine if there was an error or not
 */
func checkForErr(err error) bool {
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
*@utils getConnString(): get Connection string used for connecting to db
*@return (string): psqlconn that is used for sql.Open()
 */
func getConnString() string {
	//keys needed to connect to the database
	getSecretKey()
	var dbpassword string = os.Getenv("dbpassword")
	var dbname string = os.Getenv("dbname")
	//generate the string
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, dbpassword, dbname)
	return psqlconn
}

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
