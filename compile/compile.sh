#!/usr/bin/env bash

MODE=$1

if [[ "${MODE}" == "bin" ]]; then
   go version && go env && gcc -v && \
     CGO_ENABLED=0 GOOS=linux GOARCH=amd64  go build \
     --mod=vendor  -v -o dpos cmd/main.go
elif [[ "${MODE}" == "docker" ]]; then
    go mod vendor
    # 删除编译镜像
    docker rmi -f dpos:v1.0
    docker rmi -f $(docker images | grep "none" | awk '{print $3}')
    #编译镜像并启动
    docker build -f images/Dockerfile -t dpos:v1.0 .
    rm -rf vendor
else
  echo "para is empty bin or docker"
  exit 1
fi
