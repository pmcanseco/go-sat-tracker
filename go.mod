module github.com/pmcanseco/go-sat-tracker

go 1.18

require (
	github.com/gocarina/gocsv v0.0.0-20230123225133-763e25b40669
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/mailru/easyjson v0.7.7
	github.com/onsi/ginkgo v1.2.1-0.20160509182050-5437a97bf824
	github.com/onsi/ginkgo/v2 v2.11.0
	github.com/onsi/gomega v1.27.8
	github.com/pmcanseco/go-satellite v0.0.7
	golang.org/x/image v0.0.0-20210628002857-a66eb6448b8d
	tinygo.org/x/drivers v0.25.0
)

require (
	github.com/go-logr/logr v1.2.4 // indirect
	github.com/go-task/slim-sprig v0.0.0-20230315185526-52ccab3ef572 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/google/pprof v0.0.0-20230602150820-91b7bce49751 // indirect
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/stretchr/testify v1.8.0 // indirect
	golang.org/x/net v0.10.0 // indirect
	golang.org/x/sys v0.9.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	golang.org/x/tools v0.9.3 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace tinygo.org/x/drivers => github.com/pmcanseco/drivers v0.24.1-0.20230605010524-52b1338e91db
