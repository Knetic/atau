package atau

import (
	"strings"
	"path"
	"github.com/Knetic/presilo"
)

func GenerateGo(api *API, module string)(string, error) {

	var buffer *presilo.BufferedFormatString

	buffer = presilo.NewBufferedFormatString("\t")

	// basics
	buffer.Printfln("package %s", module)
	buffer.Printfln("\nimport (\"net/http\"\n\"encoding/json\")")

	generateGoResourceMethods(api, buffer)
	return buffer.String(), nil
}

func generateGoResourceMethods(api *API, buffer *presilo.BufferedFormatString) {

	var fullPath string
	var methodName string
	var hasResponse bool

	for resourceName, resource := range api.Resources {
		for methodPrefix, method := range resource.Methods {

			hasResponse = method.ResponseSchema != nil

			// description doc comment
			buffer.Printf("\n/*")
			buffer.AddIndentation(1)
			buffer.Printf("\n%s", method.Description)
			buffer.AddIndentation(-1)
			buffer.Printf("\n*/")

			// signature
			methodName = presilo.ToCamelCase(methodPrefix) + presilo.ToCamelCase(resourceName)
			generateGoMethodSignature(methodName, method.RequestSchema, method.ResponseSchema, buffer)
			buffer.AddIndentation(1)

			// params
			buffer.Printfln("\nvar client http.Client")
			buffer.Printfln("\nvar request *http.Request")
			buffer.Printfln("\nvar response *http.Response")
			buffer.Printfln("\nvar err error")
			if(hasResponse) {
				buffer.Printfln("var ret %s", presilo.ToCamelCase(method.ResponseSchema.GetTitle()))
			}

			buffer.Printfln("")

			// request
			fullPath = resolvePath(api, method)
			buffer.Printfln("request, err = http.NewRequest(\"%s\", \"%s\", nil)", strings.ToUpper(method.HttpMethod), fullPath)
			buffer.Printfln("request.Header.Set(\"Content-Type\", \"application/json\")")
			buffer.Printfln("response, err = client.Do(request)")
			addGoErrCheck(buffer, hasResponse)

			// marshal?
			if(hasResponse) {

				buffer.Printfln("\ndecoder := json.NewDecoder(response.Body)")
				buffer.Printfln("err = decoder.Decode(&ret)")
				addGoErrCheck(buffer, true)
			}

			// close up.
			if(hasResponse) {
				buffer.Printf("\nreturn ret, nil")
			} else {
				buffer.Printf("\nreturn nil")
			}

			buffer.AddIndentation(-1)
			buffer.Printfln("\n}")
		}
	}
}

func generateGoMethodSignature(methodName string, request presilo.TypeSchema, response presilo.TypeSchema, buffer *presilo.BufferedFormatString) {

	buffer.Printf("\nfunc %s(", methodName)
	if(request != nil) {
		buffer.Printf("params %s", presilo.ToCamelCase(request.GetTitle()))
	}

	buffer.Printf(") ")

	if(response != nil) {
		buffer.Printf("(%s, error)", presilo.ToCamelCase(response.GetTitle()))
	} else {
		buffer.Printf("error")
	}

	buffer.Printf("{")
}

/*
	Paths can have parameter placeholders baked into them.
	This provides a full path to a specific resource given an API's base path, the method's path,
	and interpolated with the correct variable names for all parameters.
*/
func resolvePath(api *API, method Method) string {
	return path.Join(api.BaseURL, method.Path)
}

func addGoErrCheck(buffer *presilo.BufferedFormatString, includeRet bool) {

	buffer.Printf("if(err != nil) {")
	buffer.AddIndentation(1)

	if(includeRet) {
		buffer.Printf("\nreturn ret, err")
	} else {
		buffer.Printf("\nreturn err")
	}
	buffer.AddIndentation(-1)
	buffer.Printf("\n}\n")
}
