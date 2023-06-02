#!/bin/bash

appname=authentication-service

docker build -f authentication-service.dockerfile -t ismail118/$appname:1.0.0 .
docker push ismail118/$appname:1.0.0