package pb

import (
	"crypto/ecdsa"
	"fmt"
)

type KeyT interface {
	Pub() (*ecdsa.PublicKey, error)
	SetImplProto(p *KeyImpl) error
}

func KeyToProto(name string, k KeyT) (*KeyImpl, error) {
	pub, err := k.Pub()
	if err != nil {
		return nil, err
	}

	impl := &KeyImpl{
		Name:      name,
		PublicKey: PubToProto(pub),
	}
	if err := k.SetImplProto(impl); err != nil {
		return nil, err
	}

	return impl, nil
}

func KeyImplTypeToString(it isKeyImpl_Impl) string {
	switch it.(type) {
	case *KeyImpl_DebugRawKey:
		return "debug_raw"
	case *KeyImpl_EmergencyBackupKey:
		return "emergency"
	case *KeyImpl_WebauthnWrappedTpm:
		return "webauthn_wrapped_tpm"
	default:
		return "<unknown>"
	}
}

func KeyImplToString(it isKeyImpl_Impl) string {
	switch v := it.(type) {
	case *KeyImpl_DebugRawKey:
		return "DebugRawKey"
	case *KeyImpl_EmergencyBackupKey:
		return fmt.Sprintf("Type: EmergencyBackupKey, Hint: %s", v.EmergencyBackupKey.Hint)
	case *KeyImpl_WebauthnWrappedTpm:
		return fmt.Sprintf("Type: WebauthnWrappedTpm, TpmKeyHandle: 0x%x, ReflectorUrl: %s",
			v.WebauthnWrappedTpm.TpmKeyHandle,
			v.WebauthnWrappedTpm.ReflectorUrl)
	default:
		return "<unknown>"
	}
}

func KeyImplSummary(k *KeyImpl) string {
	return fmt.Sprintf("%s {Name: %s, %s}",
		PubKeyPinString(ProtoToPub(k.PublicKey)),
		k.Name,
		KeyImplToString(k.Impl),
	)
}
