##Redis
` docker run --name some-redis -p 6379:6379 -d --rm redis --loglevel debug`


`docker build -t gohome .`

`docker run --name gohomecontainer -p 3000:3000 --rm gohome`


`let ws = new WebSocket("ws://localhost:3000/ws")`

`ws.onmessage = (event) => console.log(event.data)`

`ws.send("Hello, ASL?")`


## Using Dockerfile
 - ENV for using Redis is `host.docker.internal:6379`
 - `host.docker.internal` - allows to access the ports of host container