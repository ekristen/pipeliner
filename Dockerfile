FROM node:17 as ui-builder
WORKDIR /src
COPY ui/package.json /src/package.json
COPY ui/package-lock.json /src/package-lock.json
RUN npm install
COPY ui/ /src
RUN npm run build

FROM golang:1.16.4 as ui-generator
WORKDIR /src
COPY . /src
COPY --from=ui-builder /src/dist /src/ui/dist
RUN make generate

FROM golang:1.16.4 as builder
WORKDIR /src
COPY . /src
COPY --from=ui-generator /src/pkg/box/blob.go /src/pkg/box/blob.go
ENV GO111MODULE=on
ARG branch=master
ENV BRANCH=${branch}
RUN make release
RUN ls /src/release

FROM debian
ENTRYPOINT ["/usr/local/bin/pipeliner"]
COPY --from=builder /src/release/pipeliner_linux_amd64 /usr/local/bin/pipeliner
