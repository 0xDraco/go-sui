package keypair

import "github.com/0xdraco/sui-go-sdk/keychain"

type Keypair interface {
	Scheme() keychain.Scheme
	PublicKeyBytes() []byte
	PrivateKeyBytes() []byte
	SecretKeyBytes() []byte
	SuiAddress() (string, error)
	PublicKeyBase64() string
	SignPersonalMessage(message []byte) ([]byte, error)
	VerifyPersonalMessage(message []byte, signature []byte) error
}
