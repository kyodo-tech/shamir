package shamir

import (
	"crypto/rand"
	"fmt"
)

type polynomial struct {
	coefficients []uint8
}

func newPolynomial(intercept, degree uint8) (*polynomial, error) {
	p := &polynomial{
		coefficients: make([]uint8, degree+1),
	}
	p.coefficients[0] = intercept

	if _, err := rand.Read(p.coefficients[1:]); err != nil {
		return nil, err
	}

	return p, nil
}

func (p *polynomial) evaluate(x uint8) uint8 {
	if x == 0 {
		return p.coefficients[0]
	}

	out := p.coefficients[len(p.coefficients)-1]
	for i := len(p.coefficients) - 2; i >= 0; i-- {
		out = add(mult(out, x), p.coefficients[i])
	}
	return out
}

// Split divides the secret into parts shares with a threshold of minimum shares to reconstruct the secret.
func Split(secret []byte, N, T int) ([][]byte, error) {
	if N < T {
		return nil, fmt.Errorf("parts cannot be less than threshold")
	}
	if N > 255 {
		return nil, fmt.Errorf("parts cannot exceed 255")
	}
	if T < 2 {
		return nil, fmt.Errorf("threshold must be at least 2")
	}
	if T > 255 {
		return nil, fmt.Errorf("threshold cannot exceed 255")
	}
	if len(secret) == 0 {
		return nil, fmt.Errorf("cannot split an empty secret")
	}

	// Generate unique x-coordinates for each share
	xCoordinates := make([]uint8, N)
	for i := 0; i < N; i++ {
		xCoordinates[i] = uint8(i + 1)
	}

	// Initialize shares with the secret length + 1 (for the x-coordinate)
	shares := make([][]byte, N)
	for i := range shares {
		shares[i] = make([]byte, len(secret)+1)
		shares[i][len(secret)] = xCoordinates[i]
	}

	// Create a polynomial for each byte in the secret and evaluate it at each x-coordinate
	for i, b := range secret {
		p, err := newPolynomial(b, uint8(T-1))
		if err != nil {
			return nil, err
		}

		for j := 0; j < N; j++ {
			shares[j][i] = p.evaluate(xCoordinates[j])
		}
	}

	return shares, nil
}

// Combine reconstructs the secret from the provided shares.
func Combine(shares [][]byte) ([]byte, error) {
	if len(shares) < 2 {
		return nil, fmt.Errorf("less than two shares cannot be used to reconstruct the secret")
	}

	shareLength := len(shares[0])
	if shareLength < 2 {
		return nil, fmt.Errorf("shares must be at least two bytes long")
	}

	for _, share := range shares {
		if len(share) != shareLength {
			return nil, fmt.Errorf("all shares must be the same length")
		}
	}

	secret := make([]byte, shareLength-1)
	xSamples := make([]uint8, len(shares))
	ySamples := make([]uint8, len(shares))

	for i, share := range shares {
		xSamples[i] = share[shareLength-1]
	}

	for i := range secret {
		for j, share := range shares {
			ySamples[j] = share[i]
		}

		val, err := interpolatePolynomialSafe(xSamples, ySamples, 0)
		if err != nil {
			return nil, err
		}

		secret[i] = val
	}

	return secret, nil
}

func interpolatePolynomialSafe(xSamples, ySamples []uint8, x uint8) (uint8, error) {
	result := uint8(0)
	for i := range xSamples {
		num, denom := uint8(1), uint8(1)
		for j := range xSamples {
			if i != j {
				num = mult(num, add(x, xSamples[j]))
				denom = mult(denom, add(xSamples[i], xSamples[j]))
			}
		}
		term, err := div(num, denom)
		if err != nil {
			return 0, err
		}
		result = add(result, mult(ySamples[i], term))
	}
	return result, nil
}

// Helper functions for arithmetic in GF(2^8)

func add(a, b uint8) uint8 {
	return a ^ b
}

func mult(a, b uint8) uint8 {
	var p uint8
	for b > 0 {
		if b&1 == 1 {
			p ^= a
		}
		if a&0x80 > 0 {
			a = (a << 1) ^ 0x1B
		} else {
			a <<= 1
		}
		b >>= 1
	}
	return p
}

func div(a, b uint8) (uint8, error) {
	if b == 0 {
		return 0, fmt.Errorf("division by zero")
	}
	return mult(a, inverse(b)), nil
}

func inverse(a uint8) uint8 {
	var b, c uint8
	for b = 1; b != 0; b++ {
		if mult(a, b) == 1 {
			c = b
			break
		}
	}
	return c
}
