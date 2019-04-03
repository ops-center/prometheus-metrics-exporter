#!/usr/bin/env bash

pushd $GOPATH/src/github.com/nightfury1204/prometheus-remote-metric-writer

gofmt -s -w *.go promquery

goimports -w *.go promquery

popd