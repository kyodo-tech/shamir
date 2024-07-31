package main

import (
	"bufio"
	"encoding/base64"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/kyodo-tech/shamir"
)

func main() {
	mode := flag.String("mode", "split", "Mode: split or combine")
	encoding := flag.String("encoding", "base64", "Encoding: base64 or hex")
	secret := flag.String("secret", "", "The secret to split (for split mode)")
	sharesStr := flag.String("shares", "", "Comma-separated shares (for combine mode)")
	n := flag.Int("n", 5, "Total number of shares")
	t := flag.Int("t", 3, "Number of shares needed to reconstruct the secret")
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
			splitSecret(inputSecret, *n, *t, *encoding)
		} else {
			splitSecret(*secret, *n, *t, *encoding)
		}
	case "combine":
		if *sharesStr == "" {
			fmt.Println("Usage: go run main.go -mode=combine -shares=<share1>,<share2>,...,<shareN> -encoding=<base64|hex>")
			os.Exit(1)
		}
		combineShares(*sharesStr, *encoding)
	default:
		fmt.Println("Invalid mode. Use -mode=split or -mode=combine")
		os.Exit(1)
	}
}

func splitSecret(secret string, totalShares, T int, encoding string) {
	secretBytes := []byte(secret)
	shares, err := shamir.Split(secretBytes, totalShares, T)
	if err != nil {
		panic(err)
	}

	// Print share strings based on encoding
	for _, share := range shares {
		switch encoding {
		case "base64":
			fmt.Println(base64.StdEncoding.EncodeToString(share))
		case "hex":
			fmt.Println(hex.EncodeToString(share))
		default:
			fmt.Println("Invalid encoding. Use base64 or hex")
			os.Exit(1)
		}
	}
}

func combineShares(sharesStr, encoding string) {
	shareStrs := strings.Split(sharesStr, ",")

	var shares [][]byte
	for _, shareStr := range shareStrs {
		var share []byte
		var err error
		switch encoding {
		case "base64":
			share, err = base64.StdEncoding.DecodeString(shareStr)
		case "hex":
			share, err = hex.DecodeString(shareStr)
		default:
			fmt.Println("Invalid encoding. Use base64 or hex")
			os.Exit(1)
		}

		if err != nil {
			panic(err)
		}

		shares = append(shares, share)
	}

	reconstructed, err := shamir.Combine(shares)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(reconstructed))
}
