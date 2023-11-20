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
                    "V1",
                    "User"
                ],
                "summary": "Create user",
                "parameters": [
                    {
                        "description": "Request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/types.HttpCreateUserReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Response",
                        "schema": {
                            "$ref": "#/definitions/types.HttpLoginResp"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/types.HttpRespError"
                        }
                    }
                }
            }
        },
        "/api/v1/login/native": {
            "post": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "V1",
                    "login"
                ],
                "summary": "Native login",
                "parameters": [
                    {
                        "description": "Request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/types.HttpLoginReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Response",
                        "schema": {
                            "$ref": "#/definitions/types.HttpLoginResp"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/types.HttpRespError"
                        }
                    }
                }
            }
        },
        "/api/v1/user/info": {
            "get": {
                "tags": [
                    "V1",
                    "User"
                ],
                "summary": "Get user Info",
                "responses": {
                    "200": {
                        "description": "Response",
                        "schema": {
                            "$ref": "#/definitions/types.HttpUserInfoResp"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/types.HttpRespError"
                        }
                    }
                }
            },
            "put": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "V1",
                    "User"
                ],
                "summary": "Update user info",
                "parameters": [
                    {
                        "description": "Request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/types.HttpUpdateUserInfoReq"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Response",
                        "schema": {
                            "$ref": "#/definitions/types.HttpRespBase"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/types.HttpRespError"
                        }
                    }
                }
            }
        },
        "/api/v1/user/wallets": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "V1",
                    "User",
                    "Wallet"
                ],
                "summary": "Get user wallets",
                "responses": {
                    "200": {
                        "description": "Response",
                        "schema": {
                            "$ref": "#/definitions/types.HttpGetUserWalletsResp"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/types.HttpRespError"
                        }
                    }
                }
            }
        },
        "/api/v1/wallets/{walletId}": {
            "get": {
                "tags": [
                    "V1",
                    "Wallet"
                ],
                "summary": "Get wallet",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Wallet ID",
                        "name": "walletId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Response",
                        "schema": {
                            "$ref": "#/definitions/types.HttpGetWalletResp"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/types.HttpRespError"
                        }
                    }
                }
            }
        },
        "/api/v1/wallets/{walletId}/withdraw": {
            "post": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "V1",
                    "Wallet"
                ],
                "summary": "Wallet withdraw",
                "parameters": [
                    {
                        "description": "Request",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/types.HttpWalletWithdrawReq"
                        }
                    },
                    {
                        "type": "string",
                        "description": "Wallet ID",
                        "name": "walletId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Response",
                        "schema": {
                            "$ref": "#/definitions/types.HttpRespBase"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/types.HttpRespError"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
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
        "types.HttpCreateUserReq": {
            "type": "object",
            "required": [
                "password",
                "username"
            ],
            "properties": {
                "email": {
                    "type": "string",
                    "example": "example@email.com"
                },
                "nickname": {
                    "type": "string",
                    "example": "nickname"
                },
                "password": {
                    "type": "string",
                    "example": "password"
                },
                "role": {
                    "type": "string",
                    "enum": [
                        "admin"
                    ],
                    "example": "user"
                },
                "username": {
                    "type": "string",
                    "example": "username"
                }
            }
        },
        "types.HttpGetUserWalletsResp": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer",
                    "example": 3000
                },
                "data": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/types.HttpGetUserWalletsRespData"
                    }
                },
                "message": {
                    "type": "string",
                    "example": "ok"
                },
                "requestId": {
                    "type": "string",
                    "example": "b8974256-1f17-477f-8638-c7ebbac656d7"
                }
            }
        },
        "types.HttpGetUserWalletsRespData": {
            "type": "object",
            "properties": {
                "balances": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/types.HttpGetUserWalletsRespDataBalance"
                    }
                },
                "id": {
                    "type": "string",
                    "example": "af68a360-d035-469c-8ae9-a8640c2ffd19"
                }
            }
        },
        "types.HttpGetUserWalletsRespDataBalance": {
            "type": "object",
            "properties": {
                "amount": {
                    "type": "string",
                    "example": "100"
                },
                "currency": {
                    "type": "string",
                    "example": "usd"
                }
            }
        },
        "types.HttpGetWalletResp": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer",
                    "example": 3000
                },
                "data": {
                    "$ref": "#/definitions/types.HttpGetUserWalletsRespData"
                },
                "message": {
                    "type": "string",
                    "example": "ok"
                },
                "requestId": {
                    "type": "string",
                    "example": "b8974256-1f17-477f-8638-c7ebbac656d7"
                }
            }
        },
        "types.HttpLoginReq": {
            "type": "object",
            "required": [
                "password",
                "username"
            ],
            "properties": {
                "password": {
                    "type": "string",
                    "example": "password"
                },
                "username": {
                    "type": "string",
                    "example": "username"
                }
            }
        },
        "types.HttpLoginResp": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer",
                    "example": 3000
                },
                "data": {
                    "type": "object",
                    "properties": {
                        "token": {
                            "type": "string",
                            "example": "70596484-67d3-46bd-94bf-08f7c9fb7ac1"
                        }
                    }
                },
                "message": {
                    "type": "string",
                    "example": "ok"
                },
                "requestId": {
                    "type": "string",
                    "example": "b8974256-1f17-477f-8638-c7ebbac656d7"
                }
            }
        },
        "types.HttpRespBase": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer",
                    "example": 3000
                },
                "message": {
                    "type": "string",
                    "example": "ok"
                },
                "requestId": {
                    "type": "string",
                    "example": "b8974256-1f17-477f-8638-c7ebbac656d7"
                }
            }
        },
        "types.HttpRespError": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer",
                    "example": 1024
                },
                "message": {
                    "type": "string",
                    "example": "token unauthorized"
                },
                "requestId": {
                    "type": "string",
                    "example": "27c0a70e-59ab-4a94-872c-5f014aaa047f"
                }
            }
        },
        "types.HttpUpdateUserInfoReq": {
            "type": "object",
            "required": [
                "nickname"
            ],
            "properties": {
                "nickname": {
                    "type": "string",
                    "example": "nickname"
                }
            }
        },
        "types.HttpUserInfoResp": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "integer",
                    "example": 3000
                },
                "data": {
                    "type": "object",
                    "properties": {
                        "email": {
                            "type": "string",
                            "example": "example@email.com"
                        },
                        "nickname": {
                            "type": "string",
                            "example": "example"
                        },
                        "role": {
                            "type": "string",
                            "example": "user"
                        },
                        "status": {
                            "type": "string",
                            "example": "activated"
                        },
                        "userId": {
                            "type": "string",
                            "example": "userId001"
                        }
                    }
                },
                "message": {
                    "type": "string",
                    "example": "ok"
                },
                "requestId": {
                    "type": "string",
                    "example": "b8974256-1f17-477f-8638-c7ebbac656d7"
                }
            }
        },
        "types.HttpWalletWithdrawReq": {
            "type": "object",
            "required": [
                "amount"
            ],
            "properties": {
                "amount": {
                    "type": "integer",
                    "example": 100
                },
                "currency": {
                    "type": "string",
                    "example": "usd"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "Backend API",
	Description:      "API Documentation",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
