package main

import (
  "fmt"
  "os"
  "atau"
)

func main() {

	var api *atau.API
	var settings RunSettings
	var err error

	settings = ParseRunSettings()
	if(len(settings.InputPaths) < 0) {
		exitWith(1, "No input path specified")
	}

	for _, path := range settings.InputPaths {

		api, err = atau.ParseAPIFile(path)
		if(err != nil) {
			exitWith(1, "Unable to parse api file: %v\n", err.Error())
		}

		err = atau.WriteGeneratedCode(api, settings.Module, settings.OutputPath, settings.Language, false)
		if(err != nil) {
			exitWith(1, "Unable to write generated code: %v\n", err.Error())
		}
	}
}

func exitWith(code int, format string, arguments ...string) {

	fmt.Fprintf(os.Stderr, format, arguments)
	os.Exit(code)
}
