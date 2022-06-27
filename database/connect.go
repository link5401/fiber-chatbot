package database

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/storage/postgres"
	"github.com/link5401/fiber-chatbot/config"
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
	dbpassword := config.Config("dbpassword")
	dbname := config.Config("dbname")
	host := config.Config("host")
	user := config.Config("user")

	//generate the string
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, 5432, user, dbpassword, dbname)
	return psqlconn
}

/*
* Connect() Connect to database, Automigrates db
 */
func Connect() {
	// *Database setup
	var err error
	DB, err = gorm.Open(postgresDriver.Open(GetConnString()))
	if err != nil {
		panic(err)
	}
	DB.AutoMigrate(&models.User{}, &models.Intent{}, &models.Prompt{}, &models.ResponseMessage{}, &models.LogCustom{})

}

/*
 * ConfigSession(): Return configuration of DB session
 */
func ConfigSession() *postgres.Storage {
	host := config.Config("host")
	port := config.Config("port")
	user := config.Config("user")
	password := config.Config("dbpassword")
	name := config.Config("dbname")
	sshmode := config.Config("SSH")
	post, _ := strconv.Atoi(port)

	store := postgres.New(postgres.Config{
		Host:       host,
		Port:       post,
		Username:   user,
		Password:   password,
		Database:   name,
		Table:      "session",
		Reset:      false,
		GCInterval: 10 * time.Second,
		SslMode:    sshmode,
	})

	return store

}
