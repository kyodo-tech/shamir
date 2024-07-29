package main

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/kyodo-tech/shamir"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <secret>")
		os.Exit(1)
	}

	secret := []byte(os.Args[1])
	shares, err := shamir.Split(secret, 5, 3)
	if err != nil {
		panic(err)
	}

	// print share strings
	for i, share := range shares {
		fmt.Printf("Share %d: %s\n", i, base64.StdEncoding.EncodeToString(share))
	}
}
