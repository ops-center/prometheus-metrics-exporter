#!/usr/bin/env bash

pushd $GOPATH/src/github.com/searchlight/prometheus-metrics-exporter

docker build -t nightfury1204/prometheus-remote-metric-writer:canary .

popd