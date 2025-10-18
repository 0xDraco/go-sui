package secp256r1

import (
	"crypto/ecdh"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math/big"

	"github.com/0xdraco/go-sui/keychain"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
)

type Keypair struct {
	PrivateKey *ecdsa.PrivateKey
	ChainCode  []byte
	Path       keychain.DerivationPath
}

func (k Keypair) Scheme() keychain.Scheme {
	return keychain.SchemeSecp256r1
}

func (k Keypair) PrivateKeyBytes() []byte {
	b := k.PrivateKey.D.Bytes()
	if len(b) < 32 {
		padded := make([]byte, 32)
		copy(padded[32-len(b):], b)
		b = padded
	}
	return b
}

func (k Keypair) PublicKeyBytes() []byte {
	curve := elliptic.P256()
	return elliptic.MarshalCompressed(curve, k.PrivateKey.X, k.PrivateKey.Y)
}

func (k Keypair) SuiAddress() (string, error) {
	return keychain.AddressFromPublicKey(keychain.SchemeSecp256r1, k.PublicKeyBytes())
}

func (k Keypair) SecretKeyBytes() []byte {
	return k.PrivateKeyBytes()
}

func (k Keypair) PublicKeyBase64() string {
	payload := append([]byte{keychain.SchemeSecp256r1.AddressFlag()}, k.PublicKeyBytes()...)
	return base64.StdEncoding.EncodeToString(payload)
}

func Generate() (*Keypair, error) {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("secp256r1: generate key: %w", err)
	}
	return &Keypair{
		PrivateKey: priv,
	}, nil
}

func FromSecretKey(secret []byte) (*Keypair, error) {
	priv, _, err := deriveECDSAKey(secret)
	if err != nil {
		return nil, err
	}

	return &Keypair{PrivateKey: priv}, nil
}

// Mirrors Mysten's tooling by running BIP-32 derivation with secp256k1 group
// operations before converting the resulting scalar into a P-256 keypair.
func Derive(seed []byte, path keychain.DerivationPath) (*Keypair, error) {
	if err := path.ValidateForScheme(keychain.SchemeSecp256r1); err != nil {
		return nil, err
	}

	key, chain := keychain.BIP32MasterPrivateKey(seed)
	secpOrder := secp256k1.S256().N
	for _, segment := range path.Segments() {
		nextKey, nextChain, err := keychain.DeriveChildPrivateKey(key, chain, segment, func(s []byte) ([]byte, error) {
			privKey := secp256k1.PrivKeyFromBytes(s)
			if privKey == nil || privKey.Key.IsZeroBit() == 1 {
				return nil, fmt.Errorf("secp256r1: invalid intermediate key")
			}
			return privKey.PubKey().SerializeCompressed(), nil
		}, secpOrder)
		if err != nil {
			return nil, err
		}
		key = nextKey
		chain = nextChain
	}

	priv, _, err := deriveECDSAKey(key)
	if err != nil {
		return nil, err
	}

	return &Keypair{
		PrivateKey: priv,
		ChainCode:  append([]byte{}, chain...),
		Path:       path,
	}, nil
}

// Converts a raw scalar into an ECDSA private key using crypto/ecdh to obtain
// the corresponding public point and compressed SEC1 encoding.
func deriveECDSAKey(secret []byte) (*ecdsa.PrivateKey, []byte, error) {
	if len(secret) != keychain.PrivateKeySize() {
		return nil, nil, keychain.ErrInvalidSecretLength
	}
	secretCopy := make([]byte, keychain.PrivateKeySize())
	copy(secretCopy, secret)
	curve := ecdh.P256()
	priv, err := curve.NewPrivateKey(secretCopy)
	zero(secretCopy)
	if err != nil {
		return nil, nil, fmt.Errorf("secp256r1: invalid private key: %w", err)
	}
	ecCurve := elliptic.P256()
	pubBytes := priv.PublicKey().Bytes()
	if len(pubBytes) != 1+2*keychain.PrivateKeySize() || pubBytes[0] != 0x04 {
		return nil, nil, fmt.Errorf("secp256r1: unexpected public key encoding")
	}
	x := new(big.Int).SetBytes(pubBytes[1 : 1+keychain.PrivateKeySize()])
	y := new(big.Int).SetBytes(pubBytes[1+keychain.PrivateKeySize():])
	zero(pubBytes)
	compressed := elliptic.MarshalCompressed(ecCurve, x, y)
	d := new(big.Int).SetBytes(secret)
	order := ecCurve.Params().N
	if d.Sign() <= 0 || d.Cmp(order) >= 0 {
		zero(compressed)
		return nil, nil, fmt.Errorf("secp256r1: private key out of range")
	}
	ecdsaPriv := &ecdsa.PrivateKey{D: d}
	ecdsaPriv.PublicKey = ecdsa.PublicKey{Curve: ecCurve, X: x, Y: y}
	return ecdsaPriv, compressed, nil
}

func zero(b []byte) {
	for i := range b {
		b[i] = 0
	}
}
