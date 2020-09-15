# portal

![CI](https://github.com/gauravgahlot/tink-wizard/workflows/CI/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/gauravgahlot/tink-wizard)](https://goreportcard.com/report/github.com/gauravgahlot/tink-wizard)
![Experimental](https://camo.githubusercontent.com/a9257bfeb095645580d7e52bb901e033c8f479a9/68747470733a2f2f696d672e736869656c64732e696f2f62616467652f73746162696c6974792d6578706572696d656e74616c2d7265642e737667)

This repository contains a general-purpose web UI for Tinkerbell.
It allows you to manage your hardware, template, and workflows from a single place.

## Experimental

This repository is Experimental meaning that it's based on untested ideas or techniques and not yet established or finalized or involves a radically new and innovative style! This means that support is best effort (at best!) and we strongly encourage you to NOT use this in production.

## Prerequisite

-   You have already setup the Tinkerbell stack (the provisioner, as we generally know it).

## Get Started

### Clone the repository with:

```
$ git clone https://github.com/tinkerbell/portal.git && cd portal
```

### Environment settings

-   Update the `.env` file as per your environment setup.
-   If you plan to use TLS with tink-wizard, set `ALLOW_INSECURE` to `false`. Default is `true`.

### Starting server

-   In order to host a secure server, please update the CA and server configuration in respective files under `/tls`.
-   However, the definitions are preset to give you a quick start with `localhost`.
-   Once you have updated the `/tls` definitions, you can generate the certificate with:

```
$ make certs
```

-   You can now start the server by executing the following:

```
$ make redis && make run
```

-   You can now access tink-wizard at [http://localhost:7676](http://localhost:7676) or [https://localhost:7676](https://localhost:7676).

### Resources

-   [Introducing Tink-Wizard - A general purpose Web UI to manage your Tinkerbell workflows
    ](https://www.youtube.com/watch?v=SLshLxNvgC0&feature=youtu.be)
