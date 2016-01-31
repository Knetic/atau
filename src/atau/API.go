package atau

import (
	"github.com/Knetic/presilo"
)

type API struct {

	Name string `json:"name"`
	Title string `json:"title"`
	Description string `json:"description"`
	BaseURL string `json:"baseUrl"`

	Resources map[string]Resource `json:"resources"`
	Parameters ParameterList

	schemas map[string]presilo.TypeSchema
	schemaContext *presilo.SchemaParseContext
}
