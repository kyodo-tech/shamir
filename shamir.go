package shamir

import (
	"crypto/rand"
	"errors"
	"fmt"
)

var (
	ErrPartsLessThanThreshold  = errors.New("number of parts cannot be less than the threshold")
	ErrPartsExceedLimit        = errors.New("number of parts cannot exceed 255")
	ErrThresholdTooSmall       = errors.New("threshold must be at least 2")
	ErrThresholdExceedLimit    = errors.New("threshold cannot exceed 255")
	ErrEmptySecret             = errors.New("cannot split an empty secret")
	ErrInsufficientShares      = errors.New("less than two shares cannot be used to reconstruct the secret")
	ErrSharesTooShort          = errors.New("shares must be at least two bytes long")
	ErrInconsistentShareLength = errors.New("all shares must be the same length")
	ErrDivisionByZero          = errors.New("division by zero")
)

type polynomial struct {
	coeffs []uint8
}

func newPolynomial(intercept, degree uint8) (*polynomial, error) {
	p := &polynomial{
		coeffs: make([]uint8, degree+1),
	}
	p.coeffs[0] = intercept

	if _, err := rand.Read(p.coeffs[1:]); err != nil {
		return nil, fmt.Errorf("failed to generate random coefficients: %w", err)
	}

	return p, nil
}

func (p *polynomial) evaluate(x uint8) uint8 {
	if x == 0 {
		return p.coeffs[0]
	}

	result := p.coeffs[len(p.coeffs)-1]
	for i := len(p.coeffs) - 2; i >= 0; i-- {
		result = gfAdd(gfMult(result, x), p.coeffs[i])
	}
	return result
}

// Split divides the secret into n shares with a threshold t for reconstruction.
func Split(secret []byte, n, t int) ([][]byte, error) {
	if n < t {
		return nil, ErrPartsLessThanThreshold
	}
	if n > 255 {
		return nil, ErrPartsExceedLimit
	}
	if t < 2 {
		return nil, ErrThresholdTooSmall
	}
	if t > 255 {
		return nil, ErrThresholdExceedLimit
	}
	if len(secret) == 0 {
		return nil, ErrEmptySecret
	}

	xCoords := make([]uint8, n)
	for i := 0; i < n; i++ {
		xCoords[i] = uint8(i + 1)
	}

	shares := make([][]byte, n)
	for i := range shares {
		shares[i] = make([]byte, len(secret)+1)
		shares[i][len(secret)] = xCoords[i]
	}

	for i, b := range secret {
		p, err := newPolynomial(b, uint8(t-1))
		if err != nil {
			return nil, err
		}

		for j := 0; j < n; j++ {
			shares[j][i] = p.evaluate(xCoords[j])
		}
	}

	return shares, nil
}

// Combine reconstructs the secret from the shares.
func Combine(shares [][]byte) ([]byte, error) {
	if len(shares) < 2 {
		return nil, ErrInsufficientShares
	}

	shareLen := len(shares[0])
	if shareLen < 2 {
		return nil, ErrSharesTooShort
	}

	for _, share := range shares {
		if len(share) != shareLen {
			return nil, ErrInconsistentShareLength
		}
	}

	secret := make([]byte, shareLen-1)
	xSamples := make([]uint8, len(shares))
	ySamples := make([]uint8, len(shares))

	for i, share := range shares {
		xSamples[i] = share[shareLen-1]
	}

	for i := range secret {
		for j, share := range shares {
			ySamples[j] = share[i]
		}

		val, err := interpolatePolynomial(xSamples, ySamples, 0)
		if err != nil {
			return nil, err
		}

		secret[i] = val
	}

	return secret, nil
}

func interpolatePolynomial(xSamples, ySamples []uint8, x uint8) (uint8, error) {
	result := uint8(0)
	for i := range xSamples {
		num, denom := uint8(1), uint8(1)
		for j := range xSamples {
			if i != j {
				num = gfMult(num, gfAdd(x, xSamples[j]))
				denom = gfMult(denom, gfAdd(xSamples[i], xSamples[j]))
			}
		}
		term, err := gfDiv(num, denom)
		if err != nil {
			return 0, err
		}
		result = gfAdd(result, gfMult(ySamples[i], term))
	}
	return result, nil
}

// Helper functions for arithmetic in GF(2^8)

func gfAdd(a, b uint8) uint8 {
	return a ^ b
}

func gfMult(a, b uint8) uint8 {
	var product uint8
	for b > 0 {
		if b&1 == 1 {
			product ^= a
		}
		if a&0x80 > 0 {
			a = (a << 1) ^ 0x1B
		} else {
			a <<= 1
		}
		b >>= 1
	}
	return product
}

func gfDiv(a, b uint8) (uint8, error) {
	if b == 0 {
		return 0, ErrDivisionByZero
	}
	return gfMult(a, gfInverse(b)), nil
}

func gfInverse(a uint8) uint8 {
	var inv uint8
	for b := uint8(1); b != 0; b++ {
		if gfMult(a, b) == 1 {
			inv = b
			break
		}
	}
	return inv
}
