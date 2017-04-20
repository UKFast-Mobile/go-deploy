package model

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"testing"

	. "github.com/franela/goblin"
	. "github.com/onsi/gomega"
)

func TestConfig(t *testing.T) {

	g := Goblin(t)

	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })

	g.Describe("DeployServerConfig", func() {

		g.Describe("Verification function", func() {
			g.Before(func() {
				ConfigFilePath = "./../go-deploy.json"
			})

			g.It("must be able to verify configuration", func() {
				var config DeployServerConfig
				pwd := os.Getenv("PWD")
				config.PrivateKey = pwd + "/test_id_rsa"
				err := loadConfiguration("thebin", &config)
				Expect(err).To(BeNil())
				err = config.Verify()
				Expect(err).To(BeNil())
			})

			g.It("must fail if configuration isn't correct", func() {
				var config DeployServerConfig
				err := loadConfiguration("stage", &config)
				Expect(err).To(BeNil())
				err = config.Verify()
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(Equal("Value of `Host` not set"))
			})

		})

		g.Describe("Convenience functions", func() {

			var config DeployServerConfig

			g.Before(func() {
				ConfigFilePath = "./../go-deploy.json"
				loadConfiguration("thebin", &config)
			})

			g.It("should be able to generate ssh config from a given configuration", func() {
				pwd := os.Getenv("PWD")
				config.PrivateKey = pwd + "/test_id_rsa"
				sshConfig := config.SSHConfig()
				Expect(sshConfig).ToNot(BeNil())
				Expect(sshConfig.User).To(Equal("mobileteamserver"))
			})

			g.It("should be able to extract branch name from the ref", func() {
				branch := config.BranchName()
				Expect(branch).To(Equal("develop"))
			})

			g.It("should be able to extract remote name from the ref", func() {
				remote := config.RemoteName()
				Expect(remote).To(Equal("origin"))
			})
		})

	})
}

func loadConfiguration(name string, c *DeployServerConfig) error {
	file, err := openConfigFile()
	if err != nil {
		return err
	}

	configInfo := file[name]
	if configInfo == nil {
		return errors.New("Configuration not found")
	}

	data, err := json.Marshal(configInfo)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &c)
	return err
}

var ConfigFilePath = "./go-deploy.json"

func openConfigFile() (map[string]interface{}, error) {
	file, err := os.OpenFile(ConfigFilePath, os.O_RDWR|os.O_CREATE, 0644)
	defer file.Close()

	fileStat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	fileSize := fileStat.Size()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var fileJSON map[string]interface{}
	if fileSize == 0 {
		// empty file, create empty map
		fileJSON = map[string]interface{}{}
	} else {
		// get json
		err = json.Unmarshal(data, &fileJSON)
		if err != nil {
			return nil, err
		}
	}

	return fileJSON, nil
}
