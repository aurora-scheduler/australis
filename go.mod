module github.com/aurora-scheduler/australis

require (
	github.com/aurora-scheduler/gorealis/v2 v2.21.4
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/cobra v0.0.5
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.6.2
	gopkg.in/yaml.v2 v2.2.8
)

go 1.13

replace github.com/apache/thrift v0.12.0 => github.com/ridv/thrift v0.12.2
