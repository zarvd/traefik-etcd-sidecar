FROM alpine:3.8

ARG WORKSPACE=/workspace

RUN mkdir -p $WORKSPACE
COPY traefik-etcd-sidecar $WORKSPACE

WORKDIR $WORKSPACE

CMD sleep 10000
# ENTRYPOINT ["./traefik-etcd-sidecar"]
# CMD ["start"]
