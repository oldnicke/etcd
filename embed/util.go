package embed

import (
	"path/filepath"

	"github.com/oldnicke/etcd/wal"
)

func isMemberInitialized(cfg *Config) bool {
	waldir := cfg.WalDir
	if waldir == "" {
		waldir = filepath.Join(cfg.Dir, "member", "wal")
	}
	return wal.Exist(waldir)
}
