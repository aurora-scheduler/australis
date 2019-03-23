## australis start

Start a service, maintenance on a host (DRAIN), a snapshot, or a backup.

### Synopsis

Start a service, maintenance on a host (DRAIN), a snapshot, or a backup.

### Options

```
  -h, --help   help for start
```

### Options inherited from parent commands

```
  -a, --caCertsPath string      Path where CA certificates can be found.
  -c, --clientCert string       Client certificate to use to connect to Aurora.
  -k, --clientKey string        Client key to use to connect to Aurora.
      --config string           Config file to use. (default "/etc/aurora/australis.yml")
  -l, --logLevel string         Set logging level [panic fatal error warning info debug trace]. (default "info")
  -p, --password string         Password to use for API authentication
  -s, --scheduler_addr string   Aurora Scheduler's address.
  -i, --skipCertVerification    Skip CA certificate hostname verification.
      --toJSON                  Print output in JSON format.
  -u, --username string         Username to use for API authentication
  -z, --zookeeper string        Zookeeper node(s) where Aurora stores information. (comma separated list)
```

### SEE ALSO

* [australis](australis.md)	 - australis is a client for Apache Aurora
* [australis start drain](australis_start_drain.md)	 - Place a list of space separated Mesos Agents into draining mode.
* [australis start maintenance](australis_start_maintenance.md)	 - Place a list of space separated Mesos Agents into maintenance mode.
* [australis start sla-drain](australis_start_sla-drain.md)	 - Place a list of space separated Mesos Agents into maintenance mode using SLA aware strategies.

###### Auto generated by spf13/cobra on 22-Mar-2019