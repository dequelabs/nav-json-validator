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
	directory        string
}

// New returns a new NavJSON instance based on the given `data`.
func New(dir, data string) (*NavJSON, error) {
	j := NavJSON{
		directory: dir,
	}

	// Ensure JSON is parsable.
	b := []byte(data)
	err := json.Unmarshal(b, &j)
	if err != nil {
		return nil, err
	}

	// Ensure `root` is set.
	if j.Root == "" {
		return nil, fmt.Errorf("Missing `root` key")
	}

	// Ensure `AssetRoot` is set.
	if j.AssetRoot == "" {
		return nil, fmt.Errorf("Missing `assetRoot` key")
	}

	// Ensure `files` are set and have members.
	if len(j.Files) == 0 {
		return nil, fmt.Errorf("Missing or empty `files` array")
	}

	return &j, err
}

// ValidateFiles ensures all referenced files exist on disk.
func (n *NavJSON) ValidateFiles() error {
	return validateFiles(n.directory, n.Files)
}

// validateFiles returns an error if any of the given `files` do not exist within `dir`.
func validateFiles(dir string, files []NavFile) error {
	// Ensure directory exists.
	_, err := os.Stat(dir)
	if err != nil {
		return fmt.Errorf("Directory does not exist (%s)", dir)
	}

	for _, f := range files {
		// Name can be empty when:
		// - the `Path` is a placeholder for ModX
		// - the `Path` is used for nesting (this is the case for repositories like `attest-node-suite`)
		if f.Name == "" {
			// If there is no `.Name` nor any `.Files`, ignore it.
			if len(f.Files) == 0 {
				continue
			}

			// If we have files, validate them using the `.Path` as a subdirectory.
			subdir := path.Join(dir, f.Path)
			err = validateFiles(subdir, f.Files)
			if err != nil {
				return err
			}

			continue
		}

		fp := path.Join(dir, f.Name)
		i, err := os.Stat(fp)
		if err != nil {
			return fmt.Errorf("File does not exist (%s)", fp)
		}

		if i.IsDir() {
			return fmt.Errorf("Referenced file is a directory (%s)", fp)
		}

		// If the file has files underneath it, validate them.
		if len(f.Files) > 0 {
			subdir := path.Join(dir, path.Dir(f.Name))
			err = validateFiles(subdir, f.Files)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
