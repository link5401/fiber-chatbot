package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/link5401/fiber-chatbot/middleware"
	"github.com/pusher/pusher-http-go"
)

// @Tags				Intents
// @Summary      List all intents and training phrases
// @Param        token header string true  "token"
// @Description  List all intents
// @Produce      json
// @Success      200  {object}  []Intent
// @Failure      400  {object}  HTTPError
// @Failure      404  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router       /intents/listIntent [post]
func ListIntent(c *fiber.Ctx) error {

	i := queryForAllIntents()
	// middleware.WriteLogMain(c)
	// middleware.InsertLog(c)
	return c.Status(200).JSON(&fiber.Map{
		"intentList": i,
	})
}

/*
 *ReplyIntent(): A handler function for /replyIntent.
 *                This function will query Prompt and Response with a InputMessage int request body
 @param (fiber.Ctx): context of fiber. This is mainly to access request/response.
*/
var pusherClient = pusher.Client{
	AppID:   "1429753",
	Key:     "50d2887e8902ddf96951",
	Secret:  "7c8d0b7b144bd53a2394",
	Cluster: "ap1",
	Secure:  true,
}

// ReplyIntent ================================================================================
// @Tags				Intents
// @Summary      Reply to an intent
// @Description  Reply to an intent that is POST request from user
// @Param        inputMessage  body  InputMessage  true  "Message from the user"
// @Param        token header string true  "token"
// @Accept       json
// @Produce      json
// @Success      200  {object}  ResponseMessage
// @Failure      400  {object}  HTTPError
// @Failure      404  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router       /intents/ReplyIntent [post]
func ReplyIntent(c *fiber.Ctx) error {
	//Parses POST request
	inputMessage := new(InputMessage)
	if err := c.BodyParser(&inputMessage); err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	fmt.Println(inputMessage)
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
		// middleware.WriteLogMain(c)
		// middleware.InsertLog(c)
		return c.Status(200).JSON(&fiber.Map{
			"MessageContent": pResponse.MessageContent,
			"UserID":         pResponse.UserID,
		})
	}

	// middleware.WriteLogMain(c)

	middleware.InsertLog(c)
	return c.Status(200).JSON(&fiber.Map{
		"MessageContent": rResponse.MessageContent,
		"UserID":         rResponse.UserID,
	})
}

/*
 *AddIntent(): function for fiber to call when user wants to create an Intent on the DB
 @param (*fiber.Ctx) context of fiber
 ?Handling
 *Parse the request body into newIntent
 *Calls queryForInsertIntent(*newIntent)
*/
// @Tags				Intents
// @Summary             Add an intent to DB
// @Description         Add an intent to DB
// @Param               newIntent  body  Intent  true  "Name of new intent"
// @Param        		token header string true  "token"
// @deletedFlag(hidden  = true)
// @Accept              json
// @Produce             json
// @Success             200  {object}  ResponseMessage
// @Failure             400  {object}  HTTPError
// @Failure             404  {object}  HTTPError
// @Failure             500  {object}  HTTPError
// @Router              /intents/addIntent [post]
func AddIntent(c *fiber.Ctx) error {
	newIntent := new(Intent)
	if err := c.BodyParser(newIntent); err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}
	s, err := queryForInsertIntent(*newIntent)
	CheckForErr(err)
	// middleware.WriteLogMain(c)
	// middleware.InsertLog(c)
	return c.SendString(string(s))
}

/*
 *DeleteIntent(): serves DELETE request that revolves around Intent
 @param (*fiber.Ctx) context of fiber
 ?Handling
 *Parse the request body
 *Calls qeuryForDeleteIntent
*/
// @Tags				Intents
// @Summary      Delte an intent by querying intent name
// @Description  Delete an intent from DB, ===ONLY NEED  TO PASS IN IntentName===
// @Param        intentName  body  Intent  true  "Name of the intent that you want to delete from db"
// @Param        token header string true  "token"
// @Accept       json
// @Produce      json
// @Success      200  {object}  ResponseMessage
// @Failure      400  {object}  HTTPError
// @Failure      404  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router       /intents/deleteIntent [delete]
func DeleteIntent(c *fiber.Ctx) error {
	var intent Intent
	if err := c.BodyParser(&intent); err != nil {
		// fmt.Println(err)
		return c.SendStatus(http.StatusBadRequest)
	}
	// fmt.Println(intentName)
	s, err := queryForDeleteIntent(intent)
	CheckForErr(err)
	// middleware.WriteLogMain(c)
	middleware.InsertLog(c)
	return c.SendString(string(s))
}

// @Tags				Intents
// @Summary      Modify an intent
// @Description  Modify an intent to the body's intent by querying "IntentName".
// @Param        intent  body  Intent  true  "The new intent, pass in NewName to change the current name"
// @Param        token header string true  "token"
// @Accept       json
// @Produce      json
// @Success      200  {object}  Intent
// @Failure      400  {object}  HTTPError
// @Failure      404  {object}  HTTPError
// @Failure      500  {object}  HTTPError
// @Router       /intents/modifyIntent [patch]
func ModifyIntent(c *fiber.Ctx) error {
	var intent Intent
	if err := c.BodyParser(&intent); err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}
	_, err := queryForModifyIntent(intent)
	CheckForErr(err)
	// middleware.WriteLogMain(c)
	middleware.InsertLog(c)
	return c.Status(200).JSON(&fiber.Map{
		"intentName":      intent.NewName,
		"TrainingPhrases": intent.TrainingPhrases,
		"Responses":       intent.Reply.MessageContent,
		"Prompts":         intent.Prompts,
	})
}
