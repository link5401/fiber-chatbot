package main

import (
	"database/sql"
	"fmt"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

func main() {
	app := fiber.New()

	// open database
	db, err := sql.Open("postgres", getConnString())
	checkForErr(err)
	fmt.Println("Connected to db")
	// close database
	defer db.Close()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("hello")
	})

	app.Listen(":3000")
}
