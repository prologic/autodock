package config

// Config ...
type Config struct {
	Debug         bool
	Bind          string
	MsgBusURL     string
	DockerURL     string
	TLSCACert     string
	TLSCert       string
	TLSKey        string
	AllowInsecure bool
}
