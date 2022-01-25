When using some chess software like Chessbase, you may need an executable that is a proxy to a remote chess engine. This simple program is meant just for that.

Specifically, it builds an .exe that you can point to that just connects to a remote server using SSH. There is an external config file that you use so you can change the settings without needing to build a new .exe file.

This will work on linux of MacOS as well, but for those systems it may be easier to make your own script. But, if you want to use this, just build it for those environments.

More details can be found here: https://mattplayschess.com/ssh-engine/

## Setup

You will need to pull down the dependencies for Go:

```
go get ssh-engine
```

## Configuration

Create a filed called `engine.yml` in the same directory. The contents should look like this but for your configuration:

Mac/Linux:

```yml
user: "matt"
privateKeyFile: "/Users/matt/.ssh/stockfish-keypair.pem"
host: "123.45.67.8"
port: "22"
remoteCommand: "stockfish"
```

Windows:

```yml
user: "matt"
privateKeyFile: "C:\\Users\\matt\\.ssh\\stockfish-keypair.pem"
host: "123.45.67.8"
port: "22"
remoteCommand: "stockfish"
```

If you want to enable extra logging, add this to the configuration file with the name of your log file:

```yml
logFileName: "engine.log"
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

## Making a Release

Create a tag (format `0.0.0`) and the CI pipeline will automatically build a Windows .exe and create a release

```
git tag 0.0.3
git push origin --tags
```
