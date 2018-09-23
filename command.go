package main

import (
	"context"
	"log"

	"github.com/spf13/cobra"
	"go.etcd.io/etcd/clientv3"
)

var rootCmd = &cobra.Command{
	Use: "traefik-etcd-sidecar",
}

func startCmd() *cobra.Command {

	cmd := cobra.Command{
		Use: "start",
		Run: func(*cobra.Command, []string) {
			traefikKVs := fetchTraefikKVs()

			etcdConfig := fetchEtcdConf()

			etcdClient, err := clientv3.New(etcdConfig)

			if err != nil {
				log.Fatalln("failed to new etcd client", err)
			}

			registry := NewRegistry(etcdClient)
			registry.register(context.Background(), traefikKVs)
		},
	}

	return &cmd
}

func init() {
	rootCmd.AddCommand(startCmd())
}
