package atau

import (
	"strings"
	"path"
	"github.com/Knetic/presilo"
)

func GenerateGo(api *API, module string)(string, error) {

	var buffer *presilo.BufferedFormatString

	buffer = presilo.NewBufferedFormatString("\t")

	generateGoImports(api, module, buffer)
	generateGoResourceMethods(api, buffer)
	return buffer.String(), nil
}

/*
	Generates all methods for all defined api resources.
*/
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
			buffer.Printf("\nvar client http.Client")
			buffer.Printf("\nvar request *http.Request")
			buffer.Printf("\nvar response *http.Response")
			buffer.Printf("\nvar err error")
			if(hasResponse) {
				buffer.Printf("\nvar ret %s", presilo.ToCamelCase(method.ResponseSchema.GetTitle()))
			}

			buffer.Printfln("\n")

			// request
			fullPath = resolvePath(api, method)
			buffer.Printfln("request, err = http.NewRequest(\"%s\", \"%s\", nil)", strings.ToUpper(method.HttpMethod), fullPath)
			buffer.Printfln("request.Header.Set(\"Content-Type\", \"application/json\")")
			buffer.Printfln("response, err = client.Do(request)")
			addGoErrCheck(buffer, hasResponse)

			// check for non-2xx
			buffer.Printf("if(response.StatusCode >= 400) {")
			buffer.AddIndentation(1)
			buffer.Printf("\nerrorMsg := fmt.Sprintf(\"Unable to complete request, server returned %%s\", response.Status)")
			buffer.Printf("\nreturn ret, errors.New(errorMsg)")
			buffer.AddIndentation(-1)
			buffer.Printf("\n}\n")

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

func generateGoImports(api *API, module string, buffer *presilo.BufferedFormatString) {

	buffer.Printfln("package %s", module)
	buffer.Printf("\nimport (")
	buffer.AddIndentation(1)
	buffer.Printf("\n\"fmt\"")
	buffer.Printf("\n\"net/http\"")
	buffer.Printf("\n\"encoding/json\"")
	buffer.Printf("\n\"errors\"")
	buffer.AddIndentation(-1)
	buffer.Printfln(")")
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
