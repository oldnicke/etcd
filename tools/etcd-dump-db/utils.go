package main

import "os"

func existFileOrDir(name string) bool {
	_, err := os.Stat(name)
	return err == nil
}
