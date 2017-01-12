#!/bin/sh

# compile sources with linux target.
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o statspout-linux github.com/mijara/statspout

# build the docker image.
docker build -t mijara/statspout .

