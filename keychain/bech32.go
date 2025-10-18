// This file implements Sui's SIP-15 Bech32 private key format.
// Ref: https://github.com/sui-foundation/sips/pull/15
package keychain

import (
	"errors"
	"fmt"

	"github.com/btcsuite/btcutil/bech32"
)

const (
	suiPrivateKeyHRP = "suiprivkey"
	privateKeySize   = 32
	flagSize         = 1
)

var ErrInvalidSecretLength = errors.New("secret key must be 32 bytes")

type ParsedPrivateKey struct {
	Scheme    Scheme
	SecretKey []byte
}

// EncodePrivateKey serializes flag||secret as Bech32 with HRP "suiprivkey".
func EncodePrivateKey(s Scheme, secretKey []byte) (string, error) {
	if len(secretKey) != privateKeySize {
		return "", fmt.Errorf("bech32: private key must be %d bytes, got %d", privateKeySize, len(secretKey))
	}

	flag := s.AddressFlag()
	if flag == flagUnspecified {
		return "", fmt.Errorf("bech32: unsupported scheme %d", s)
	}

	payload := make([]byte, flagSize+privateKeySize)
	payload[0] = flag
	copy(payload[1:], secretKey)
	defer zeroBytes(payload)

	words, err := bech32.ConvertBits(payload, 8, 5, true)
	if err != nil {
		return "", fmt.Errorf("bech32: convert bits: %w", err)
	}

	encoded, err := bech32.Encode(suiPrivateKeyHRP, words)
	if err != nil {
		return "", fmt.Errorf("bech32: encode: %w", err)
	}

	return encoded, nil
}

// DecodePrivateKey reverses EncodePrivateKey, yielding the scheme and secret bytes.
func DecodePrivateKey(encoded string) (*ParsedPrivateKey, error) {
	hrp, words, err := bech32.Decode(encoded)
	if err != nil {
		return nil, fmt.Errorf("bech32: decode: %w", err)
	}

	if hrp != suiPrivateKeyHRP {
		return nil, fmt.Errorf("bech32: unexpected hrp %q", hrp)
	}

	payload, err := bech32.ConvertBits(words, 5, 8, false)
	if err != nil {
		return nil, fmt.Errorf("bech32: convert bits: %w", err)
	}

	if len(payload) != flagSize+privateKeySize {
		return nil, fmt.Errorf("bech32: invalid payload length %d", len(payload))
	}

	scheme, err := SchemeFromFlag(payload[0])
	if err != nil {
		return nil, err
	}

	secret := make([]byte, privateKeySize)
	copy(secret, payload[1:])
	zeroBytes(payload)

	return &ParsedPrivateKey{Scheme: scheme, SecretKey: secret}, nil
}

func PrivateKeySize() int {
	return privateKeySize
}

func zeroBytes(b []byte) {
	for i := range b {
		b[i] = 0
	}
}
