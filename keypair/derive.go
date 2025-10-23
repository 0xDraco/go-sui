package keypair

import (
	"fmt"

	ed25519keys "github.com/0xdraco/sui-go-sdk/cryptography/ed25519"
	secp256k1keys "github.com/0xdraco/sui-go-sdk/cryptography/secp256k1"
	secp256r1keys "github.com/0xdraco/sui-go-sdk/cryptography/secp256r1"
	"github.com/0xdraco/sui-go-sdk/keychain"
)

func DeriveFromMnemonic(s keychain.Scheme, mnemonic, passphrase, path string) (Keypair, error) {
	parsed, err := keychain.ParseDerivationPath(path)
	if err != nil {
		return nil, err
	}
	seed, err := keychain.SeedFromMnemonic(mnemonic, passphrase)
	if err != nil {
		return nil, err
	}

	switch s {
	case keychain.SchemeEd25519:
		return ed25519keys.Derive(seed, parsed)
	case keychain.SchemeSecp256k1:
		return secp256k1keys.Derive(seed, parsed)
	case keychain.SchemeSecp256r1:
		return secp256r1keys.Derive(seed, parsed)
	default:
		return nil, fmt.Errorf("derive: unsupported scheme %d", s)
	}
}

func Generate(s keychain.Scheme) (Keypair, error) {
	switch s {
	case keychain.SchemeEd25519:
		return ed25519keys.Generate()
	case keychain.SchemeSecp256k1:
		return secp256k1keys.Generate()
	case keychain.SchemeSecp256r1:
		return secp256r1keys.Generate()
	default:
		return nil, fmt.Errorf("generate: unsupported scheme %d", s)
	}
}

func FromSecretKey(s keychain.Scheme, secret []byte) (Keypair, error) {
	switch s {
	case keychain.SchemeEd25519:
		return ed25519keys.FromSecretKey(secret)
	case keychain.SchemeSecp256k1:
		return secp256k1keys.FromSecretKey(secret)
	case keychain.SchemeSecp256r1:
		return secp256r1keys.FromSecretKey(secret)
	default:
		return nil, fmt.Errorf("from secret: unsupported scheme %d", s)
	}
}

func FromBech32(encoded string) (Keypair, error) {
	parsed, err := keychain.DecodePrivateKey(encoded)
	if err != nil {
		return nil, err
	}

	kp, err := FromSecretKey(parsed.Scheme, parsed.SecretKey)
	zero(parsed.SecretKey)
	return kp, err
}

func ToBech32(k Keypair) (string, error) {
	secret := k.SecretKeyBytes()
	if len(secret) != keychain.PrivateKeySize() {
		return "", fmt.Errorf("export: expected %d secret bytes, got %d", keychain.PrivateKeySize(), len(secret))
	}

	encoded, err := keychain.EncodePrivateKey(k.Scheme(), secret)
	zero(secret)
	return encoded, err
}

func zero(b []byte) {
	for i := range b {
		b[i] = 0
	}
}
