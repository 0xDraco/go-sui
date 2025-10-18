package ed25519

import (
	cryptoed25519 "crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"fmt"

	"github.com/0xdraco/go-sui/keychain"
)

const seedSize = 32

// The master key HMAC salt for Ed25519 according to SLIP-0010.
var slip10Key = []byte("ed25519 seed")

type Keypair struct {
	PrivateKey cryptoed25519.PrivateKey
	PublicKey  cryptoed25519.PublicKey
	ChainCode  []byte
	Path       keychain.DerivationPath
}

func (k Keypair) Scheme() keychain.Scheme {
	return keychain.SchemeEd25519
}

func (k Keypair) PrivateKeyBytes() []byte {
	return append(cryptoed25519.PrivateKey(nil), k.PrivateKey...)
}

func (k Keypair) PublicKeyBytes() []byte {
	return append(cryptoed25519.PublicKey(nil), k.PublicKey...)
}

func (k Keypair) SuiAddress() (string, error) {
	return keychain.AddressFromPublicKey(keychain.SchemeEd25519, k.PublicKey)
}

func (k Keypair) PublicKeyBase64() string {
	pub := k.PublicKeyBytes()
	payload := append([]byte{keychain.SchemeEd25519.AddressFlag()}, pub...)
	return base64.StdEncoding.EncodeToString(payload)
}

// Extracts the 32-byte seed component from the expanded
// Ed25519 private key.
func (k Keypair) SecretKeyBytes() []byte {
	seed := k.PrivateKey.Seed()
	out := make([]byte, len(seed))

	copy(out, seed)
	return out
}

// Generate produces a fresh Ed25519 keypair (based on crypto/rand).
func Generate() (*Keypair, error) {
	pub, priv, err := cryptoed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("ed25519: generate key: %w", err)
	}

	return &Keypair{
		PrivateKey: priv,
		PublicKey:  pub,
	}, nil
}

// Rebuilds the expanded Ed25519 keypair from a 32-byte
// SLIP-0010 seed.
func FromSecretKey(seed []byte) (*Keypair, error) {
	if len(seed) != keychain.PrivateKeySize() {
		return nil, keychain.ErrInvalidSecretLength
	}

	priv := cryptoed25519.NewKeyFromSeed(seed)
	pub := cryptoed25519.PublicKey(priv[32:])
	return &Keypair{
		PrivateKey: priv,
		PublicKey:  pub,
	}, nil
}

// Traverses the SLIP-0010 hardened derivation path starting from the
// supplied master seed. Ed25519 only supports hardened segments; any
// non-hardened index results in an error.
func Derive(seed []byte, path keychain.DerivationPath) (*Keypair, error) {
	if err := path.ValidateForScheme(keychain.SchemeEd25519); err != nil {
		return nil, err
	}

	key, chain := slip10MasterKey(seed)
	for _, segment := range path.Segments() {
		if !segment.Hardened {
			return nil, fmt.Errorf("ed25519: slip-0010 only supports hardened segments")
		}

		// SLIP-0010 hardened child derivation: HMAC-SHA512 with parent chain code
		// over [0x00 || parent_secret || child_index_with_high_bit]. The left half
		// becomes the child secret scalar; the right half becomes the next chain
		// code.
		data := make([]byte, 1+seedSize+4)
		data[0] = 0x00
		copy(data[1:], key)
		binary.BigEndian.PutUint32(data[1+seedSize:], segment.HardenedIndex())
		digest := keychain.HMACSHA512(chain, data)
		key = digest[:seedSize]
		chain = digest[seedSize:]
	}

	privateKey := cryptoed25519.NewKeyFromSeed(key)
	publicKey := cryptoed25519.PublicKey(privateKey[32:])
	return &Keypair{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
		ChainCode:  append([]byte{}, chain...),
		Path:       path,
	}, nil
}

// Initializes the root private key and chain code according to
// SLIP-0010's Ed25519 specification.
func slip10MasterKey(seed []byte) ([]byte, []byte) {
	digest := keychain.HMACSHA512(slip10Key, seed)
	return digest[:seedSize], digest[seedSize:]
}
