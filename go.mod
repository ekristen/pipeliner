module github.com/ekristen/pipeliner

go 1.16

require (
	github.com/Azure/go-ansiterm v0.0.0-20210617225240-d185dfc1b5a1 // indirect
	github.com/Microsoft/go-winio v0.5.1 // indirect
	github.com/buildkite/terminal-to-html/v3 v3.6.1
	github.com/bwmarrin/snowflake v0.3.0
	github.com/certifi/gocertifi v0.0.0-20200922220541-2c3bb06c6054 // indirect
	github.com/chartmuseum/storage v0.11.0
	github.com/cockroachdb/datadriven v0.0.0-20200714090401-bf6692d28da5 // indirect
	github.com/containerd/containerd v1.4.12 // indirect
	github.com/coreos/go-systemd/v22 v22.3.2 // indirect
	github.com/docker/docker v20.10.11+incompatible // indirect
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/machine v0.16.2 // indirect
	github.com/go-http-utils/etag v0.0.0-20161124023236-513ea8f21eb1
	github.com/go-http-utils/fresh v0.0.0-20161124030543-7231e26a4b27 // indirect
	github.com/go-http-utils/headers v0.0.0-20181008091004-fed159eddc2a // indirect
	github.com/google/uuid v1.3.0
	github.com/gorhill/cronexpr v0.0.0-20180427100037-88b0669f7d75
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/sessions v1.2.1
	github.com/gorilla/websocket v1.4.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.16.0 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/jinzhu/copier v0.3.2
	github.com/jonboulle/clockwork v0.2.2 // indirect
	github.com/mattn/go-sqlite3 v2.0.3+incompatible // indirect
	github.com/moby/term v0.0.0-20210610120745-9d4ed1856297 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.0.2 // indirect
	github.com/prometheus/common v0.10.0 // indirect
	github.com/prometheus/procfs v0.6.0 // indirect
	github.com/rancher/lasso v0.0.0-20210616224652-fc3ebd901c08
	github.com/rancher/wrangler v0.8.9
	github.com/sirupsen/logrus v1.8.1
	github.com/soheilhy/cmux v0.1.5 // indirect
	github.com/stevenle/topsort v0.2.0 // indirect
	github.com/tevino/abool v1.2.0 // indirect
	github.com/tmc/grpc-websocket-proxy v0.0.0-20201229170055-e5319fda7802 // indirect
	github.com/urfave/cli v1.22.2 // indirect
	github.com/urfave/cli/v2 v2.3.0
	gitlab.com/gitlab-org/gitlab-runner v12.5.0+incompatible
	go.etcd.io/bbolt v1.3.6 // indirect
	go.etcd.io/etcd v0.5.0-alpha.5.0.20200910180754-dd1b699fc489 // indirect
	go.uber.org/zap v1.17.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/square/go-jose.v2 v2.5.1 // indirect
	gopkg.in/yaml.v1 v1.0.0-20140924161607-9f9df34309c0
	gopkg.in/yaml.v2 v2.4.0
	gorm.io/datatypes v1.0.3
	gorm.io/driver/mysql v1.2.0
	gorm.io/driver/sqlite v1.2.6
	gorm.io/gorm v1.22.3
	gotest.tools/v3 v3.0.3 // indirect
	k8s.io/api v0.22.4
	k8s.io/apiextensions-apiserver v0.19.0-alpha.3 // indirect
	k8s.io/apimachinery v0.22.4
	k8s.io/client-go v0.22.4
	k8s.io/code-generator v0.22.4 // indirect
	k8s.io/utils v0.0.0-20211116205334-6203023598ed // indirect
)

//replace github.com/docker/docker v1.4.2-0.20190822180741-9552f2b2fdde => github.com/docker/engine v1.4.2-0.20190822180741-9552f2b2fdde

replace github.com/minio/go-homedir v0.0.0-20190425115525-017018655514 => gitlab.com/steveazz/go-homedir v0.0.0-20190425115525-017018655514

// Added this for chartmuseum/storage
replace github.com/NetEase-Object-Storage/nos-golang-sdk => github.com/karuppiah7890/nos-golang-sdk v0.0.0-20191116042345-0792ba35abcc

// Added this for chartmuseum/storage
//replace go.etcd.io/etcd => github.com/eddycjy/etcd v0.5.0-alpha.5.0.20200218102753-4258cdd2efdf
replace go.etcd.io/etcd => github.com/eddycjy/etcd v0.5.0-alpha.5.0.20200129063535-5d07a202ae30

// Added because of https://github.com/etcd-io/etcd/issues/12124
replace google.golang.org/grpc => google.golang.org/grpc v1.29.1

// Added because of https://github.com/go-gorm/gorm/issues/4719, can be removed after this is fixed
replace gorm.io/gorm => gorm.io/gorm v1.21.14

// Added because of https://github.com/rancher/wrangler/pull/186 (off of branch 0.8.7 because 0.8.8 introduced breaking changes :()
replace github.com/rancher/wrangler => github.com/ekristen/wrangler v0.4.1-0.20211125193741-094c487f3bb2

replace gitlab.com/gitlab-org/gitlab-runner => gitlab.com/gitlab-org/gitlab-runner v1.11.1-0.20211121164806-f0a95a76c6db
