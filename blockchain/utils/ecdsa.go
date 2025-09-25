package utils

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"fmt"
	"math/big"
)

type Signature struct {
	R *big.Int
	S *big.Int
}

func (s *Signature) String() string {
	return fmt.Sprintf("%064x%064x", s.R, s.S)
}

func String2BigIntTuple(s string) (big.Int, big.Int) {
	bx, _ := hex.DecodeString(s[:64])
	by, _ := hex.DecodeString(s[64:])
	var r, t big.Int
	r.SetBytes(bx)
	t.SetBytes(by)
	return r, t
}

func SignatureFromString(s string) *Signature {
	r, t := String2BigIntTuple(s)
	return &Signature{&r, &t}
}

func PublicKeyFromString(s string) *ecdsa.PublicKey {
	x, y := String2BigIntTuple(s)
	return &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     &x,
		Y:     &y,
	}
}

func PrivateKeyFromString(s string, publicKey *ecdsa.PublicKey) *ecdsa.PrivateKey {
	b, _ := hex.DecodeString(s)
	var d big.Int
	d.SetBytes(b)
	return &ecdsa.PrivateKey{
		PublicKey: *publicKey,
		D:         &d,
	}
}
