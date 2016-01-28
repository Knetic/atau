package main

import (
	"flag"
)

type RunSettings struct {

	InputPaths []string
	OutputPath string

	Language string
	Module string
}

func ParseRunSettings() RunSettings {

	var ret RunSettings

	flag.StringVar(&ret.OutputPath, "o", "./", "Output directory to which generated files should be written")
	flag.StringVar(&ret.Language, "l", "go", "Language that generated files ought to be")
	flag.StringVar(&ret.Module, "m", "main", "Module (or package) path which generated files ought to use")

	flag.Parse()
	ret.InputPaths = flag.Args()
	return ret
}
