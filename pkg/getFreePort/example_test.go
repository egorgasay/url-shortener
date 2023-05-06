package getfreeport

import (
	"fmt"
	"log"
)

func Example() {
	var flag bool
	var err error

	port := ""
	for !flag {
		port, err = GetFreePort()
		if err != nil {
			log.Fatalf("Oops, couldn't get a free port: %v", err)
		}
	}

	fmt.Println("Free port received:", port)
}
