#!/bin/bash

appname=broker-service

docker build -f broker-service.dockerfile -t ismail118/$appname:1.0.0 .
docker push ismail118/$appname:1.0.0