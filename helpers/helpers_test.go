package helpers

import (
	"testing"

	"os"

	"io/ioutil"

	. "github.com/franela/goblin"
	. "github.com/onsi/gomega"
)

func TestHelpers(t *testing.T) {
	g := Goblin(t)

	// Gomega fail handler
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("Helpers", func() {
		g.Describe("OpenConfigFile function", func() {
			g.Before(func() {
				// setup
				ConfigFilePath = "./testsConfig.json"
				// make sure file doesn't exists at the start
				_ = os.Remove(ConfigFilePath)
			})

			g.It("must create a .json file when run for the first time and return an empty map[string]interface{}", func() {
				fileJSON, err := OpenConfigFile()
				Expect(err).To(BeNil(), "OpenConfigFile resulted in error")
				expected := new(map[string]interface{})
				Expect(fileJSON).To(Equal(expected), "Not an empty map")
			})

			g.It("must exctract json from file correctly", func() {
				err := ioutil.WriteFile(ConfigFilePath, []byte(`{"text":"hello","id":1,"bool":true}`), 0644)
				Expect(err).To(BeNil(), "Failed to prepare file (write to file)")
				fileJSON, err := OpenConfigFile()
				Expect(err).To(BeNil(), "OpenConfigFile resulted in error")
				expected := &map[string]interface{}{
					"text": "hello",
					"id":   1,
					"bool": true,
				}
				Expect(fileJSON["text"]).To(Equal("hello"))
				Expect(fileJSON["id"]).To(Equal(1))
				Expect(fileJSON["bool"]).To(BeTrue())
				Expect(fileJSON).To(HaveLen(3))
			})

			g.After(func() {
				// delete test file
				_ = os.Remove(ConfigFilePath)
			})
		})

		g.Describe("WriteToFile function", func() {

		})
	})
}
