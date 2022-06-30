package middleware

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/link5401/fiber-chatbot/config"
)

/*
 *GenerateJWToken(): returns a JWT
 @param user_name, user_agent, ip_address
 ?Handling
 *Create a MapClaims instance with the parameters
 * calls NewWiwhtClaims() to generate a new JWT
*/
func GenerateJWToken(user_name, user_agent, ip_address string) (string, error) {
	secret_key := config.Config("JWT_SECRET_KEY")
	expire_time := config.Config("JWT_EXPIRED_TIME")

	minutesCount, _ := strconv.Atoi(expire_time)
	claims := jwt.MapClaims{}

	claims["user_name"] = user_name
	claims["user_agent"] = user_agent
	claims["ip_address"] = ip_address
	claims["createdat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(time.Minute * time.Duration(minutesCount)).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(secret_key))
	if err != nil {
		return "", err
	}
	return t, err
}
