// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        (unknown)
// source: pb/h132.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type P256PublicKey struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	X []byte `protobuf:"bytes,1,opt,name=x,proto3" json:"x,omitempty"`
	Y []byte `protobuf:"bytes,2,opt,name=y,proto3" json:"y,omitempty"`
}

func (x *P256PublicKey) Reset() {
	*x = P256PublicKey{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_h132_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *P256PublicKey) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*P256PublicKey) ProtoMessage() {}

func (x *P256PublicKey) ProtoReflect() protoreflect.Message {
	mi := &file_pb_h132_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use P256PublicKey.ProtoReflect.Descriptor instead.
func (*P256PublicKey) Descriptor() ([]byte, []int) {
	return file_pb_h132_proto_rawDescGZIP(), []int{0}
}

func (x *P256PublicKey) GetX() []byte {
	if x != nil {
		return x.X
	}
	return nil
}

func (x *P256PublicKey) GetY() []byte {
	if x != nil {
		return x.Y
	}
	return nil
}

type P256PrivateKey struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	D []byte `protobuf:"bytes,1,opt,name=d,proto3" json:"d,omitempty"`
}

func (x *P256PrivateKey) Reset() {
	*x = P256PrivateKey{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_h132_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *P256PrivateKey) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*P256PrivateKey) ProtoMessage() {}

func (x *P256PrivateKey) ProtoReflect() protoreflect.Message {
	mi := &file_pb_h132_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use P256PrivateKey.ProtoReflect.Descriptor instead.
func (*P256PrivateKey) Descriptor() ([]byte, []int) {
	return file_pb_h132_proto_rawDescGZIP(), []int{1}
}

func (x *P256PrivateKey) GetD() []byte {
	if x != nil {
		return x.D
	}
	return nil
}

type EncryptedSymmetricKey struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The public key of the recipient used to encrypt the symmetric key
	RecipientPublicKey *P256PublicKey `protobuf:"bytes,1,opt,name=recipient_public_key,json=recipientPublicKey,proto3" json:"recipient_public_key,omitempty"`
	// The ephemeral public key of the sender used to encrypt the symmmetric key
	SenderEphemeralPublicKey *P256PublicKey `protobuf:"bytes,2,opt,name=sender_ephemeral_public_key,json=senderEphemeralPublicKey,proto3" json:"sender_ephemeral_public_key,omitempty"`
	// The signature of the ephemeral public key by the sender key
	SenderEphemeralPublicKeySign []byte `protobuf:"bytes,3,opt,name=sender_ephemeral_public_key_sign,json=senderEphemeralPublicKeySign,proto3" json:"sender_ephemeral_public_key_sign,omitempty"`
	// The symmetric key encrypted with the recipient's public key
	EncryptedSymmetricKey []byte `protobuf:"bytes,4,opt,name=encrypted_symmetric_key,json=encryptedSymmetricKey,proto3" json:"encrypted_symmetric_key,omitempty"`
	// The nonce used to encrypt the symmetric key
	Nonce []byte `protobuf:"bytes,5,opt,name=nonce,proto3" json:"nonce,omitempty"`
}

func (x *EncryptedSymmetricKey) Reset() {
	*x = EncryptedSymmetricKey{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_h132_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EncryptedSymmetricKey) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EncryptedSymmetricKey) ProtoMessage() {}

