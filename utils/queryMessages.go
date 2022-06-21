package utils

/*
 * This is a string for handling replyIntent()
 */
const findResponseMessageQuery string = `
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
const findPromptMessageQuery string = `
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
const currentIDQuery = `SELECT MAX(id) FROM intent`
const insertIntentQuery = `INSERT INTO intent("intent_name", "training_phrases") VALUES ($1,ARRAY[$2])`
const insertResponseQuery = `INSERT INTO response_message("intent_id", "message_content") VALUES ($1,$2)`
const insertPromptQuery = `INSERT INTO prompt("intent_id","prompt_question") VALUES ($1,ARRAY[$2])`

const findIntentIDQuery = `SELECT id FROM intent WHERE intent_name = $1`

const deleteIntentQuery = `DELETE FROM intent WHERE id = $1`
const deleteResponseQuery = `DELETE FROM response_message WHERE intent_id = $1`
const deletePromptQuery = `DELETE FROM prompt WHERE intent_id = $1`

// const restartSequenceQuery = `ALTER SEQUENCE intent_id_seq RESTART WITH $1`
