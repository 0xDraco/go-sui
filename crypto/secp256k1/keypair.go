package secp256k1

import (
	"encoding/base64"
	"fmt"
	"math/big"

	"github.com/0xdraco/go-sui/keychain"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
)

type Keypair struct {
	PrivateKey *secp256k1.PrivateKey
	PublicKey  *secp256k1.PublicKey
	ChainCode  []byte
	Path       keychain.DerivationPath
}

func (k Keypair) Scheme() keychain.Scheme {
	return keychain.SchemeSecp256k1
}

func (k Keypair) PrivateKeyBytes() []byte {
	return k.PrivateKey.Serialize()
}

func (k Keypair) PublicKeyBytes() []byte {
	return k.PublicKey.SerializeCompressed()
}

func (k Keypair) SuiAddress() (string, error) {
	return keychain.AddressFromPublicKey(keychain.SchemeSecp256k1, k.PublicKeyBytes())
}

func (k Keypair) SecretKeyBytes() []byte {
	return k.PrivateKey.Serialize()
}

func (k Keypair) PublicKeyBase64() string {
	payload := append([]byte{keychain.SchemeSecp256k1.AddressFlag()}, k.PublicKeyBytes()...)
	return base64.StdEncoding.EncodeToString(payload)
}

func Generate() (*Keypair, error) {
	priv, err := secp256k1.GeneratePrivateKey()
	if err != nil {
		return nil, fmt.Errorf("secp256k1: generate key: %w", err)
	}

	return &Keypair{
		PrivateKey: priv,
		PublicKey:  priv.PubKey(),
	}, nil
}

func FromSecretKey(secret []byte) (*Keypair, error) {
	if len(secret) != keychain.PrivateKeySize() {
		return nil, keychain.ErrInvalidSecretLength
	}

	order := secp256k1.S256().N
	k := new(big.Int).SetBytes(secret)
	if k.Sign() <= 0 || k.Cmp(order) >= 0 {
		return nil, fmt.Errorf("secp256k1: private key out of range")
	}

	priv := secp256k1.PrivKeyFromBytes(secret)
	if priv == nil || priv.Key.IsZeroBit() == 1 {
		return nil, fmt.Errorf("secp256k1: invalid private key")
	}

	return &Keypair{
		PrivateKey: priv,
		PublicKey:  priv.PubKey(),
	}, nil
}

// Walks the BIP-32 derivation path (allowing both hardened and unhardened
// steps) starting from the provided seed and returns the resulting keypair and
// chain code.
func Derive(seed []byte, path keychain.DerivationPath) (*Keypair, error) {
	if err := path.ValidateForScheme(keychain.SchemeSecp256k1); err != nil {
		return nil, err
	}

	key, chain := keychain.BIP32MasterPrivateKey(seed)
	segments := path.Segments()
	for _, segment := range segments {
		nextKey, nextChain, err := keychain.DeriveChildPrivateKey(key, chain, segment, func(priv []byte) ([]byte, error) {
			privKey := secp256k1.PrivKeyFromBytes(priv)
			if privKey == nil || privKey.Key.IsZeroBit() == 1 {
				return nil, fmt.Errorf("secp256k1: invalid private key")
			}

			return privKey.PubKey().SerializeCompressed(), nil
		}, secp256k1.S256().N)
		if err != nil {
			return nil, err
		}

		key = nextKey
		chain = nextChain
	}

	privKey := secp256k1.PrivKeyFromBytes(key)
	if privKey == nil || privKey.Key.IsZeroBit() == 1 {
		return nil, fmt.Errorf("secp256k1: failed to build private key")
	}

	pubKey := privKey.PubKey()
	return &Keypair{
		PrivateKey: privKey,
		PublicKey:  pubKey,
		ChainCode:  append([]byte{}, chain...),
		Path:       path,
	}, nil
}
