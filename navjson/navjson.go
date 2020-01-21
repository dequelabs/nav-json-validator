// Package navjson provides a `nav.json` parser and validator.
package navjson

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
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

// EnsureFilesExist returns an error if any referenced files do not exist or are directories.
func EnsureFilesExist(dir string, files []NavFile) error {
	for _, f := range files {
		p := path.Join(dir, f.Path, f.Name)
		i, err := os.Stat(p)
		if err != nil {
			return fmt.Errorf("Referenced file does not exist: %s", p)
		}

		if i.IsDir() {
			return fmt.Errorf("Referenced file is directory: %s", p)
		}

		// Recursively check all nested files.
		if len(f.Files) > 0 {
			subdir := path.Join(dir, f.Path)
			if err = EnsureFilesExist(subdir, f.Files); err != nil {
				return err
			}
		}
	}

	return nil
}

// IsValid checks if the given data is a valid NavJSON.
func IsValid(data string) bool {
	_, err := Parse(data)
	if err != nil {
		return false
	}
	return true
}
