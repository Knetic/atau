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

	api, err = atau.ParseAPIFile(settings.InputPaths[0])
	if(err != nil) {
		exitWith(1, "Unable to parse api file: %v\n", err.Error())
	}

	atau.WriteGeneratedCode(api, settings.Module, settings.OutputPath, settings.Language, false)
}

func exitWith(code int, format string, arguments ...string) {

	fmt.Fprintf(os.Stderr, format, arguments)
	os.Exit(code)
}
