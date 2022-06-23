package routes

import (
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/link5401/fiber-chatbot/controllers"
)

func helloWorld(c *fiber.Ctx) error {
	return c.SendString("Hello World!")
}

func SetupFiberRoute(app *fiber.App) {

	// Default config
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PATCH,PUT,DELETE",
	}))

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
	app.Post("/replyIntent", controllers.ReplyIntent)
	app.Post("/addIntent", controllers.AddIntent)
	app.Get("/listIntent", controllers.ListIntent)

	app.Delete("/deleteIntent", controllers.DeleteIntent)
	app.Patch("/modifyIntent", controllers.ModifyIntent)
}
