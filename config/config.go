package config

// Config ...
type Config struct {
	Debug         bool
	Bind          string
	DockerURL     string
	TLSCACert     string
	TLSCert       string
	TLSKey        string
	AllowInsecure bool
}
