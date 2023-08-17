// Package osutil implements operating system-related utility functions.
package osutil

import (
	"os"
	"strings"

	"github.com/coreos/pkg/capnslog"
)

var (
	plog = capnslog.NewPackageLogger("go.etcd.io/etcd", "pkg/osutil")

	// support to override setting SIG_DFL so tests don't terminate early
	setDflSignal = dflSignal
)

func Unsetenv(key string) error {
	envs := os.Environ()
	os.Clearenv()
	for _, e := range envs {
		strs := strings.SplitN(e, "=", 2)
		if strs[0] == key {
			continue
		}
		if err := os.Setenv(strs[0], strs[1]); err != nil {
			return err
		}
	}
	return nil
}
