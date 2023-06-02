#!/bin/bash

appname=logger-service

docker build -f logger-service.dockerfile -t ismail118/$appname:1.0.0 .
docker push ismail118/$appname:1.0.0