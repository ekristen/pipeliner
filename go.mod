module github.com/ekristen/pipeliner

go 1.16

require (
	github.com/Microsoft/go-winio v0.5.1 // indirect
	github.com/buildkite/terminal-to-html/v3 v3.6.1
	github.com/bwmarrin/snowflake v0.3.0
	github.com/chartmuseum/storage v0.11.0
	github.com/containerd/containerd v1.5.8 // indirect
	github.com/docker/cli v20.10.11+incompatible // indirect
	github.com/docker/docker v20.10.11+incompatible // indirect
	github.com/docker/docker-credential-helpers v0.6.4 // indirect
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
	github.com/jinzhu/copier v0.3.2
	github.com/mattn/go-sqlite3 v2.0.3+incompatible // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/onsi/ginkgo v1.14.2 // indirect
	github.com/onsi/gomega v1.10.4 // indirect
	github.com/opencontainers/image-spec v1.0.2 // indirect
	github.com/rancher/lasso v0.0.0-20210616224652-fc3ebd901c08
	github.com/rancher/wrangler v0.8.9
	github.com/sirupsen/logrus v1.8.1
	github.com/tevino/abool v1.2.0 // indirect
	github.com/urfave/cli/v2 v2.3.0
	gitlab.com/ayufan/golang-cli-helpers v0.0.0-20210929155855-70bef318ae0a // indirect
	gitlab.com/gitlab-org/gitlab-runner v12.5.0+incompatible
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v1 v1.0.0-20140924161607-9f9df34309c0
	gopkg.in/yaml.v2 v2.4.0
	gorm.io/datatypes v1.0.3
	gorm.io/driver/mysql v1.2.0
	gorm.io/driver/sqlite v1.2.6
	gorm.io/gorm v1.22.3
	k8s.io/apiextensions-apiserver v0.22.4 // indirect
	k8s.io/apimachinery v0.22.4
	k8s.io/client-go v0.22.4
	k8s.io/utils v0.0.0-20211116205334-6203023598ed // indirect
)

replace github.com/docker/docker v1.4.2-0.20190822180741-9552f2b2fdde => github.com/docker/engine v1.4.2-0.20190822180741-9552f2b2fdde

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
