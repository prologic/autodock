package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	pkgver "github.com/prologic/autodock/version"

	"github.com/namsral/flag"
	"github.com/prologic/autodock/config"
	"github.com/prologic/autodock/server"
)

func main() {
	var (
		dockerurl string
		msgbusurl string

		tlsverify bool
		tlscacert string
		tlscert   string
		tlskey    string
		tls       bool

		debug   bool
		version bool

		bind string
	)

	flag.String(flag.DefaultConfigFlagname, "", "path to config file")
	flag.BoolVar(&debug, "debug", false, "enable debug logging")
	flag.BoolVar(&version, "v", false, "display version information")

	flag.StringVar(
		&bind, "bind", "0.0.0.0:8000",
		"[int]:<port> to bind to for HTTP",
	)

	flag.StringVar(&dockerurl, "dockerurl", "", "Docker URL to connect to")
	flag.StringVar(&msgbusurl, "msgbusurl", "", "MessageBus URL to connect to")

	flag.BoolVar(&tls, "tls", false, "Use TLS; implied by --tlsverify")
	flag.StringVar(
		&tlscacert, "tlscacert", "",
		"Trust certs signed only by this CA",
	)
	flag.StringVar(
		&tlscert, "tlscert", "",
		"Path to TLS certificate file",
	)
	flag.StringVar(
		&tlskey, "tlskey", "",
		"Path to TLS key file",
	)
	flag.BoolVar(
		&tlsverify, "tlsverify", true,
		"Use TLS and verify the remote",
	)

	flag.Parse()

	if version {
		fmt.Printf("autodock v%s", pkgver.FullVersion())
		os.Exit(0)
	}

	if debug {
		log.SetLevel(log.DebugLevel)
	}

	cfg := &config.Config{
		Debug: debug,

		Bind: bind,

		DockerURL:     dockerurl,
		MsgBusURL:     msgbusurl,
		TLSCACert:     tlscacert,
		TLSCert:       tlscert,
		TLSKey:        tlskey,
		AllowInsecure: !tlsverify,
	}

	srv, err := server.NewServer(cfg)
	if err != nil {
		log.Fatal(err)
	}

	err = srv.EnableCollector()
	if err != nil {
		log.Fatalf("error enabling collector: %s", err)
	}

	if cfg.MsgBusURL == "" {
		srv.EnableMessageBus()
	}

	err = srv.EnableProxy()
	if err != nil {
		log.Fatalf("error enabling proxy: %s", err)
	}

	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
