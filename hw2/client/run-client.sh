#! /bin/bash
# build projects
echo build projects ... 
docker build -t mahmoudreza/cloud-hw2-client ./client


# push images to docker hub
echo psuh images to docker hub
docker push mahmoudreza/cloud-hw2-client 


# create volumes 
echo create volumes 
docker volume create clientvol

# create network 
echo create netowrk
docker network create cloud-hw2

# run client 
echo run client
docker run -v clientvol:/clientdata --network cloud-hw2 --network-alias client -e SERVER_HOST=server SERVER_PORT=3000 SAVE_FILE_PATH=/clientdata mahmoudreza/cloud-hw2-client 
