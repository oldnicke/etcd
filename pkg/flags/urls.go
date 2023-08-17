package flags

import (
	"flag"
	"net/url"
	"strings"

	"go.etcd.io/etcd/pkg/types"
)

// URLsValue wraps "types.URLs".
type URLsValue types.URLs

// Set parses a command line set of URLs formatted like:
// http://127.0.0.1:2380,http://10.1.1.2:80
// Implements "flag.Value" interface.
func (us *URLsValue) Set(s string) error {
	ss, err := types.NewURLs(strings.Split(s, ","))
	if err != nil {
		return err
	}
	*us = URLsValue(ss)
	return nil
}

// String implements "flag.Value" interface.
func (us *URLsValue) String() string {
	all := make([]string, len(*us))
	for i, u := range *us {
		all[i] = u.String()
	}
	return strings.Join(all, ",")
}

// NewURLsValue implements "url.URL" slice as flag.Value interface.
// Given value is to be separated by comma.
func NewURLsValue(s string) *URLsValue {
	if s == "" {
		return &URLsValue{}
	}
	v := &URLsValue{}
	if err := v.Set(s); err != nil {
		plog.Panicf("new URLsValue should never fail: %v", err)
	}
	return v
}

// URLsFromFlag returns a slices from url got from the flag.
func URLsFromFlag(fs *flag.FlagSet, urlsFlagName string) []url.URL {
	return []url.URL(*fs.Lookup(urlsFlagName).Value.(*URLsValue))
}
