package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/link5401/fiber-chatbot/models"
	postgresDriver "gorm.io/driver/postgres"

	"gorm.io/gorm"
)

var DB *gorm.DB

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
 *getSecretKey(): locate secret.env and load database secrets
 */
func GetSecretKey() {
	err := godotenv.Load("./dev/secret.env")
	if err != nil {
		log.Fatal("Cant load Secret file")
	}
}

func Connect() {
	// *Database setup
	var err error
	DB, err = gorm.Open(postgresDriver.Open(GetConnString()))
	if err != nil {
		panic(err)
	}
	DB.AutoMigrate(&models.User{}, &models.Intent{}, &models.Prompt{}, &models.ResponseMessage{}, &models.LogCustom{})

}
