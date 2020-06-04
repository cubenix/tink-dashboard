# tink-wizard 

![CI](https://github.com/gauravgahlot/tink-wizard/workflows/CI/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/gauravgahlot/tink-wizard)](https://goreportcard.com/report/github.com/gauravgahlot/tink-wizard)

Tink-Wizard is a general-purpose web UI for Tinkerbell.
It allows you to manage your hardware, template, and workflows from a single place.

## Prerequisite
 - You have already setup the Tinkerbell stack (the provisioner, as we generally know it).

## Get Started

### Clone the repository with:

```
$ git clone https://github.com/gauravgahlot/tink-wizard.git && cd tink-wizard
```

### Environment settings

 - Update the `.env` file as per your environment setup.
 - If you plan to use TLS with tink-wizard, set `ALLOW_INSECURE` to `false`. Default is `true`.
 
### Starting server 

 - In order to host a secure server, please update the CA and server configuration in respective files under `/tls`.
 - However, the definitions are preset to give you a quick start with `localhost`.
 - Once you have updated the `/tls` definitions, you can generate the certificate with:
```
$ make certs
```

 - You can now start the server by executing the following:
```
$ make redis && make run
```
 - You can now access tink-wizard at [http://localhost:7676](http://localhost:7676) or [https://localhost:7676](https://localhost:7676).
