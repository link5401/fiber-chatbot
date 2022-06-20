package main

import (
	"database/sql"
	"fmt"

	"chatbot_fiber.com/app/utils"
	"github.com/gofiber/fiber/v2"

	_ "github.com/lib/pq"
)

func helloWorld(c *fiber.Ctx) error {
	return c.SendString("Hello World!")
}
func setupFiberRoute(app *fiber.App) {
	app.Get("/", helloWorld)
	app.Post("/replyIntent", utils.ReplyIntent)
	// app.Post("/promptTest", utils.PromptTest)
}
func main() {
	// *Database setup
	db, err := sql.Open("postgres", utils.GetConnString())
	if !utils.CheckForErr(err) {
		fmt.Println("Connected to db")
	}
	// close database
	defer db.Close()
	utils.DB = db
	// *Fiber setup
	//Fiber routing
	app := fiber.New()
	setupFiberRoute(app)
	app.Listen(":3000")
}
