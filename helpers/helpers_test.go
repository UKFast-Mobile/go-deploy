package helpers

import (
	"errors"
	"testing"

	"os"

	"io/ioutil"

	"encoding/json"

	"github.com/UKFast-Mobile/go-deploy/model"
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

		g.Describe("LoadConfiguration fucntion", func() {
			g.Before(func() {
				ConfigFilePath = "./../go-deploy.json"
			})

			g.It("Should be able to load `test` configuration without an error", func() {
				var config model.DeployServerConfig
				err := LoadConfiguration("test", &config)
				Expect(err).To(BeNil())
				Expect(config).ToNot(BeNil())
			})

			g.It("Should result in error for an unregistered configuration name", func() {
				var config model.DeployServerConfig
				err := LoadConfiguration("this does not exists", &config)
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("Configuration not found"))
			})

			g.It("Should result in error for unexisting configuration file", func() {
				ConfigFilePath = "./this_does_not_exists.json"
				var config model.DeployServerConfig
				err := LoadConfiguration("test", &config)
				Expect(err).ToNot(BeNil())
			})
		})

		g.Describe("GetEnv function", func() {
			g.It("Should be able to get env variable", func() {
				err := os.Setenv("TEST_VAR", "TEST")
				Expect(err).To(BeNil())
				var myVar string
				GetEnv("TEST_VAR", &myVar)
				Expect(myVar).To(Equal("TEST"))
			})

			g.It("Should fail silently if env varialbe doesn't exists, allowing for a default var implementation", func() {
				myVar := "doesn't exist"
				GetEnv("TEST_VAR_NOT_SET", &myVar)
				Expect(myVar).To(Equal("doesn't exist"))
			})
		})

		g.Describe("GetEnvBool function", func() {
			g.It("Should be able to get boolean environment varialbe", func() {
				err := os.Setenv("TEST_VAR", "true")
				Expect(err).To(BeNil())
				var myVar bool
				GetEnvBool("TEST_VAR", &myVar)
				Expect(myVar).To(BeTrue())
				err = os.Setenv("TEST_VAR", "1")
				Expect(err).To(BeNil())
				myVar = false
				GetEnvBool("TEST_VAR", &myVar)
				Expect(myVar).To(BeTrue())
			})

			g.It("Should fail silently if var not set", func() {
				myVar := true
				GetEnvBool("TEST_VAR_NOT_SET", &myVar)
				Expect(myVar).To(BeTrue())
			})

			g.It("Should fall back to default if var isn't a boolean", func() {
				err := os.Setenv("TEST_VAR", "Hello")
				Expect(err).To(BeNil())
				myVar := false
				GetEnvBool("TEST_VAR", &myVar)
				Expect(myVar).To(BeFalse())
			})
		})

		g.Describe("Set env vars function", func() {
			g.It("Should be able to set env vars from the configuration", func() {
				config := new(model.DeployServerConfig)
				config.EnvVars = map[string]string{
					"test_env1": "1",
					"test_env2": "hello",
				}

				SetEnvVars(config)
				env1 := os.Getenv("test_env1")
				env2 := os.Getenv("test_env2")

				Expect(env1).To(Equal("1"), "test_env1 isn't set")
				Expect(env2).To(Equal("hello"), "test_env2 isn't set")
			})
		})

	})
}
