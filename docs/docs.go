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
            "name": "MIT License",
            "url": "https://github.com/brajine/metatrader-live/blob/master/LICENSE"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/stats": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "summary": "Provide actual list of connected accounts",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/metatrader.StateData"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/rest/{page}": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "summary": "Provide actual data on connected account",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Account Page name",
                        "name": "page",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/metatrader.Account"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/wss/{page}": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "summary": "Provide actual data on connected account via WebSocket connection",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Account Page name",
                        "name": "page",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/metatrader.Account"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "metatrader.Account": {
            "type": "object",
            "properties": {
                "balance": {
                    "type": "string",
                    "example": "1000.00"
                },
                "clientversion": {
                    "type": "string",
                    "example": "1.0"
                },
                "company": {
                    "type": "string",
                    "example": "My own company"
                },
                "equity": {
                    "type": "string",
                    "example": "1000.0"
                },
                "freemargin": {
                    "type": "string",
                    "example": "1000.0"
                },
                "login": {
                    "type": "string",
                    "example": "010203"
                },
                "margin": {
                    "type": "string",
                    "example": "1000.0"
                },
                "marginlevel": {
                    "type": "string",
                    "example": "100.0"
                },
                "name": {
                    "type": "string",
                    "example": "Alexandre Dumas"
                },
                "orders": {
                    "description": "Ticket is used as Order key",
                    "type": "object",
                    "additionalProperties": {
                        "$ref": "#/definitions/metatrader.Order"
                    }
                },
                "orderscount": {
                    "type": "integer",
                    "example": 3
                },
                "page": {
                    "type": "string",
                    "example": "my-test-page"
                },
                "profittotal": {
                    "type": "string",
                    "example": "0.0"
                },
                "server": {
                    "type": "string",
                    "example": "Metatrader test server"
                },
                "started": {
                    "type": "string",
                    "example": "2020-12-20 23:10:01"
                },
                "updated": {
                    "type": "string",
                    "example": "2020-12-20 23:10:01"
                },
                "updatefreq": {
                    "type": "string",
                    "example": "minute"
                }
            }
        },
        "metatrader.Order": {
            "type": "object",
            "properties": {
                "curvolume": {
                    "type": "string",
                    "example": "0.1"
                },
                "initvolume": {
                    "type": "string",
                    "example": "0.1"
                },
                "priceopen": {
                    "type": "string",
                    "example": "1.13234"
                },
                "pricesl": {
                    "type": "string",
                    "example": "0.0"
                },
                "profit": {
                    "type": "string",
                    "example": "-10.23"
                },
                "sl": {
                    "type": "string",
                    "example": "0.0"
                },
                "swap": {
                    "type": "string",
                    "example": "0.1"
                },
                "symbol": {
                    "type": "string",
                    "example": "EURUSD"
                },
                "timeopen": {
                    "type": "string",
                    "example": "2020-12-20 23:10:01"
                },
                "tp": {
                    "type": "string",
                    "example": "0.0"
                },
                "type": {
                    "type": "string",
                    "example": "OP_BUY"
                }
            }
        },
        "metatrader.StateData": {
            "type": "object",
            "properties": {
                "accounts": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/metatrader.StateEntry"
                    }
                },
                "online": {
                    "type": "integer",
                    "example": 1
                }
            }
        },
        "metatrader.StateEntry": {
            "type": "object",
            "properties": {
                "page": {
                    "type": "string",
                    "example": "my-test-page"
                },
                "started": {
                    "type": "string",
                    "example": "2020-12-20 23:10:01"
                },
                "updateFreq": {
                    "type": "string",
                    "example": "minute"
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
	Host:        "metatrader.live",
	BasePath:    "/api",
	Schemes:     []string{},
	Title:       "Metatrader.live API",
	Description: "Swagger API doc for Metatrader.live.",
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