func (x *EncryptedSymmetricKey) ProtoReflect() protoreflect.Message {
	mi := &file_pb_h132_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EncryptedSymmetricKey.ProtoReflect.Descriptor instead.
func (*EncryptedSymmetricKey) Descriptor() ([]byte, []int) {
	return file_pb_h132_proto_rawDescGZIP(), []int{2}
}

func (x *EncryptedSymmetricKey) GetRecipientPublicKey() *P256PublicKey {
	if x != nil {
		return x.RecipientPublicKey
	}
	return nil
}

func (x *EncryptedSymmetricKey) GetSenderEphemeralPublicKey() *P256PublicKey {
	if x != nil {
		return x.SenderEphemeralPublicKey
	}
	return nil
}

func (x *EncryptedSymmetricKey) GetSenderEphemeralPublicKeySign() []byte {
	if x != nil {
		return x.SenderEphemeralPublicKeySign
	}
	return nil
}

func (x *EncryptedSymmetricKey) GetEncryptedSymmetricKey() []byte {
	if x != nil {
		return x.EncryptedSymmetricKey
	}
	return nil
}

func (x *EncryptedSymmetricKey) GetNonce() []byte {
	if x != nil {
		return x.Nonce
	}
	return nil
}

type Letter struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ciphertext      []byte                   `protobuf:"bytes,1,opt,name=ciphertext,proto3" json:"ciphertext,omitempty"`
	Nonce           []byte                   `protobuf:"bytes,2,opt,name=nonce,proto3" json:"nonce,omitempty"`
	RecipientKeys   []*EncryptedSymmetricKey `protobuf:"bytes,3,rep,name=recipient_keys,json=recipientKeys,proto3" json:"recipient_keys,omitempty"`
	SenderPublicKey *P256PublicKey           `protobuf:"bytes,4,opt,name=sender_public_key,json=senderPublicKey,proto3" json:"sender_public_key,omitempty"`
}

func (x *Letter) Reset() {
	*x = Letter{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_h132_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Letter) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Letter) ProtoMessage() {}

func (x *Letter) ProtoReflect() protoreflect.Message {
	mi := &file_pb_h132_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Letter.ProtoReflect.Descriptor instead.
func (*Letter) Descriptor() ([]byte, []int) {
	return file_pb_h132_proto_rawDescGZIP(), []int{3}
}

func (x *Letter) GetCiphertext() []byte {
	if x != nil {
		return x.Ciphertext
	}
	return nil
}

func (x *Letter) GetNonce() []byte {
	if x != nil {
		return x.Nonce
	}
	return nil
}

func (x *Letter) GetRecipientKeys() []*EncryptedSymmetricKey {
	if x != nil {
		return x.RecipientKeys
	}
	return nil
}

func (x *Letter) GetSenderPublicKey() *P256PublicKey {
	if x != nil {
		return x.SenderPublicKey
	}
	return nil
}

type Envelope struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	LetterProto []byte `protobuf:"bytes,1,opt,name=letter_proto,json=letterProto,proto3" json:"letter_proto,omitempty"`
	Signature   []byte `protobuf:"bytes,2,opt,name=signature,proto3" json:"signature,omitempty"`
}

func (x *Envelope) Reset() {
	*x = Envelope{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_h132_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Envelope) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Envelope) ProtoMessage() {}

func (x *Envelope) ProtoReflect() protoreflect.Message {
	mi := &file_pb_h132_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Envelope.ProtoReflect.Descriptor instead.
func (*Envelope) Descriptor() ([]byte, []int) {
	return file_pb_h132_proto_rawDescGZIP(), []int{4}
}

func (x *Envelope) GetLetterProto() []byte {
	if x != nil {
		return x.LetterProto
	}
	return nil
}

func (x *Envelope) GetSignature() []byte {
	if x != nil {
		return x.Signature
	}
	return nil
}

// The emergency backup key, used to recover envelopes in case of a TPM failure.
// The private key is encoded as a bip39 mneumonic, written down on a piece of
// paper, and stored in a safe.
type EmergencyBackupKey struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Hint about the physical location of the key mneumonic written down.
	Hint string `protobuf:"bytes,1,opt,name=hint,proto3" json:"hint,omitempty"`
}

func (x *EmergencyBackupKey) Reset() {
	*x = EmergencyBackupKey{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_h132_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EmergencyBackupKey) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EmergencyBackupKey) ProtoMessage() {}

