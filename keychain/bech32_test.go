package keychain

import "testing"

func TestEncodeDecodePrivateKey(t *testing.T) {
	secret := make([]byte, PrivateKeySize())
	for i := range secret {
		secret[i] = byte(i)
	}

	encoded, err := EncodePrivateKey(SchemeSecp256k1, secret)
	if err != nil {
		t.Fatalf("encode: %v", err)
	}

	parsed, err := DecodePrivateKey(encoded)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}

	if parsed.Scheme != SchemeSecp256k1 {
		t.Fatalf("scheme mismatch: got %v want %v", parsed.Scheme, SchemeSecp256k1)
	}
	if len(parsed.SecretKey) != len(secret) {
		t.Fatalf("secret length mismatch: got %d want %d", len(parsed.SecretKey), len(secret))
	}
	for i := range secret {
		if parsed.SecretKey[i] != secret[i] {
			t.Fatalf("secret differs at %d", i)
		}
	}
}

func TestDecodeRejectsBadHrp(t *testing.T) {
	_, err := DecodePrivateKey("invalid1deadbeef")
	if err == nil {
		t.Fatalf("expected error for bad encoding")
	}
}
