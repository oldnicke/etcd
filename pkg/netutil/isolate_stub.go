//go:build !linux
// +build !linux

package netutil

func DropPort(port int) error { return nil }

func RecoverPort(port int) error { return nil }

func SetLatency(ms, rv int) error { return nil }

func RemoveLatency() error { return nil }
