package utils

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
)

/*
 *ReplyIntent(): A handler function for /replyIntent.
 *                This function will query Prompt and Response with a InputMessage int request body
 @param (fiber.Ctx): context of fiber. This is mainly to access request/response.
*/

func ReplyIntent(c *fiber.Ctx) error {
	//Parses POST request
	inputMessage := new(InputMessage)
	if err := c.BodyParser(inputMessage); err != nil {
		return c.SendStatus(200)
	}

	//Query
	p, p_err := queryForPrompt(*inputMessage)
	CheckForErr(p_err)
	r, r_err := queryForResponse(*inputMessage)
	CheckForErr(r_err)

	//parsing []byte to ResponseMessage to read result
	pResponse := ResponseMessage{}
	rResponse := ResponseMessage{}
	json.Unmarshal(p, &pResponse)
	json.Unmarshal(r, &rResponse)

	//decide if the bot should prompt or response
	if pResponse.MessageContent != "" {
		return c.SendString(string(p))
	}
	return c.SendString(string(r))
}

/*
 *AddIntent(): function for fiber to call when user wants to create an Intent on the DB
 @param (*fiber.Ctx) context of fiber
 ?Handling
 *Parse the request body into newIntent
 *Calls queryForInsertIntent(*newIntent)
*/
func AddIntent(c *fiber.Ctx) error {
	newIntent := new(Intent)
	if err := c.BodyParser(newIntent); err != nil {
		c.SendStatus(200)
	}
	s := queryForInsertIntent(*newIntent)
	return c.SendString(s)
}
