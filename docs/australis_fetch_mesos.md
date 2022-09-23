## australis fetch mesos

Fetch information from Mesos.

### Synopsis

Fetch information from Mesos.

### Options

```
  -h, --help   help for mesos
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
  -t, --timeout duration        Gorealis timeout. (default 20s)
      --toJSON                  Print output in JSON format.
  -u, --username string         Username to use for API authentication
  -z, --zookeeper string        Zookeeper node(s) where Aurora stores information. (comma separated list)
```

### SEE ALSO

* [australis fetch](australis_fetch.md)	 - Fetch information from Aurora
* [australis fetch mesos leader](australis_fetch_mesos_leader.md)	 - Fetch current Mesos-master leader given Zookeeper nodes.
* [australis fetch mesos master](australis_fetch_mesos_master.md)	 - Fetch current Mesos-master nodes/leader given Zookeeper nodes.

###### Auto generated by spf13/cobra on 8-Sep-2022