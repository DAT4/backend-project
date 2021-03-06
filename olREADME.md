# Backend game

## Go routines

Go routines or coroutines are used for parallel programming and is used a lot in the application

* Each client will have his own
* The hub will have one


## Websocket

The websocket is used for live streaming the data of the game between the clients.

Websocket is nice because it allows the server to make callbacks. And provides a much more reactive program.

The websocket consists of a hub which has a number of clients who can register and unregister. When a client sends his coordinates to the server then the hub will register that a client has sent a message and therefore notify/broadcast all the other clients about the new coordinates of the player.


### Websocket over ssl

To connect to the websocket we have to connect with `ws://servername` but because the server is https only then we will write `wss://servername

### Websocket through nginx reverse proxy

We have to set the nginx configuration file up so that websocket is allowed. This is done by adding the following to  the nginx servers config file.
```
proxy_http_version 1.1;
proxy_set_header Upgrade $http_upgrade;
proxy_set_header Connection "Upgrade";
proxy_set_header Host $host;
```

_TODO: find what this upgrade thing is all about_

### Websocket behind JWT

Luckily for us the websocket can easily be wrapped with a jwt authentication function in go. Because the websocket request is for Go actually just understood a normal HTTP GET request, then we can add the header with the token that we got from an earlier request.

```go
r.Handle("/game", AuthMiddleware(http.HandlerFunc(Game))).Methods("GET")
```

_Note: it is important that the request is marked as a GET request in the gorilla router, else it will make an error._

### **(TODO)** Setting up a new hub for each game. 
Right now only one hub exists for the whole application. This is not how it should be in the future.

The plan is to create hubs dynamically for each 4 players who connect to a game.
