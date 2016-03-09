package atau

import (
	"strings"
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
			fullPath = interpolateGoPath(api, method, fullPath)

			// body, if applicable.
			if(method.RequestSchema != nil) {
				buffer.Printfln("marshalledBody, err := json.Marshal(requestContents)")
				addGoErrCheck(buffer, hasResponse)

				buffer.Printfln("request, err = http.NewRequest(\"%s\", %s, bytes.NewReader(marshalledBody))", strings.ToUpper(method.HttpMethod), fullPath)
			} else {
				buffer.Printfln("request, err = http.NewRequest(\"%s\", %s, nil)", strings.ToUpper(method.HttpMethod), fullPath)
			}

			buffer.Printfln("request.Header.Set(\"Content-Type\", \"application/json\")")
			for key, _ := range method.Headers.Parameters {
				buffer.Printfln("request.Header.Set(\"%s\", fmt.Sprintf(\"%%v\", %s))", key, presilo.ToStrictJavaCase(key))
			}

			// make request
			buffer.Printfln("response, err = client.Do(request)")
			addGoErrCheck(buffer, hasResponse)

			// check for non-2xx
			buffer.Printf("if(response.StatusCode >= 400) {")
			buffer.AddIndentation(1)
			buffer.Printf("\nresponseBodyString, _ := ioutil.ReadAll(response.Body)")
			buffer.Printf("\nerrorMsg := fmt.Sprintf(\"Unable to complete request, server returned %%s: %%s\", response.Status, responseBodyString)")

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

	// method-specific header parameters
	for _, parameterName := range method.Headers.GetOrderedParameters() {

		parameterSchema = method.Headers.Parameters[parameterName]
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

	var modules []string

	buffer.Printfln("package %s", module)

	modules = []string{"fmt", "net/http", "encoding/json", "errors", "io/ioutil"}

	// we only import "bytes" if we marshal a request
	ModulesDefined:
	for _, resource := range api.Resources {
		for _, method := range resource.Methods {
			if(method.RequestSchema != nil) {

				modules = append(modules, "bytes")
				break ModulesDefined
			}
		}
	}

	buffer.Printf("\nimport (")
	buffer.AddIndentation(1)

	for _, module := range modules {
		buffer.Printf("\n\"%s\"", module)
	}

	buffer.AddIndentation(-1)
	buffer.Printfln("\n)")
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

func interpolateGoPath(api *API, method Method, fullPath string) string {

	var querystrings []string
	var queryKeys []string
	var placeholder, replacement string

	// replace parameters as referenced in paths
	for _, parameter := range method.PathParameters {

		placeholder = fmt.Sprintf("{%s}", parameter)
		replacement = fmt.Sprintf("%%v")
		fullPath = strings.Replace(fullPath, placeholder, replacement, -1)
	}

	// set querystring
	for _, key := range method.QueryParameters {
		querystrings = append(querystrings, fmt.Sprintf("%s=%%v", key))
		queryKeys = append(queryKeys, presilo.ToStrictJavaCase(key))
	}

	if(len(querystrings) > 0) {
		fullPath = fullPath + "?" + strings.Join(querystrings, "&")
	}

	fullPath = "fmt.Sprintf(\"" + fullPath + "\", "
	fullPath = fullPath + strings.Join(method.PathParameters, ", ")
	fullPath = fullPath + strings.Join(queryKeys, ", ") + ")"

	return fullPath
}
