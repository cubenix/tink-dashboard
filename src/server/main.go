package main

import (
	"log"
	"net/http"

	"github.com/gauravgahlot/frawn/src/pkg"
	"github.com/gauravgahlot/frawn/src/pkg/controller"
)

const applicationPort = ":5000"

func main() {
	templates := pkg.PopulateTemplates()
	controller.Startup(templates)

	server := http.Server{
		Addr: applicationPort,
	}
	log.Println("server listening at port", applicationPort)
	log.Fatal(server.ListenAndServe())
}
