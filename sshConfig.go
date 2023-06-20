package main

type sshConfig struct {
	localPort  string
	hostUser   string
	targetPort string
	targetHost string
	sshHost    string
}

func NewSSHConfig() *sshConfig {
	return &sshConfig{}
}