func (x *EmergencyBackupKey) ProtoReflect() protoreflect.Message {
	mi := &file_pb_h132_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EmergencyBackupKey.ProtoReflect.Descriptor instead.
func (*EmergencyBackupKey) Descriptor() ([]byte, []int) {
	return file_pb_h132_proto_rawDescGZIP(), []int{5}
}

func (x *EmergencyBackupKey) GetHint() string {
	if x != nil {
		return x.Hint
	}
	return ""
}

// A TPM key, with a passphrase derived from WebAuthn PRF.
type WebAuthnWrappedTPMKey struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The URL used to trigger WebAuthn interaction.
	ReflectorUrl string `protobuf:"bytes,1,opt,name=reflector_url,json=reflectorUrl,proto3" json:"reflector_url,omitempty"`
	// The key handle used to identify the key in the TPM.
	TpmKeyHandle uint32 `protobuf:"varint,2,opt,name=tpm_key_handle,json=tpmKeyHandle,proto3" json:"tpm_key_handle,omitempty"`
	// The salt given to retrieve PRF output from WebAuthn. Must be 32 bytes.
	PrfSalt []byte `protobuf:"bytes,3,opt,name=prf_salt,json=prfSalt,proto3" json:"prf_salt,omitempty"`
	// The salt given to HKDF to derive the TPM passphrase from the PRF output.
	HkdfSalt []byte `protobuf:"bytes,4,opt,name=hkdf_salt,json=hkdfSalt,proto3" json:"hkdf_salt,omitempty"`
	// The username used to derive the WebAuthn credential {id, name, displayName}.
	WebauthnUsername string `protobuf:"bytes,5,opt,name=webauthn_username,json=webauthnUsername,proto3" json:"webauthn_username,omitempty"`
	// The credential JSON returned by WebAuthn.
	WebauthnCredentialJson []byte `protobuf:"bytes,6,opt,name=webauthn_credential_json,json=webauthnCredentialJson,proto3" json:"webauthn_credential_json,omitempty"`
}

func (x *WebAuthnWrappedTPMKey) Reset() {
	*x = WebAuthnWrappedTPMKey{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_h132_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WebAuthnWrappedTPMKey) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WebAuthnWrappedTPMKey) ProtoMessage() {}

