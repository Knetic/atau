package atau

import (
	"strings"
	"path"
	"fmt"
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
			methodName = presilo.ToStrictCamelCase(methodPrefix) + presilo.ToStrictCamelCase(resourceName)
			generateGoMethodSignature(methodName, api.optionsSchema, method, buffer)
			buffer.AddIndentation(1)

			// params
			buffer.Printf("\nvar client http.Client")
			buffer.Printf("\nvar request *http.Request")
			buffer.Printf("\nvar response *http.Response")
			buffer.Printf("\nvar err error")
			if(hasResponse) {
				buffer.Printf("\nvar ret %s", presilo.ToStrictCamelCase(method.ResponseSchema.GetTitle()))
			}

			buffer.Printfln("\n")

			// request
			fullPath = resolvePath(api, method)
			fullPath = generateGoInterpolatedPath(api, method, fullPath)
			buffer.Printfln("request, err = http.NewRequest(\"%s\", \"%s\", nil)", strings.ToUpper(method.HttpMethod), fullPath)
			buffer.Printfln("request.Header.Set(\"Content-Type\", \"application/json\")")
			buffer.Printfln("response, err = client.Do(request)")
			addGoErrCheck(buffer, hasResponse)

			// check for non-2xx
			buffer.Printf("if(response.StatusCode >= 400) {")
			buffer.AddIndentation(1)
			buffer.Printf("\nerrorMsg := fmt.Sprintf(\"Unable to complete request, server returned %%s\", response.Status)")

			if(hasResponse) {
				buffer.Printf("\nreturn ret, errors.New(errorMsg)")
			} else {
				buffer.Printf("\nreturn errors.New(errorMsg)")
			}
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

func generateGoMethodSignature(methodName string, optionsSchema *presilo.ObjectSchema, method Method, buffer *presilo.BufferedFormatString) {

	var arguments []string
	var parameterSchema presilo.TypeSchema
	var parameterType string

	buffer.Printf("\nfunc %s(", methodName)

	// method-specific parameters
	for _, parameterName := range method.Parameters.GetOrderedParameters() {

		parameterSchema = method.Parameters.Parameters[parameterName]
		parameterType = presilo.GenerateGoTypeForSchema(parameterSchema)
		parameterName = presilo.ToStrictJavaCase(parameterName)

		arguments = append(arguments, fmt.Sprintf("%s %s", parameterName, parameterType))
	}

	if(method.RequestSchema != nil) {
		arguments = append(arguments, fmt.Sprintf("requestContents %s", presilo.ToStrictCamelCase(method.RequestSchema.GetTitle())))
	}
	if(optionsSchema != nil) {
		arguments = append(arguments, fmt.Sprintf("options %s", presilo.ToStrictCamelCase(optionsSchema.GetTitle())))
	}

	buffer.Printf("%s) ", strings.Join(arguments, ", "))

	if(method.ResponseSchema != nil) {
		buffer.Printf("(%s, error)", presilo.ToStrictCamelCase(method.ResponseSchema.GetTitle()))
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
	Interpolates querystring parameters for the given [fullPath].
*/
func generateGoInterpolatedPath(api *API, method Method, fullPath string) string {

	var placeholder, replacement string

	for key, _ := range method.Parameters.Parameters {

		placeholder = fmt.Sprintf("{%s}", key)
		replacement = fmt.Sprintf("\"+%s+\"", key)
		fullPath = strings.Replace(fullPath, placeholder, replacement, -1)
	}

	return fullPath
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
