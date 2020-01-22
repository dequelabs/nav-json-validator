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

func TestNew(t *testing.T) {
	t.Run("full valid doc", func(t *testing.T) {
		data, err := readTestdataFile("simple-valid.json")
		assert.NoError(t, err)

		j, err := New(".", data)
		assert.NoError(t, err)
		assert.NotNil(t, j)
		assert.Equal(t, j.Root, "guide/attest/2.7-experiment/")
	})

	t.Run("valid doc no packages", func(t *testing.T) {
		data, err := readTestdataFile("no-packages.json")
		assert.NoError(t, err)

		j, err := New(".", data)
		assert.NoError(t, err)
		assert.NotNil(t, j)
		assert.Nil(t, j.Packages)
	})

	t.Run("packages", func(t *testing.T) {
		data, err := readTestdataFile("with-packages.json")
		assert.NoError(t, err)

		j, err := New(".", data)
		assert.NoError(t, err)
		assert.NotNil(t, j)
		assert.NotNil(t, j.Packages)
		assert.Equal(t, j.Packages["attest-js"], "path")
	})

	t.Run("files with name", func(t *testing.T) {
		data, err := readTestdataFile("files-with-names.json")
		assert.NoError(t, err)

		j, err := New(".", data)
		assert.NoError(t, err)
		assert.NotNil(t, j)
		assert.Len(t, j.Files, 1)

		f := j.Files[0]
		assert.Equal(t, f.Name, "hello")
	})

	t.Run("invalid JSON", func(t *testing.T) {
		_, err := New(".", `|)(*(&DF(*FUDF)))}`)
		assert.Error(t, err)
	})

	t.Run("missing root", func(t *testing.T) {
		data, err := readTestdataFile("missing-root.json")
		assert.NoError(t, err)

		_, err = New(".", data)
		assert.Error(t, err)
	})

	t.Run("missing assetRoot", func(t *testing.T) {
		data, err := readTestdataFile("missing-assetRoot.json")
		assert.NoError(t, err)

		_, err = New(".", data)
		assert.Error(t, err)
	})

	t.Run("missing files", func(t *testing.T) {
		data, err := readTestdataFile("missing-files.json")
		assert.NoError(t, err)

		_, err = New(".", data)
		assert.Error(t, err)
	})

	t.Run("skipMenuOrdering as string", func(t *testing.T) {
		data, err := readTestdataFile("invalid-skipMenuOrdering.json")
		assert.NoError(t, err)

		_, err = New(".", data)
		assert.Error(t, err)
	})
}

func TestParseExamples(t *testing.T) {
	files := []string{"attest-docs.json", "attest-node-suite.json"}
	for _, filename := range files {
		t.Run(filename, func(t *testing.T) {
			data, err := readTestdataFile(filename)
			assert.NoError(t, err)
			_, err = New("./testdata", data)
			assert.NoError(t, err)
		})
	}
}

func TestValidateFiles(t *testing.T) {
	cwd, err := os.Getwd()
	assert.NoError(t, err)

	t.Run("missing directory", func(t *testing.T) {
		n, err := New("/nope", `
			{
				"root": "guide/attest/2.7-experiment/",
				"assetRoot": "assets/images/attest/2.7-experiment/",
				"files": [
					{
						"name": "foo.html",
						"path": "foo",
						"files": [
							{ "name": "bar.html", "path": "bar" }
						]
					}
				]
			}
		`)
		assert.NoError(t, err)

		err = n.ValidateFiles()
		assert.Error(t, err)
	})

	t.Run("skip missing Name and empty Files", func(t *testing.T) {
		n, err := New(cwd, `
			{
				"root": "guide/attest/2.7-experiment/",
				"assetRoot": "assets/images/attest/2.7-experiment/",
				"files": [
					{ "path": "foo" }
				]
			}
		`)
		assert.NoError(t, err)

		assert.NoError(t, n.ValidateFiles())
	})

	t.Run("nested files no name", func(t *testing.T) {
		n, err := New(cwd, `
			{
				"root": "guide/attest/2.7-experiment/",
				"assetRoot": "assets/images/attest/2.7-experiment/",
				"files": [
					{
						"path": "testdata",
						"files": [
							{ "name": "files/bar", "path": "bar" }
						]
					}
				]
			}
		`)
		assert.NoError(t, err)
		assert.NoError(t, n.ValidateFiles())
	})

	t.Run("file is directory", func(t *testing.T) {
		n, err := New(cwd, `
			{
				"root": "guide/attest/2.7-experiment/",
				"assetRoot": "assets/images/attest/2.7-experiment/",
				"files": [
					{
						"path": "testdata/files/qux",
						"name": "testdata/files/qux"
					}
				]
			}
		`)
		assert.NoError(t, err)
		assert.Error(t, n.ValidateFiles())
	})

	t.Run("nested files missing", func(t *testing.T) {
		n, err := New(cwd, `
			{
				"root": "guide/attest/2.7-experiment/",
				"assetRoot": "assets/images/attest/2.7-experiment/",
				"files": [
					{
						"path": "testdata",
						"files": [
							{ "name": "files/nope", "path": "bar" }
						]
					}
				]
			}
		`)
		assert.NoError(t, err)
		assert.Error(t, n.ValidateFiles())
	})

	t.Run("nested files no name (multiple levels)", func(t *testing.T) {
		n, err := New(cwd, `
			{
				"root": "guide/attest/2.7-experiment/",
				"assetRoot": "assets/images/attest/2.7-experiment/",
				"files": [
					{
						"path": "testdata",
						"files": [
							{ "name": "files/bar", "path": "bar" },
							{
								"path": "files/qux",
								"files": [
									{ "name": "1", "path": "1" },
									{ "name": "2", "path": "2" }
								]
							}
						]
					}
				]
			}
		`)
		assert.NoError(t, err)
		assert.NoError(t, n.ValidateFiles())
	})

	t.Run("missing nested files no name (multiple levels)", func(t *testing.T) {
		n, err := New(cwd, `
			{
				"root": "guide/attest/2.7-experiment/",
				"assetRoot": "assets/images/attest/2.7-experiment/",
				"files": [
					{
						"path": "testdata",
						"files": [
							{ "name": "files/bar", "path": "bar" },
							{
								"path": "files/qux",
								"files": [
									{ "name": "potato", "path": "banana" }
								]
							}
						]
					}
				]
			}
		`)
		assert.NoError(t, err)
		assert.Error(t, n.ValidateFiles())
	})

	t.Run("nested files with name", func(t *testing.T) {
		n, err := New(cwd, `
			{
				"root": "guide/attest/2.7-experiment/",
				"assetRoot": "assets/images/attest/2.7-experiment/",
				"files": [
					{
						"path": "testdata",
						"name": "testdata/files/bar",
						"files": [
							{ "name": "baz", "path": "baz" }
						]
					}
				]
			}
		`)
		assert.NoError(t, err)
		assert.NoError(t, n.ValidateFiles())
	})

	t.Run("missing nested files with name", func(t *testing.T) {
		n, err := New(cwd, `
			{
				"root": "guide/attest/2.7-experiment/",
				"assetRoot": "assets/images/attest/2.7-experiment/",
				"files": [
					{
						"path": "testdata",
						"name": "testdata/files/bar",
						"files": [
							{ "name": "nope", "path": "baz" }
						]
					}
				]
			}
		`)
		assert.NoError(t, err)
		assert.Error(t, n.ValidateFiles())
	})
}
