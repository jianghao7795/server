#!/bin/bash

cp conf/config.yaml config.yaml
sed -i -e "s/127.0.0.1/host.containers.internal/g" config.yaml
podman stop fiber
podman rm fiber
podman rmi fiber

podman build --progress=plain -t fiber .
podman run --name fiber -d -v ./logs/:/app/logs/:z -v ./uploads:/app/uploads:z -p 3100:3100 fiber
