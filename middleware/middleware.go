package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/link5401/fiber-chatbot/database"
	"github.com/link5401/fiber-chatbot/models"
	services "github.com/link5401/fiber-chatbot/services"
)

func AppAuthen(c *fiber.Ctx) error {
	db := database.DB
	user := new(models.User)
	store := database.ConfigSession()
	is_error := 0

	// Check token valid
	tokenData, err := ExtractTokenData(c)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// Check user exists
	if res := db.Where("username = ? and deleted_at is null", tokenData.Username).First(&user); res.RowsAffected <= 0 {
		is_error = 1
	}
	// Check session Exist and comparse token
	sess, err := store.Get(tokenData.Username)
	authen := c.Get("token")
	if err != nil || len(sess) == 0 || string(sess) != authen {
		is_error = 1
	}

	// Check User Agent / IP
	ip_address := c.IP()
	if ip_address != tokenData.IPAdress {
		is_error = 1
	}

	if is_error == 1 {
		return c.Status(400).JSON(fiber.Map{
			"message": "Token is incorrect",
		})
	}

	// Set data log
	services.UserName = tokenData.Username
	WriteLogMidleWare(c)

	return c.Next()
}
