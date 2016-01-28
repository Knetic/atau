package atau

import (
	"encoding/json"
	"strings"
	"errors"
	"bytes"
	"path/filepath"
	"io"
	"os"
	"github.com/Knetic/presilo"
)

/*
	Intermediate struct while parsing a real API
*/
type marshalledAPI struct {

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
		return nil, err
	}

	// parse parameters in order
	ret.Parameters, err = parseSchemaBlock(intermediate.Parameters, schemaContext)
	if(err != nil) {
		return nil, err
	}

	ret.schemas = schemaContext.SchemaDefinitions
	return ret, nil
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
			return nil, err
		}

		schema, err = presilo.ParseSchemaStreamContinue(bytes.NewReader(rawBody), name, schemaContext)
		if(err != nil) {
			return nil, err
		}

		ret[name] = schema
	}

	return ret, nil
}

func translateAPIStructs(intermediate marshalledAPI) *API {

	var ret *API

	ret = new(API)
	ret.Title = intermediate.Title
	ret.BaseURL = intermediate.BaseURL
	ret.Description = intermediate.Description
	ret.Resources = intermediate.Resources
	ret.orderedParameters = intermediate.OrderedParameters

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
