module github.com/ekristen/pipeliner

go 1.16

require (
	github.com/Microsoft/go-winio v0.4.14 // indirect
	github.com/Nvveen/Gotty v0.0.0-20120604004816-cd527374f1e5 // indirect
	github.com/buildkite/terminal-to-html/v3 v3.4.1-0.20200917034757-9577c9f1146b
	github.com/bwmarrin/snowflake v0.3.0
	github.com/chartmuseum/storage v0.11.0
	github.com/containerd/continuity v0.0.0-20200928162600-f2cc35102c2a // indirect
	github.com/docker/cli v0.0.0-20181219132003-336b2a5cac7f // indirect
	github.com/docker/docker v1.4.2-0.20190822180741-9552f2b2fdde // indirect
	github.com/docker/docker-credential-helpers v0.6.3 // indirect
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-metrics v0.0.1 // indirect
	github.com/docker/libtrust v0.0.0-20160708172513-aabc10ec26b7 // indirect
	github.com/docker/machine v0.16.2 // indirect
	github.com/go-http-utils/etag v0.0.0-20161124023236-513ea8f21eb1
	github.com/go-http-utils/fresh v0.0.0-20161124030543-7231e26a4b27 // indirect
	github.com/go-http-utils/headers v0.0.0-20181008091004-fed159eddc2a // indirect
	github.com/google/uuid v1.3.0
	github.com/gorhill/cronexpr v0.0.0-20180427100037-88b0669f7d75
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/sessions v1.2.1
	github.com/gorilla/websocket v1.4.2
	github.com/jinzhu/copier v0.2.3
	github.com/kr/text v0.2.0 // indirect
	github.com/mattn/go-sqlite3 v2.0.3+incompatible // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/opencontainers/image-spec v1.0.1 // indirect
	github.com/opencontainers/runc v0.1.1 // indirect
	github.com/prometheus/client_golang v1.9.0 // indirect
	github.com/rancher/wrangler v0.8.6
	github.com/sirupsen/logrus v1.8.1
	github.com/tevino/abool v1.2.0 // indirect
	github.com/urfave/cli/v2 v2.3.0
	gitlab.com/ayufan/golang-cli-helpers v0.0.0-20171103152739-a7cf72d604cd // indirect
	gitlab.com/gitlab-org/gitlab-runner v12.5.0+incompatible
	golang.org/x/time v0.0.0-20200630173020-3af7569d3a1e // indirect
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	gopkg.in/yaml.v1 v1.0.0-20140924161607-9f9df34309c0
	gopkg.in/yaml.v2 v2.4.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
	gorm.io/datatypes v0.0.0-20200924071644-3967db6857cf
	gorm.io/driver/mysql v1.1.2
	gorm.io/driver/sqlite v1.1.5
	gorm.io/gorm v1.21.15
	k8s.io/api v0.19.2 // indirect
)

replace github.com/docker/docker v1.4.2-0.20190822180741-9552f2b2fdde => github.com/docker/engine v1.4.2-0.20190822180741-9552f2b2fdde

replace github.com/minio/go-homedir v0.0.0-20190425115525-017018655514 => gitlab.com/steveazz/go-homedir v0.0.0-20190425115525-017018655514

// Added this for chartmuseum/storage
replace github.com/NetEase-Object-Storage/nos-golang-sdk => github.com/karuppiah7890/nos-golang-sdk v0.0.0-20191116042345-0792ba35abcc

// Added this for chartmuseum/storage
replace go.etcd.io/etcd => github.com/eddycjy/etcd v0.5.0-alpha.5.0.20200218102753-4258cdd2efdf
