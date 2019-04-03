#!/usr/bin/env bash

pushd $GOPATH/src/github.com/searchlight/prometheus-metrics-exporter

gofmt -s -w *.go promquery

goimports -w *.go promquery

popd