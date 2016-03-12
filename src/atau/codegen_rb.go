package atau

import (
	"strings"
	"github.com/Knetic/presilo"
)

func GenerateRB(api *API, module string)(string, error) {

	var buffer *presilo.BufferedFormatString

	buffer = presilo.NewBufferedFormatString("\t")

	generateRBImports(api, buffer)

	buffer.Printfln("module %s", presilo.ToStrictCamelCase(module))
	buffer.AddIndentation(1)

	generateRBResourceMethods(api, buffer)

	buffer.AddIndentation(-1)
	buffer.Printfln("\nend")
	return buffer.String(), nil
}

func generateRBImports(api *API, buffer *presilo.BufferedFormatString) {

	buffer.Printfln("require 'uri'")
	buffer.Printfln("require 'net/http'")
	buffer.Printfln("")
}

func generateRBResourceMethods(api *API, buffer *presilo.BufferedFormatString) {

	var fullPath string
	var methodName string

	for resourceName, resource := range api.Resources {
		for methodPrefix, method := range resource.Methods {

			// description doc comment
			buffer.Printf("\n# %s", method.Description)

			// signature
			methodName = presilo.ToSnakeCase(methodPrefix) + "_" + presilo.ToSnakeCase(resourceName)
			generateRBMethodSignature(methodName, api.optionsSchema, method, buffer)
			buffer.AddIndentation(1)

			// request
			fullPath = resolvePath(api, method)
			fullPath = interpolatePath(api, method, fullPath)
			fullPath = appendQuerystringPath(api, method, fullPath)

			buffer.Printfln("\nuri = URI.parse(\"%s\")", fullPath)
			buffer.Printfln("http = Net::HTTP.new(uri.host, uri.port)")
			buffer.Printfln("request = Net::HTTP::%s.new(uri.path)", presilo.ToJavaCase(strings.ToLower(method.HttpMethod)))

			if(method.RequestSchema != nil) {
				buffer.Printfln("request.body = request_contents.to_json()")
			}

			// headers
			for key, _ := range method.Headers.Parameters {
				buffer.Printfln("request[\"%s\"] = %s", key, presilo.ToStrictSnakeCase(key))
			}

			buffer.Printfln("response = http.request(request)")
			buffer.Printf("if(Integer(response.code) >= 400)")
			buffer.AddIndentation(1)
			buffer.Printf("\nraise StandardError.new(\"Unable to reach resource, returned http\" + response.code)")
			buffer.AddIndentation(-1)
			buffer.Printfln("\nend")
			buffer.Printf("return JSON.parse(response.body)")

			buffer.AddIndentation(-1)
			buffer.Printfln("\nend")
		}
	}
}

func generateRBMethodSignature(methodName string, optionsSchema *presilo.ObjectSchema, method Method, buffer *presilo.BufferedFormatString) {

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

	buffer.Printfln("%s) ", strings.Join(arguments, ", "))
}
