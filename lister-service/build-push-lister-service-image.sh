#!/bin/bash

appname=lister-service

docker build -f lister-service.dockerfile -t ismail118/$appname:1.0.0 .
docker push ismail118/$appname:1.0.0