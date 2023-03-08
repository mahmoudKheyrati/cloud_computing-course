#! /bin/bash
# build projects
docker build -t mahmoudreza/cloud-hw2-server ./server

# push images to docker hub
echo psuh images to docker hub
docker push mahmoudreza/cloud-hw2-server 

# create volumes 
echo create volumes 
docker volume create servervol

# create network 
echo create netowrk
docker network create cloud-hw2

# run server
echo run server
docker run -v servervol:/serverdata --network cloud-hw2 --network-alias server -e MONGO_USERNAME=user MONGO_PASSWORD=pass MONGO_HOST=mongo MONGO_PORT=27017 -p 3000:3000 mahmoudreza/cloud-hw2-server
