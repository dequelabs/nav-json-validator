package navjson

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func readTestdataFile(filename string) (string, error) {
	data, err := ioutil.ReadFile("./testdata/" + filename)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func TestParse(t *testing.T) {
	t.Run("full valid doc", func(t *testing.T) {
		data, err := readTestdataFile("simple-valid.json")
		assert.NoError(t, err)

		j, err := Parse(data)
		assert.NoError(t, err)
		assert.NotNil(t, j)
		assert.Equal(t, j.Root, "guide/attest/2.7-experiment/")
	})

	t.Run("valid doc no packages", func(t *testing.T) {
		data, err := readTestdataFile("no-packages.json")
		assert.NoError(t, err)

		j, err := Parse(data)
		assert.NoError(t, err)
		assert.NotNil(t, j)
		assert.Nil(t, j.Packages)
	})

	t.Run("packages", func(t *testing.T) {
		data, err := readTestdataFile("with-packages.json")
		assert.NoError(t, err)

		j, err := Parse(data)
		assert.NoError(t, err)
		assert.NotNil(t, j)
		assert.NotNil(t, j.Packages)
		assert.Equal(t, j.Packages["attest-js"], "path")
	})

	t.Run("files with name", func(t *testing.T) {
		data, err := readTestdataFile("files-with-names.json")
		assert.NoError(t, err)

		j, err := Parse(data)
		assert.NoError(t, err)
		assert.NotNil(t, j)
		assert.Len(t, j.Files, 1)

		f := j.Files[0]
		assert.Equal(t, f.Name, "hello")
	})

	t.Run("invalid JSON", func(t *testing.T) {
		_, err := Parse(`|)(*(&DF(*FUDF)))}`)
		assert.Error(t, err)
	})

	t.Run("missing root", func(t *testing.T) {
		data, err := readTestdataFile("missing-root.json")
		assert.NoError(t, err)

		_, err = Parse(data)
		assert.Error(t, err)
	})

	t.Run("missing assetRoot", func(t *testing.T) {
		data, err := readTestdataFile("missing-assetRoot.json")
		assert.NoError(t, err)

		_, err = Parse(data)
		assert.Error(t, err)
	})

	t.Run("missing files", func(t *testing.T) {
		data, err := readTestdataFile("missing-files.json")
		assert.NoError(t, err)

		_, err = Parse(data)
		assert.Error(t, err)
	})

	t.Run("skipMenuOrdering as string", func(t *testing.T) {
		data, err := readTestdataFile("invalid-skipMenuOrdering.json")
		assert.NoError(t, err)

		_, err = Parse(data)
		assert.Error(t, err)
	})
}

func TestParseExamples(t *testing.T) {
	files := []string{"attest-docs.json", "attest-node-suite.json"}
	for _, filename := range files {
		t.Run(filename, func(t *testing.T) {
			data, err := readTestdataFile(filename)
			assert.NoError(t, err)
			_, err = Parse(data)
			assert.NoError(t, err)
		})
	}
}

func TestIsValid(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		data, err := readTestdataFile("simple-valid.json")
		assert.NoError(t, err)
		assert.True(t, IsValid(data))
	})

	t.Run("invalid", func(t *testing.T) {
		data, err := readTestdataFile("missing-files.json")
		assert.NoError(t, err)
		assert.False(t, IsValid(data))
	})
}

func TestEnsureFilesExist(t *testing.T) {
	cwd, err := os.Getwd()
	assert.NoError(t, err)

	t.Run("missing path", func(t *testing.T) {
		err := EnsureFilesExist(cwd, []NavFile{
			{Name: "foo", Path: "testdata/dsfds"},
		})
		assert.Error(t, err)
	})

	t.Run("missing file", func(t *testing.T) {
		err := EnsureFilesExist(cwd, []NavFile{
			{Name: "dsfds", Path: "testdata/files"},
		})
		assert.Error(t, err)
	})

	t.Run("exists", func(t *testing.T) {
		err := EnsureFilesExist(cwd, []NavFile{
			{Name: "foo", Path: "testdata/files"},
			{Name: "bar", Path: "testdata/files"},
			{Name: "baz", Path: "testdata/files"},
		})
		assert.NoError(t, err)
	})

	t.Run("directory", func(t *testing.T) {
		err := EnsureFilesExist(cwd, []NavFile{
			{Name: "qux", Path: "testdata/files"},
		})
		assert.Error(t, err)
	})

	t.Run("nested exists", func(t *testing.T) {
		err := EnsureFilesExist(cwd, []NavFile{
			{Name: "foo", Path: "testdata/files", Files: []NavFile{
				{Name: "1", Path: "qux"}, {Name: "2", Path: "qux"}},
			},
			{Name: "baz", Path: "testdata/files"},
		})
		assert.NoError(t, err)
	})

	t.Run("nested missing", func(t *testing.T) {
		err := EnsureFilesExist(cwd, []NavFile{
			{Name: "foo", Path: "testdata/files", Files: []NavFile{
				{Name: "1", Path: "qux"},
				{Name: "3", Path: "qux"}},
			},
		})
		assert.Error(t, err)
	})
}
