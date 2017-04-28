package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"bytes"

	"log"

	"strconv"

	"github.com/UKFast-Mobile/go-deploy/model"
	"github.com/ttacon/chalk"
	"golang.org/x/crypto/ssh"
)

// ConfigFilePath default config file name
var ConfigFilePath = "./go-deploy.json"

// Debug holds information if CLI run with the debug flag
var Debug bool

// LoadConfiguration loads deployment configuration from a file with the given name and sets config struct properties accordingly
func LoadConfiguration(name string, c *model.DeployServerConfig) error {
	file, err := OpenConfigFile()
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

// ExecuteCmd establishes a secure shell session with the remote and executes provided commands.
// Note - context isn't preserved between commands
func ExecuteCmd(cmd []string, deployConfig *model.DeployServerConfig) (string, error) {
	conn, err := ssh.Dial("tcp", deployConfig.Host+":"+deployConfig.Port, deployConfig.SSHConfig())
	if err != nil {
		return "", err
	}
	defer conn.Close()

	var stdoutBuf bytes.Buffer
	var debug bool
	GetEnvBool("DEBUG", &debug)

	for _, command := range cmd {
		LogDebug(fmt.Sprintf("Running command: %s", command))
		session, err := conn.NewSession()
		if err != nil {
			return "", err
		}
		defer session.Close()
		session.Stdout = &stdoutBuf
		err = session.Run(command)
		if err != nil {
			LogDebug(fmt.Sprintf("Failed, stdoutbuf: %s", stdoutBuf.String()))
			return "", err
		}
	}

	return stdoutBuf.String(), nil
}

// GetEnv populates your variable with the env var value
func GetEnv(name string, target *string) {
	value := os.Getenv(name)
	if value != "" {
		*target = value
	}
}

// GetEnvBool parses environment variable to a bool
func GetEnvBool(name string, target *bool) {
	var result string

	GetEnv(name, &result)
	if result != "" {
		b, err := strconv.ParseBool(result)
		if err == nil {
			*target = b
		}
	}
}

// SetEnvVars sets environment variables as per configuration
func SetEnvVars(config *model.DeployServerConfig) {
	for k := range config.EnvVars {
		os.Setenv(k, config.EnvVars[k])
	}
}

// LogDebug logs only if debug passed is true
func LogDebug(msg string) {
	if Debug {
		dlog := chalk.Magenta.NewStyle().
			WithTextStyle(chalk.Italic).Style
		log.Printf(dlog(msg))
	}
}

// LogDebugf logs only if debug passed is true, you can use format here
func LogDebugf(format string, a ...interface{}) {
	LogDebug(fmt.Sprintf(format, a))
}
