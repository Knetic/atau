package atau

import (
	"github.com/Knetic/presilo"
)

func GenerateCSharp(api *API, module string)(string, error) {

	var buffer *presilo.BufferedFormatString

	buffer = presilo.NewBufferedFormatString("\t")

	generateCSharpImports(api, module, buffer)
	generateCSharpResourceMethods(api, buffer)
	return buffer.String(), nil
}

func generateCSharpImports(api *API, module string, buffer *presilo.BufferedFormatString) {

}

func generateCSharpResourceMethods(api *API, buffer *presilo.BufferedFormatString) {

	buffer.Printfln("using System;")
	buffer.Printfln("using System.Text.RegularExpressions;")
	buffer.Printfln("using System.Runtime.Serialization;")
}
