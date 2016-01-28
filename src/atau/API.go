package atau

import (
	"github.com/Knetic/presilo"
)

type API struct {

	Title string `json:"title"`
	Description string `json:"description"`
	BaseURL string `json:"baseUrl"`

	Resources map[string]Resource `json:"resources"`

	parameters []presilo.TypeSchema
	requiredParameters []string
	orderedParameters []string

	schemas map[string]presilo.TypeSchema
}
