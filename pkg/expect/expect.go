// Package expect implements a small expect-style interface
package expect

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"

	"github.com/creack/pty"
)

type ExpectProcess struct {
	cmd  *exec.Cmd
	fpty *os.File
	wg   sync.WaitGroup

	cond  *sync.Cond // for broadcasting updates are available
	mu    sync.Mutex // protects lines and err
	lines []string
	count int // increment whenever new line gets added
	err   error

	// StopSignal is the signal Stop sends to the process; defaults to SIGKILL.
	StopSignal os.Signal
}

// NewExpect creates a new process for expect testing.
func NewExpect(name string, arg ...string) (ep *ExpectProcess, err error) {
	// if env[] is nil, use current system env
	return NewExpectWithEnv(name, arg, nil)
}

// NewExpectWithEnv creates a new process with user defined env variables for expect testing.
func NewExpectWithEnv(name string, args []string, env []string) (ep *ExpectProcess, err error) {
	cmd := exec.Command(name, args...)
	cmd.Env = env
	ep = &ExpectProcess{
		cmd:        cmd,
		StopSignal: syscall.SIGKILL,
	}
	ep.cond = sync.NewCond(&ep.mu)
	ep.cmd.Stderr = ep.cmd.Stdout
	ep.cmd.Stdin = nil

	if ep.fpty, err = pty.Start(ep.cmd); err != nil {
		return nil, err
	}

	ep.wg.Add(1)
	go ep.read()
	return ep, nil
}

func (ep *ExpectProcess) read() {
	defer ep.wg.Done()
	printDebugLines := os.Getenv("EXPECT_DEBUG") != ""
	r := bufio.NewReader(ep.fpty)
	for ep.err == nil {
		l, rerr := r.ReadString('\n')
		ep.mu.Lock()
		ep.err = rerr
		if l != "" {
			if printDebugLines {
				fmt.Printf("%s-%d: %s", ep.cmd.Path, ep.cmd.Process.Pid, l)
			}
			ep.lines = append(ep.lines, l)
			ep.count++
			if len(ep.lines) == 1 {
				ep.cond.Signal()
			}
		}
		ep.mu.Unlock()
	}
	ep.cond.Signal()
}

// ExpectFunc returns the first line satisfying the function f.
func (ep *ExpectProcess) ExpectFunc(f func(string) bool) (string, error) {
	ep.mu.Lock()
	for {
		for len(ep.lines) == 0 && ep.err == nil {
			ep.cond.Wait()
		}
		if len(ep.lines) == 0 {
			break
		}
		l := ep.lines[0]
		ep.lines = ep.lines[1:]
		if f(l) {
			ep.mu.Unlock()
			return l, nil
		}
	}
	ep.mu.Unlock()
	return "", ep.err
}

// Expect returns the first line containing the given string.
func (ep *ExpectProcess) Expect(s string) (string, error) {
	return ep.ExpectFunc(func(txt string) bool { return strings.Contains(txt, s) })
}

// LineCount returns the number of recorded lines since
// the beginning of the process.
func (ep *ExpectProcess) LineCount() int {
	ep.mu.Lock()
	defer ep.mu.Unlock()
	return ep.count
}

// Stop kills the expect process and waits for it to exit.
func (ep *ExpectProcess) Stop() error { return ep.close(true) }

// Signal sends a signal to the expect process
func (ep *ExpectProcess) Signal(sig os.Signal) error {
	return ep.cmd.Process.Signal(sig)
}

// Close waits for the expect process to exit.
func (ep *ExpectProcess) Close() error { return ep.close(false) }

func (ep *ExpectProcess) close(kill bool) error {
	if ep.cmd == nil {
		return ep.err
	}
	if kill {
		ep.Signal(ep.StopSignal)
	}

	err := ep.cmd.Wait()
	ep.fpty.Close()
	ep.wg.Wait()

	if err != nil {
		ep.err = err
		if !kill && strings.Contains(err.Error(), "exit status") {
			// non-zero exit code
			err = nil
		} else if kill && strings.Contains(err.Error(), "signal:") {
			err = nil
		}
	}
	ep.cmd = nil
	return err
}

func (ep *ExpectProcess) Send(command string) error {
	_, err := io.WriteString(ep.fpty, command)
	return err
}
