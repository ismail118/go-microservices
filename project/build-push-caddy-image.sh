#!/bin/bash

docker build -f caddy.dockerfile -t ismail118/micro-caddy:1.0.0 .
docker push ismail118/micro-caddy:1.0.0