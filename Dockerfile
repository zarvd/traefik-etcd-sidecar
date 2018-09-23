FROM golang:1.11-alpine

WORKDIR /go/src/traefik-etcd-sidecar
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["traefik-etcd-sidecar"]