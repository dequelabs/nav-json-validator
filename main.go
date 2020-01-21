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

func init() {
	silent = flag.Bool("silent", false, "Silence output")
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

	n, err := navjson.Parse(string(data))
	check(err)

	check(navjson.EnsureFilesExist(cwd, n.Files))

	if *silent == false {
		fmt.Printf("File `%s` is valid", *file)
	}
}
