FROM ubuntu:latest

RUN mkdir /myapp

ADD bin/wxmanager /myapp

WORKDIR /myapp