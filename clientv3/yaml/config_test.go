package yaml

import (
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"testing"

	"sigs.k8s.io/yaml"
)

var (
	certPath       = "../../integration/fixtures/server.crt"
	privateKeyPath = "../../integration/fixtures/server.key.insecure"
	caPath         = "../../integration/fixtures/ca.crt"
)

func TestConfigFromFile(t *testing.T) {
	tests := []struct {
		ym *yamlConfig

		werr bool
	}{
		{
			&yamlConfig{},
			false,
		},
		{
			&yamlConfig{
				InsecureTransport: true,
			},
			false,
		},
		{
			&yamlConfig{
				Keyfile:               privateKeyPath,
				Certfile:              certPath,
				TrustedCAfile:         caPath,
				InsecureSkipTLSVerify: true,
			},
			false,
		},
		{
			&yamlConfig{
				Keyfile:  "bad",
				Certfile: "bad",
			},
			true,
		},
		{
			&yamlConfig{
				Keyfile:       privateKeyPath,
				Certfile:      certPath,
				TrustedCAfile: "bad",
			},
			true,
		},
	}

	for i, tt := range tests {
		tmpfile, err := ioutil.TempFile("", "clientcfg")
		if err != nil {
			log.Fatal(err)
		}

		b, err := yaml.Marshal(tt.ym)
		if err != nil {
			t.Fatal(err)
		}

		_, err = tmpfile.Write(b)
		if err != nil {
			t.Fatal(err)
		}
		err = tmpfile.Close()
		if err != nil {
			t.Fatal(err)
		}

		cfg, cerr := NewConfig(tmpfile.Name())
		if cerr != nil && !tt.werr {
			t.Errorf("#%d: err = %v, want %v", i, cerr, tt.werr)
			continue
		}
		if cerr != nil {
			continue
		}

		if !reflect.DeepEqual(cfg.Endpoints, tt.ym.Endpoints) {
			t.Errorf("#%d: endpoint = %v, want %v", i, cfg.Endpoints, tt.ym.Endpoints)
		}

		if tt.ym.InsecureTransport != (cfg.TLS == nil) {
			t.Errorf("#%d: insecureTransport = %v, want %v", i, cfg.TLS == nil, tt.ym.InsecureTransport)
		}

		if !tt.ym.InsecureTransport {
			if tt.ym.Certfile != "" && len(cfg.TLS.Certificates) == 0 {
				t.Errorf("#%d: failed to load in cert", i)
			}
			if tt.ym.TrustedCAfile != "" && cfg.TLS.RootCAs == nil {
				t.Errorf("#%d: failed to load in ca cert", i)
			}
			if cfg.TLS.InsecureSkipVerify != tt.ym.InsecureSkipTLSVerify {
				t.Errorf("#%d: skipTLSVeify = %v, want %v", i, cfg.TLS.InsecureSkipVerify, tt.ym.InsecureSkipTLSVerify)
			}
		}

		os.Remove(tmpfile.Name())
	}
}
