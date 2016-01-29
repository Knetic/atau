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

	RawRequestSchema *json.RawMessage `json:"request"`
	RawResponseSchema *json.RawMessage `json:"response"`
	RawParameters map[string]*json.RawMessage `json:"parameters"`
}
