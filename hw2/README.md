# cloud hw2 (client and server with mongodb)
## reduce size of the images 
I use multi-stage build and scatch base image to delete unnecessary dependencies. 

## volumes 
there are two types of volumes: Named Volumes (docker managed) and Bind Mounts (bind directory you decide and bind it to the container)
create Named Volumes with `docker volume create <VOLUME_NAME>`. 
you can bind directory or files easily with Bind Mounts. you just neeed to bind it and don't need to create volume. 

## some commands 
access interactive shell in the container: `docker exec -it <CONTAINER_ID> /bin/bash `
list of volumes: `docker volume ls`
inspect volume and find MountPoint (you can access the volume on your hard-disk). it is usefull especially when we use Named Volume: `docker volume inspect <VOLUME_NAME>`
login to docker: `docker login`
save docker image to tar file: `docker save <IMAGE_NAME> SOME_NANME.tar `
load docker image from tar file: `docker load < SOME_NAME.tar `
run docker compose: `docker compose up -d`
view logs of the container: `docker logs <CONTAINER_ID>`
view live logs of the container: `docker logs -f <CONTAINER_ID>`
