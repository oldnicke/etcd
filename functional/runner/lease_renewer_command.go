package runner

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"oldnicke/etcd/clientv3"

	"github.com/spf13/cobra"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	leaseTTL int64
)

// NewLeaseRenewerCommand returns the cobra command for "lease-renewer runner".
func NewLeaseRenewerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lease-renewer",
		Short: "Performs lease renew operation",
		Run:   runLeaseRenewerFunc,
	}
	cmd.Flags().Int64Var(&leaseTTL, "ttl", 5, "lease's ttl")
	return cmd
}

func runLeaseRenewerFunc(cmd *cobra.Command, args []string) {
	if len(args) > 0 {
		ExitWithError(ExitBadArgs, errors.New("lease-renewer does not take any argument"))
	}

	eps := endpointsFromFlag(cmd)
	c := newClient(eps, dialTimeout)
	ctx := context.Background()

	for {
		var (
			l   *clientv3.LeaseGrantResponse
			lk  *clientv3.LeaseKeepAliveResponse
			err error
		)
		for {
			l, err = c.Lease.Grant(ctx, leaseTTL)
			if err == nil {
				break
			}
		}
		expire := time.Now().Add(time.Duration(l.TTL-1) * time.Second)

		for {
			lk, err = c.Lease.KeepAliveOnce(ctx, l.ID)
			if ev, ok := status.FromError(err); ok && ev.Code() == codes.NotFound {
				if time.Since(expire) < 0 {
					log.Fatalf("bad renew! exceeded: %v", time.Since(expire))
					for {
						lk, err = c.Lease.KeepAliveOnce(ctx, l.ID)
						fmt.Println(lk, err)
						time.Sleep(time.Second)
					}
				}
				log.Fatalf("lost lease %d, expire: %v\n", l.ID, expire)
				break
			}
			if err != nil {
				continue
			}
			expire = time.Now().Add(time.Duration(lk.TTL-1) * time.Second)
			log.Printf("renewed lease %d, expire: %v\n", lk.ID, expire)
			time.Sleep(time.Duration(lk.TTL-2) * time.Second)
		}
	}
}
