# Coding challenge

## Goal

The goal of this exercise is to create a dead simple "chat" system. This system must be built in two parts:

* `Server` Receiving messages from a network interface (any kind) and forwarding them to all the connected clients
* `Client` A process reading a string on `STDIN` and forwarding it to the server, and also receiving messages from the same server and writing them to `STDOUT`.

## How to run this exercise

The code is written in `Golang`, and uses modules.

Therefore, you must have `Golang` version > 1.11 installed on your system to be able to build and execute this code.

* To launch the **Server**, issue the following command inside the `server/cmd` folder:
```
go run main.go
```

If the command is successful, you should be able to see the following output in your console:
```
2021/01/24 15:58:59 Chat Server listening on port 6697
```

* To run a new **Client**, issue the following command inside the `client/cmd` folder:
```
go run main.go
```

Each new client will have a username automatically generated, so if the command is successful, you will see this kind of output in your console:
```
021/01/24 16:01:45 Client connected to server / port 6697
Connected as Elizabeth Thomas
>
```

Then you can enter the text you want after the prompt (>) and see how the Server and the other Clients react.