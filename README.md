# Australis

A light-weight client for [Apache Aurora](https://aurora.apache.org/) built using [gorealis](https://github.com/paypal/gorealis).

## Usage 

```
A light-weight command line client for use with Apache Aurora built using gorealis.

Usage:
  australis [command]

Available Commands:
  create      Create an Aurora Job
  fetch       Fetch information from Aurora
  help        Help about any command
  kill        Kill an Aurora Job
  start       Start a service, maintenance on a host (DRAIN), a snapshot, or a backup.
  stop        Stop a service or maintenance on a host (DRAIN).

Flags:
  -a, --caCertsPath string      CA certificates path to use.
  -c, --clientCert string       Client certificate to use to connect to Aurora.
  -k, --clientKey string        Client key to use to connect to Aurora.
  -h, --help                    help for australis
  -i, --insecureSkipVerify      Skip verification.
  -p, --password string         Password to use for API authentication
  -s, --scheduler_addr string   Aurora Scheduler's address.
  -u, --username string         Username to use for API authentication
  -z, --zookeeper string        Zookeeper node(s) where Aurora stores information.

Use "australis [command] --help" for more information about a command.
```

## Sample commands:

### Fetching current leader
`australis fetch leader [ZK NODE 1] [ZK NODE 2]...[ZK NODE N]`

### Fetching status/mode of a agent
`australis fetch status [HOST 1] [HOST 2]...[HOST N]`

### Setting host to DRAIN:
`australis start drain [HOST 1] [HOST 2]...[HOST N]`

### Taking hosts out of DRAIN (End maintenance):
`australis stop drain [HOST 1] [HOST 2]...[HOST N]`

### Force a snapshot
`australis force snapshot`

### Force a backup
`australis force backup`

## Status
Australis is a work in progress and does not support all the features of Apache Aurora.
