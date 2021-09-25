# Pipeliner

**Status: ALPHA** - Everything works, but decisions could be made that cause breaking changes still. Definitely still rough edges.

The initial goals of this project is to provide a way to run arbitrary pipelines of jobs in a fast and sane manor. This project uses golang and Vue.JS to create a single compiled binary with the UI embedded. By default is leverages SQLite and it can be configured to connect to any SQL server that [GORM](https://gorm.io/index.html) supports.

## Getting Started

The simplest way to get started is to download the binary for your OS and run the following:

```bash
./pipeliner api-server
```

Then open the UI at http://localhost:4444, navigate to the runners tab and follow the in UI instructions for adding a runner to the API.

**Note:** there currently isn't any authentication to for registering runners, but once the runner is registered it's token it uses is unique. Registration auth is in the works.

There are bunch of workflows in `examples/workerflows`, you can manually load them or use the `seed-workflows` by the following command:

```bash
./pipeliner seed-workflows --directory examples/workflows
```

## Acknowledgements, Credits and Thanks

Thank you to GitLab for building a wonderful CI/CD system and for open sourcing their GitLab Runner.

Thank you to Alloy CI from which inspiration and some better understanding of the GitLab backend works.

## Design

This project utilizes the [GitLab Runner](https://docs.gitlab.com/runner/) and it's extensive capabilities. The GitLab Runner is a superb piece of software that can run in a number of different ways, by re-using this, the project is 1000x easier to maintain.

A workflow is a [GitLab CI YAML](https://docs.gitlab.com/ee/ci/yaml/) file, see the Support Matrix below for more information on what is and is not supported.

A pipeline is an instance of a workflow, a workflow is made up of one or more jobs.

Variables can be defined globally or at the pipeline level.

**Important:** SCM integration has been **PURPOSEFULLY** left out at this point. It's possible this project will be expanded to handle SCMs.

If you are looking for an open source CI/CD system, please check out [Alloy CI](https://github.com/AlloyCI/alloy_ci), Alloy CI is based on GitLab Runner as well, but is designed to integration with your SCM.

## Goals

- Easy to deploy
- Easy to maintian
- Easy to scale
- API Driven
- Build in UI
- Provide a simple and effective way to run workflows reliably at scale
- Remain compatible with GitLab Runner
- Maintain a GitLab Runner Fork that exposes additional privileges/scheduling that normally would not be wanted in a multi-tenant environment
- Webhooks after pipelines and jobs run

## Non-Goals

**Note:** Subject to change at any time.

- Integrated CI/CD with SCM (see [Alloy CI](https://github.com/AlloyCI/alloy_ci))
- Multi-tenancy

## Outstanding Tasks

- [ ] High Availability for the API Server
- [ ] UI Authentication
- [ ] Runner Registration Auth
- [ ] Support include (Gitlab CI, or YAML include of some sort)
- [ ] Support trigger (GitLab CI)
- [ ] Support rules (GitLab CI terminology for when jobs execute)
- [ ] Support needs
- [ ] Webhooks

## Features

- Most [GitLab CI Features](https://docs.gitlab.com/runner/#features)
- Dynamic Kubernetes Scheduling

### Dynamic Kubernetes Scheduling

When running a job there are many occassions where a job (or workload) would need to be run on a specific node in a Kubernetes cluster, maybe because of unique hardware.

Using the GitLab Runner natively this is not possible without running a different executor PER node, while the GitLab Runner makes this possibly by configuring multiple executors per runner. If nodes change frequently the reconfiguring of the GitLab Runner gets difficult, enter the **Dynamic Kubernetes Scheduling** feature.

This feature allows using a first class feature of GitLab Runner to pass environment variables between jobs and by leveraging a specific environment variable and it's corresponding value to dynamically alter the scheduling of a job to a specific node.

`PIPELINER_K8S_AFFINITY` allows for the dynamic altering of downstream dependent jobs. The value of this command takes the following format `node:required|preferred:expression|field:test.io/key:in|exists|gt|lt:[value]`

For example, to schedule to a specific node based on an annotation the value of the environment variable would be `node:required:expression:ekristen.github.io/pipeliner:exists`

## GitLab CI YAML Support Matrix

- Yes, it is currently supported.
- No, it is not and isn't planned likely due to not being relevant.
- Planned, it is not yet implemented, but will be added at some point.

| Keyword | Supported |
| ------- | --------- |
| script | YES |
| after_script | YES |
| allow_failure | YES |
| allow_failure:exit_codes | PLANNED |
| artifacts | YES |
| before_script | YES |
| cache | YES |
| coverage | NO |
| dependencies | YES |
| environment | NO |
| except | NO |
| extends | NO (yaml anchors are supported) |
| image | YES |
| include | PLANNED |
| include:local | NO |
| include:file | PLANNED |
| include:remote | PLANNED |
| include:template | NO |
| interruptible | NO |
| only | NO |
| needs| PLANNED |
| pages | NO |
| parallel | YES |
| release | NO |
| resource_group | PLANNED |
| retry | PLANNED |
| rules | NO |
| secrets | PLANNED |
| services | YES |
| stage | YES |
| tags | YES |
| timeout | YES |
| trigger | PLANNED |
| variables | YES |
| when | YES |
| when:delayed | PLANNED |

## Development

There are two components to this project, the API and the UI.

The UI is build with Node.JS 12.x, Vue.JS 2.x, and Vuetify.

The API is based entirely on golang and embeds the at compile time.

### Builds

Production builds are a single binary that can be compiled for any architecture.

### Developing the API

To do a production build run the following command

```bash
make build
```

To do active development run the following command

```bash
go run main.go api-server -l debug
```

### Developing the UI

To do a production build run the following command

```bash
npm run build
```

This will result in a `ui/dist` folder which is git ignored.

To do active development run the following command

```bash
npm run serve
```

This will start a webpack development server on port at [http://0.0.0.0:4445](http://0.0.0.0:4445)
