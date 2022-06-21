package utils

/*
 * This is a string for handling replyIntent()
 */
var findResponseMessageQuery string = `
 SELECT "message_content" 
 FROM "response_message"
 WHERE "intent_id" = (
	 SELECT "id"
	 FROM (
		 SELECT "id", unnest("training_phrases") as "phrase"
		 FROM "intent"
	 ) search_training_phrase 
	 WHERE lower("phrase") LIKE lower($1) LIMIT 1
 )`

/*
 * String for finding prompts
 */
var findPromptMessageQuery string = `
	 SELECT prompt_question
		 FROM "prompt"
		 WHERE "intent_id" = (
			 SELECT "id"
			 FROM (
				 SELECT "id", unnest("training_phrases") as "phrase"
				 FROM "intent" 
			 ) search_training_phrase 
			 WHERE lower("phrase") LIKE lower($1)LIMIT 1
	 )`
var currentIDQuery = `SELECT MAX(id) FROM "intent"`
var insertIntentQuery = `INSERT INTO intent("intent_name", "training_phrases") VALUES ($1,ARRAY[$2])`
var insertResponseQuery = `INSERT INTO response_message("intent_id", "message_content") VALUES ($1,$2)`
var insertPromptQuery = `INSERT INTO prompt("intent_id","prompt_question") VALUES ($1,ARRAY[$2])`
