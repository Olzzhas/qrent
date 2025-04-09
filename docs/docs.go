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
        "/organizations": {
            "get": {
                "description": "Возвращает все организации",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "organizations"
                ],
                "summary": "Возвращает список организаций",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.OrganizationListResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            },
            "post": {
                "description": "Создаёт организацию с заданными name и location",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "organizations"
                ],
                "summary": "Создаёт новую организацию",
                "parameters": [
                    {
                        "description": "Organization Data",
                        "name": "organization",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.CreateOrganizationRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/main.OrganizationResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/organizations/{id}": {
            "get": {
                "description": "Возвращает организацию по переданному идентификатору",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "organizations"
                ],
                "summary": "Получает организацию по ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Organization ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.OrganizationResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            },
            "put": {
                "description": "Обновляет данные организации. Обновляются только переданные поля.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "organizations"
                ],
                "summary": "Обновляет организацию по ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Organization ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Organization Data",
                        "name": "organization",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.UpdateOrganizationRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.OrganizationResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            },
            "delete": {
                "description": "Удаляет организацию с заданным идентификатором",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "organizations"
                ],
                "summary": "Удаляет организацию по ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Organization ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.MessageResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/powerbanks": {
            "get": {
                "description": "Возвращает все повербанки",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "powerbanks"
                ],
                "summary": "Возвращает список повербанков",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.PowerbankListResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            },
            "post": {
                "description": "Создаёт повербанк с заданными current_station_id и status",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "powerbanks"
                ],
                "summary": "Создаёт новый повербанк",
                "parameters": [
                    {
                        "description": "Powerbank Data",
                        "name": "powerbank",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.CreatePowerbankRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/main.PowerbankResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/powerbanks/{id}": {
            "get": {
                "description": "Возвращает повербанк по переданному идентификатору",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "powerbanks"
                ],
                "summary": "Получает повербанк по ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Powerbank ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.PowerbankResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            },
            "put": {
                "description": "Обновляет данные повербанка. Обновляются только переданные поля.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "powerbanks"
                ],
                "summary": "Обновляет повербанк по ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Powerbank ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Powerbank Data",
                        "name": "powerbank",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.UpdatePowerbankRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.PowerbankResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            },
            "delete": {
                "description": "Удаляет повербанк с заданным идентификатором",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "powerbanks"
                ],
                "summary": "Удаляет повербанк по ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Powerbank ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.MessageResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/stations": {
            "get": {
                "description": "Возвращает все станции",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "stations"
                ],
                "summary": "Возвращает список станций",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.StationListResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            },
            "post": {
                "description": "Создаёт станцию, привязанную к организации (org_id)",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "stations"
                ],
                "summary": "Создаёт новую станцию",
                "parameters": [
                    {
                        "description": "Station Data",
                        "name": "station",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.CreateStationRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/main.StationResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        },
        "/stations/{id}": {
            "get": {
                "description": "Возвращает станцию по идентификатору",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "stations"
                ],
                "summary": "Получает станцию по ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Station ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.StationResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            },
            "put": {
                "description": "Обновляет данные станции. Обновляются только переданные поля.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "stations"
                ],
                "summary": "Обновляет станцию по ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Station ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Station Data",
                        "name": "station",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/main.UpdateStationRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.StationResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "422": {
                        "description": "Unprocessable Entity",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            },
            "delete": {
                "description": "Удаляет станцию с заданным идентификатором",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "stations"
                ],
                "summary": "Удаляет станцию по ID",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "Station ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/main.MessageResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "object",
                            "additionalProperties": {
                                "type": "string"
                            }
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "data.Organization": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "location": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "data.Powerbank": {
            "type": "object",
            "properties": {
                "current_station_id": {
                    "type": "integer"
                },
                "id": {
                    "type": "integer"
                },
                "status": {
                    "$ref": "#/definitions/data.PowerbankStatus"
                }
            }
        },
        "data.PowerbankStatus": {
            "type": "string",
            "enum": [
                "rented",
                "available",
                "charging"
            ],
            "x-enum-varnames": [
                "PowerbankStatusRented",
                "PowerbankStatusAvailable",
                "PowerbankStatusCharging"
            ]
        },
        "data.Station": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "org_id": {
                    "type": "integer"
                }
            }
        },
        "main.CreateOrganizationRequest": {
            "type": "object",
            "properties": {
                "location": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "main.CreatePowerbankRequest": {
            "type": "object",
            "properties": {
                "current_station_id": {
                    "type": "integer"
                },
                "status": {
                    "type": "string"
                }
            }
        },
        "main.CreateStationRequest": {
            "type": "object",
            "properties": {
                "org_id": {
                    "type": "integer"
                }
            }
        },
        "main.MessageResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "main.OrganizationListResponse": {
            "type": "object",
            "properties": {
                "organizations": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/data.Organization"
                    }
                }
            }
        },
        "main.OrganizationResponse": {
            "type": "object",
            "properties": {
                "organization": {
                    "$ref": "#/definitions/data.Organization"
                }
            }
        },
        "main.PowerbankListResponse": {
            "type": "object",
            "properties": {
                "powerbanks": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/data.Powerbank"
                    }
                }
            }
        },
        "main.PowerbankResponse": {
            "type": "object",
            "properties": {
                "powerbank": {
                    "$ref": "#/definitions/data.Powerbank"
                }
            }
        },
        "main.StationListResponse": {
            "type": "object",
            "properties": {
                "stations": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/data.Station"
                    }
                }
            }
        },
        "main.StationResponse": {
            "type": "object",
            "properties": {
                "station": {
                    "$ref": "#/definitions/data.Station"
                }
            }
        },
        "main.UpdateOrganizationRequest": {
            "type": "object",
            "properties": {
                "location": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "main.UpdatePowerbankRequest": {
            "type": "object",
            "properties": {
                "current_station_id": {
                    "type": "integer"
                },
                "status": {
                    "type": "string"
                }
            }
        },
        "main.UpdateStationRequest": {
            "type": "object",
            "properties": {
                "org_id": {
                    "type": "integer"
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
