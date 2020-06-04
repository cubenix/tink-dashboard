#!/bin/sh

cat /tmp/hardware-one.json | tink hardware push
cat /tmp/hardware-two.json | tink hardware push
tink template create -n hello -p /tmp/hello-world.tmpl
tink template create -n sample -p /tmp/sample.tmpl

tink workflow create -r '{"worker_1": "ec:0d:9a:bf:ff:dc"}' -t 
tink workflow create -r '{"worker_1": "ec:0d:9a:bf:ff:dc", "worker_2": "9a:00:15:98:0d:a0"}' -t 
