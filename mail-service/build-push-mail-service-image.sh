#!/bin/bash

appname=mail-service

docker build -f mail-service.dockerfile -t ismail118/$appname:1.0.0 .
docker push ismail118/$appname:1.0.0