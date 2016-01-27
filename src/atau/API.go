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

type API struct {

	Title string `json:"title"`
	Description string `json:"description"`
	BaseURL string `json:"baseUrl"`
	Parameters []presilo.TypeSchema

	Resources map[string]Resource `json:"resources"`
}

/*
	Intermediate struct while parsing a real API
*/
type marshalledAPI struct {

	Title string `json:"title"`
	Description string `json:"description"`
	BaseURL string `json:"baseUrl"`

	Parameters *json.RawMessage `json:"parameters"`
	Resources map[string]Resource `json:"resources"`
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

	// parse parameters in order
	ret.Parameters, err = parseParameters(intermediate.Parameters)
	if(err != nil) {
		return nil, err
	}

	return ret, nil
}

func parseParameters(parameters *json.RawMessage) ([]presilo.TypeSchema, error) {

	var ret []presilo.TypeSchema
	var decoder *json.Decoder
	var token json.Token
	/*var delimiter json.Delim
	var number float64
	var str string
	var boolean bool*/
	var err error

	decoder = json.NewDecoder(bytes.NewReader(*parameters))

	for decoder.More() {

		token, err = decoder.Token()
		if(err != nil) {
			return ret, err
		}

		if token == nil {
			return ret, errors.New("Encountered null token while parsing parameter schemas")
		}
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
