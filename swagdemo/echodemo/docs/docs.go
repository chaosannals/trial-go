// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/accounts/add": {
            "put": {
                "produces": [
                    "application/json"
                ],
                "summary": "Some Endpoint",
                "parameters": [
                    {
                        "minimum": 1,
                        "type": "integer",
                        "example": 1,
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "maximum": 20,
                        "minimum": 1,
                        "type": "integer",
                        "example": 10,
                        "name": "page_size",
                        "in": "query"
                    }
                ],
                "responses": {}
            }
        },
        "/accounts/delete": {
            "delete": {
                "produces": [
                    "application/json"
                ],
                "summary": "Some Endpoint",
                "parameters": [
                    {
                        "description": "Query Params",
                        "name": "request",
                        "in": "body",
                        "schema": {
                            "$ref": "#/definitions/controller.IndexAddParam"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/accounts/list": {
            "get": {
                "description": "get string by ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "accounts"
                ],
                "summary": "Show an account",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Account ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/util.HTTPError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/util.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/util.HTTPError"
                        }
                    }
                }
            }
        },
        "/accounts/{id}": {
            "get": {
                "description": "get string by ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "accounts"
                ],
                "summary": "Show an account",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Account ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "object"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/util.HTTPError"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/util.HTTPError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/util.HTTPError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "controller.IndexAddParam": {
            "type": "object",
            "properties": {
                "page": {
                    "type": "integer",
                    "minimum": 1,
                    "example": 1
                },
                "page_size": {
                    "type": "integer",
                    "maximum": 20,
                    "minimum": 1,
                    "example": 10
                }
            }
        },
        "util.HTTPError": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer",
                    "example": 400
                },
                "message": {
                    "type": "string",
                    "example": "status bad request"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
