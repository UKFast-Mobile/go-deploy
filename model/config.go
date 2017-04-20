package model

import (
	"fmt"
	"io/ioutil"
	"reflect"

	"strings"

	"golang.org/x/crypto/ssh"
)

// DeployServerConfig configuration for the server communication
type DeployServerConfig struct {
	Host       string `json:"host" cli_q:"Host: "`
	Port       string `json:"port,omitempty" cli_q:"Port: "`
	Username   string `json:"username" cli_q:"Username: "`
	Repo       string `json:"repo" cli_q:"Repo : "`
	Ref        string `json:"refs" cli_q:"Ref: "`
	Path       string `json:"path" cli_q:"Deployment path: "`
	Cmd        string `json:"cmd" cli_q:"Command: "`
	PrivateKey string `json:"privateKey,omitempty" cli_q:"Private key path: "`
}

// Verify verfies if the config is of correct format
func (c *DeployServerConfig) Verify() error {
	// TODO: verify config setup is correct.
	v := reflect.ValueOf(c).Elem()

	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i).Interface()
		if f == "" || f == nil {
			return fmt.Errorf("Value of `%s` not set", v.Type().Field(i).Name)
		}
	}

	if len(strings.Split(c.Ref, "/")) != 2 {
		return fmt.Errorf("Incorrect ref value, must be of `origin/master` type")
	}

	return nil
}

// SSHConfig returns an ssh config for the deployment configuration
func (c *DeployServerConfig) SSHConfig() *ssh.ClientConfig {
	return &ssh.ClientConfig{
		User: c.Username,
		Auth: []ssh.AuthMethod{
			publicKeyFile(c.PrivateKey),
		},
	}
}

// BranchName returns a branch name from a given refs
func (c *DeployServerConfig) BranchName() string {
	return strings.Split(c.Ref, "/")[1]
}

// RemoteName returns a remote name from a given refs
func (c *DeployServerConfig) RemoteName() string {
	return strings.Split(c.Ref, "/")[0]
}

func publicKeyFile(file string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		errMsg := fmt.Sprintf("%s: %s", "Private key not found", err.Error())
		panic(errMsg)
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		errMsg := fmt.Sprintf("%s: %s", "Private key not found", err.Error())
		panic(errMsg)
	}

	return ssh.PublicKeys(key)
}
