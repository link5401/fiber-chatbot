// Package docs GENERATED BY SWAG; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "API Support",
            "email": "fiber@swagger.io"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/ReplyIntent": {
            "post": {
                "description": "Reply to an intent that is POST request from user",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Reply to an intent",
                "parameters": [
                    {
                        "description": "user id",
                        "name": "inputMessage",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controllers.InputMessage"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controllers.ResponseMessage"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/controllers.HTTPError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/controllers.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/controllers.HTTPError"
                        }
                    }
                }
            }
        },
        "/addIntent": {
            "post": {
                "description": "Add an intent to DB",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Add an intent to DB",
                "parameters": [
                    {
                        "description": "Name of new intent",
                        "name": "newIntent",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controllers.Intent"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controllers.ResponseMessage"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/controllers.HTTPError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/controllers.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/controllers.HTTPError"
                        }
                    }
                }
            }
        },
        "/deleteIntent": {
            "delete": {
                "description": "Delete an intent from DB, ===ONLY NEED  TO PASS IN IntentName===",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Delte an intent by querying intent name",
                "parameters": [
                    {
                        "description": "Name of the intent that you want to delete from db",
                        "name": "intentName",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controllers.Intent"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/controllers.ResponseMessage"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/controllers.HTTPError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/controllers.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/controllers.HTTPError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "controllers.HTTPError": {
            "type": "object"
        },
        "controllers.InputMessage": {
            "type": "object",
            "properties": {
                "messageContent": {
                    "type": "string"
                },
                "userID": {
                    "type": "string"
                }
            }
        },
        "controllers.Intent": {
            "type": "object",
            "properties": {
                "intentName": {
                    "type": "string"
                },
                "prompt": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/models.Prompt"
                    }
                },
                "reply": {
                    "$ref": "#/definitions/models.ResponseMessage"
                },
                "trainingPhrases": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "controllers.ResponseMessage": {
            "type": "object",
            "properties": {
                "messageContent": {
                    "type": "string"
                },
                "userID": {
                    "type": "string"
                }
            }
        },
        "models.Prompt": {
            "type": "object",
            "properties": {
                "paramName": {
                    "type": "string"
                },
                "paramType": {
                    "type": "string"
                },
                "promptQuestion": {
                    "type": "string"
                }
            }
        },
        "models.ResponseMessage": {
            "type": "object",
            "properties": {
                "messageContent": {
                    "type": "string"
                },
                "userID": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:3000",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Chatbot API",
	Description:      "Chatbot API with Fiber",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
