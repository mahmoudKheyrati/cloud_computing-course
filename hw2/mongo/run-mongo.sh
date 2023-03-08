#! /bin/bash

# create volumes 
docker volume create mongodbvol

# create network 
echo create netowrk
docker network create cloud-hw2

# run mongodb 
echo run mongodb
docker run  --network cloud-hw2 --network-alias mongo  -v mongodbvol:/data/db/ -e MONGO_INITDB_ROOT_USERNAME=user -e MONGO_INITDB_ROOT_PASSWORD=pass mongo:latest