func (x *WebAuthnWrappedTPMKey) ProtoReflect() protoreflect.Message {
	mi := &file_pb_h132_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WebAuthnWrappedTPMKey.ProtoReflect.Descriptor instead.
func (*WebAuthnWrappedTPMKey) Descriptor() ([]byte, []int) {
	return file_pb_h132_proto_rawDescGZIP(), []int{6}
}

func (x *WebAuthnWrappedTPMKey) GetReflectorUrl() string {
	if x != nil {
		return x.ReflectorUrl
	}
	return ""
}

func (x *WebAuthnWrappedTPMKey) GetTpmKeyHandle() uint32 {
	if x != nil {
		return x.TpmKeyHandle
	}
	return 0
}

func (x *WebAuthnWrappedTPMKey) GetPrfSalt() []byte {
	if x != nil {
		return x.PrfSalt
	}
	return nil
}

func (x *WebAuthnWrappedTPMKey) GetHkdfSalt() []byte {
	if x != nil {
		return x.HkdfSalt
	}
	return nil
}

func (x *WebAuthnWrappedTPMKey) GetWebauthnUsername() string {
	if x != nil {
		return x.WebauthnUsername
	}
	return ""
}

func (x *WebAuthnWrappedTPMKey) GetWebauthnCredentialJson() []byte {
	if x != nil {
		return x.WebauthnCredentialJson
	}
	return nil
}

type KeyImpl struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The name of the key, specified by the user when starting a h132 session.
	Name      string         `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	PublicKey *P256PublicKey `protobuf:"bytes,2,opt,name=public_key,json=publicKey,proto3" json:"public_key,omitempty"`
	// Types that are assignable to Impl:
	//
	//	*KeyImpl_DebugRawKey
	//	*KeyImpl_EmergencyBackupKey
	//	*KeyImpl_WebauthnWrappedTpm
	Impl isKeyImpl_Impl `protobuf_oneof:"impl"`
}

func (x *KeyImpl) Reset() {
	*x = KeyImpl{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_h132_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *KeyImpl) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*KeyImpl) ProtoMessage() {}

func (x *KeyImpl) ProtoReflect() protoreflect.Message {
	mi := &file_pb_h132_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use KeyImpl.ProtoReflect.Descriptor instead.
func (*KeyImpl) Descriptor() ([]byte, []int) {
	return file_pb_h132_proto_rawDescGZIP(), []int{7}
}

func (x *KeyImpl) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *KeyImpl) GetPublicKey() *P256PublicKey {
	if x != nil {
		return x.PublicKey
	}
	return nil
}

func (m *KeyImpl) GetImpl() isKeyImpl_Impl {
	if m != nil {
		return m.Impl
	}
	return nil
}

func (x *KeyImpl) GetDebugRawKey() *P256PrivateKey {
	if x, ok := x.GetImpl().(*KeyImpl_DebugRawKey); ok {
		return x.DebugRawKey
	}
	return nil
}

func (x *KeyImpl) GetEmergencyBackupKey() *EmergencyBackupKey {
	if x, ok := x.GetImpl().(*KeyImpl_EmergencyBackupKey); ok {
		return x.EmergencyBackupKey
	}
	return nil
}

func (x *KeyImpl) GetWebauthnWrappedTpm() *WebAuthnWrappedTPMKey {
	if x, ok := x.GetImpl().(*KeyImpl_WebauthnWrappedTpm); ok {
		return x.WebauthnWrappedTpm
	}
	return nil
}

type isKeyImpl_Impl interface {
	isKeyImpl_Impl()
}

type KeyImpl_DebugRawKey struct {
	// Private key presented in the clear. Only used for debugging.
	// h132 refuses to operate unless the debug flag is set.
	DebugRawKey *P256PrivateKey `protobuf:"bytes,10,opt,name=debug_raw_key,json=debugRawKey,proto3,oneof"`
}

type KeyImpl_EmergencyBackupKey struct {
	EmergencyBackupKey *EmergencyBackupKey `protobuf:"bytes,11,opt,name=emergency_backup_key,json=emergencyBackupKey,proto3,oneof"`
}

type KeyImpl_WebauthnWrappedTpm struct {
	WebauthnWrappedTpm *WebAuthnWrappedTPMKey `protobuf:"bytes,12,opt,name=webauthn_wrapped_tpm,json=webauthnWrappedTpm,proto3,oneof"`
}

func (*KeyImpl_DebugRawKey) isKeyImpl_Impl() {}

func (*KeyImpl_EmergencyBackupKey) isKeyImpl_Impl() {}

func (*KeyImpl_WebauthnWrappedTpm) isKeyImpl_Impl() {}

// LetterWritingSet is a h132 configuration that is stored as
// "h132_letter_writing_set.binpb" file in the directory with the envelopes.
type LetterWritingSet struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The name of the letter writing set - presented to the user when operating
	// under the letter writing set.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// The keys that can decrypt the envelopes.
	Keys []*KeyImpl `protobuf:"bytes,2,rep,name=keys,proto3" json:"keys,omitempty"`
	// Shell command before editing a envelope. The envelope path is passed as an argument.
	EditPreconditionShellCommand string `protobuf:"bytes,10,opt,name=edit_precondition_shell_command,json=editPreconditionShellCommand,proto3" json:"edit_precondition_shell_command,omitempty"`
	// Shell command after new envelope is created or edited. The envelope path is passed as an argument.
	EditPostconditionShellCommand string `protobuf:"bytes,11,opt,name=edit_postcondition_shell_command,json=editPostconditionShellCommand,proto3" json:"edit_postcondition_shell_command,omitempty"`
}

func (x *LetterWritingSet) Reset() {
	*x = LetterWritingSet{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_h132_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LetterWritingSet) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LetterWritingSet) ProtoMessage() {}

func (x *LetterWritingSet) ProtoReflect() protoreflect.Message {
	mi := &file_pb_h132_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LetterWritingSet.ProtoReflect.Descriptor instead.
func (*LetterWritingSet) Descriptor() ([]byte, []int) {
	return file_pb_h132_proto_rawDescGZIP(), []int{8}
}

func (x *LetterWritingSet) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *LetterWritingSet) GetKeys() []*KeyImpl {
	if x != nil {
		return x.Keys
	}
	return nil
}

func (x *LetterWritingSet) GetEditPreconditionShellCommand() string {
	if x != nil {
		return x.EditPreconditionShellCommand
	}
	return ""
}

func (x *LetterWritingSet) GetEditPostconditionShellCommand() string {
	if x != nil {
		return x.EditPostconditionShellCommand
	}
	return ""
}

var File_pb_h132_proto protoreflect.FileDescriptor

var file_pb_h132_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x70, 0x62, 0x2f, 0x68, 0x31, 0x33, 0x32, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x02, 0x70, 0x62, 0x22, 0x2b, 0x0a, 0x0d, 0x50, 0x32, 0x35, 0x36, 0x50, 0x75, 0x62, 0x6c, 0x69,
	0x63, 0x4b, 0x65, 0x79, 0x12, 0x0c, 0x0a, 0x01, 0x78, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x01, 0x78, 0x12, 0x0c, 0x0a, 0x01, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x01, 0x79,
	0x22, 0x1e, 0x0a, 0x0e, 0x50, 0x32, 0x35, 0x36, 0x50, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x4b,
	0x65, 0x79, 0x12, 0x0c, 0x0a, 0x01, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x01, 0x64,
	0x22, 0xc4, 0x02, 0x0a, 0x15, 0x45, 0x6e, 0x63, 0x72, 0x79, 0x70, 0x74, 0x65, 0x64, 0x53, 0x79,
	0x6d, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x4b, 0x65, 0x79, 0x12, 0x43, 0x0a, 0x14, 0x72, 0x65,
	0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x5f, 0x6b,
	0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x70, 0x62, 0x2e, 0x50, 0x32,
	0x35, 0x36, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x4b, 0x65, 0x79, 0x52, 0x12, 0x72, 0x65, 0x63,
	0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x4b, 0x65, 0x79, 0x12,
	0x50, 0x0a, 0x1b, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x5f, 0x65, 0x70, 0x68, 0x65, 0x6d, 0x65,
	0x72, 0x61, 0x6c, 0x5f, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x70, 0x62, 0x2e, 0x50, 0x32, 0x35, 0x36, 0x50, 0x75,
	0x62, 0x6c, 0x69, 0x63, 0x4b, 0x65, 0x79, 0x52, 0x18, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x45,
	0x70, 0x68, 0x65, 0x6d, 0x65, 0x72, 0x61, 0x6c, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x4b, 0x65,
	0x79, 0x12, 0x46, 0x0a, 0x20, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x5f, 0x65, 0x70, 0x68, 0x65,
	0x6d, 0x65, 0x72, 0x61, 0x6c, 0x5f, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x5f, 0x6b, 0x65, 0x79,
	0x5f, 0x73, 0x69, 0x67, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x1c, 0x73, 0x65, 0x6e,
	0x64, 0x65, 0x72, 0x45, 0x70, 0x68, 0x65, 0x6d, 0x65, 0x72, 0x61, 0x6c, 0x50, 0x75, 0x62, 0x6c,
	0x69, 0x63, 0x4b, 0x65, 0x79, 0x53, 0x69, 0x67, 0x6e, 0x12, 0x36, 0x0a, 0x17, 0x65, 0x6e, 0x63,
	0x72, 0x79, 0x70, 0x74, 0x65, 0x64, 0x5f, 0x73, 0x79, 0x6d, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63,
	0x5f, 0x6b, 0x65, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x15, 0x65, 0x6e, 0x63, 0x72,
	0x79, 0x70, 0x74, 0x65, 0x64, 0x53, 0x79, 0x6d, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x4b, 0x65,
	0x79, 0x12, 0x14, 0x0a, 0x05, 0x6e, 0x6f, 0x6e, 0x63, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0c,
	0x52, 0x05, 0x6e, 0x6f, 0x6e, 0x63, 0x65, 0x22, 0xbf, 0x01, 0x0a, 0x06, 0x4c, 0x65, 0x74, 0x74,
	0x65, 0x72, 0x12, 0x1e, 0x0a, 0x0a, 0x63, 0x69, 0x70, 0x68, 0x65, 0x72, 0x74, 0x65, 0x78, 0x74,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0a, 0x63, 0x69, 0x70, 0x68, 0x65, 0x72, 0x74, 0x65,
	0x78, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x6e, 0x6f, 0x6e, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x05, 0x6e, 0x6f, 0x6e, 0x63, 0x65, 0x12, 0x40, 0x0a, 0x0e, 0x72, 0x65, 0x63, 0x69,
	0x70, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x6b, 0x65, 0x79, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x19, 0x2e, 0x70, 0x62, 0x2e, 0x45, 0x6e, 0x63, 0x72, 0x79, 0x70, 0x74, 0x65, 0x64, 0x53,
	0x79, 0x6d, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x4b, 0x65, 0x79, 0x52, 0x0d, 0x72, 0x65, 0x63,
	0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x4b, 0x65, 0x79, 0x73, 0x12, 0x3d, 0x0a, 0x11, 0x73, 0x65,
	0x6e, 0x64, 0x65, 0x72, 0x5f, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x5f, 0x6b, 0x65, 0x79, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x70, 0x62, 0x2e, 0x50, 0x32, 0x35, 0x36, 0x50,
	0x75, 0x62, 0x6c, 0x69, 0x63, 0x4b, 0x65, 0x79, 0x52, 0x0f, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72,
	0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x4b, 0x65, 0x79, 0x22, 0x4b, 0x0a, 0x08, 0x45, 0x6e, 0x76,
	0x65, 0x6c, 0x6f, 0x70, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x6c, 0x65, 0x74, 0x74, 0x65, 0x72, 0x5f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0b, 0x6c, 0x65, 0x74,
	0x74, 0x65, 0x72, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1c, 0x0a, 0x09, 0x73, 0x69, 0x67, 0x6e,
	0x61, 0x74, 0x75, 0x72, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x73, 0x69, 0x67,
	0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x22, 0x28, 0x0a, 0x12, 0x45, 0x6d, 0x65, 0x72, 0x67, 0x65,
	0x6e, 0x63, 0x79, 0x42, 0x61, 0x63, 0x6b, 0x75, 0x70, 0x4b, 0x65, 0x79, 0x12, 0x12, 0x0a, 0x04,
	0x68, 0x69, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x68, 0x69, 0x6e, 0x74,
	0x22, 0x81, 0x02, 0x0a, 0x15, 0x57, 0x65, 0x62, 0x41, 0x75, 0x74, 0x68, 0x6e, 0x57, 0x72, 0x61,
	0x70, 0x70, 0x65, 0x64, 0x54, 0x50, 0x4d, 0x4b, 0x65, 0x79, 0x12, 0x23, 0x0a, 0x0d, 0x72, 0x65,
	0x66, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x5f, 0x75, 0x72, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0c, 0x72, 0x65, 0x66, 0x6c, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x55, 0x72, 0x6c, 0x12,
	0x24, 0x0a, 0x0e, 0x74, 0x70, 0x6d, 0x5f, 0x6b, 0x65, 0x79, 0x5f, 0x68, 0x61, 0x6e, 0x64, 0x6c,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0c, 0x74, 0x70, 0x6d, 0x4b, 0x65, 0x79, 0x48,
	0x61, 0x6e, 0x64, 0x6c, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x70, 0x72, 0x66, 0x5f, 0x73, 0x61, 0x6c,
	0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x70, 0x72, 0x66, 0x53, 0x61, 0x6c, 0x74,
	0x12, 0x1b, 0x0a, 0x09, 0x68, 0x6b, 0x64, 0x66, 0x5f, 0x73, 0x61, 0x6c, 0x74, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x08, 0x68, 0x6b, 0x64, 0x66, 0x53, 0x61, 0x6c, 0x74, 0x12, 0x2b, 0x0a,
	0x11, 0x77, 0x65, 0x62, 0x61, 0x75, 0x74, 0x68, 0x6e, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x10, 0x77, 0x65, 0x62, 0x61, 0x75, 0x74,
	0x68, 0x6e, 0x55, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x38, 0x0a, 0x18, 0x77, 0x65,
	0x62, 0x61, 0x75, 0x74, 0x68, 0x6e, 0x5f, 0x63, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x61,
	0x6c, 0x5f, 0x6a, 0x73, 0x6f, 0x6e, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x16, 0x77, 0x65,
	0x62, 0x61, 0x75, 0x74, 0x68, 0x6e, 0x43, 0x72, 0x65, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x61, 0x6c,
	0x4a, 0x73, 0x6f, 0x6e, 0x22, 0xac, 0x02, 0x0a, 0x07, 0x4b, 0x65, 0x79, 0x49, 0x6d, 0x70, 0x6c,
	0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x12, 0x30, 0x0a, 0x0a, 0x70, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x5f, 0x6b,
	0x65, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x70, 0x62, 0x2e, 0x50, 0x32,
	0x35, 0x36, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x4b, 0x65, 0x79, 0x52, 0x09, 0x70, 0x75, 0x62,
	0x6c, 0x69, 0x63, 0x4b, 0x65, 0x79, 0x12, 0x38, 0x0a, 0x0d, 0x64, 0x65, 0x62, 0x75, 0x67, 0x5f,
	0x72, 0x61, 0x77, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e,
	0x70, 0x62, 0x2e, 0x50, 0x32, 0x35, 0x36, 0x50, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x4b, 0x65,
	0x79, 0x48, 0x00, 0x52, 0x0b, 0x64, 0x65, 0x62, 0x75, 0x67, 0x52, 0x61, 0x77, 0x4b, 0x65, 0x79,
	0x12, 0x4a, 0x0a, 0x14, 0x65, 0x6d, 0x65, 0x72, 0x67, 0x65, 0x6e, 0x63, 0x79, 0x5f, 0x62, 0x61,
	0x63, 0x6b, 0x75, 0x70, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16,
	0x2e, 0x70, 0x62, 0x2e, 0x45, 0x6d, 0x65, 0x72, 0x67, 0x65, 0x6e, 0x63, 0x79, 0x42, 0x61, 0x63,
	0x6b, 0x75, 0x70, 0x4b, 0x65, 0x79, 0x48, 0x00, 0x52, 0x12, 0x65, 0x6d, 0x65, 0x72, 0x67, 0x65,
	0x6e, 0x63, 0x79, 0x42, 0x61, 0x63, 0x6b, 0x75, 0x70, 0x4b, 0x65, 0x79, 0x12, 0x4d, 0x0a, 0x14,
	0x77, 0x65, 0x62, 0x61, 0x75, 0x74, 0x68, 0x6e, 0x5f, 0x77, 0x72, 0x61, 0x70, 0x70, 0x65, 0x64,
	0x5f, 0x74, 0x70, 0x6d, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x70, 0x62, 0x2e,
	0x57, 0x65, 0x62, 0x41, 0x75, 0x74, 0x68, 0x6e, 0x57, 0x72, 0x61, 0x70, 0x70, 0x65, 0x64, 0x54,
	0x50, 0x4d, 0x4b, 0x65, 0x79, 0x48, 0x00, 0x52, 0x12, 0x77, 0x65, 0x62, 0x61, 0x75, 0x74, 0x68,
	0x6e, 0x57, 0x72, 0x61, 0x70, 0x70, 0x65, 0x64, 0x54, 0x70, 0x6d, 0x42, 0x06, 0x0a, 0x04, 0x69,
	0x6d, 0x70, 0x6c, 0x22, 0xd7, 0x01, 0x0a, 0x10, 0x4c, 0x65, 0x74, 0x74, 0x65, 0x72, 0x57, 0x72,
	0x69, 0x74, 0x69, 0x6e, 0x67, 0x53, 0x65, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1f, 0x0a, 0x04,
	0x6b, 0x65, 0x79, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x70, 0x62, 0x2e,
	0x4b, 0x65, 0x79, 0x49, 0x6d, 0x70, 0x6c, 0x52, 0x04, 0x6b, 0x65, 0x79, 0x73, 0x12, 0x45, 0x0a,
	0x1f, 0x65, 0x64, 0x69, 0x74, 0x5f, 0x70, 0x72, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x69, 0x74, 0x69,
	0x6f, 0x6e, 0x5f, 0x73, 0x68, 0x65, 0x6c, 0x6c, 0x5f, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64,
	0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x1c, 0x65, 0x64, 0x69, 0x74, 0x50, 0x72, 0x65, 0x63,
	0x6f, 0x6e, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x68, 0x65, 0x6c, 0x6c, 0x43, 0x6f, 0x6d,
	0x6d, 0x61, 0x6e, 0x64, 0x12, 0x47, 0x0a, 0x20, 0x65, 0x64, 0x69, 0x74, 0x5f, 0x70, 0x6f, 0x73,
	0x74, 0x63, 0x6f, 0x6e, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x73, 0x68, 0x65, 0x6c, 0x6c,
	0x5f, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x52, 0x1d,
	0x65, 0x64, 0x69, 0x74, 0x50, 0x6f, 0x73, 0x74, 0x63, 0x6f, 0x6e, 0x64, 0x69, 0x74, 0x69, 0x6f,
	0x6e, 0x53, 0x68, 0x65, 0x6c, 0x6c, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x42, 0x07, 0x5a,
	0x05, 0x2e, 0x2f, 0x3b, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pb_h132_proto_rawDescOnce sync.Once
	file_pb_h132_proto_rawDescData = file_pb_h132_proto_rawDesc
)

func file_pb_h132_proto_rawDescGZIP() []byte {
	file_pb_h132_proto_rawDescOnce.Do(func() {
		file_pb_h132_proto_rawDescData = protoimpl.X.CompressGZIP(file_pb_h132_proto_rawDescData)
	})
	return file_pb_h132_proto_rawDescData
}

var file_pb_h132_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_pb_h132_proto_goTypes = []any{
	(*P256PublicKey)(nil),         // 0: pb.P256PublicKey
	(*P256PrivateKey)(nil),        // 1: pb.P256PrivateKey
	(*EncryptedSymmetricKey)(nil), // 2: pb.EncryptedSymmetricKey
	(*Letter)(nil),                // 3: pb.Letter
	(*Envelope)(nil),              // 4: pb.Envelope
	(*EmergencyBackupKey)(nil),    // 5: pb.EmergencyBackupKey
	(*WebAuthnWrappedTPMKey)(nil), // 6: pb.WebAuthnWrappedTPMKey
	(*KeyImpl)(nil),               // 7: pb.KeyImpl
	(*LetterWritingSet)(nil),      // 8: pb.LetterWritingSet
}
var file_pb_h132_proto_depIdxs = []int32{
	0, // 0: pb.EncryptedSymmetricKey.recipient_public_key:type_name -> pb.P256PublicKey
	0, // 1: pb.EncryptedSymmetricKey.sender_ephemeral_public_key:type_name -> pb.P256PublicKey
	2, // 2: pb.Letter.recipient_keys:type_name -> pb.EncryptedSymmetricKey
	0, // 3: pb.Letter.sender_public_key:type_name -> pb.P256PublicKey
	0, // 4: pb.KeyImpl.public_key:type_name -> pb.P256PublicKey
	1, // 5: pb.KeyImpl.debug_raw_key:type_name -> pb.P256PrivateKey
	5, // 6: pb.KeyImpl.emergency_backup_key:type_name -> pb.EmergencyBackupKey
	6, // 7: pb.KeyImpl.webauthn_wrapped_tpm:type_name -> pb.WebAuthnWrappedTPMKey
	7, // 8: pb.LetterWritingSet.keys:type_name -> pb.KeyImpl
	9, // [9:9] is the sub-list for method output_type
	9, // [9:9] is the sub-list for method input_type
	9, // [9:9] is the sub-list for extension type_name
	9, // [9:9] is the sub-list for extension extendee
	0, // [0:9] is the sub-list for field type_name
}

func init() { file_pb_h132_proto_init() }
func file_pb_h132_proto_init() {
	if File_pb_h132_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pb_h132_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*P256PublicKey); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pb_h132_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*P256PrivateKey); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pb_h132_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*EncryptedSymmetricKey); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pb_h132_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*Letter); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pb_h132_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*Envelope); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pb_h132_proto_msgTypes[5].Exporter = func(v any, i int) any {
			switch v := v.(*EmergencyBackupKey); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pb_h132_proto_msgTypes[6].Exporter = func(v any, i int) any {
			switch v := v.(*WebAuthnWrappedTPMKey); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pb_h132_proto_msgTypes[7].Exporter = func(v any, i int) any {
			switch v := v.(*KeyImpl); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pb_h132_proto_msgTypes[8].Exporter = func(v any, i int) any {
			switch v := v.(*LetterWritingSet); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_pb_h132_proto_msgTypes[7].OneofWrappers = []any{
		(*KeyImpl_DebugRawKey)(nil),
		(*KeyImpl_EmergencyBackupKey)(nil),
		(*KeyImpl_WebauthnWrappedTpm)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_pb_h132_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_pb_h132_proto_goTypes,
		DependencyIndexes: file_pb_h132_proto_depIdxs,
		MessageInfos:      file_pb_h132_proto_msgTypes,
	}.Build()
	File_pb_h132_proto = out.File
	file_pb_h132_proto_rawDesc = nil
	file_pb_h132_proto_goTypes = nil
	file_pb_h132_proto_depIdxs = nil
}
