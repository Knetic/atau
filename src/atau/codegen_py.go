package atau

import (
	"strings"
	"github.com/Knetic/presilo"
)

func GeneratePython(api *API, module string)(string, error) {

	var buffer *presilo.BufferedFormatString

	buffer = presilo.NewBufferedFormatString("\t")

	generatePythonImports(api, buffer)
	generatePythonResourceMethods(api, buffer)
	return buffer.String(), nil
}

func generatePythonImports(api *API, buffer *presilo.BufferedFormatString) {

	buffer.Printfln("import urllib2")
	buffer.Printfln("import json")
	buffer.Printfln("")
}

func generatePythonResourceMethods(api *API, buffer *presilo.BufferedFormatString) {

	var fullPath string
	var methodName string

	for resourceName, resource := range api.Resources {
		for methodPrefix, method := range resource.Methods {

			// description doc comment
			buffer.Printf("\n# %s", method.Description)

			// signature
			methodName = presilo.ToSnakeCase(methodPrefix) + "_" + presilo.ToSnakeCase(resourceName)
			generatePythonMethodSignature(methodName, api.optionsSchema, method, buffer)
			buffer.AddIndentation(1)

			// request
			fullPath = resolvePath(api, method)
			fullPath = interpolatePath(api, method, fullPath)
			fullPath = appendQuerystringPath(api, method, fullPath)

			buffer.Printfln("\nhandler = urllib2.HTTPHandler()")

			if(method.RequestSchema != nil) {

				buffer.Printfln("request = urllib2.Request(\"%s\", data=request_contents.to_json())", fullPath)
				buffer.Printfln("request.add_header(\"Content-Type\", \"application/json\")")
			} else {

				buffer.Printfln("request = urllib2.Request(\"%s\")", fullPath)
			}

			buffer.Printfln("request.get_method = lambda: \"%s\"", method.HttpMethod)
			buffer.Printfln("opener = urllib2.build_opener(handler)")
			buffer.Printfln("connection = opener.open(request)")

			buffer.Printf("if connection.code >= 400:")
			buffer.AddIndentation(1)
			buffer.Printf("\nraise Exception(\"Failed to make request: http\" + connection.code + \": \" + connection.read())")
			buffer.AddIndentation(-1)

			if(method.ResponseSchema != nil) {

				buffer.Printf("\nhash = json.loads(connection.read())")
				buffer.Printf("\nreturn %s.deserialize_from(hash)", presilo.ToStrictCamelCase(method.ResponseSchema.GetTitle()))
			}

			buffer.AddIndentation(-1)
			buffer.Printfln("")
		}
	}
}

func generatePythonMethodSignature(methodName string, optionsSchema *presilo.ObjectSchema, method Method, buffer *presilo.BufferedFormatString) {

	var arguments []string

	buffer.Printf("\ndef %s(", methodName)

	// method-specific parameters
	for _, parameterName := range method.Parameters.GetOrderedParameters() {
		arguments = append(arguments, presilo.ToStrictSnakeCase(parameterName))
	}

	// method-specific header parameters
	for _, parameterName := range method.Headers.GetOrderedParameters() {
		arguments = append(arguments, presilo.ToStrictSnakeCase(parameterName))
	}

	if(method.RequestSchema != nil) {
		arguments = append(arguments, "request_contents")
	}
	if(optionsSchema != nil) {
		arguments = append(arguments, "options")
	}

	buffer.Printfln("%s):", strings.Join(arguments, ", "))
}
