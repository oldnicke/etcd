package netutil

import (
	"fmt"
	"os/exec"
)

// DropPort drops all tcp packets that are received from the given port and sent to the given port.
func DropPort(port int) error {
	cmdStr := fmt.Sprintf("sudo iptables -A OUTPUT -p tcp --destination-port %d -j DROP", port)
	if _, err := exec.Command("/bin/sh", "-c", cmdStr).Output(); err != nil {
		return err
	}
	cmdStr = fmt.Sprintf("sudo iptables -A INPUT -p tcp --destination-port %d -j DROP", port)
	_, err := exec.Command("/bin/sh", "-c", cmdStr).Output()
	return err
}

// RecoverPort stops dropping tcp packets at given port.
func RecoverPort(port int) error {
	cmdStr := fmt.Sprintf("sudo iptables -D OUTPUT -p tcp --destination-port %d -j DROP", port)
	if _, err := exec.Command("/bin/sh", "-c", cmdStr).Output(); err != nil {
		return err
	}
	cmdStr = fmt.Sprintf("sudo iptables -D INPUT -p tcp --destination-port %d -j DROP", port)
	_, err := exec.Command("/bin/sh", "-c", cmdStr).Output()
	return err
}

// SetLatency adds latency in millisecond scale with random variations.
func SetLatency(ms, rv int) error {
	ifces, err := GetDefaultInterfaces()
	if err != nil {
		return err
	}

	if rv > ms {
		rv = 1
	}
	for ifce := range ifces {
		cmdStr := fmt.Sprintf("sudo tc qdisc add dev %s root netem delay %dms %dms distribution normal", ifce, ms, rv)
		_, err = exec.Command("/bin/sh", "-c", cmdStr).Output()
		if err != nil {
			// the rule has already been added. Overwrite it.
			cmdStr = fmt.Sprintf("sudo tc qdisc change dev %s root netem delay %dms %dms distribution normal", ifce, ms, rv)
			_, err = exec.Command("/bin/sh", "-c", cmdStr).Output()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// RemoveLatency resets latency configurations.
func RemoveLatency() error {
	ifces, err := GetDefaultInterfaces()
	if err != nil {
		return err
	}
	for ifce := range ifces {
		_, err = exec.Command("/bin/sh", "-c", fmt.Sprintf("sudo tc qdisc del dev %s root netem", ifce)).Output()
		if err != nil {
			return err
		}
	}
	return nil
}
