package main

import (
	"fmt"
	"log"
	"os"
	"time"

	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	_ "github.com/lib/pq"

	"github.com/link5401/fiber-chatbot/config"
	"github.com/link5401/fiber-chatbot/controllers"
	"github.com/link5401/fiber-chatbot/middleware"
	"github.com/link5401/fiber-chatbot/routes"

	database "github.com/link5401/fiber-chatbot/database"
	_ "github.com/link5401/fiber-chatbot/docs"
)

// @title           Chatbot API
// @version         1.0
// @description     Chatbot API with Fiber
// @termsOfService  http://swagger.io/terms/
// @contact.name    Linh
// @contact.email   linh.ldx@vn-cubesystem.com
// @license.name    Apache 2.0
// @license.url     http://www.apache.org/licenses/LICENSE-2.0.html
// @host            localhost:3000
// @BasePath        /
func main() {
	database.Connect()
	controllers.DB = database.DB

	// *Fiber setup and logging

	dir_path := "./assets/log/system"
	file_name := fmt.Sprintf("%s/%s.txt", dir_path, time.Now().Format("2006-01-02"))
	file, err := os.OpenFile(file_name, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer file.Close()
	app := fiber.New()

	app.Static("/assets", "./assets")
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PATCH,PUT,DELETE",
	}))

	app.Use(logger.New(logger.Config{
		Format:     "${time} | ${ip} | ${method} | ${status} - ${error} | ${path} \n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   config.Config("APP_TIME_ZONE"),
		Output:     file,
	}))
	app.Get("/swagger/*", swagger.HandlerDefault)

	app.Use(func(c *fiber.Ctx) error {
		return middleware.WriteLogMain(c)

	})

	// app.Use("/users/*", func(c *fiber.Ctx) error {
	// 	middleware.WriteLogMiddleWare(c)
	// 	return err
	// })
	routes.SetupFiberRoute(app)
	app.Listen(":3000")
}
