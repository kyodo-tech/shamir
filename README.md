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

Also see the two example programs in the `example` directory:

```sh
go run example/split.go "my secret"
# Output:
# Share 0: lK8LyefHjxzOAQ==
# Share 1: 1Ws3IGJh/lVQAg==
# Share 2: LL0cmuDFAyzqAw==
# Share 3: Wf8RJQs43iALBA==
# Share 4: oCk6n4mcI1mxBQ==
go run ./example/combine.go lK8LyefHjxzOAQ==,LL0cmuDFAyzqAw==,oCk6n4mcI1mxBQ==
# Output:
# Reconstructed secret: my secret
```
