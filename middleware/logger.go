package middleware

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func Logging(app *fiber.App) {
	file, err := os.OpenFile("./request.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	app.Use(logger.New(logger.Config{
		TimeFormat: "2 Jan 2006 15:04:05",
		TimeZone:   "Asia/Ho_Chi_Minh",
		Output:     file,
	}))

}
