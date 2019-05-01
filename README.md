# nav-json-validator

> A tiny Golang program for validating ModX `nav.json` files.

## Installation

First, install that [`go` thing](https://golang.org/).

Next, do:

```
$ go get github.com/dequelabs/nav-json-validator
```

## Usage

```
./nav-json-validator --help
Usage of ./nav-json-validator:
  -file string
    	Path to nav.json file (default "docs/nav.json")
  -silent
    	Silence output
```

### Example Usage in CI

The below is an CircleCI job which uses `nav-json-validator`:

```yml
validate_navjson:
  docker:
    - image: circleci/golang:1.11
  steps:
    - checkout
    - run: go get github.com/dequelabs/nav-json-validator
    - run: nav-json-validator --file=/path/to/nav.json
```

## License

MPL-2.0
