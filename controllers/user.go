package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/link5401/fiber-chatbot/middleware"
	models "github.com/link5401/fiber-chatbot/models"
)

// @Tags		 users
// @Summary      Add an user
// @Description  Modify an user.
// @Param        user body models.User true "user1"
// @Accept       json
// @Produce      json
// @Success 200 {object} 	ResponseMessage
// @Failure 400 {object}   	HTTPError
// @Failure 500 {object} 	HTTPError
// @Failure 501 {object}   	HTTPError
// @Router       /users/addUser [post]
func AddUser(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}
	user.Password = models.HashPassword(user.Password)
	user.CreatedAt = time.Now()
	// Check exists
	if res := DB.Where("username = ? and deleted_at is null", user.Username).Find(&user); res.RowsAffected > 0 {
		return c.Status(501).JSON(&fiber.Map{
			"userID":          "admin",
			"message_content": "Username already existed",
		})
	}
	err := DB.Exec(addNewUser, user.Username, user.Password, time.Now().Format(TimeFormat)).Error
	if err != nil {
		return c.Status(501).JSON(&fiber.Map{
			"userID":          "admin",
			"message_content": "System error",
		})
	}

	return c.Status(200).JSON(&fiber.Map{
		"userID":          "admin",
		"message_content": "Successfully added a user",
	})

}

// @Tags		 users
// @Summary      Login
// @Description  Login
// @Param        user body models.User true "user1"
// @Accept       json
// @Produce      json
// @Success 200 {object} 	string
// @Failure 400 {object}   	HTTPError
// @Failure 500 {object} 	HTTPError
// @Failure 501 {object}   	HTTPError
// @Router       /users/Login [post]
func Login(c *fiber.Ctx) error {
	var input models.User
	if err := c.BodyParser(&input); err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	var user models.User
	// Check if username exists
	if res := DB.Where("username = ? and deleted_at is null", input.Username).First(&user); res.RowsAffected <= 0 {
		return c.Status(501).JSON(fiber.Map{
			"UserID":         "admin",
			"MessageContent": "Username doesnt exist",
		})
	}

	//Check if passwordm matches
	fmt.Println(user.Password)
	fmt.Println(models.HashPassword(input.Password))
	if err := models.CheckPasswordHash(user.Password, input.Password); err != nil {
		return c.Status(501).JSON(fiber.Map{
			"UserID":         "admin",
			"MessageContent": "Password is incorrect",
		})
	}
	token, err := middleware.GenerateJWToken(user.Username, c.Get("User-Agent"), c.IP())
	fmt.Println(token)
	if err != nil {
		return c.Status(501).JSON(fiber.Map{"message": "system_error"})
	}
	return c.JSON(fiber.Map{
		"UserName":       input.Username,
		"MessageContent": "success",
		"token":          token,
	})

}
