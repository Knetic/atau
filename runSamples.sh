#!/bin/bash

make

./.output/atau -o "./.temp/go" -l go -m "sample" ./samples/urlshortener.json

pushd ./.temp/go
go build .
popd
