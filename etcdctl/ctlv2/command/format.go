package command

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/oldnicke/etcd/client"
)

// printResponseKey only supports to print key correctly.
func printResponseKey(resp *client.Response, format string) {
	// Format the result.
	switch format {
	case "simple":
		if resp.Action != "delete" {
			fmt.Println(resp.Node.Value)
		} else {
			fmt.Println("PrevNode.Value:", resp.PrevNode.Value)
		}
	case "extended":
		// Extended prints in a rfc2822 style format
		fmt.Println("Key:", resp.Node.Key)
		fmt.Println("Created-Index:", resp.Node.CreatedIndex)
		fmt.Println("Modified-Index:", resp.Node.ModifiedIndex)

		if resp.PrevNode != nil {
			fmt.Println("PrevNode.Value:", resp.PrevNode.Value)
		}

		fmt.Println("TTL:", resp.Node.TTL)
		fmt.Println("Index:", resp.Index)
		if resp.Action != "delete" {
			fmt.Println("")
			fmt.Println(resp.Node.Value)
		}
	case "json":
		b, err := json.Marshal(resp)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(b))
	default:
		fmt.Fprintln(os.Stderr, "Unsupported output format:", format)
	}
}
