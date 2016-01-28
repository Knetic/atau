package atau

import (
	"encoding/json"
	"github.com/Knetic/presilo"
)

type Method struct {

	HttpMethod string `json:"httpMethod`
	Path string `json:"path"`
	Description string `json:"description"`

	RequestSchema presilo.TypeSchema
	ResponseSchema presilo.TypeSchema
	Parameters ParameterList

	RawRequestSchema *json.RawMessage
	RawResponseSchema *json.RawMessage
}
