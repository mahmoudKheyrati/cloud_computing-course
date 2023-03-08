#! /bin/bash
# build projects
echo build projects ... 
# docker build -t mahmoudreza/cloud-hw2-client ./client
# docker build -t mahmoudreza/cloud-hw2-server ./server

# push images to docker hub
echo psuh images to docker hub
docker push mahmoudreza/cloud-hw2-client 
docker push mahmoudreza/cloud-hw2-server

# create volumes 
echo create volumes 
docker volume create clientvol
docker volume create servervol
docker volume create mongodbvol

# create network 
echo create netowrk
docker network create cloud-hw2

# run mongodb 
echo run mongodb
docker run  --network cloud-hw2 --network-alias mongo  -v mongodbvol:/data/db/ -e MONGO_INITDB_ROOT_USERNAME=user -e MONGO_INITDB_ROOT_PASSWORD=pass mongo:latest

# run server
echo run server
docker run -v servervol:/serverdata --network cloud-hw2 --network-alias server -e MONGO_USERNAME=user -e MONGO_PASSWORD=pass -e MONGO_HOST=mongo -e MONGO_PORT=27017 -p 3000:3000 mahmoudreza/cloud-hw2-server

# run client 
echo run client
docker run -v clientvol:/clientdata --network cloud-hw2 --network-alias client -e SERVER_HOST=server -e SERVER_PORT=3000 -e SAVE_FILE_PATH=/clientdata mahmoudreza/cloud-hw2-client
