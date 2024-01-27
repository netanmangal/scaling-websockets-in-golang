## Go Program
    - `go run main.go`
    - Used to 

## Dockerfile
    - `docker build -t golang-docker .`
    - `docker run -p 3000:3000 golang-docker`
    - Used to build a docker image
    - `ENV PORT=3000` has been set in the Dockerfile
    - `EXPOSE 3000` has been set in the Dockerfile. However, dockerfile can't publish ports. It can only expose them.

## DockerCompose
    - `ports: - "3000:3000"` has been published in the docker-compose.yaml file
    - Used to build a docker image and run it
    - `docker-compose up`
