Use as a proxy to a UCI engine that is running on a remote server.

## Setup

You will need to pull down the dependencies for Go:

```
go get ssh-engine
```

## Configuration

Create a filed called `engine.yml` in the same directory. The contents should look like this but for your configuration:

```yml
user: "matt"
privateKeyFile: "/Users/nohr/.ssh/id_rsa"
host: "127.0.0.1"
port: "22"
remoteCommand: "/home/matt/bin/stockfish"
```

## Running

Run the proxy:

```
go run SshEngine.go
```

## Building

To build an executable for Windows:

```
env GOOS=windows GOARCH=386 go build SshEngine.go
```

This will create SshEngine.exe. Copy that along with the config file to a suitable directory on Windows.