package blockchain

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
)

type Signature string

func NewSignature(r, s *big.Int) Signature {
	signature := append(r.Bytes(), s.Bytes()...)
	return Signature(hex.EncodeToString(signature))
}

func (signature Signature) RS() (r, s *big.Int, err error) {
	signatureBytes, err := hex.DecodeString(string(signature))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decode string %s: %v", signature, err)
	}

	curveOrderByteSize := elliptic.P256().Params().BitSize / 8

	r = new(big.Int).SetBytes(signatureBytes[:curveOrderByteSize])
	s = new(big.Int).SetBytes(signatureBytes[curveOrderByteSize:])

	return r, s, nil
}

func (signature Signature) Verify(address Address, data []byte) (bool, error) {
	publicKey, err := address.PublicKey()
	if err != nil {
		return false, fmt.Errorf("failed to get public key from address %s: %v", address, err)
	}

	r, s, err := signature.RS()
	if err != nil {
		return false, fmt.Errorf("failed to get RS %s: %v", s, err)
	}

	return ecdsa.Verify(publicKey, data, r, s), nil
}

type Address string

func NewAddressFromPublicKey(publicKey *ecdsa.PublicKey) Address {
	publicKeyBytes := append(publicKey.X.Bytes(), publicKey.Y.Bytes()...)
	return Address(hex.EncodeToString(publicKeyBytes[:]))
}

func (a Address) PublicKey() (*ecdsa.PublicKey, error) {
	publicKeyBytes, err := hex.DecodeString(string(a))
	if err != nil {
		return nil, fmt.Errorf("failed to get bytes %s: %v", a, err)
	}

	curve := elliptic.P256()
	curveSize := curve.Params().BitSize / 8

	x := new(big.Int).SetBytes(publicKeyBytes[:curveSize])
	y := new(big.Int).SetBytes(publicKeyBytes[curveSize:])

	publicKey := ecdsa.PublicKey{
		Curve: curve,
		X:     x,
		Y:     y,
	}

	return &publicKey, nil
}

type Wallet string

func NewRandomWallet() Wallet {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}

	dBytes := privateKey.D.Bytes()
	return Wallet(hex.EncodeToString(dBytes))
}

func NewWalletFromString(wallet string) Wallet {
	return Wallet(wallet)
}

func (w Wallet) PrivateKey() (*ecdsa.PrivateKey, error) {
	dBytes, err := hex.DecodeString(string(w))
	if err != nil {
		return nil, fmt.Errorf("failed to decode string %s: %v", w, err)
	}

	publicKeyX, publicKeyY := elliptic.P256().ScalarBaseMult(dBytes)

	publicKey := ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     publicKeyX,
		Y:     publicKeyY,
	}

	d := new(big.Int).SetBytes(dBytes)

	return &ecdsa.PrivateKey{
		PublicKey: publicKey,
		D:         d,
	}, nil
}

func (w Wallet) Address() (Address, error) {
	privateKey, err := w.PrivateKey()
	if err != nil {
		return "", fmt.Errorf("failed to get private key from %s: %v", w, err)
	}

	return NewAddressFromPublicKey(&privateKey.PublicKey), nil
}

func (w Wallet) Sign(data []byte) (Signature, error) {
	privateKey, err := w.PrivateKey()
	if err != nil {
		return "", fmt.Errorf("failed to get private key %s: %v", w, err)
	}

	r, s, err := ecdsa.Sign(rand.Reader, privateKey, data)
	if err != nil {
		return "", fmt.Errorf("failed to sign block: %v", err)
	}

	signature := NewSignature(r, s)
	return signature, nil
}
