package crypto

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"io"
	"strconv"
)

const (
	privKeyLen = 64
	signLen    = 64
	pubKeyLen  = 32
	seedLen    = 32
	addressLen = 20
)

type PrivateKey struct {
	key ed25519.PrivateKey
}

func NewPrivateKeyFromSeed(seed []byte) *PrivateKey {
	if len(seed) != seedLen {
		panic("seed len must be " + strconv.Itoa(seedLen))
	}

	return &PrivateKey{
		key: ed25519.NewKeyFromSeed(seed),
	}
}

func NewPrivateKetFromString(s string) *PrivateKey {
	seed, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return NewPrivateKeyFromSeed(seed)
}

func GeneratePrivateKey() *PrivateKey {
	seed := make([]byte, seedLen)
	_, err := io.ReadFull(rand.Reader, seed)

	if err != nil {
		panic(err)
	}

	return NewPrivateKeyFromSeed(seed)
}

func (p *PrivateKey) Bytes() []byte {
	return p.key
}

func (p *PrivateKey) Sign(msg []byte) *Signature {
	return &Signature{
		value: ed25519.Sign(p.key, msg),
	}
}

func (p *PrivateKey) Public() *PublicKey {
	b := make([]byte, pubKeyLen)
	copy(b, p.key[32:])

	return &PublicKey{
		key: b,
	}
}

type PublicKey struct {
	key ed25519.PublicKey
}

func NewPublicKeyFromBytes(b []byte) *PublicKey {
	if len(b) != pubKeyLen {
		panic("public key len must be " + strconv.Itoa(pubKeyLen))
	}

	return &PublicKey{
		key: ed25519.PublicKey(b),
	}
}

func (p *PublicKey) Address() Address {
	return Address{
		value: p.key[len(p.key)-addressLen:],
	}
}

func (p *PublicKey) Bytes() []byte {
	return p.key
}

type Signature struct {
	value []byte
}

func NewSignatureFromBytes(b []byte) *Signature {
	if len(b) != signLen {
		panic("signature len must be " + strconv.Itoa(signLen))
	}

	return &Signature{
		value: b,
	}
}

func (s *Signature) Bytes() []byte {
	return s.value
}

func (s *Signature) Verify(pubKey *PublicKey, msg []byte) bool {
	return ed25519.Verify(pubKey.Bytes(), msg, s.Bytes())
}

type Address struct {
	value []byte
}

func (a Address) String() string {
	return hex.EncodeToString(a.value)
}

func (a Address) Bytes() []byte {
	return a.value
}
