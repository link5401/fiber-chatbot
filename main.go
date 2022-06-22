package main

import (
	"database/sql"
	"fmt"

	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
	_ "github.com/link5401/fiber-chatbot/docs"
	"github.com/link5401/fiber-chatbot/utils"
)

func helloWorld(c *fiber.Ctx) error {
	return c.SendString("Hello World!")
}

func setupFiberRoute(app *fiber.App) {
	app.Get("/", helloWorld)
	app.Get("/swagger/*", swagger.HandlerDefault)     // default
	app.Get("/swagger/*", swagger.New(swagger.Config{ // custom
		URL:         "http://example.com/doc.json",
		DeepLinking: false,
		// Expand ("list") or Collapse ("none") tag groups by default
		DocExpansion: "none",
		// Prefill OAuth ClientId on Authorize popup
		OAuth: &swagger.OAuthConfig{
			AppName:  "OAuth Provider",
			ClientId: "21bb4edc-05a7-4afc-86f1-2e151e4ba6e2",
		},
		// Ability to change OAuth2 redirect uri location
		OAuth2RedirectUrl: "http://localhost:3000/swagger/oauth2-redirect.html",
	}))
	app.Post("/replyIntent", utils.ReplyIntent)
	app.Post("/addIntent", utils.AddIntent)

	app.Delete("/deleteIntent", utils.DeleteIntent)
}

// @title Fiber Example API
// @version 1.0
// @description This is a sample swagger for Fiber
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email fiber@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:3000
// @BasePath /
func main() {
	// *Database setup
	db, err := sql.Open("postgres", utils.GetConnString())
	if !utils.CheckForErr(err) {
		fmt.Println("Connected to db")
	}
	defer db.Close()
	utils.DB = db
	// *Fiber setup
	app := fiber.New()

	setupFiberRoute(app)
	app.Listen(":3000")
}
