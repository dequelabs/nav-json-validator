package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/dequelabs/nav-json-validator/navjson"
)

var silent *bool
var file *string
var skipFileCheck *bool

func init() {
	silent = flag.Bool("silent", false, "Silence output")
	// Skip checking for file existence by default to avoid a breaking change.
	// TODO: make the default `false` once all projects are setup to support this.
	skipFileCheck = flag.Bool("skip-file-check", true, "Skip file existence check")
	file = flag.String("file", "docs/nav.json", "Path to nav.json file")
	flag.Parse()
}

func check(err error) {
	if err == nil {
		return
	}

	if *silent {
		os.Exit(1)
	}

	panic(err)
}

func main() {
	cwd := path.Dir(*file)

	data, err := ioutil.ReadFile(*file)
	check(err)

	n, err := navjson.New(cwd, string(data))
	check(err)

	if !*skipFileCheck {
		check(n.ValidateFiles())
	}

	if *silent == false {
		fmt.Printf("File `%s` is valid", *file)
	}
}
