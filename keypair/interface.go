package keypair

import "github.com/0xdraco/go-sui/keychain"

type Keypair interface {
	Scheme() keychain.Scheme
	PublicKeyBytes() []byte
	PrivateKeyBytes() []byte
	SecretKeyBytes() []byte
	SuiAddress() (string, error)
	PublicKeyBase64() string
}
