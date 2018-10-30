package main

import (
	"context"
	"log"
	"os"
	"syscall"
	"os/signal"

	"github.com/spf13/cobra"
	"go.etcd.io/etcd/clientv3"

	"github.com/0x636363/traefik-etcd-sidecar/readiness"
	"github.com/0x636363/traefik-etcd-sidecar/traefik"
)

var rootCmd = &cobra.Command{
	Use: "traefik-etcd-sidecar",
}

type httpReadinessOptions struct {
	Host     string
	Port     uint
	Path     string
	Interval uint
}

type traefikOptions struct {
	EtcdPrefix string
	Backend    traefik.Backend
}

func startCmd() *cobra.Command {

	var etcdConfig clientv3.Config
	var traefikOptions traefikOptions
	var httpReadinessOptions httpReadinessOptions

	cmd := cobra.Command{
		Use: "start",
		Run: func(*cobra.Command, []string) {
			etcdClient, err := clientv3.New(etcdConfig)

			if err != nil {
				log.Fatalln("failed to new etcd client", err)
			}

			ready := readiness.NewHTTPReadiness(
				readiness.HTTPHost(httpReadinessOptions.Host, httpReadinessOptions.Port),
				readiness.HTTPPath(httpReadinessOptions.Path),
				readiness.HTTPInterval(httpReadinessOptions.Interval),
			)

			traefikClient := traefik.NewClient(
				traefik.EtcdClient(etcdClient),
				traefik.EtcdPrefix(traefikOptions.EtcdPrefix),
			)

			ctx, cancel := context.WithCancel(context.Background())

			c := make(chan os.Signal)
			signal.Notify(c, syscall.SIGTERM, syscall.SIGKILL)

			go traefikClient.RegisterBackend(ctx, traefikOptions.Backend,
				traefik.Ready(ready),
				traefik.KeepAlive(5),
			)

			select {
			case <-c:
				cancel()
				log.Printf("Gracefully shutting down...")
			}
		},
	}

	cmd.PersistentFlags().StringSliceVar(&etcdConfig.Endpoints, "etcd-endpoints", []string{}, "etcd endpoints")
	cmd.PersistentFlags().StringVar(&etcdConfig.Username, "etcd-username", "", "etcd username")
	cmd.PersistentFlags().StringVar(&etcdConfig.Password, "etcd-password", "", "etcd password")

	cmd.PersistentFlags().StringVar(&traefikOptions.EtcdPrefix, "traefik-etcd-prefix", "/traefik", "traefik etcd prefix")

	cmd.PersistentFlags().StringVar(&traefikOptions.Backend.Name, "traefik-backend-name", "", "traefik backend name")
	cmd.PersistentFlags().StringVar(&traefikOptions.Backend.Node, "traefik-backend-node", "", "traefik backend node")
	cmd.PersistentFlags().StringVar(&traefikOptions.Backend.URL, "traefik-backend-url", "", "traefik backend url")
	cmd.PersistentFlags().UintVar(&traefikOptions.Backend.Weight, "traefik-backend-weight", 1, "traefik backend weight")

	cmd.PersistentFlags().StringVar(&httpReadinessOptions.Host, "service-http-readiness-host", "localhost", "backend service HTTP readiness host")
	cmd.PersistentFlags().UintVar(&httpReadinessOptions.Port, "service-http-readiness-port", 80, "backend service HTTP readiness port")
	cmd.PersistentFlags().StringVar(&httpReadinessOptions.Path, "service-http-readiness-path", "/", "backend service HTTP readiness path")
	cmd.PersistentFlags().UintVar(&httpReadinessOptions.Interval, "service-http-readiness-interval", 5, "backend service HTTP readiness interval [Seconds]")

	return &cmd
}

func init() {
	rootCmd.AddCommand(startCmd())
}
