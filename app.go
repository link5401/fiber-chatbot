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
}
func main() {
	// *Database setup
	// open database
	db, err := sql.Open("postgres", utils.GetConnString())
	if !utils.CheckForErr(err) {
		fmt.Println("Connected to db")
	}
	// close database
	defer db.Close()

	// *Fiber setup
	//Fiber routing
	app := fiber.New()
	setupFiberRoute(app)
	app.Listen(":3000")
}
