package runner

import (
	"context"
	"errors"
	"fmt"

	"github.com/oldnicke/etcd/clientv3/concurrency"

	"github.com/spf13/cobra"
)

// NewElectionCommand returns the cobra command for "election runner".
func NewElectionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "election [election name (defaults to 'elector')]",
		Short: "Performs election operation",
		Run:   runElectionFunc,
	}
	cmd.Flags().IntVar(&totalClientConnections, "total-client-connections", 10, "total number of client connections")
	return cmd
}

func runElectionFunc(cmd *cobra.Command, args []string) {
	election := "elector"
	if len(args) == 1 {
		election = args[0]
	}
	if len(args) > 1 {
		ExitWithError(ExitBadArgs, errors.New("election takes at most one argument"))
	}

	rcs := make([]roundClient, totalClientConnections)
	validatec := make(chan struct{}, len(rcs))
	// nextc closes when election is ready for next round.
	nextc := make(chan struct{})
	eps := endpointsFromFlag(cmd)

	for i := range rcs {
		v := fmt.Sprintf("%d", i)
		observedLeader := ""
		validateWaiters := 0
		var rcNextc chan struct{}
		setRcNextc := func() {
			rcNextc = nextc
		}

		rcs[i].c = newClient(eps, dialTimeout)
		var (
			s   *concurrency.Session
			err error
		)
		for {
			s, err = concurrency.NewSession(rcs[i].c)
			if err == nil {
				break
			}
		}

		e := concurrency.NewElection(s, election)
		rcs[i].acquire = func() (err error) {
			ctx, cancel := context.WithCancel(context.Background())
			donec := make(chan struct{})
			go func() {
				defer close(donec)
				for ctx.Err() == nil {
					if ol, ok := <-e.Observe(ctx); ok {
						observedLeader = string(ol.Kvs[0].Value)
						break
					}
				}
				if observedLeader != v {
					cancel()
				}
			}()
			err = e.Campaign(ctx, v)
			cancel()
			<-donec
			if err == nil {
				observedLeader = v
			}
			if observedLeader == v {
				validateWaiters = len(rcs)
			}
			select {
			case <-ctx.Done():
				return nil
			default:
				return err
			}
		}
		rcs[i].validate = func() error {
			l, err := e.Leader(context.TODO())
			if err == nil && string(l.Kvs[0].Value) != observedLeader {
				return fmt.Errorf("expected leader %q, got %q", observedLeader, l.Kvs[0].Value)
			}
			if err != nil {
				return err
			}
			setRcNextc()
			validatec <- struct{}{}
			return nil
		}
		rcs[i].release = func() error {
			for validateWaiters > 0 {
				select {
				case <-validatec:
					validateWaiters--
				default:
					return fmt.Errorf("waiting on followers")
				}
			}
			if err := e.Resign(context.TODO()); err != nil {
				return err
			}
			if observedLeader == v {
				oldNextc := nextc
				nextc = make(chan struct{})
				close(oldNextc)

			}
			<-rcNextc
			observedLeader = ""
			return nil
		}
	}
	// each client creates 1 key from Campaign() and delete it from Resign()
	// a round involves in 2*len(rcs) requests.
	doRounds(rcs, rounds, 2*len(rcs))
}
