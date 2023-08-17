package command

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/oldnicke/etcd/clientv3/snapshot"
)

type jsonPrinter struct{ printer }

func newJSONPrinter() printer {
	return &jsonPrinter{
		&printerRPC{newPrinterUnsupported("json"), printJSON},
	}
}

func (p *jsonPrinter) EndpointHealth(r []epHealth) { printJSON(r) }
func (p *jsonPrinter) EndpointStatus(r []epStatus) { printJSON(r) }
func (p *jsonPrinter) EndpointHashKV(r []epHashKV) { printJSON(r) }
func (p *jsonPrinter) DBStatus(r snapshot.Status)  { printJSON(r) }

func printJSON(v interface{}) {
	b, err := json.Marshal(v)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	fmt.Println(string(b))
}
