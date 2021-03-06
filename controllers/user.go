package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/link5401/fiber-chatbot/config"
	"github.com/link5401/fiber-chatbot/database"
	"github.com/link5401/fiber-chatbot/middleware"
	models "github.com/link5401/fiber-chatbot/models"
)

// @Tags         users
// @Summary      Add an user
// @Description  Register a user to the database.
// @Param        user  body  models.User  true  "User information"
// @Accept       json
// @Produce      json
// @Success      200  {object}    ResponseMessage
// @Failure      400  {object}    HTTPError
// @Failure      500  {object}    HTTPError
// @Failure      501  {object}    HTTPError
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

// @Tags         users
// @Summary      Login
// @Description  Login to get the token
// @Param        user  body  models.User  true  "User information"
// @Accept       json
// @Produce      json
// @Success      200  {object}    string
// @Failure      400  {object}    HTTPError
// @Failure      500  {object}    HTTPError
// @Failure      501  {object}    HTTPError
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
	if err := models.CheckPasswordHash(user.Password, input.Password); err != nil {
		return c.Status(501).JSON(fiber.Map{
			"UserID":         "admin",
			"MessageContent": "Password is incorrect",
		})
	}
	token, err := middleware.GenerateJWToken(user.Username, c.Get("User-Agent"), c.IP())

	if err != nil {
		return c.Status(501).JSON(fiber.Map{"message": "system_error"})
	}
	store := database.ConfigSession()
	time_expire := config.Config("JWT_EXPIRED_TIME")
	minutesCount, _ := strconv.Atoi(time_expire)
	store.Set(user.Username, []byte(token), time.Duration(minutesCount)*time.Minute)

	return c.JSON(fiber.Map{
		"UserName":       input.Username,
		"MessageContent": "success",
		"token":          token,
	})

}

// @Tags         users
// @Summary      CheckToken
// @Description  CheckToken
// @Param        token header string true  "token"
// @Accept       json
// @Produce      json
// @Success      200  {object}    string
// @Failure      400  {object}    HTTPError
// @Failure      500  {object}    HTTPError
// @Failure      501  {object}    HTTPError
// @Router       /users/CheckToken [post]
func CheckToken(c *fiber.Ctx) error {
	TokenData, _ := middleware.ExtractTokenData(c)
	// fmt.Println(TokenData.Username)
	store := database.ConfigSession()
	return c.Status(200).JSON(fiber.Map{
		"TokenData": TokenData,
		"storeData": store,
	})
}
