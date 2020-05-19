package client

import (
	"crypto/x509"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gauravgahlot/tink-wizard/src/pkg/redis"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/tinkerbell/tink/protos/hardware"
	"github.com/tinkerbell/tink/protos/template"
	"github.com/tinkerbell/tink/protos/workflow"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// gRPC clients
var (
	cache *redis.Cache

	templateClient template.TemplateClient
	hardwareClient hardware.HardwareServiceClient
	workflowClient workflow.WorkflowSvcClient
)

// Init initializes a gRPC connection with server
func Init() {
	conn, err := getConnection()
	if err != nil {
		log.Fatal(err)
	}
	templateClient = template.NewTemplateClient(conn)
	hardwareClient = hardware.NewHardwareServiceClient(conn)
	workflowClient = workflow.NewWorkflowSvcClient(conn)
}

// GetConnection returns a gRPC client connection
func getConnection() (*grpc.ClientConn, error) {
	certURL := os.Getenv("TINKERBELL_CERT_URL")
	if certURL == "" {
		return nil, errors.New("undefined TINKERBELL_CERT_URL")
	}

	resp, err := http.Get(certURL)
	if err != nil {
		return nil, errors.Wrap(err, "fetch cert")
	}
	defer resp.Body.Close()

	certs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "read cert")
	}
	cp := x509.NewCertPool()
	ok := cp.AppendCertsFromPEM(certs)
	if !ok {
		return nil, errors.Wrap(err, "parse cert")
	}

	grpcAuthority := os.Getenv("TINKERBELL_GRPC_AUTHORITY")
	if grpcAuthority == "" {
		return nil, errors.New("undefined TINKERBELL_GRPC_AUTHORITY")
	}
	creds := credentials.NewClientTLSFromCert(cp, "")
	conn, err := grpc.Dial(grpcAuthority, grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, errors.Wrap(err, "connect to tinkerbell server")
	}
	return conn, nil
}

func init() {
	cache = redis.Instance()
}
