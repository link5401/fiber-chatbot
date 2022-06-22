package controllers

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

// ReplyIntent ================================================================================
// @Summary Reply to an intent
// @Description Reply to an intent that is POST request from user
// @Param inputMessage body InputMessage true "user id"
// @Accept  json
// @Produce  json
// @Success 200 {object} ResponseMessage
// @Failure 400 {object} HTTPError
// @Failure 404 {object} HTTPError
// @Failure 500 {object} HTTPError
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

// @Summary Add an intent to DB
// @Description Add an intent to DB
// @Param newIntent  body Intent true "Name of new intent"
// @Accept  json
// @Produce  json
// @Success 200 {object} ResponseMessage
// @Failure 400 {object} HTTPError
// @Failure 404 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Router /addIntent [post]
func AddIntent(c *fiber.Ctx) error {
	newIntent := new(Intent)
	if err := c.BodyParser(newIntent); err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}
	s, err := queryForInsertIntent(*newIntent)
	CheckForErr(err)

	return c.SendString(string(s))
}

/*
 *DeleteIntent(): serves DELETE request that revolves around Intent
 @param (*fiber.Ctx) context of fiber
 ?Handling
 *Parse the request body
 *Calls qeuryForDeleteIntent
*/
// @Summary Delte an intent by querying intent name
// @Description Delete an intent from DB, ===ONLY NEED  TO PASS IN IntentName===
// @Param intentName body Intent true "Name of the intent that you want to delete from db"
// @Accept  json
// @Produce  json
// @Success 200 {object} ResponseMessage
// @Failure 400 {object} HTTPError
// @Failure 404 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Router /deleteIntent [delete]
func DeleteIntent(c *fiber.Ctx) error {
	var intent Intent
	if err := c.BodyParser(&intent); err != nil {
		fmt.Println(err)
		return c.SendStatus(http.StatusBadRequest)
	}
	// fmt.Println(intentName)
	s, err := queryForDeleteIntent(intent.IntentName)
	CheckForErr(err)
	return c.SendString(string(s))
}

// @Summary List all intents and training phrases
// @Description List all intents
// @Produce  json
// @Success 200 {object} string
// @Failure 400 {object} HTTPError
// @Failure 404 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Router /listIntent [get]
func ListIntent(c *fiber.Ctx) error {

	i, err := queryForAllIntents()
	CheckForErr(err)
	return c.SendString(string(i))
}
