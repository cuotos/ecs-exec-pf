> continuation of winebarrel/ecs-exec-pf

# ecs-exec-pf

Port forwarding using the ECS task container. (aws-cli wrapper)

## Usage

```
ecs-exec-pf - Port forwarding using the ECS task container. (aws-cli wrapper)

  Flags:
       --version      Displays the program version string.
    -h --help         Displays help with available flag, subcommand, and positional value parameters.
    -c --cluster      ECS cluster name.
    -t --task         ECS task ID.
    -n --container    Container name in ECS task.
    -p --port         Target remote port.
    -l --local-port   Client local port.
    -d --debug        Only print the commands that would be run.
```

## Installation

```sh
brew tap cuotos/ecs-exec-pf
brew install ecs-exec-pf
```

## Execution Example

```sh
$ ecs-exec-pf -c my-cluster -t 0113f61a4b1044d99c627daeee8c0d0c -p 80 -l 8080
Starting session with SessionId: root-03f56652a5f120d48
Port 8080 opened for sessionId root-03f56652a5f120d48.
Waiting for connections...
```

```
$ curl -s localhost:8080 | grep title
<title>Welcome to nginx!</title>
```

It is possible to connect to multiple ports at the same time but passing multiple `-l and -p` pairs. there must be the same number of them and they must be in the same order.

```sh
$ $ ecs-exec-pf -c my-cluster -t 0113f61a4b1044d99c627daeee8c0d0c -p 80 -l 8080 -p 81 -l 8081
Starting session with SessionId: root-03f56652a5f120d48

Starting session with SessionId: root-03f56652a5f120d49
Port 8080 opened for sessionId root-03f56652a5f120d48.
Waiting for connections...
Port 8081 opened for sessionId root-03f56652a5f120d49.
Waiting for connections...
```
