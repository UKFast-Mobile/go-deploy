package model

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
	return nil
}
