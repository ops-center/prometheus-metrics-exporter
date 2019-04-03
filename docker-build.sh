#!/usr/bin/env bash

pushd $GOPATH/src/github.com/nightfury1204/prometheus-remote-metric-writer

docker build -t nightfury1204/prometheus-remote-metric-writer:canary .

popd