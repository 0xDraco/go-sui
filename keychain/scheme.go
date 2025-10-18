package keychain

import "fmt"

type Scheme uint8

const (
	SchemeEd25519 Scheme = iota
	SchemeSecp256k1
	SchemeSecp256r1
)

const (
	flagUnspecified byte = 0xff
)

func (s Scheme) AddressFlag() byte {
	switch s {
	case SchemeEd25519:
		return 0x00
	case SchemeSecp256k1:
		return 0x01
	case SchemeSecp256r1:
		return 0x02
	default:
		return flagUnspecified
	}
}

func (s Scheme) Purpose() uint32 {
	switch s {
	case SchemeEd25519:
		return 44
	case SchemeSecp256k1:
		return 54
	case SchemeSecp256r1:
		return 74
	default:
		return 0
	}
}

func SchemeFromFlag(flag byte) (Scheme, error) {
	switch flag {
	case 0x00:
		return SchemeEd25519, nil
	case 0x01:
		return SchemeSecp256k1, nil
	case 0x02:
		return SchemeSecp256r1, nil
	default:
		return 0, fmt.Errorf("unknown scheme flag 0x%02x", flag)
	}
}
