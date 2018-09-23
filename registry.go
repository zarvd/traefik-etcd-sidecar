package main

import (
	"context"
	"log"

	"go.etcd.io/etcd/clientv3"
)

type Registry struct {
	etcdClient      *clientv3.Client
	generalLeaseTTL int64
}

func NewRegistry(etcdClient *clientv3.Client) *Registry {
	return &Registry{
		etcdClient:      etcdClient,
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

func (r *Registry) register(ctx context.Context, traefikKVs map[string]string) {
	resp, err := r.etcdClient.Grant(ctx, r.generalLeaseTTL)

	if err != nil {
		log.Fatalln("failed to grant general lease", err)
	}

	leaseID := resp.ID

	r.etcdClient.KeepAlive(context.Background(), leaseID)

	for k, v := range traefikKVs {
		if r.isRegister(ctx, k) {
			log.Printf("WARNING key %s is registered, it will be override\n", k)
		}

		_, err := r.etcdClient.Put(ctx, k, v, clientv3.WithLease(leaseID))
		if err != nil {
			log.Fatalln("failed to register", err)
		}
	}
}
