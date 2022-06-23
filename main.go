package main

import (
	"database/sql"
	"fmt"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
	controllers "github.com/link5401/fiber-chatbot/controllers"
	_ "github.com/link5401/fiber-chatbot/docs"
	middleware "github.com/link5401/fiber-chatbot/middleware"
	routes "github.com/link5401/fiber-chatbot/routes"
)

// @title Chatbot API
// @version 1.0
// @description Chatbot API with Fiber
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email fiber@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:3000
// @BasePath /
func main() {

	// *Database setup
	db, err := sql.Open("postgres", controllers.GetConnString())
	if !controllers.CheckForErr(err) {
		fmt.Println("Connected to db")
	}
	defer db.Close()
	controllers.DB = db

	// *Fiber setup
	app := fiber.New()
	middleware.Logging(app)
	routes.SetupFiberRoute(app)
	app.Listen(":3000")
}
