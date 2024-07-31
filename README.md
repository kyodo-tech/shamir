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

```sh
Usage of sss:
  -encoding string
        Encoding: base64 or hex (default "base64")
  -mode string
        Mode: split or combine (default "split")
  -secret string
        The secret to split (for split mode)
  -shares string
        Comma-separated shares (for combine mode)
  -n int
        Total number of shares (default 5)
  -t int
        Number of shares needed to reconstruct the secret (default 3)
```

Then run `sss` to split a secret into shares and combine them to reconstruct the secret:

```sh
# Example:
sss -secret "my secret" -encoding hex -t 2 -n 5
# 4de08f03d23ab8360d01
# 2d50659310d1fdc38602
# 0dc9cae3a7883790ff03
# ed2baaa88f1c77328b04
# cdb205d83845bd61f205
sss -mode=combine -encoding hex -shares 4de08f03d23ab8360d01,0dc9cae3a7883790ff03
# my secret
```
