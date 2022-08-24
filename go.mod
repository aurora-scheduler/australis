module github.com/aurora-scheduler/australis

go 1.15

require (
	github.com/aurora-scheduler/gorealis/v2 v2.29.0
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/cobra v1.0.0
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.6.3
	github.com/stretchr/testify v1.5.0
	gopkg.in/yaml.v2 v2.2.8
)

replace github.com/apache/thrift v0.13.0 => github.com/ridv/thrift v0.13.2
