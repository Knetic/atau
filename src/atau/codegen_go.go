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
			fullPath = interpolatePath(api, method, fullPath)
			fullPath = appendGoQuerystringPath(api, method, fullPath)

			// body, if applicable.
			if(method.RequestSchema != nil) {
				buffer.Printfln("marshalledBody, err := json.Marshal(requestContents)")
				addGoErrCheck(buffer, hasResponse)

				buffer.Printfln("request, err = http.NewRequest(\"%s\", %s, bytes.NewReader(marshalledBody))", strings.ToUpper(method.HttpMethod), fullPath)
			} else {
				buffer.Printfln("request, err = http.NewRequest(\"%s\", %s, nil)", strings.ToUpper(method.HttpMethod), fullPath)
			}

			buffer.Printfln("request.Header.Set(\"Content-Type\", \"application/json\")")

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
	buffer.Printf("\n\"bytes\"")
	buffer.Printf("\n\"io/ioutil\"")
	buffer.AddIndentation(-1)
	buffer.Printfln(")")
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

func appendGoQuerystringPath(api *API, method Method, fullPath string) string {

	var querystrings []string
	var queryKeys []string

	// set querystring
	for key, _ := range method.Parameters.Parameters {
		querystrings = append(querystrings, fmt.Sprintf("%s=%%v", key))
		queryKeys = append(queryKeys, key)
	}

	if(len(querystrings) > 0) {

		fullPath = fullPath + "?" + strings.Join(querystrings, "&")
		fullPath = "fmt.Sprintf(\"" + fullPath + "\", " + strings.Join(queryKeys, ", ") + ")"
	}

	return fullPath
}
