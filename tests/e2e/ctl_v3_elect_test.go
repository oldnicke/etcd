package e2e

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/oldnicke/etcd/pkg/expect"
)

func TestCtlV3Elect(t *testing.T) {
	oldenv := os.Getenv("EXPECT_DEBUG")
	defer os.Setenv("EXPECT_DEBUG", oldenv)
	os.Setenv("EXPECT_DEBUG", "1")

	testCtl(t, testElect)
}

func testElect(cx ctlCtx) {
	name := "a"

	holder, ch, err := ctlV3Elect(cx, name, "p1")
	if err != nil {
		cx.t.Fatal(err)
	}

	l1 := ""
	select {
	case <-time.After(2 * time.Second):
		cx.t.Fatalf("timed out electing")
	case l1 = <-ch:
		if !strings.HasPrefix(l1, name) {
			cx.t.Errorf("got %q, expected %q prefix", l1, name)
		}
	}

	// blocked process that won't win the election
	blocked, ch, err := ctlV3Elect(cx, name, "p2")
	if err != nil {
		cx.t.Fatal(err)
	}
	select {
	case <-time.After(100 * time.Millisecond):
	case <-ch:
		cx.t.Fatalf("should block")
	}

	// overlap with a blocker that will win the election
	blockAcquire, ch, err := ctlV3Elect(cx, name, "p2")
	if err != nil {
		cx.t.Fatal(err)
	}
	defer blockAcquire.Stop()
	select {
	case <-time.After(100 * time.Millisecond):
	case <-ch:
		cx.t.Fatalf("should block")
	}

	// kill blocked process with clean shutdown
	if err = blocked.Signal(os.Interrupt); err != nil {
		cx.t.Fatal(err)
	}
	if err = closeWithTimeout(blocked, time.Second); err != nil {
		cx.t.Fatal(err)
	}

	// kill the holder with clean shutdown
	if err = holder.Signal(os.Interrupt); err != nil {
		cx.t.Fatal(err)
	}
	if err = closeWithTimeout(holder, time.Second); err != nil {
		cx.t.Fatal(err)
	}

	// blockAcquire should win the election
	select {
	case <-time.After(time.Second):
		cx.t.Fatalf("timed out from waiting to holding")
	case l2 := <-ch:
		if l1 == l2 || !strings.HasPrefix(l2, name) {
			cx.t.Fatalf("expected different elect name, got l1=%q, l2=%q", l1, l2)
		}
	}
}

// ctlV3Elect creates a elect process with a channel listening for when it wins the election.
func ctlV3Elect(cx ctlCtx, name, proposal string) (*expect.ExpectProcess, <-chan string, error) {
	cmdArgs := append(cx.PrefixArgs(), "elect", name, proposal)
	proc, err := spawnCmd(cmdArgs)
	outc := make(chan string, 1)
	if err != nil {
		close(outc)
		return proc, outc, err
	}
	go func() {
		s, xerr := proc.ExpectFunc(func(string) bool { return true })
		if xerr != nil {
			cx.t.Errorf("expect failed (%v)", xerr)
		}
		outc <- s
	}()
	return proc, outc, err
}
