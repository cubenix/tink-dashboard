#!/usr/bin/env sh

if ! { [[ -r certs/ca.pem ]] && [[ -r certs/ca-key.pem ]]; }; then
	cfssl gencert -initca ca-csr.json | cfssljson -bare ca
fi
if ! { [[ -r server-csr.json ]] && [[ -r server.pem ]] && [[ -r server-key.pem ]]; }; then
	cfssl gencert \
		-ca=ca.pem \
		-ca-key=ca-key.pem \
		-config=ca-config.json \
		-profile=server \
		server-csr.json | cfssljson -bare server
fi
cp *.pem /certs
