## Shamir's Secret Sharing

This repository provides a simple implementation of Shamir's Secret Sharing in Go, allowing to split a secret into multiple shares and reconstruct it using a subset of those shares.

- Split a secret into `N` shares with a threshold of `T` shares required to reconstruct the secret.
- Arithmetic operations in Galois Field (GF(`2^8`)).
- Polynomial creation and evaluation.

### Usage

```go
package main

import (
	"fmt"
	"github.com/kyodo-tech/shamir"
)

func main() {
	secret := []byte("my secret")
	shares, err := shamir.Split(secret, 5, 3)
	if err != nil {
		panic(err)
	}

	// Use any 3 out of 5 shares to reconstruct the secret
	reconstructed, err := shamir.Combine(shares[:3])
	if err != nil {
		panic(err)
	}

	fmt.Printf("Reconstructed secret: %s\n", reconstructed)
}
```

Install `cmd/sss` to use the command line interface:

	go install github.com/kyodo-tech/shamir/cmd/sss@latest

Then run `sss` to split a secret into shares and combine them to reconstruct the secret:

```sh
Usage: sss -mode=split -secret=<secret> -n=<number of shares> -T=<threshold>
# Example:
sss -secret "my secret"
# Output:
# N8QLmd6YcLu2AQ==
# Nnubhbvx52y8Ag==
# bMawbwAK5bJ+Aw==
# SurUYDi7Mfm4BA==
# EFf/ioNAMyd6BQ==
sss -mode combine -shares N8QLmd6YcLu2AQ==,bMawbwAK5bJ+Aw==,EFf/ioNAMyd6BQ==
# Output:
# my secret
```
