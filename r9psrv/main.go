
package main

import (
	"fmt"
	"os"
	
	"github.com/cmars/replican-sync/replican/fs"
	"github.com/cmars/replican-sync/replican/r9p/srv"
)

func main() {
	if (len(os.Args) < 2) {
		fmt.Printf("Usage: %s <path>\n", os.Args[0])
		os.Exit(1)
	}
	
	path := os.Args[1]
	
	store, err := fs.NewLocalStore(path)
	if err != nil { die(fmt.Sprintf("Failed to read source %s", path), err) }
	
	srv, err := srv.NewStoreSrv(store)
	if err != nil { die("Failed to create server", err) }
	
	srv.Start(srv)
	err = srv.StartNetListener("tcp", ":5640")
	if err != nil {
		die("Failed to start server", err)
	}
}

func die(message string, err os.Error) {
	fmt.Fprintf(os.Stderr, "%s: %v\n", message, err)
	os.Exit(1)
}



