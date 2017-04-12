package helpers

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// ConfigFilePath default config file name
var ConfigFilePath = "./go-deploy.json"

// OpenConfigFile opens an existing config file or creates a new one if it doesn't exists
func OpenConfigFile() (*map[string]interface{}, error) {
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

	var fileJSON *map[string]interface{}
	if fileSize == 0 {
		// empty file, create empty map
		fileJSON = new(map[string]interface{})
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
