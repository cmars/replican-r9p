
package main

import (
	"fmt"
	"os"
	
	"github.com/cmars/replican-sync/replican/fs"
	"github.com/cmars/replican-sync/replican/r9p/clnt"
)

func main() {
	store, err := clnt.Connect("127.0.0.1:5640")
	if err != nil { 
		die("Failed to connect", err);
	}
	
	err = store.Refresh()
	if err != nil { 
		die("Failed to refresh", err);
	}
	
	fs.Walk(store.Root(), func(node fs.Node) bool {
		switch node.(type) {
		case *fs.Dir:
			fmt.Printf("d\t%s\t%s\n", node.Strong(), fs.RelPath(node.(fs.FsNode)))
			return true
		case *fs.File:
			fmt.Printf("f\t%s\t%s\n", node.Strong(), fs.RelPath(node.(fs.FsNode)))
		}
		return false
	})
}

func die(message string, err os.Error) {
	fmt.Fprintf(os.Stderr, "%s: %v\n", message, err)
	os.Exit(1)
}


