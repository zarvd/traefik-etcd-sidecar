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

	var etcdConfig clientv3.Config
	var traefikBackend TraefikBackend
	var traefikConf TreafikConf

	cmd := cobra.Command{
		Use: "start",
		Run: func(*cobra.Command, []string) {
			etcdClient, err := clientv3.New(etcdConfig)

			if err != nil {
				log.Fatalln("failed to new etcd client", err)
			}

			registry := NewRegistry(etcdClient, traefikConf)
			registry.register(context.Background(), traefikBackend)

			select {} // sleep forever
		},
	}

	cmd.PersistentFlags().StringSliceVar(&etcdConfig.Endpoints, "etcd-endpoints", []string{}, "etcd endpoints")
	cmd.PersistentFlags().StringVar(&etcdConfig.Username, "etcd-username", "", "etcd username")
	cmd.PersistentFlags().StringVar(&etcdConfig.Password, "etcd-password", "", "etcd password")

	cmd.PersistentFlags().StringVar(&traefikConf.EtcdPrefix, "traefik-etcd-prefix", "/traefik", "traefik etcd prefix")

	cmd.PersistentFlags().StringVar(&traefikBackend.Name, "traefik-backend-name", "", "traefik backend name")
	cmd.PersistentFlags().StringVar(&traefikBackend.Node, "traefik-backend-node", "", "traefik backend node")
	cmd.PersistentFlags().StringVar(&traefikBackend.URL, "traefik-backend-url", "", "traefik backend url")
	cmd.PersistentFlags().UintVar(&traefikBackend.Weight, "traefik-backend-weight", 1, "traefik backend weight")

	return &cmd
}

func init() {
	rootCmd.AddCommand(startCmd())
}
