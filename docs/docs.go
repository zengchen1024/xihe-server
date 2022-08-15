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
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/": {
            "get": {
                "description": "callback of authentication by authing",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "Login"
                ],
                "parameters": [
                    {
                        "type": "string",
                        "description": "authing code",
                        "name": "code",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/app.UserDTO"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "system_error"
                        }
                    },
                    "501": {
                        "description": "Not Implemented",
                        "schema": {
                            "type": "duplicate_creating"
                        }
                    }
                }
            }
        },
        "/v1/dataset": {
            "post": {
                "description": "create dataset",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "Dataset"
                ],
                "summary": "Create",
                "parameters": [
                    {
                        "description": "body of creating dataset",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controller.datasetCreateRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/app.DatasetDTO"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "bad_request_param"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "duplicate_creating"
                        }
                    }
                }
            }
        },
        "/v1/dataset/{owner}": {
            "get": {
                "description": "list dataset",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Dataset"
                ],
                "summary": "List",
                "responses": {}
            }
        },
        "/v1/dataset/{owner}/{id}": {
            "get": {
                "description": "get dataset",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "Dataset"
                ],
                "summary": "Get",
                "parameters": [
                    {
                        "type": "string",
                        "description": "id of dataset",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/app.DatasetDTO"
                        }
                    }
                }
            }
        },
        "/v1/model": {
            "post": {
                "description": "create model",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "Model"
                ],
                "summary": "Create",
                "parameters": [
                    {
                        "description": "body of creating model",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controller.modelCreateRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/app.ModelDTO"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "bad_request_param"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "duplicate_creating"
                        }
                    }
                }
            }
        },
        "/v1/model/{owner}": {
            "get": {
                "description": "list model",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Model"
                ],
                "summary": "List",
                "responses": {}
            }
        },
        "/v1/model/{owner}/{id}": {
            "get": {
                "description": "get model",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "Model"
                ],
                "summary": "Get",
                "parameters": [
                    {
                        "type": "string",
                        "description": "id of model",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/app.ModelDTO"
                        }
                    }
                }
            }
        },
        "/v1/project": {
            "post": {
                "description": "create project",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Project"
                ],
                "summary": "Create",
                "parameters": [
                    {
                        "description": "body of creating project",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controller.projectCreateRequest"
                        }
                    }
                ],
                "responses": {}
            }
        },
        "/v1/project/{owner}": {
            "get": {
                "description": "list project",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Project"
                ],
                "summary": "List",
                "responses": {}
            }
        },
        "/v1/project/{owner}/{id}": {
            "get": {
                "description": "get project",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Project"
                ],
                "summary": "Get",
                "parameters": [
                    {
                        "type": "string",
                        "description": "id of project",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {}
            },
            "put": {
                "description": "update project",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Project"
                ],
                "summary": "Update",
                "parameters": [
                    {
                        "type": "string",
                        "description": "id of project",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "body of updating project",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controller.projectUpdateRequest"
                        }
                    }
                ],
                "responses": {}
            },
            "post": {
                "description": "fork project",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Project"
                ],
                "summary": "Fork",
                "parameters": [
                    {
                        "type": "string",
                        "description": "id of project",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {}
            }
        },
        "/v1/user": {
            "get": {
                "description": "get user",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Get",
                "parameters": [
                    {
                        "type": "string",
                        "description": "account",
                        "name": "account",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/app.UserDTO"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "bad_request_param"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "resource_not_exists"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "system_error"
                        }
                    }
                }
            },
            "put": {
                "description": "update user basic info",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Update",
                "responses": {}
            },
            "post": {
                "description": "create user",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "summary": "Create",
                "parameters": [
                    {
                        "description": "body of creating user",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controller.userCreateRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/app.UserDTO"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "bad_request_param"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "duplicate_creating"
                        }
                    }
                }
            }
        },
        "/v1/user/following": {
            "get": {
                "description": "list followings",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "Following"
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/app.FollowDTO"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "system_error"
                        }
                    }
                }
            },
            "post": {
                "description": "add a following",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "Following"
                ],
                "parameters": [
                    {
                        "description": "body of creating following",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/controller.followingCreateRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": ""
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "bad_request_body"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "bad_request_param"
                        }
                    },
                    "402": {
                        "description": "Payment Required",
                        "schema": {
                            "type": "not_allowed"
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "type": "resource_not_exists"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "duplicate_creating"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "system_error"
                        }
                    }
                }
            }
        },
        "/v1/user/following/{account}": {
            "delete": {
                "description": "remove a following",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "Following"
                ],
                "parameters": [
                    {
                        "type": "string",
                        "description": "the account of following",
                        "name": "account",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": ""
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "bad_request_param"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "not_allowed"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "system_error"
                        }
                    }
                }
            }
        },
        "/{account}": {
            "get": {
                "description": "get info of login",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "Login"
                ],
                "parameters": [
                    {
                        "type": "string",
                        "description": "account",
                        "name": "account",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/app.LoginDTO"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "bad_request_param"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "type": "not_allowed"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "system_error"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "app.DatasetDTO": {
            "type": "object",
            "properties": {
                "desc": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "owner": {
                    "type": "string"
                },
                "protocol": {
                    "type": "string"
                },
                "repo_id": {
                    "type": "string"
                },
                "repo_type": {
                    "type": "string"
                },
                "tags": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "app.FollowDTO": {
            "type": "object",
            "properties": {
                "account": {
                    "type": "string"
                },
                "avatar_id": {
                    "type": "string"
                },
                "bio": {
                    "type": "string"
                }
            }
        },
        "app.LoginDTO": {
            "type": "object",
            "properties": {
                "info": {
                    "type": "string"
                }
            }
        },
        "app.ModelDTO": {
            "type": "object",
            "properties": {
                "desc": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "owner": {
                    "type": "string"
                },
                "protocol": {
                    "type": "string"
                },
                "repo_id": {
                    "type": "string"
                },
                "repo_type": {
                    "type": "string"
                },
                "tags": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                }
            }
        },
        "app.UserDTO": {
            "type": "object",
            "properties": {
                "account": {
                    "type": "string"
                },
                "avatar_id": {
                    "type": "string"
                },
                "bio": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "follower_count": {
                    "type": "integer"
                },
                "following_count": {
                    "type": "integer"
                },
                "id": {
                    "type": "string"
                }
            }
        },
        "controller.datasetCreateRequest": {
            "type": "object",
            "properties": {
                "desc": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "owner": {
                    "type": "string"
                },
                "protocol": {
                    "type": "string"
                },
                "repo_type": {
                    "type": "string"
                }
            }
        },
        "controller.followingCreateRequest": {
            "type": "object",
            "properties": {
                "account": {
                    "type": "string"
                }
            }
        },
        "controller.modelCreateRequest": {
            "type": "object",
            "properties": {
                "desc": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "owner": {
                    "type": "string"
                },
                "protocol": {
                    "type": "string"
                },
                "repo_type": {
                    "type": "string"
                }
            }
        },
        "controller.projectCreateRequest": {
            "type": "object",
            "properties": {
                "cover_id": {
                    "type": "string"
                },
                "desc": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "owner": {
                    "type": "string"
                },
                "protocol": {
                    "type": "string"
                },
                "repo_type": {
                    "type": "string"
                },
                "training": {
                    "type": "string"
                },
                "type": {
                    "type": "string"
                }
            }
        },
        "controller.projectUpdateRequest": {
            "type": "object",
            "properties": {
                "cover_id": {
                    "type": "string"
                },
                "desc": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "tags": {
                    "description": "json [] will be converted to []string",
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "type": {
                    "type": "string"
                }
            }
        },
        "controller.userCreateRequest": {
            "type": "object",
            "properties": {
                "account": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
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
