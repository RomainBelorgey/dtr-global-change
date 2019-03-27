# Change configuration for DTR on all replicas

## Usage

1 - Log in on a node with a DTR replica

2 - Get the replica id:
```
REPLICA_ID=$(docker ps -lf name='^/dtr-rethinkdb-.{12}$' --format '{{.Names}}' | cut -d- -f3)
```
3 - Execute this command line:
```
docker run -i --rm --net dtr-ol -v dtr-ca-$REPLICA_ID:/ca -e DTR_REPLICA_ID=$REPLICA_ID romainbelorgey/dtr-global-change
```
4 - You will need to make a `dtr reconfigure` to apply the changes to the containers

## Help

```
This command will change the configuration of all replicas.

/!\ Please do a backup before using it !

You will need to do a "dtr reconfigure" to apply globally these changes.

Usage:
  dtr-global-change [flags]

Flags:
  -h, --help                     help for dtr-global-change
      --http-port int            Http port that will use all replicas
      --https-port int           Https port that will use all replicas
      --replica-id string        Replica-id to connect (default "4d1e2c7382c5")
      --rethinkdb-cache-mb int   Max rethinkdb memory cache that will use all replicas | 0 = auto (default -1)
  -t, --toggle                   Help message for toggle
```

## To compile

```
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build
```
