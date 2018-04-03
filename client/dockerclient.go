package client

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/docker/docker/client"
	"github.com/prologic/autodock/version"
	log "github.com/sirupsen/logrus"
)

const (
	apiVersion       = "1.35"
	defaultDockerURL = "unix:///var/run/docker.sock"
)

// NewTLSConfig ...
func NewTLSConfig(caCert, cert, key []byte, allowInsecure bool) (*tls.Config, error) {
	// TLS config
	var tlsConfig tls.Config
	tlsConfig.InsecureSkipVerify = true
	certPool := x509.NewCertPool()

	certPool.AppendCertsFromPEM(caCert)
	tlsConfig.RootCAs = certPool
	keypair, err := tls.X509KeyPair(cert, key)
	if err != nil {
		return &tlsConfig, err
	}
	tlsConfig.Certificates = []tls.Certificate{keypair}
	if allowInsecure {
		tlsConfig.InsecureSkipVerify = true
	}

	return &tlsConfig, nil
}

// GetDockerURL ...
func GetDockerURL(dockerURL string) string {
	// check environment for docker client config
	envDockerHost := os.Getenv("DOCKER_HOST")
	if dockerURL == "" && envDockerHost != "" {
		dockerURL = envDockerHost
	}

	if dockerURL == "" {
		dockerURL = defaultDockerURL
	}

	return dockerURL
}

// GetDockerTLSConfig ...
func GetDockerTLSConfig(tlsCaCert, tlsCert, tlsKey string, allowInsecure bool) *tls.Config {
	envDockerCertPath := os.Getenv("DOCKER_CERT_PATH")
	envDockerTLSVerify := os.Getenv("DOCKER_TLS_VERIFY")
	if tlsCaCert == "" && envDockerCertPath != "" && envDockerTLSVerify != "" {
		tlsCaCert = filepath.Join(envDockerCertPath, "ca.pem")
		tlsCert = filepath.Join(envDockerCertPath, "cert.pem")
		tlsKey = filepath.Join(envDockerCertPath, "key.pem")
	}

	if tlsCaCert != "" && tlsCert != "" && tlsKey != "" {
		log.Debug("using tls for communication with docker")
		caCert, err := ioutil.ReadFile(tlsCaCert)
		if err != nil {
			log.Fatalf("error loading tls ca cert: %s", err)
		}

		cert, err := ioutil.ReadFile(tlsCert)
		if err != nil {
			log.Fatalf("error loading tls cert: %s", err)
		}

		key, err := ioutil.ReadFile(tlsKey)
		if err != nil {
			log.Fatalf("error loading tls key: %s", err)
		}

		cfg, err := NewTLSConfig(caCert, cert, key, allowInsecure)
		if err != nil {
			log.Fatalf("error configuring tls: %s", err)
		}
		cfg.InsecureSkipVerify = envDockerTLSVerify == ""
		return cfg
	}

	return nil
}

// GetDockerClient ...
func GetDockerClient(dockerURL, tlsCaCert, tlsCert, tlsKey string, allowInsecure bool) (*client.Client, error) {
	dockerURL = GetDockerURL(dockerURL)

	var httpClient *http.Client

	// load tlsconfig
	tlsConfig := GetDockerTLSConfig(tlsCaCert, tlsCert, tlsKey, allowInsecure)

	if tlsConfig != nil {
		httpClient = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: tlsConfig,
			},
		}
	}

	defaultHeaders := map[string]string{
		"User-Agent": fmt.Sprintf("autodock-%s", version.Version),
	}
	c, err := client.NewClient(
		dockerURL, apiVersion, httpClient, defaultHeaders,
	)
	if err != nil {
		return nil, err
	}

	return c, nil
}
