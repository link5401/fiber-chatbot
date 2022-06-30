package routes

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/link5401/fiber-chatbot/controllers"
	"github.com/link5401/fiber-chatbot/middleware"
	pusher "github.com/pusher/pusher-http-go"
)

func helloWorld(c *fiber.Ctx) error {
	return c.SendString("Hello World!")
}

var pusherClient = pusher.Client{
	AppID:   "1429753",
	Key:     "50d2887e8902ddf96951",
	Secret:  "7c8d0b7b144bd53a2394",
	Cluster: "ap1",
	Secure:  true,
}

func SetupFiberRoute(app *fiber.App) {

	//cors setup

	//mainpage
	app.Get("/", helloWorld)

	// intents tags
	intents := app.Group("/intents", middleware.AppAuthen)
	intents.Post("/replyIntent", controllers.ReplyIntent)
	intents.Post("/addIntent", controllers.AddIntent)
	intents.Post("/listIntent", controllers.ListIntent)
	intents.Delete("/deleteIntent", controllers.DeleteIntent)
	intents.Patch("/modifyIntent", controllers.ModifyIntent)

	//users tags
	app.Post("/users/Login", controllers.Login)
	users := app.Group("/users", middleware.AppAuthen)
	users.Post("/addUser", controllers.AddUser)
	users.Post("/CheckToken", controllers.CheckToken)

	app.Post("/api/messages", func(c *fiber.Ctx) error {
		var data controllers.InputMessage

		if err := c.BodyParser(&data); err != nil {
			return err
		}
		fmt.Println(data)
		pusherClient.Trigger("chat", "message", data)

		return c.JSON([]string{})
	})
}
