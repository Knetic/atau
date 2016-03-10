#!/bin/bash

make

./.output/atau -o "./.temp/go" -l go -m "sample" ./samples/urlshortener.json
./.output/atau -o "./.temp/cs" -l cs -m "sample" ./samples/urlshortener.json

pushd ./.temp/go
go build .
popd
pushd ./.temp/cs
mcs $(find . -name "*.cs") --parse
popd
