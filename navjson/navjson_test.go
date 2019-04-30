package navjson

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	t.Run("full valid doc", func(t *testing.T) {
		j, err := Parse(`
			{
				"root": "guide/attest/2.7-experiment/",
				"assetRoot": "assets/images/attest/2.7-experiment/",
				"skipMenuOrdering": true,
				"packages": {
					"attest-js": "path",
					"attest-puppeteer": "path"
				},
				"files": [
					{
						"path": "api-integrations",
						"files": [
							{
								"name": "attest-js/browser-js.html",
								"path": "browser-js"
							},
							{
								"name": "attest-puppeteer/index.html",
								"path": "node-integrations",
								"files": [
									{
										"name": "attest-puppeteer/attest-puppeteer.html",
										"path": "attest-puppeteer"
									}
								]
							}
						]
					}
				]
			}
		`)

		assert.NoError(t, err)
		assert.NotNil(t, j)
		assert.Equal(t, j.Root, "guide/attest/2.7-experiment/")
	})

	t.Run("valid doc no packages", func(t *testing.T) {
		j, err := Parse(`
			{
				"root": "guide/attest/2.7-experiment/",
				"assetRoot": "assets/images/attest/2.7-experiment/",
				"skipMenuOrdering": true,
				"files": [
					{
						"path": "api-integrations"
					}
				]
			}
		`)

		assert.NoError(t, err)
		assert.NotNil(t, j)
		assert.Nil(t, j.Packages)
	})

	t.Run("packages", func(t *testing.T) {
		j, err := Parse(`
			{
				"root": "guide/attest/2.7-experiment/",
				"assetRoot": "assets/images/attest/2.7-experiment/",
				"skipMenuOrdering": true,
				"packages": {
					"attest-js": "path",
					"attest-puppeteer": "path"
				},
				"files": [
					{
						"path": "api-integrations"
					}
				]
			}
		`)

		assert.NoError(t, err)
		assert.NotNil(t, j)
		assert.NotNil(t, j.Packages)
		assert.Equal(t, j.Packages["attest-js"], "path")
	})

	t.Run("files with name", func(t *testing.T) {
		j, err := Parse(`
			{
				"root": "guide/attest/2.7-experiment/",
				"assetRoot": "assets/images/attest/2.7-experiment/",
				"skipMenuOrdering": true,
				"files": [
					{
						"path": "api-integrations",
						"name": "hello"
					}
				]
			}
		`)

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
		_, err := Parse(`
			{
				"assetRoot": "assets/images/attest/2.7-experiment/",
				"files": [
					{
						"path": "api-integrations"
					}
				]
			}
		`)
		assert.Error(t, err)
	})

	t.Run("missing assetRoot", func(t *testing.T) {
		_, err := Parse(`
			{
				"root": "guide/attest/2.7-experiment/",
				"files": [
					{
						"path": "api-integrations"
					}
				]
			}
		`)
		assert.Error(t, err)
	})

	t.Run("missing files", func(t *testing.T) {
		_, err := Parse(`
			{
				"root": "guide/attest/2.7-experiment/",
				"assetRoot": "assets/images/attest/2.7-experiment/"
			}
		`)
		assert.Error(t, err)
	})

	t.Run("skipMenuOrdering as string", func(t *testing.T) {
		_, err := Parse(`
			{
				"root": "guide/attest/2.7-experiment/",
				"assetRoot": "assets/images/attest/2.7-experiment/",
				"skipMenuOrdering": "true",
				"files": [
					{
						"path": "api-integrations"
					}
				]
			}
		`)
		assert.Error(t, err)
	})
}

func TestParseExamples(t *testing.T) {
	files := []string{"attest-docs.json", "attest-node-suite.json"}
	for _, filename := range files {
		t.Run(filename, func(t *testing.T) {
			data, err := ioutil.ReadFile("./testdata/" + filename)
			assert.NoError(t, err)
			_, err = Parse(string(data))
			assert.NoError(t, err)
		})
	}
}

func TestIsValid(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		data := `
			{
				"root": "guide/attest/2.7-experiment/",
				"assetRoot": "assets/images/attest/2.7-experiment/",
				"skipMenuOrdering": true,
				"packages": {
					"attest-js": "path",
					"attest-puppeteer": "path"
				},
				"files": [
					{
						"path": "api-integrations",
						"files": [
							{
								"name": "attest-js/browser-js.html",
								"path": "browser-js"
							},
							{
								"name": "attest-puppeteer/index.html",
								"path": "node-integrations",
								"files": [
									{
										"name": "attest-puppeteer/attest-puppeteer.html",
										"path": "attest-puppeteer"
									}
								]
							}
						]
					}
				]
			}
		`
		assert.True(t, IsValid(data))
	})

	t.Run("invalid", func(t *testing.T) {
		data := `
			{
				"root": "guide/attest/2.7-experiment/",
				"assetRoot": "assets/images/attest/2.7-experiment/",
			}
		`
		assert.False(t, IsValid(data))
	})
}
