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
	//cors setup
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PATCH,PUT,DELETE",
	}))

	//mainpage
	app.Get("/", helloWorld)

	//Swagger configration

	// intents tags
	app.Post("/intents/replyIntent", controllers.ReplyIntent)
	app.Post("/intents/addIntent", controllers.AddIntent)
	app.Get("/intents/listIntent", controllers.ListIntent)
	app.Delete("/intents/deleteIntent", controllers.DeleteIntent)
	app.Patch("/intents/modifyIntent", controllers.ModifyIntent)

	//users tags
	app.Post("/users/addUser", controllers.AddUser)
	app.Post("/users/Login", controllers.Login)
}
