package intent

import (
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/fardream/go-bcs/bcs"
	"golang.org/x/crypto/blake2b"
)

type IntentVersion uint8

const (
	IntentVersionV0 IntentVersion = 0
)

var (
	errInvalidIntentLength = errors.New("intent: invalid serialized intent length")
	errInvalidIntentScope  = errors.New("intent: invalid scope byte")
	errInvalidIntentAppID  = errors.New("intent: invalid app id byte")
	errInvalidIntentVers   = errors.New("intent: invalid version byte")
)

func (v IntentVersion) Validate() error {
	if v != IntentVersionV0 {
		return errInvalidIntentVers
	}
	return nil
}

type AppID uint8

const (
	AppIDSui AppID = 0
)

func (a AppID) Validate() error {
	if a != AppIDSui {
		return errInvalidIntentAppID
	}
	return nil
}

type IntentScope uint8

const (
	IntentScopeTransactionData         IntentScope = 0
	IntentScopeTransactionEffects      IntentScope = 1
	IntentScopeCheckpointSummary       IntentScope = 2
	IntentScopePersonalMessage         IntentScope = 3
	IntentScopeSenderSignedTransaction IntentScope = 4
	IntentScopeProofOfPossession       IntentScope = 5
)

func (s IntentScope) Validate() error {
	switch s {
	case IntentScopeTransactionData,
		IntentScopeTransactionEffects,
		IntentScopeCheckpointSummary,
		IntentScopePersonalMessage,
		IntentScopeSenderSignedTransaction,
		IntentScopeProofOfPossession:
		return nil
	default:
		return errInvalidIntentScope
	}
}

type Intent struct {
	Scope   IntentScope
	Version IntentVersion
	AppID   AppID
}

func DefaultIntent() Intent {
	return Intent{
		Scope:   IntentScopeTransactionData,
		Version: IntentVersionV0,
		AppID:   AppIDSui,
	}
}

func (i Intent) WithAppID(appID AppID) Intent {
	i.AppID = appID
	return i
}

func (i Intent) WithScope(scope IntentScope) Intent {
	i.Scope = scope
	return i
}

func (i Intent) Bytes() [3]byte {
	return [3]byte{byte(i.Scope), byte(i.Version), byte(i.AppID)}
}

func (i Intent) Validate() error {
	if err := i.Scope.Validate(); err != nil {
		return err
	}
	if err := i.Version.Validate(); err != nil {
		return err
	}
	return i.AppID.Validate()
}

func ParseIntent(hexEncoded string) (Intent, error) {
	raw, err := hex.DecodeString(hexEncoded)
	if err != nil {
		return Intent{}, fmt.Errorf("intent: decode hex: %w", err)
	}
	intent, err := IntentFromBytes(raw)
	if err != nil {
		return Intent{}, err
	}
	return intent, nil
}

func IntentFromBytes(raw []byte) (Intent, error) {
	if len(raw) != 3 {
		return Intent{}, errInvalidIntentLength
	}
	intent := Intent{
		Scope:   IntentScope(raw[0]),
		Version: IntentVersion(raw[1]),
		AppID:   AppID(raw[2]),
	}
	if err := intent.Validate(); err != nil {
		return Intent{}, err
	}
	return intent, nil
}

type IntentMessage[T any] struct {
	Intent Intent
	Value  T
}

func NewIntentMessage[T any](intent Intent, value T) IntentMessage[T] {
	return IntentMessage[T]{
		Intent: intent,
		Value:  value,
	}
}

func (m IntentMessage[T]) MarshalBCS() ([]byte, error) {
	if err := m.Intent.Validate(); err != nil {
		return nil, err
	}
	valueBytes, err := bcs.Marshal(m.Value)
	if err != nil {
		return nil, fmt.Errorf("intent: marshal value: %w", err)
	}
	intentBytes := m.Intent.Bytes()
	encoded := make([]byte, 0, len(valueBytes)+len(intentBytes))
	encoded = append(encoded, intentBytes[:]...)
	encoded = append(encoded, valueBytes...)
	return encoded, nil
}

func HashIntentMessage[T any](message IntentMessage[T]) ([32]byte, error) {
	serialized, err := message.MarshalBCS()
	if err != nil {
		return [32]byte{}, err
	}
	return blake2b.Sum256(serialized), nil
}

type PersonalMessage struct {
	Message []byte
}
