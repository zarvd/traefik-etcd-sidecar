package traefik

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"go.etcd.io/etcd/clientv3"

	"github.com/0x636363/traefik-etcd-sidecar/readiness"
)

type Client struct {
	etcdClient *clientv3.Client
	etcdPrefix string
}

type ClientOption func(*Client)

func EtcdClient(client *clientv3.Client) ClientOption {
	return func(c *Client) {
		c.etcdClient = client
	}
}

func EtcdPrefix(prefix string) ClientOption {
	return func(c *Client) {
		c.etcdPrefix = prefix
	}
}

var defaultClient = Client{
	etcdPrefix: "/traefik",
}

func NewClient(opts ...ClientOption) *Client {

	c := defaultClient

	for _, opt := range opts {
		opt(&c)
	}

	if c.etcdClient == nil {
		panic("etcd client is not init")
	}

	return &c
}

type registerBackendOptions struct {
	readiness    readiness.Readiness
	keepAliveTTL int64
}

type RegisterBackendOption func(*registerBackendOptions)

func Ready(readiness readiness.Readiness) RegisterBackendOption {
	return func(opts *registerBackendOptions) {
		opts.readiness = readiness
	}
}

func KeepAlive(ttl int64) RegisterBackendOption {
	return func(opts *registerBackendOptions) {
		opts.keepAliveTTL = ttl
	}
}

var defaultRegisterBackendOptions = registerBackendOptions{
	readiness:    nil,
	keepAliveTTL: 5,
}

func (c *Client) registerBackendWithLease(ctx context.Context, backend Backend, leaseTTL int64) clientv3.LeaseID {
	var leaseID clientv3.LeaseID

	if resp, err := c.etcdClient.Grant(ctx, leaseTTL); err != nil {
		log.Fatalln("failed to grant etcd lease", err)
	} else {
		leaseID = resp.ID
	}

	putOptions := []clientv3.OpOption{
		clientv3.WithLease(leaseID),
	}

	if err := c.SetBackend(ctx, backend, putOptions...); err != nil {
		log.Fatalln("failed to set traefik backend", err)
	}

	return leaseID
}

// register backend and keep alive
// stop keeping alive by cancel context
func (c *Client) RegisterBackend(ctx context.Context, backend Backend, opts ...RegisterBackendOption) {
	opt := defaultRegisterBackendOptions

	for _, o := range opts {
		o(&opt)
	}

	// FIXME re-register if keep alive fail
	if opt.readiness == nil {
		// without readiness check
		log.Println("register backend without readiness check", backend)
		leaseID := c.registerBackendWithLease(ctx, backend, opt.keepAliveTTL)
		c.etcdClient.KeepAlive(context.Background(), leaseID)

		for {
			select {
			case <-ctx.Done():
				log.Println("unregister backend by context cancellation", backend)
				c.etcdClient.Revoke(context.Background(), leaseID)
				return
			}
		}
	} else {
		var leaseID clientv3.LeaseID

		ready := readiness.Ready(opt.readiness)

		for {
			select {
			case isReady := <-ready:
				if isReady {
					log.Println("register backend by readiness", backend)
					// register and keep alive
					leaseID = c.registerBackendWithLease(ctx, backend, opt.keepAliveTTL)
					c.etcdClient.KeepAlive(ctx, leaseID)
				} else {
					log.Println("unregister backend by readiness", backend)
					// revoke
					c.etcdClient.Revoke(context.Background(), leaseID)
				}
			case <-ctx.Done():
				log.Println("unregister backend by context cancellation", backend)
				c.etcdClient.Revoke(context.Background(), leaseID)
				return
			}
		}
	}
}

func (c *Client) SetBackend(ctx context.Context, backend Backend, etcdPutOptions ...clientv3.OpOption) error {

	// write url
	urlRecord := fmt.Sprintf("%s/backends/%s/servers/%s/url", c.etcdPrefix, backend.Name, backend.Node)
	if _, err := c.etcdClient.Put(ctx, urlRecord, backend.URL, etcdPutOptions...); err != nil {
		return err
	}

	// write weigth
	weightRecord := fmt.Sprintf("%s/backends/%s/servers/%s/weight", c.etcdPrefix, backend.Name, backend.Node)
	if _, err := c.etcdClient.Put(ctx, weightRecord, strconv.Itoa(int(backend.Weight)), etcdPutOptions...); err != nil {
		// FIXME: rollback url
		return err
	}

	return nil
}
