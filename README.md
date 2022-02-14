# Blogo

This is a blog app with frontend and backend.

Frontend is a Flask app in Python

Backend is a Go API

Data is stored in Redis

## API
Please create a **multistage** Dockerfile to:
- build go code in /api:
    - use image: `golang:alpine`
    - commands `got get -d -v && go build -o api main.go`
- package the resulting `api` binary in a vanilla `alpine` Docker image
- run the container as user `appuser`
- make the packaged binary the entrypoint of the container

## Frontend
- Create a Dockerfile in `/front`
    - Based on `python:alpine`
    - Add all code to /usr/src/app
    - Install required modules : `pip install -r requirements.txt`
    - run the container as user `appuser`
    - Define the command to run: `python3 front.py`

## Compose
Write a docker-compose file that would contain 4 services:
- front - serving the fronted app on port 80 of your host
    - Depends on `api`
    - Note - the frontend app needs the following env variables:
        - PORT (which port it should listen on)
        - API_PORT (which port the api is listening on)
    - Frontend is looking for api on hostname `api`
    - Configure a healthcheck to call the app at `/healthz`
- api  - serving the api on docker bridge network
    - Depends on `redis`
    - Note - the api needs the following env variables:
        - PORT (which port it should listen on)
        - REDIS_HOST (which hostname redis is available on)
        - REDIS_PORT (which port redis is listening on)
    - Configure a healthcheck to call the api at `/healthz`
- redis - use the official redis image
    - mount a host directory into `/data` to persist redis data if container is deleted
- cadvisor - run the official cadvisor image mounting all needed Docker host directories
    - expose cadvisor on port 9999 of your host

**NOTE** - create 2 networks in compose:
  - front - for front, api and cadvisor
  - back - for api and redis
## Bonus: 
Run all containers only with needed capabilities
