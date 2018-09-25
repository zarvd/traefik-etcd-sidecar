Sidecar for registering treafik backend service by etcd
----

## Example

docker-compose:

```yaml
---
version: "3"

services:
  api:
    image: emilevauge/whoami
    restart: always

  api-traefik-sidecar:
    image: 0x636363/traefik-etcd-sidecar
    container_name: api_traefik_sidecar
    command:
      - "start"
      - "--etcd-endpoints=127.0.0.1:2379"
      - "--etcd-username=traefik"
      - "--etcd-password=traefik"
      - "--traefik-backend-name=api"
      - "--traefik-backend-node=node1"
      - "--traefik-backend-url=http://api:80"
      - "--traefik-backend-weight=1"
      - "--traefik-etcd-prefix=/traefik"
    depends_on:
      - api

```

## License

MIT
