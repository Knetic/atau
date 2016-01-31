package atau

import (
	"encoding/json"
	"strings"
	"errors"
	"bytes"
	"path/filepath"
	"io"
	"os"
	"fmt"
	"github.com/Knetic/presilo"
)

/*
	Intermediate struct while parsing a real API
*/
type marshalledAPI struct {

	Name string `json:"name"`
	Title string `json:"title"`
	Description string `json:"description"`
	BaseURL string `json:"baseUrl"`

	Resources map[string]Resource `json:"resources"`

	Schemas map[string]*json.RawMessage `json:"schemas"`
	Parameters map[string]*json.RawMessage `json:"parameters"`
	OrderedParameters []string `json:"orderedParameters"`
}

func ParseAPIFile(path string) (*API, error) {

	var sourceFile *os.File
	var baseName string
	var err error

	baseName = filepath.Base(path)
 	baseName = strings.TrimSuffix(baseName, filepath.Ext(baseName))

	sourceFile, err = os.Open(path)
	if(err != nil) {
		return nil, err
	}
	defer sourceFile.Close()

	return ParseAPIStream(sourceFile, baseName)
}

func ParseAPIStream(reader io.Reader, defaultTitle string) (*API, error) {

	var intermediate marshalledAPI
	var ret *API
	var decoder *json.Decoder
	var schemaContext *presilo.SchemaParseContext
	var parameters map[string]presilo.TypeSchema
	var err error

	decoder = json.NewDecoder(reader)
	err = decoder.Decode(&intermediate)
	if(err != nil) {
		return nil, err
	}

	err = intermediate.Validate()
	if(err != nil) {
		return nil, err
	}

	ret = translateAPIStructs(intermediate)
	schemaContext = presilo.NewSchemaParseContext()

	// parse schemas
	_, err = parseSchemaBlock(intermediate.Schemas, schemaContext)
	if(err != nil) {
		errorMsg := fmt.Sprintf("Unable to parse api schemas: %v", err)
		return nil, errors.New(errorMsg)
	}

	err = presilo.LinkSchemas(schemaContext)
	if(err != nil) {
		return nil, err
	}

	// parse parameters
	parameters, err = parseSchemaBlock(intermediate.Parameters, schemaContext)
	if(err != nil) {
		errorMsg := fmt.Sprintf("Unable to parse api parameters: %v", err)
		return nil, errors.New(errorMsg)
	}

	err = presilo.LinkSchemas(schemaContext)
	if(err != nil) {
		return nil, err
	}
	ret.Parameters = ParameterList{Parameters: parameters}

	// deal with resources/methods.
	for resourceKey, resource := range ret.Resources {
		for methodKey, method := range resource.Methods {

			err = parseResourceMethod(resource, &method, schemaContext)
			if(err != nil) {
				errorMsg := fmt.Sprintf("Unable to parse request parameter schemas for %s%s: %v", methodKey, resourceKey, err)
				return nil, errors.New(errorMsg)
			}

			resource.Methods[methodKey] = method
		}
		ret.Resources[resourceKey] = resource
	}

	// synthetically generate an api-wide parameters schema.
	generateAPIOptions(ret, schemaContext)

	ret.schemas = schemaContext.SchemaDefinitions
	ret.schemaContext = schemaContext
	return ret, nil
}

func generateAPIOptions(api *API, schemaContext *presilo.SchemaParseContext) {

	var schema *presilo.ObjectSchema
	var name string

	schema = presilo.NewObjectSchema()
	name = presilo.ToCamelCase(api.Name) + "Options"
	schema.Title = name
	schema.ID = name

	for key, propertySchema := range api.Parameters.Parameters {
		schema.AddProperty(key, propertySchema)
	}

	api.optionsSchema = schema
	schemaContext.SchemaDefinitions[name] = schema
}

/*
	Resource Methods require a little extra parsing around parameters and schemas.
*/
func parseResourceMethod(resource Resource, method *Method, schemaContext *presilo.SchemaParseContext) error {

	var parameters ParameterList
	var err error

	// request/response
	if(method.RawRequestSchema != nil) {

		method.RequestSchema, err = unmarshalSchema(method.RawRequestSchema, "", schemaContext)
		if(err != nil) {
			errorMsg := fmt.Sprintf("Unable to parse request schema for %s: %v", resource.Name, err)
			return errors.New(errorMsg)
		}
	}

	if(method.RawResponseSchema != nil) {

		method.ResponseSchema, err = unmarshalSchema(method.RawResponseSchema, "", schemaContext)
		if(err != nil) {
			errorMsg := fmt.Sprintf("Unable to parse response schema for %s: %v", resource.Name, err)
			return errors.New(errorMsg)
		}
	}

	if(method.RawParameters != nil) {

		parameters.Parameters, err = parseSchemaBlock(method.RawParameters, schemaContext)
		if(err != nil) {
			errorMsg := fmt.Sprintf("Unable to parse parameters for %s: %v", resource.Name, err)
			return errors.New(errorMsg)
		}

		method.Parameters = parameters
	}

	return nil
}

func parseSchemaBlock(parameters map[string]*json.RawMessage, schemaContext *presilo.SchemaParseContext) (map[string]presilo.TypeSchema, error) {

	var ret map[string]presilo.TypeSchema
	var schema presilo.TypeSchema
	var rawBody []byte
	var err error

	ret = make(map[string]presilo.TypeSchema)

	for name, body := range parameters {

		rawBody, err = body.MarshalJSON()
		if(err != nil) {
			errorMsg := fmt.Sprintf("Unable unmarshal JSON for %s: %v", name, err)
			return nil, errors.New(errorMsg)
		}

		schema, err = presilo.ParseSchemaStreamContinue(bytes.NewReader(rawBody), name, schemaContext)
		if(err != nil) {
			errorMsg := fmt.Sprintf("Unable to parse schema for %s: %v", name, err)
			return nil, errors.New(errorMsg)
		}

		ret[name] = schema
	}

	return ret, nil
}

func unmarshalSchema(rawSchema *json.RawMessage, defaultTitle string, schemaContext *presilo.SchemaParseContext) (presilo.TypeSchema, error) {

	var rawBody []byte
	var err error

	rawBody, err = rawSchema.MarshalJSON()
	if(err != nil) {
		return nil, err
	}

	return presilo.ParseSchemaStreamContinue(bytes.NewReader(rawBody), defaultTitle, schemaContext)
}

func translateAPIStructs(intermediate marshalledAPI) *API {

	var ret *API

	ret = new(API)
	ret.Name = intermediate.Name
	ret.Title = intermediate.Title
	ret.BaseURL = intermediate.BaseURL
	ret.Description = intermediate.Description
	ret.Resources = intermediate.Resources

	return ret
}

/*
	Validates the first parse of an API object.
*/
func (this marshalledAPI) Validate() error {

	if(this.Title == "") {
		return errors.New("'title' was empty or not present")
	}
	if(this.BaseURL == "") {
		return errors.New("'baseUrl' was empty or not present")
	}
	if(!strings.HasPrefix(this.BaseURL, "http://") && !strings.HasPrefix(this.BaseURL, "https://")) {
		return errors.New("'baseUrl' did not refer to an http path")
	}
	if(len(this.Resources) <= 0) {
		return errors.New("'resources' section did not contain any valid resources")
	}
	return nil
}
