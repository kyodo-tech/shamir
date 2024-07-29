package main

import (
	"bufio"
	"encoding/base64"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/kyodo-tech/shamir"
)

func main() {
	mode := flag.String("mode", "split", "Mode: split or combine")
	secret := flag.String("secret", "", "The secret to split (for split mode)")
	sharesStr := flag.String("shares", "", "Comma-separated shares (for combine mode)")
	N := flag.Int("n", 5, "Total number of shares")
	T := flag.Int("t", 3, "Number of shares needed to reconstruct the secret")
	flag.Parse()

	switch *mode {
	case "split":
		if *secret == "" {
			reader := bufio.NewReader(os.Stdin)
			inputSecret, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Error reading input:", err)
				os.Exit(1)
			}
			inputSecret = strings.TrimSpace(inputSecret)
			splitSecret(inputSecret, *N, *T)
		} else {
			splitSecret(*secret, *N, *T)
		}
	case "combine":
		if *sharesStr == "" {
			fmt.Println("Usage: go run main.go -mode=combine -shares=<share1>,<share2>,...,<shareN>")
			os.Exit(1)
		}
		combineShares(*sharesStr)
	default:
		fmt.Println("Invalid mode. Use -mode=split or -mode=combine")
		os.Exit(1)
	}
}

func splitSecret(secret string, totalShares, T int) {
	secretBytes := []byte(secret)
	shares, err := shamir.Split(secretBytes, totalShares, T)
	if err != nil {
		panic(err)
	}

	// Print share strings
	for _, share := range shares {
		fmt.Println(base64.StdEncoding.EncodeToString(share))
	}
}

func combineShares(sharesStr string) {
	shareStrs := strings.Split(sharesStr, ",")

	var shares [][]byte
	for _, shareStr := range shareStrs {
		share, err := base64.StdEncoding.DecodeString(shareStr)
		if err != nil {
			panic(err)
		}

		shares = append(shares, share)
	}

	reconstructed, err := shamir.Combine(shares)
	if err != nil {
		panic(err)
	}

	fmt.Println(reconstructed)
}
