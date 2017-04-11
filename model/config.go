package model

type DeployServerConfig struct {
	Host       string            `json:"host"`
	Port       string            `json:"port,omitempty"`
	Username   string            `json:"username"`
	Repo       string            `json:"repo"`
	Ref        string            `json:"refs"`
	Path       string            `json:"path"`
	Cmd        string            `json:"cmd"`
	PrivateKey string            `json:"privateKey,omitempty"`
	Env        map[string]string `json:"env,omitempty"`
}
