# CLI Usage
```
NAME:
   dlm - Docker load monitor

USAGE:
   dlm [global options] command [command options] [arguments...]

VERSION:
   0.1.0

COMMANDS:
     stopped   view stopped containers
     launched  view launched containers
     logs      view container logs
     metrics   view containers metrics
     status    view API status
     help, h   Shows a list of commands or help for one command

GLOBAL OPTIONS:
   -v, --verbose           enable verbose output
   -a value, --addr value  set API address (default: "http://localhost:4222")
   --help, -h              show help
   --version, -V           print the version
```
## Container logs
```
$ dlm logs <container_id>
```
## Stopped contaiers
```
$ dlm stopped
```
## Launched containers
```
$ dlm launched
```
## Containers metrics
```
$ dlm metrics <container_id> <container_id>
$ dlm --addr http://localhost:4222 metrics all
```
## API status
```
$ dlm status
```
