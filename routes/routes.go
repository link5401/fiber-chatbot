package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/link5401/fiber-chatbot/controllers"
)

func helloWorld(c *fiber.Ctx) error {
	return c.SendString("Hello World!")
}

func SetupFiberRoute(app *fiber.App) {

	//cors setup

	//mainpage
	app.Get("/", helloWorld)

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
