Use as a proxy to a UCI engine that is running on a remote server.

## Setup

```
go get ssh-engine
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