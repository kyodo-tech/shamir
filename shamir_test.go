package shamir

import (
	"bytes"
	"testing"
)

func TestSplitInvalid(t *testing.T) {
	secret := []byte("test")

	tests := []struct {
		parts     int
		threshold int
	}{
		{0, 0},
		{2, 3},
		{1000, 3},
		{10, 1},
		{3, 256},
	}

	for _, tt := range tests {
		if _, err := Split(secret, tt.parts, tt.threshold); err == nil {
			t.Fatalf("expected error for parts: %d, threshold: %d", tt.parts, tt.threshold)
		}
	}

	if _, err := Split(nil, 3, 2); err == nil {
		t.Fatalf("expected error for nil secret")
	}
}

func TestSplit(t *testing.T) {
	secret := []byte("test")

	out, err := Split(secret, 5, 3)
	if err != nil {
		t.Fatalf("Split error: %v", err)
	}

	if len(out) != 5 {
		t.Fatalf("expected 5 shares, got %d", len(out))
	}

	for _, share := range out {
		if len(share) != len(secret)+1 {
			t.Fatalf("expected share length %d, got %d", len(secret)+1, len(share))
		}
	}
}

func TestCombineInvalid(t *testing.T) {
	tests := [][][]byte{
		nil,
		{[]byte("foo"), []byte("ba")},
		{[]byte("f"), []byte("b")},
		{[]byte("foo"), []byte("foo")},
	}

	for _, parts := range tests {
		if _, err := Combine(parts); err == nil {
			t.Fatalf("expected error for parts: %v", parts)
		}
	}
}

func TestCombine(t *testing.T) {
	secret := []byte("test")

	out, err := Split(secret, 5, 3)
	if err != nil {
		t.Fatalf("Split error: %v", err)
	}

	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			if j == i {
				continue
			}
			for k := 0; k < 5; k++ {
				if k == i || k == j {
					continue
				}
				parts := [][]byte{out[i], out[j], out[k]}
				recomb, err := Combine(parts)
				if err != nil {
					t.Fatalf("Combine error: %v", err)
				}

				if !bytes.Equal(recomb, secret) {
					t.Fatalf("expected %v, got %v", secret, recomb)
				}
			}
		}
	}
}

func TestFieldOperations(t *testing.T) {
	tests := []struct {
		a, b, expected uint8
		op             func(uint8, uint8) (uint8, error)
	}{
		{16, 16, 0, func(a, b uint8) (uint8, error) { return gfAdd(a, b), nil }},
		{3, 4, 7, func(a, b uint8) (uint8, error) { return gfAdd(a, b), nil }},
		{3, 7, 9, func(a, b uint8) (uint8, error) { return gfMult(a, b), nil }},
		{3, 0, 0, func(a, b uint8) (uint8, error) { return gfMult(a, b), nil }},
		{0, 3, 0, func(a, b uint8) (uint8, error) { return gfMult(a, b), nil }},
		{0, 7, 0, gfDiv},
		{3, 3, 1, gfDiv},
		{6, 3, 2, gfDiv},
	}

	for _, tt := range tests {
		out, err := tt.op(tt.a, tt.b)
		if err != nil {
			t.Fatalf("operation error: %v", err)
		}
		if out != tt.expected {
			t.Fatalf("expected %d, got %d", tt.expected, out)
		}
	}
}

func TestPolynomialCreationAndEvaluation(t *testing.T) {
	p, err := newPolynomial(42, 1)
	if err != nil {
		t.Fatalf("NewPolynomial error: %v", err)
	}

	if p.coeffs[0] != 42 {
		t.Fatalf("expected intercept 42, got %d", p.coeffs[0])
	}

	if out := p.evaluate(0); out != 42 {
		t.Fatalf("expected 42, got %d", out)
	}

	x := uint8(1)
	expected := gfAdd(42, gfMult(x, p.coeffs[1]))
	if out := p.evaluate(x); out != expected {
		t.Fatalf("expected %d, got %d", expected, out)
	}
}

func TestPolynomialInterpolation(t *testing.T) {
	for i := 0; i < 256; i++ {
		p, err := newPolynomial(uint8(i), 2)
		if err != nil {
			t.Fatalf("NewPolynomial error: %v", err)
		}

		xVals := []uint8{1, 2, 3}
		yVals := []uint8{p.evaluate(1), p.evaluate(2), p.evaluate(3)}
		out, err := interpolatePolynomial(xVals, yVals, 0)
		if err != nil {
			t.Fatalf("interpolatePolynomial error: %v", err)
		}
		if out != uint8(i) {
			t.Fatalf("expected %d, got %d", uint8(i), out)
		}
	}
}
