FROM alpine:3.8

ARG WORKSPACE=/workspace

RUN mkdir -p $WORKSPACE
COPY traefik-etcd-sidecar $WORKSPACE

WORKDIR $WORKSPACE

ENTRYPOINT ["./traefik-etcd-sidecar"]
CMD ["start"]
