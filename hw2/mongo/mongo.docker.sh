#! /bin/bash
docker run -p 27017:27017 -e MONGO_INITDB_ROOT_USERNAME=user -e MONGO_INITDB_ROOT_PASSWORD=pass -v ./data:/data/db/ mongo:latest
