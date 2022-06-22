package utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

/*
 *ReplyIntent(): A handler function for /replyIntent.
 *                This function will query Prompt and Response with a InputMessage int request body
 @param (fiber.Ctx): context of fiber. This is mainly to access request/response.
*/

// Check token ================================================================================
// @Tags User
// @Summary Check token
// @Description Check token
// @Param appKey header string true "Contact Admin to get"
// @Accept  json
// @Produce  json
// @Success 200 {object} string "{message: 'Token is correct'}"
// @Failure 400 {object} string "{message: string}"
// @Failure 408 {object} string "{message: string}"
// @Failure 500 {object} string "{message: string}"
// @Router /ReplyIntent [post]
func ReplyIntent(c *fiber.Ctx) error {
	//Parses POST request
	inputMessage := new(InputMessage)
	if err := c.BodyParser(inputMessage); err != nil {
		return c.SendStatus(http.StatusBadRequest)
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
		return c.SendStatus(http.StatusBadRequest)
	}
	s := queryForInsertIntent(*newIntent)
	return c.SendString(s)
}

/*
 *DeleteIntent(): serves DELETE request that revolves around Intent
 @param (*fiber.Ctx) context of fiber
 ?Handling
 *Parse the request body
 *Calls qeuryForDeleteIntent
*/
func DeleteIntent(c *fiber.Ctx) error {
	intent := new(Intent)
	if err := c.BodyParser(intent); err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}
	fmt.Println(intent.IntentName)
	queryForDeleteIntent(intent.IntentName)
	return c.SendStatus(200)
}
