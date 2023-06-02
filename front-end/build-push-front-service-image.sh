#!/bin/bash

docker build -f front-end.dockerfile -t ismail118/front-service:1.0.0 .
docker push ismail118/front-service:1.0.0