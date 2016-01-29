#!/bin/bash

make

./.output/atau -o "./.output/go" -l go -m "sample" ./samples/urlshortener.json

pushd ./.output/go
go build .
popd
