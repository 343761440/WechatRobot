#! /bin/bash

docker rmi wechatrobot:latest
docker image build -t wechatrobot:latest . 