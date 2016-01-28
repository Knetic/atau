package atau

import (
	"github.com/Knetic/presilo"
)

func GenerateGo(api *API, module string)(string, error) {

	var buffer *presilo.BufferedFormatString

	buffer = presilo.NewBufferedFormatString("\t")
	buffer.Printfln("package %s", module)
	generateGoImports(api, buffer)

	return buffer.String(), nil
}

func generateGoImports(api *API, buffer *presilo.BufferedFormatString) {
	buffer.Printfln("\nimport \"net/http\"")
}
