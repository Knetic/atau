package atau

import (
	"fmt"
	"strings"
	"github.com/Knetic/presilo"
)

func GenerateCSharp(api *API, module string)(string, error) {

	var buffer *presilo.BufferedFormatString

	buffer = presilo.NewBufferedFormatString("\t")

	generateCSharpImports(api, buffer)

	generateCSharpClassSignature(api, module, buffer)
	generateCSharpResourceMethods(api, buffer)

	// close up namespace/class
	buffer.AddIndentation(-1)
	buffer.Printfln("\n}")
	buffer.AddIndentation(-1)
	buffer.Printfln("\n}")
	return buffer.String(), nil
}

func generateCSharpImports(api *API, buffer *presilo.BufferedFormatString) {

	buffer.Printfln("using System;")
	buffer.Printfln("using System.Text.RegularExpressions;")
	buffer.Printfln("using System.Runtime.Serialization;")
	buffer.Printfln("using System.Net;")
	buffer.Printfln("using System.Runtime.Serialization.Json;")

	buffer.Printfln("")
}

func generateCSharpClassSignature(api *API, module string, buffer *presilo.BufferedFormatString) {

	buffer.Printf("namespace %s\n{", module)
	buffer.AddIndentation(1)
	buffer.Printf("\npublic class %s\n{", presilo.ToStrictCamelCase(api.Name))
	buffer.AddIndentation(1)
}

func generateCSharpResourceMethods(api *API, buffer *presilo.BufferedFormatString) {

	var fullPath string
	var methodName string
	var responseTypeName string
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
			generateCSharpMethodSignature(methodName, api.optionsSchema, method, buffer)

			// params
			buffer.Printfln("\nHttpWebRequest request;")
			buffer.Printfln("HttpWebResponse response;")
			buffer.Printfln("")

			// request
			fullPath = resolvePath(api, method)
			fullPath = interpolatePath(api, method, fullPath)
			fullPath = appendQuerystringPath(api, method, fullPath)

			buffer.Printfln("request = (HttpWebRequest)WebRequest.Create(\"%s\");", fullPath)
			buffer.Printfln("request.Method = \"%s\";", method.HttpMethod)

			// headers
			for key, _ := range method.Headers.Parameters {
				buffer.Printfln("request.Headers.Add(\"%s\") = %s;", key, key)
			}

			buffer.Printfln("response = (HttpWebResponse)request.GetResponse();");

			// check for bad status and hey, since this is C#, throw an exception!
			buffer.Printf("if((int)response.StatusCode >= 400)")
			buffer.AddIndentation(1)
			buffer.Printf("\nthrow new Exception(\"Server returned status \"+response.StatusCode+\".\");")
			buffer.AddIndentation(-1)

			// read response and unmarshal
			if(hasResponse) {

				responseTypeName = presilo.ToStrictCamelCase(method.ResponseSchema.GetTitle())
				buffer.Printf("\nDataContractJsonSerializer deserializer = new DataContractJsonSerializer(typeof(%s));", responseTypeName)
				buffer.Printfln("\nreturn (deserializer.ReadObject(response.GetResponseStream()) as %s);", responseTypeName)
			}

			buffer.AddIndentation(-1)
			buffer.Printf("\n}\n")
		}
	}

}

func generateCSharpMethodSignature(methodName string, optionsSchema *presilo.ObjectSchema, method Method, buffer *presilo.BufferedFormatString) {

	var arguments []string
	var parameterSchema presilo.TypeSchema
	var parameterType string
	var returnType string

	if(method.ResponseSchema != nil) {
		returnType = presilo.ToStrictCamelCase(method.ResponseSchema.GetTitle())
	} else {
		returnType = "void"
	}

	buffer.Printf("\npublic static %s %s(", returnType, methodName)

	// method-specific parameters
	for _, parameterName := range method.Parameters.GetOrderedParameters() {

		parameterSchema = method.Parameters.Parameters[parameterName]
		parameterType = presilo.GenerateGoTypeForSchema(parameterSchema)
		parameterName = presilo.ToStrictJavaCase(parameterName)

		arguments = append(arguments, fmt.Sprintf("%s %s", parameterType, parameterName))
	}

	// method-specific header parameters
	for _, parameterName := range method.Headers.GetOrderedParameters() {

		parameterSchema = method.Headers.Parameters[parameterName]
		parameterType = presilo.GenerateGoTypeForSchema(parameterSchema)
		parameterName = presilo.ToStrictJavaCase(parameterName)

		arguments = append(arguments, fmt.Sprintf("%s %s", parameterType, parameterName))
	}

	if(method.RequestSchema != nil) {
		arguments = append(arguments, fmt.Sprintf("%s requestContents", presilo.ToStrictCamelCase(method.RequestSchema.GetTitle())))
	}
	if(optionsSchema != nil) {
		arguments = append(arguments, fmt.Sprintf("%s options", presilo.ToStrictCamelCase(optionsSchema.GetTitle())))
	}

	buffer.Printfln("%s) ", strings.Join(arguments, ", "))
	buffer.Printf("{")
	buffer.AddIndentation(1)
}
