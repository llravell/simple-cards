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
        "/api/modules/": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "modules"
                ],
                "summary": "Get all user's modules",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/entity.Module"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "modules"
                ],
                "summary": "Create new module",
                "parameters": [
                    {
                        "description": "Module params",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/modules.createOrUpdateModuleRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/entity.Module"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/api/modules/{module_uuid}/": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "modules"
                ],
                "summary": "Get module with cards",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Module UUID",
                        "name": "module_uuid",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/entity.ModuleWithCards"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            },
            "put": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "modules"
                ],
                "summary": "Update module",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Module UUID",
                        "name": "module_uuid",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Module params",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/modules.createOrUpdateModuleRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/entity.Module"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "tags": [
                    "modules"
                ],
                "summary": "Delete module",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Module UUID",
                        "name": "module_uuid",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "202": {
                        "description": "Accepted"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/api/user/login": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Verify user creds and login",
                "parameters": [
                    {
                        "description": "User creds",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/auth.authRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "invalid data"
                    },
                    "401": {
                        "description": "verification failed"
                    },
                    "500": {
                        "description": "token building error"
                    }
                }
            }
        },
        "/api/user/register": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "auth"
                ],
                "summary": "Register new user",
                "parameters": [
                    {
                        "description": "User creds",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/auth.authRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "invalid data"
                    },
                    "409": {
                        "description": "user with same login already exists"
                    },
                    "500": {
                        "description": "token building error"
                    }
                }
            }
        },
        "/ping": {
            "get": {
                "tags": [
                    "health"
                ],
                "summary": "Check database connection",
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        }
    },
    "definitions": {
        "auth.authRequest": {
            "type": "object",
            "properties": {
                "login": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "entity.Card": {
            "type": "object",
            "properties": {
                "meaning": {
                    "type": "string"
                },
                "module_uuid": {
                    "type": "string"
                },
                "term": {
                    "type": "string"
                },
                "uuid": {
                    "type": "string"
                }
            }
        },
        "entity.Module": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string"
                },
                "user_uuid": {
                    "type": "string"
                },
                "uuid": {
                    "type": "string"
                }
            }
        },
        "entity.ModuleWithCards": {
            "type": "object",
            "properties": {
                "cards": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/entity.Card"
                    }
                },
                "name": {
                    "type": "string"
                },
                "user_uuid": {
                    "type": "string"
                },
                "uuid": {
                    "type": "string"
                }
            }
        },
        "modules.createOrUpdateModuleRequest": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Simple Cards API",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
