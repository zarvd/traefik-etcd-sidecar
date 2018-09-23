package main

import (
	"os"
	"strings"

	"go.etcd.io/etcd/clientv3"
)

func fetchTraefikKVs() map[string]string {
	traefikKVs := make(map[string]string)

	for _, pair := range os.Environ() {

		if strings.HasPrefix(pair, "traefik_") {
			kv := strings.Split(pair, "=")
			k := strings.Replace(kv[0], "_", "/", -1)

			traefikKVs[k] = kv[1]
		}
	}

	return traefikKVs
}

func fetchEtcdConf() (config clientv3.Config) {
	const EtcdEndpointsKey = "etcd_endpoints"
	const EtcdUsernameKey = "etcd_username"
	const EtcdPasswordKey = "etcd_password"

	for _, pair := range os.Environ() {
		kv := strings.Split(pair, "=")

		switch kv[0] {
		case EtcdEndpointsKey:
			config.Endpoints = strings.Split(kv[1], ",")
		case EtcdUsernameKey:
			config.Username = kv[1]
		case EtcdPasswordKey:
			config.Password = kv[1]
		}
	}
	return
}
