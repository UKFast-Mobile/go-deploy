package model

import (
	"fmt"
	"io/ioutil"

	"strings"

	"golang.org/x/crypto/ssh"
)

// DeployServerConfig configuration for the server communication
type DeployServerConfig struct {
	Host       string            `json:"host" cli_q:"Host: "`
	Port       string            `json:"port,omitempty" cli_q:"Port: "`
	Username   string            `json:"username" cli_q:"Username: "`
	Repo       string            `json:"repo" cli_q:"Repo : "`
	Ref        string            `json:"refs" cli_q:"Ref: "`
	Path       string            `json:"path" cli_q:"Deployment path: "`
	Cmd        string            `json:"cmd" cli_q:"Command: "`
	PrivateKey string            `json:"privateKey,omitempty" cli_q:"Private key path: "`
	Env        map[string]string `json:"env,omitempty"`
}

// Verify verfies if the config is of correct format
func (c *DeployServerConfig) Verify() error {
	// TODO: verify config setup is correct.
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
