package main

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/tinkerbell/portal/src/client"
	"github.com/tinkerbell/portal/src/pkg"
	"github.com/tinkerbell/portal/src/pkg/controller"
	_ "github.com/tinkerbell/portal/src/pkg/redis"
)

const (
	applicationPort  = ":7676"
	serverCRT        = "tls/server.pem"
	serverKey        = "tls/server-key.pem"
	envAllowInsecure = "ALLOW_INSECURE"

	infoServerListening = "server listening at port %v"
	errAllowInsecure    = "failed to parse env variable ALLOW_INSECURE"
	errTLSSetup         = "failed to setup TLS"
	errServer           = "failed to start the server"
	warnHostingInsecure = "hosting an insecure server"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
	client.Init()
}

func main() {
	templates := pkg.PopulateTemplates()
	controller.Startup(templates)
	startListening()
}

func startListening() {
	allowInsecure := false
	if env := os.Getenv(envAllowInsecure); env != "" {
		insecure, err := strconv.ParseBool(env)
		if err != nil {
			log.Fatal(errAllowInsecure)
		}
		allowInsecure = insecure
	}

	pwd, _ := os.Getwd()
	log.Infof(infoServerListening, applicationPort)
	if !allowInsecure {
		log.Fatal(
			http.ListenAndServeTLS(
				applicationPort,
				filepath.Join(pwd, serverCRT),
				filepath.Join(pwd, serverKey),
				nil,
			))
	} else {
		log.Warn(warnHostingInsecure)
		log.Fatal(http.ListenAndServe(applicationPort, nil))
	}
}
