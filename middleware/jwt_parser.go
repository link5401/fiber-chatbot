package middleware

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"

	"github.com/golang-jwt/jwt/v4"
	"github.com/link5401/fiber-chatbot/config"
)

type TokenData struct {
	Username  string
	Useragent string
	IPAdress  string
	Createdat int64
	Expires   int64
}

func ExtractTokenData(c *fiber.Ctx) (*TokenData, error) {
	token, err := VerifyToken(c)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		return &TokenData{
			Username:  fmt.Sprint(claims["user_name"]),
			Useragent: fmt.Sprint(claims["user_agent"]),
			IPAdress:  fmt.Sprint(claims["ip_address"]),
			Createdat: int64((claims["createdat"]).(float64)),
			Expires:   int64((claims["exp"]).(float64)),
		}, nil
	}
	return nil, err
}

func VerifyToken(c *fiber.Ctx) (*jwt.Token, error) {
	token := c.Get("token")

	if len(token) == 0 {
		return nil, errors.New("token is empty")
	}

	JWToken, err := jwt.Parse(token, jwtFunc)
	if err != nil {
		return nil, err
	}
	return JWToken, nil
}

func jwtFunc(token *jwt.Token) (interface{}, error) {
	return []byte(config.Config("JWT_SECRET_KEY")), nil
}
