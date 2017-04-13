package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/UKFast-Mobile/go-deploy/model"
	"golang.org/x/crypto/ssh"
)

// ConfigFilePath default config file name
var ConfigFilePath = "./go-deploy.json"

// LoadConfiguration loads deployment configuration from a file with the given name and sets config struct properties accordingly
func LoadConfiguration(name string, c *model.DeployServerConfig) error {
	file, err := OpenConfigFile()
	FailOnError(err, "Failed to open configuration file, see help for usage.")
	configInfo := file[name]
	if configInfo == nil {
		FailOnError(errors.New("Configuration not found"), "Failed to load given configuration")
	}

	data, err := json.Marshal(configInfo)
	FailOnError(err, "Failed to parse configuraiton JSON")

	err = json.Unmarshal(data, &c)
	return err
}

// OpenConfigFile opens an existing config file or creates a new one if it doesn't exists
func OpenConfigFile() (map[string]interface{}, error) {
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

// WriteToFile writes json map to a config file
func WriteToFile(j map[string]interface{}) error {
	data, err := json.MarshalIndent(j, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(ConfigFilePath, data, 0644)
	return err
}

// FailOnError panics and displays the message if an error passed
// is not nil
func FailOnError(err error, msg string) {
	if err != nil {
		errMsg := fmt.Sprintf("%s: %s", msg, err.Error())
		panic(errMsg)
	}
}

func PublicKeyFile(file string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	FailOnError(err, "Private key not found")

	key, err := ssh.ParsePrivateKey(buffer)
	FailOnError(err, "Failed to parse private key")

	return ssh.PublicKeys(key)
}
