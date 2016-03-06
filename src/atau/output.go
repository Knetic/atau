package atau

import (
	"os"
	"errors"
	"io/ioutil"
	"path/filepath"
	"github.com/Knetic/presilo"
)

func WriteGeneratedCode(api *API, module string, targetPath string, language string, unsafeModule bool) error {

	var generator func(*API, string)(string, error)
	var targetName string
	var contents string
	var err error

	// get final file name
	targetPath, err = prepareOutputPath(targetPath)
	if(err != nil) {
		return err
	}

	targetName = filepath.Join(targetPath, module + "." + language)

	switch language {

	case "go":
		generator = GenerateGo
	case "cs":
		generator = GenerateCSharp
	default:
		return errors.New("Invalid output language specified")
	}

	// first, write presilo schemas
	err = presilo.WriteGeneratedCode(api.schemaContext, module, targetPath, language, "\t", false, true)
	if(err != nil) {
		return err
	}

	// now write atau wiring
	contents, err = generator(api, module)
	if(err != nil) {
		return err
	}

	return ioutil.WriteFile(targetName, []byte(contents), os.ModePerm)
}

/*
  Given the output path, returns the absolute value of it,
  and ensures that the given path exists.
*/
func prepareOutputPath(targetPath string) (string, error) {

	var err error

	targetPath, err = filepath.Abs(targetPath)
	if err != nil {
		return "", err
	}

	err = os.MkdirAll(targetPath, os.ModePerm)
	if err != nil {
		return "", err
	}

	return targetPath, nil
}
