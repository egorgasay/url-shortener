package shortenalgorithm

import (
	"fmt"
	"log"
)

func Example() {
	short, err := GetShortName(100)
	if err != nil {
		log.Fatalf("Oops, unexpected error: %v", short)
	}

	fmt.Println("Short name received:", short)
}
