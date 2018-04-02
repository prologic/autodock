package server

import (
	"github.com/prologic/autodock/client"
	"github.com/prologic/autodock/proxy"
)

func (s *Server) getDockerURL() string {
	return client.GetDockerURL(s.cfg.DockerURL)
}

func (s *Server) getDockerProxy() (*proxy.Proxy, error) {
	tlsConfig := client.GetDockerTLSConfig(
		s.cfg.TLSCACert,
		s.cfg.TLSCert,
		s.cfg.TLSKey,
		s.cfg.AllowInsecure,
	)

	return proxy.NewProxy(client.GetDockerURL(s.cfg.DockerURL), tlsConfig)
}
