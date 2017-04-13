package helpers

import (
	"errors"
	"testing"

	"os"

	"io/ioutil"

	"encoding/json"

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
				expected := map[string]interface{}{}
				Expect(fileJSON).To(Equal(expected), "Not an empty map")
			})

			g.It("must exctract json from file correctly", func() {
				err := ioutil.WriteFile(ConfigFilePath, []byte(`{"text":"hello","id":1,"bool":true}`), 0644)
				Expect(err).To(BeNil(), "Failed to prepare file (write to file)")
				fileJSON, err := OpenConfigFile()

				Expect(err).To(BeNil(), "OpenConfigFile resulted in error")

				Expect((fileJSON)["text"]).To(Equal("hello"))
				Expect((fileJSON)["id"]).To(BeEquivalentTo(1))
				Expect((fileJSON)["bool"]).To(BeTrue())
				Expect(fileJSON).To(HaveLen(3))
			})

			g.It("Must return an error if file containes malformed json", func() {
				err := ioutil.WriteFile(ConfigFilePath, []byte(`hello`), 0644)
				Expect(err).To(BeNil(), "Failed to prepare file (write to file)")
				fileJSON, err := OpenConfigFile()
				Expect(err).To(Not(BeNil()))
				Expect(fileJSON).To(BeNil())
			})

			g.After(func() {
				// delete test file
				_ = os.Remove(ConfigFilePath)
			})
		})

		g.Describe("WriteToFile function", func() {
			g.Before(func() {
				// setup
				ConfigFilePath = "./testsConfig.json"
				// make sure file doesn't exists at the start
				_ = os.Remove(ConfigFilePath)
			})

			g.It("Should be able to write to file", func() {
				jOrig := map[string]interface{}{
					"text": "hello",
				}

				err := WriteToFile(jOrig)
				Expect(err).To(BeNil())

				file, err := os.OpenFile(ConfigFilePath, os.O_RDONLY, 0644)
				Expect(err).To(BeNil())
				Expect(file).To(Not(BeNil()))

				data, err := ioutil.ReadAll(file)
				Expect(err).To(BeNil())
				Expect(data).To(Not(BeNil()))

				var jExt map[string]interface{}
				err = json.Unmarshal(data, &jExt)

				Expect(err).To(BeNil())
				Expect(jExt).To(BeEquivalentTo(jOrig))

			})

			g.After(func() {
				// delete test file
				_ = os.Remove(ConfigFilePath)
			})
		})

		g.Describe("FailOneError", func() {
			g.It("should panic if an error is passed", func() {
				defer func() {
					if r := recover(); r == nil {
						g.Fail("Did not panic")
					}
				}()
				FailOnError(errors.New("Panic error"), "Something went wrong")
			})
			g.It("Should not panic if an error isn't passed", func() {
				defer func() {
					if r := recover(); r != nil {
						g.Fail("Panicked")
					}
					FailOnError(nil, "Nothing to see here")
				}()
			})
		})
	})
}
