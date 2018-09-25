package main

import (
	"context"
	"log"

	"fmt"
	"go.etcd.io/etcd/clientv3"
	"strconv"
)

type Registry struct {
	etcdClient      *clientv3.Client
	traefikConf     TreafikConf
	generalLeaseTTL int64
}

func NewRegistry(etcdClient *clientv3.Client, traefikConf TreafikConf) *Registry {
	return &Registry{
		etcdClient:      etcdClient,
		traefikConf:     traefikConf,
		generalLeaseTTL: 10,
	}
}

func (r *Registry) isRegister(ctx context.Context, key string) bool {
	resp, err := r.etcdClient.Get(ctx, key)
	if err != nil {
		log.Fatalln("failed to fetch value", err)
	}

	return len(resp.Kvs) != 0
}

func (r *Registry) register(ctx context.Context, backend TraefikBackend) {
	resp, err := r.etcdClient.Grant(ctx, r.generalLeaseTTL)

	if err != nil {
		log.Fatalln("failed to grant general lease", err)
	}

	leaseID := resp.ID

	r.etcdClient.KeepAlive(context.Background(), leaseID)

	urlRecord := fmt.Sprintf("%s/backends/%s/servers/%s/url", r.traefikConf.EtcdPrefix, backend.Name, backend.Node)
	weightRecord := fmt.Sprintf("%s/backends/%s/servers/%s/weight", r.traefikConf.EtcdPrefix, backend.Name, backend.Node)

	if _, err := r.etcdClient.Put(ctx, urlRecord, backend.URL, clientv3.WithLease(leaseID)); err != nil {
		log.Fatalln("failed to register", err)
	}

	if _, err := r.etcdClient.Put(ctx, weightRecord, strconv.Itoa(int(backend.Weight)), clientv3.WithLease(leaseID)); err != nil {
		log.Fatalln("failed to register", err)
	}
}
