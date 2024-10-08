package bls

import (
	"math/big"

	bn256 "github.com/umbracle/go-eth-bn256"
)

// PrivateKey holds private key for bls implementation
type PrivateKey struct {
	s *big.Int
}

// NewZeroPrivateKey is the constructor of an empty PrivateKey
func NewZeroPrivateKey() *PrivateKey {
	return &PrivateKey{
		s: new(big.Int),
	}
}

// UnmarshalPrivateKey unmarshals the private key from the given byte slice.
// This function supports both raw big int string and hex-encoded big int string.
func UnmarshalPrivateKey(data []byte) (*PrivateKey, error) {
	pk := NewZeroPrivateKey()

	// First trying to use a default unmarshaling function of big int.
	// It works for the raw big int string format and for hex with 0x prefix.
	if err := pk.s.UnmarshalText(data); err == nil {
		return pk, nil
	}

	// Otherwise, trying to assume the given data is a hex encoded big int represented as a bytes array.
	// This is needed in order to be compatible with the currently stored polybft BLS keys.
	pk.s.SetBytes(data)

	return pk, nil
}

// PublicKey returns the public key from the PrivateKey
func (p *PrivateKey) PublicKey() *PublicKey {
	g2 := new(bn256.G2).ScalarMult(g2Point, p.s)

	return &PublicKey{g2: g2}
}

// Marshal marshals private key hex (without 0x prefix) represented as a byte slice
func (p *PrivateKey) Marshal() ([]byte, error) {
	return p.s.Bytes(), nil
}

// Sign generates a simple BLS signature of the given message
func (p *PrivateKey) Sign(message, domain []byte) (*Signature, error) {
	point, err := hashToPoint(message, domain)
	if err != nil {
		return nil, err
	}

	g1 := new(bn256.G1).ScalarMult(point, p.s)

	return &Signature{g1: g1}, nil
}
