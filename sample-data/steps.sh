#!/bin/sh

cat /tmp/hardware-one.json | tink hardware push
cat /tmp/hardware-two.json | tink hardware push
tink template create -n hello -p /tmp/hello-world.tmpl
tink template create -n sample -p /tmp/sample.tmpl

tink workflow create -t d87d2c0c-0dc8-4e3c-83dc-af69eb688d36 -r '{"worker_1": "ec:0d:9a:bf:ff:dc"}'
tink workflow create -t e4243328-b06f-4ccc-823a-4df05c8157e7 -r '{"worker_1": "ec:0d:9a:bf:ff:dc", "worker_2": "9a:00:15:98:0d:a0"}'
