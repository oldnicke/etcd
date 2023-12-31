package flags

import (
	"flag"
	"os"
	"strings"
	"testing"
)

func TestSetFlagsFromEnv(t *testing.T) {
	fs := flag.NewFlagSet("testing", flag.ExitOnError)
	fs.String("a", "", "")
	fs.String("b", "", "")
	fs.String("c", "", "")
	fs.Parse([]string{})

	os.Clearenv()
	// flags should be settable using env vars
	os.Setenv("ETCD_A", "foo")
	// and command-line flags
	if err := fs.Set("b", "bar"); err != nil {
		t.Fatal(err)
	}

	// first verify that flags are as expected before reading the env
	for f, want := range map[string]string{
		"a": "",
		"b": "bar",
	} {
		if got := fs.Lookup(f).Value.String(); got != want {
			t.Fatalf("flag %q=%q, want %q", f, got, want)
		}
	}

	// now read the env and verify flags were updated as expected
	err := SetFlagsFromEnv("ETCD", fs)
	if err != nil {
		t.Errorf("err=%v, want nil", err)
	}
	for f, want := range map[string]string{
		"a": "foo",
		"b": "bar",
	} {
		if got := fs.Lookup(f).Value.String(); got != want {
			t.Errorf("flag %q=%q, want %q", f, got, want)
		}
	}
}

func TestSetFlagsFromEnvBad(t *testing.T) {
	// now verify that an error is propagated
	fs := flag.NewFlagSet("testing", flag.ExitOnError)
	fs.Int("x", 0, "")
	os.Setenv("ETCD_X", "not_a_number")
	if err := SetFlagsFromEnv("ETCD", fs); err == nil {
		t.Errorf("err=nil, want != nil")
	}
}

func TestSetFlagsFromEnvParsingError(t *testing.T) {
	fs := flag.NewFlagSet("etcd", flag.ContinueOnError)
	var tickMs uint
	fs.UintVar(&tickMs, "heartbeat-interval", 0, "Time (in milliseconds) of a heartbeat interval.")

	if oerr := os.Setenv("ETCD_HEARTBEAT_INTERVAL", "100 # ms"); oerr != nil {
		t.Fatal(oerr)
	}
	defer os.Unsetenv("ETCD_HEARTBEAT_INTERVAL")

	err := SetFlagsFromEnv("ETCD", fs)
	for _, v := range []string{"invalid syntax", "parse error"} {
		if strings.Contains(err.Error(), v) {
			err = nil
			break
		}
	}
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
}
