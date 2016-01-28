package atau

import (
	"github.com/Knetic/presilo"
)

func GenerateGo(api *API, module string)(string, error) {

	var buffer *presilo.BufferedFormatString

	buffer = presilo.NewBufferedFormatString("\t")

	// basics
	buffer.Printfln("package %s", module)
	buffer.Printfln("\nimport \"net/http\"")

	generateGoResourceMethods(api, buffer)
	return buffer.String(), nil
}

func generateGoResourceMethods(api *API, buffer *presilo.BufferedFormatString) {

	
}
