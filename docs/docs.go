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
        "/context/all": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Context"
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/dtos.PunqContext"
                            }
                        }
                    }
                }
            }
        },
        "/user": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID of the user",
                        "name": "userid",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dtos.PunqUser"
                        }
                    }
                }
            },
            "post": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "parameters": [
                    {
                        "description": "PunqUser",
                        "name": "body",
                        "in": "body",
                        "schema": {
                            "$ref": "#/definitions/dtos.PunqUser"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dtos.PunqUser"
                        }
                    }
                }
            },
            "delete": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "parameters": [
                    {
                        "type": "string",
                        "description": "ID of the user",
                        "name": "userid",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dtos.PunqUser"
                        }
                    }
                }
            },
            "patch": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "parameters": [
                    {
                        "description": "PunqUser",
                        "name": "body",
                        "in": "body",
                        "schema": {
                            "$ref": "#/definitions/dtos.PunqUser"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dtos.PunqUser"
                        }
                    }
                }
            }
        },
        "/user/all": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "User"
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/dtos.PunqUser"
                            }
                        }
                    }
                }
            }
        },
        "/version": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Misc"
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/structs.Version"
                        }
                    }
                }
            }
        },
        "/workload/available-resources": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "General"
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/workload/templates": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "General"
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/kubernetes.K8sNewWorkload"
                            }
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "dtos.AccessLevel": {
            "type": "integer",
            "enum": [
                0,
                1,
                2
            ],
            "x-enum-varnames": [
                "READER",
                "USER",
                "ADMIN"
            ]
        },
        "dtos.PunqAccess": {
            "type": "object",
            "required": [
                "level",
                "userId"
            ],
            "properties": {
                "level": {
                    "$ref": "#/definitions/dtos.AccessLevel"
                },
                "userId": {
                    "type": "string"
                }
            }
        },
        "dtos.PunqContext": {
            "type": "object",
            "required": [
                "access",
                "contextBase64",
                "id",
                "name"
            ],
            "properties": {
                "access": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/dtos.PunqAccess"
                    }
                },
                "contextBase64": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "dtos.PunqUser": {
            "type": "object",
            "required": [
                "accessLevel",
                "createdAt",
                "displayName",
                "email",
                "id",
                "password"
            ],
            "properties": {
                "accessLevel": {
                    "$ref": "#/definitions/dtos.AccessLevel"
                },
                "createdAt": {
                    "type": "string"
                },
                "displayName": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "kubernetes.K8sNewWorkload": {
            "type": "object",
            "properties": {
                "description": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "yamlString": {
                    "type": "string"
                }
            }
        },
        "structs.Version": {
            "type": "object",
            "properties": {
                "branch": {
                    "type": "string"
                },
                "buildTimestamp": {
                    "type": "string"
                },
                "gitCommitHash": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "version": {
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
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
