# Cloud Computing hw#1 ( mahmoud reza kheyrati fard - stdNo: 9832110 )

### dockerhub image link: [link](https://hub.docker.com/r/mahmoudreza/cloud-hw1-vidly)

## Dockerfile
``` dockerfile
FROM node:alpine
LABEL "maintainer"="@mahmoud_dev"
WORKDIR /application
COPY package.json .
RUN npm install
COPY . .
EXPOSE 3000
CMD npm start
```
I use node:alpine as base image. the base image must have the dependencies and the compiler for the project.
with `LABEL` command we can add some metadata to the image. 
we can set a directory as work directory using `WORKDIR`  command. every command that executes after this command run at the workdir without need to use `cd` command or absoloute path. 
first I copy the `package.json` that have all libraries and dependencies list. 
after that we should install the dependencies using `npm install` command. 
after installing dependencies, I copy the source-code into working directory, then exposing port=3000. 
after all, running the application using `npm start` command. 

## docker-compose
we can use docker-compose to run container: 
``` yaml
version: '3'
services:
  vidly:
    hostname: "vidly-the-simple-react-application"
    build:
      context: .
    ports:
      - 3000:3000
```
## some useful commands to work with docker
``` bash
#!/bin/bash
# build docker image
docker build . -t mahmoudreza/cloud-hw1-vidly

# login to dockerhub
docker login 

# push image to dockerhub or another docker registery
docker image push mahmoudreza/cloud-hw1-vidly

# pull image form registry
docker pull mahmoudreza/cloud-hw1-vidly

# run docker container
docker run -p 3000:3000 mahmoudreza/cloud-hw1-vidly 

# run using docker compose
docker compose up 

# if you want to run the container in the background you can use -d flag in the commands above

# save docker image to tar file ( if you don't want to push to any registry )
docker save docker.io/mahmoudreza/cloud-hw1-vidly > cloud-hw1-vidly.docker.tar

# load image to docker from tar file
docker load < cloud-hw1-vidly.docker.tar

```

