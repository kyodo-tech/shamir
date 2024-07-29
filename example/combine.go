package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"github.com/kyodo-tech/shamir"
)

func main() {
	// read comma separated shares from command line
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <share1>,<share2>,...,<shareN>")
		os.Exit(1)
	}

	// split shares
	sharesStr := strings.Split(os.Args[1], ",")

	var shares [][]byte
	for _, shareStr := range sharesStr {
		share, err := base64.StdEncoding.DecodeString(shareStr)
		if err != nil {
			panic(err)
		}

		shares = append(shares, share)
	}

	// Use any 3 out of 5 shares to reconstruct the secret
	reconstructed, err := shamir.Combine(shares[:3])
	if err != nil {
		panic(err)
	}

	fmt.Printf("Reconstructed secret: %s\n", reconstructed)
}
