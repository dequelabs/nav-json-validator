// Package navjson provides a `nav.json` parser and validator.
package navjson

import (
	"encoding/json"
	"fmt"
)

// NavFile represents a file/path association in a `nav.json` document.
type NavFile struct {
	Name  string    `json:"name,omitempty"`
	Path  string    `json:"path"`
	Files []NavFile `json:"files,omitempty"`
}

// NavPackage represents a list of packages in a `nav.json` document. This is only used for monorepos.
type NavPackage map[string]string

// NavJSON represents a `nav.json` document.
type NavJSON struct {
	Root             string     `json:"root"`
	AssetRoot        string     `json:"assetRoot"`
	SkipMenuOrdering bool       `json:"skipMenuOrdering,omitempty"`
	Packages         NavPackage `json:"packages,omitempty"`
	Files            []NavFile  `json:"files"`
}

// Parse attempts to parse the given data into a NavJSON object.
func Parse(data string) (NavJSON, error) {
	j := NavJSON{}
	b := []byte(data)
	err := json.Unmarshal(b, &j)
	if err != nil {
		return j, err
	}

	if j.Root == "" {
		return j, fmt.Errorf("Missing `root` key")
	}

	if j.AssetRoot == "" {
		return j, fmt.Errorf("Missing `assetRoot` key")
	}

	if len(j.Files) == 0 {
		return j, fmt.Errorf("Missing or empty `files` array")
	}

	return j, err
}

// IsValid checks if the given data is a valid NavJSON.
func IsValid(data string) bool {
	_, err := Parse(data)
	if err != nil {
		return false
	}
	return true
}
