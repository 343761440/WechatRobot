#!/bin/bash

cd ./master
docker-compose up -d 
cd ..

cd ./sentinel
docker-compose up -d 
cd ..

docker ps 
