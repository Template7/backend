// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Hello"
                ],
                "summary": "Hello Page",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/handler.hello"
                        }
                    }
                }
            }
        },
        "/admin/v1/user": {
            "post": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "v1",
                    "user",
                    "admin"
                ],
                "summary": "Create user",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Access token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "User data",
                        "name": "userData",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/collection.User"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "User object",
                        "schema": {
                            "$ref": "#/definitions/collection.User"
                        }
                    },
                    "400": {
                        "description": "Error object",
                        "schema": {
                            "$ref": "#/definitions/t7Error.Error"
                        }
                    },
                    "401": {
                        "description": "Error object",
                        "schema": {
                            "$ref": "#/definitions/t7Error.Error"
                        }
                    }
                }
            }
        },
        "/admin/v1/users": {
            "delete": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Delete user",
                "responses": {
                    "200": {
                        "description": ""
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/t7Error.Error"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/t7Error.Error"
                        }
                    }
                }
            }
        },
        "/api/v1/users/{UserId}": {
            "get": {
                "tags": [
                    "v1",
                    "user"
                ],
                "summary": "Get user Info",
                "parameters": [
                    {
                        "type": "string",
                        "description": "User ID",
                        "name": "UserId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Access token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/collection.User"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/t7Error.Error"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/t7Error.Error"
                        }
                    }
                }
            },
            "put": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Update user",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Access token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "User basic info",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/collection.UserInfo"
                        }
                    },
                    {
                        "type": "string",
                        "description": "User ID",
                        "name": "UserId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": ""
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/t7Error.Error"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/t7Error.Error"
                        }
                    }
                }
            }
        },
        "/api/v1/users/{UserId}/token": {
            "put": {
                "tags": [
                    "v1",
                    "token"
                ],
                "summary": "Refresh access token",
                "parameters": [
                    {
                        "description": "Token",
                        "name": "token",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/collection.Token"
                        }
                    },
                    {
                        "type": "string",
                        "description": "User ID",
                        "name": "UserId",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Access token",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/collection.Token"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/t7Error.Error"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/t7Error.Error"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "collection.LoginInfo": {
            "type": "object",
            "required": [
                "device"
            ],
            "properties": {
                "channel": {
                    "type": "integer"
                },
                "channel_user_id": {
                    "description": "user id of the channel",
                    "type": "string"
                },
                "device": {
                    "description": "iPhoneN, PixelN, NoteN, ...",
                    "type": "string"
                },
                "os": {
                    "type": "integer"
                }
            }
        },
        "collection.Token": {
            "type": "object",
            "properties": {
                "access_token": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "refresh_token": {
                    "type": "string"
                }
            }
        },
        "collection.User": {
            "type": "object",
            "properties": {
                "basic_info": {
                    "$ref": "#/definitions/collection.UserInfo"
                },
                "email": {
                    "type": "string",
                    "example": "username@mail.com"
                },
                "last_update": {
                    "description": "unix time in second",
                    "type": "integer"
                },
                "login_info": {
                    "$ref": "#/definitions/collection.LoginInfo"
                },
                "mobile": {
                    "description": "+886987654321",
                    "type": "string",
                    "example": "+886987654321"
                },
                "status": {
                    "type": "integer"
                }
            }
        },
        "collection.UserInfo": {
            "type": "object",
            "required": [
                "nick_name"
            ],
            "properties": {
                "Avatar": {
                    "description": "s3 object url",
                    "type": "string"
                },
                "bio": {
                    "type": "string"
                },
                "birthday": {
                    "type": "integer"
                },
                "gender": {
                    "description": "configable?",
                    "type": "integer"
                },
                "hobbies": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "nick_name": {
                    "type": "string"
                },
                "profile_pictures": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "handler.hello": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string",
                    "example": "Hello"
                },
                "timestamp": {
                    "type": "string",
                    "example": "2021-07-24T20:01:25.874565+08:00"
                }
            }
        },
        "t7Error.Error": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "string"
                },
                "detail": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "1.0",
	Host:        "",
	BasePath:    "",
	Schemes:     []string{},
	Title:       "Backend API",
	Description: "API Documentation",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
