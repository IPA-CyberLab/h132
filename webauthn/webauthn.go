package webauthn

import (
	"bytes"
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/url"

	"github.com/go-jose/go-jose/v4"
	"github.com/go-webauthn/webauthn/protocol"
	gowebauthn "github.com/go-webauthn/webauthn/webauthn"
	"golang.org/x/crypto/hkdf"
)

func genUrl(reflectorUrl *url.URL, jsonbs []byte) (string, error) {
	var gzipped bytes.Buffer
	zw := gzip.NewWriter(&gzipped)
	if _, err := zw.Write(jsonbs); err != nil {
		return "", fmt.Errorf("failed to gzip.Write: %w", err)
	}
	if err := zw.Close(); err != nil {
		return "", fmt.Errorf("failed to gzip.Close: %w", err)
	}

	fragment := base64.RawStdEncoding.WithPadding(base64.StdPadding).EncodeToString(gzipped.Bytes())
	registrationUrl := url.URL{
		Scheme:      reflectorUrl.Scheme,
		User:        reflectorUrl.User,
		Host:        reflectorUrl.Host,
		Path:        reflectorUrl.Path,
		RawPath:     reflectorUrl.RawPath,
		RawQuery:    "",
		Fragment:    fragment,
		RawFragment: "",
	}

	return registrationUrl.String(), nil
}

type Credential struct {
	UserName        string
	ReflectorUrlStr string

	// WebAuthn Credential JSON
	WacJson []byte
}

type wauser struct {
	Name  string
	Creds []gowebauthn.Credential
}

var _ = gowebauthn.User(&wauser{})

func userFromCredential(cred *Credential) (*wauser, error) {
	var wacred gowebauthn.Credential
	if err := json.Unmarshal(cred.WacJson, &wacred); err != nil {
		return nil, fmt.Errorf("Failed to unmarshal Credential: %w", err)
	}

	return &wauser{
		Name:  cred.UserName,
		Creds: []gowebauthn.Credential{wacred},
	}, nil
}

func (u *wauser) WebAuthnID() []byte {
	return []byte(u.Name)
}

func (u *wauser) WebAuthnName() string {
	return u.Name
}

func (u *wauser) WebAuthnDisplayName() string {
	return u.Name
}

func (u *wauser) WebAuthnCredentials() []gowebauthn.Credential {
	return u.Creds
}

type WebAuthnRegistrationSession struct {
	RegistrationUrlStr string

	reflectorUrl *url.URL
	wan          *gowebauthn.WebAuthn
	wau          *wauser
	was          *gowebauthn.SessionData
}

func gowebauthnFromReflectorUrlStr(reflectorUrl *url.URL) (*gowebauthn.WebAuthn, error) {
	originUrl := url.URL{
		Scheme: reflectorUrl.Scheme,
		Host:   reflectorUrl.Host,
	}

	cfg := gowebauthn.Config{
		RPID:          reflectorUrl.Hostname(),
		RPDisplayName: "h132",
		RPOrigins:     []string{originUrl.String()},
	}
	return gowebauthn.New(&cfg)
}

func StartRegistration(reflectorUrlStr, userName string) (*WebAuthnRegistrationSession, error) {
	reflectorUrl, err := url.Parse(reflectorUrlStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse reflector URL: %w", err)
	}

	wan, err := gowebauthnFromReflectorUrlStr(reflectorUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to webauthn.New: %w", err)
	}

	first := make([]byte, 32)
	for i := range first {
		first[i] = byte(i)
	}

	wau := &wauser{Name: userName}

	TRUE := true
	creation, was, err := wan.BeginRegistration(
		wau,
		gowebauthn.WithAuthenticatorSelection(protocol.AuthenticatorSelection{
			AuthenticatorAttachment: protocol.CrossPlatform,
			// Require ResidentKey (PassKey) - Windows Hello requires this for prf extension
			RequireResidentKey: &TRUE,
			ResidentKey:        protocol.ResidentKeyRequirementRequired,
		}),
		gowebauthn.WithExtensions(protocol.AuthenticationExtensions{
			"prf": map[string]any{
				"eval": map[string]any{
					"first": protocol.URLEncodedBase64(first),
				},
			},
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to webauthn.BeginRegistration: %w", err)
	}

	registrationData := struct {
		Action   string                       `json:"action"`
		Creation *protocol.CredentialCreation `json:"creation"`
	}{
		Action:   "registration",
		Creation: creation,
	}

	jsonbs, err := json.Marshal(registrationData)
	if err != nil {
		return nil, fmt.Errorf("failed to json.Marshal(registrationOptions): %w", err)
	}
	registrationUrlStr, err := genUrl(reflectorUrl, jsonbs)
	if err != nil {
		return nil, fmt.Errorf("failed to genUrl: %w", err)
	}

	return &WebAuthnRegistrationSession{
		RegistrationUrlStr: registrationUrlStr,

		reflectorUrl: reflectorUrl,
		wan:          wan,
		wau:          wau,
		was:          was,
	}, nil
}

func B64ToBytes(b64 []byte) ([]byte, error) {
	for len(b64) > 0 {
		last := b64[len(b64)-1]
		if last == '.' || last == '\r' || last == '\n' || last == ' ' {
			b64 = b64[:len(b64)-1]
			continue
		}
		break
	}
	if len(b64) == 0 {
		return nil, fmt.Errorf("no input given")
	}

	enc := base64.RawStdEncoding.WithPadding(base64.StdPadding)

	blen := enc.DecodedLen(len(b64))
	bs := make([]byte, blen)
	blen, err := enc.Decode(bs, b64)
	if err != nil {
		return nil, fmt.Errorf("failed to base64 decode: %w", err)
	}
	bs = bs[:blen]

	return bs, nil
}

func Zb64ToBytes(zb64 []byte) ([]byte, error) {
	z, err := B64ToBytes(zb64)
	if err != nil {
		return nil, err
	}

	zreader, err := gzip.NewReader(bytes.NewBuffer(z))
	if err != nil {
		return nil, fmt.Errorf("failed to gzip decode (new): %w", err)
	}
	bs, err := io.ReadAll(zreader)
	if err != nil {
		return nil, fmt.Errorf("failed to gzip decode (read): %w", err)
	}

	return bs, nil
}

func (sess *WebAuthnRegistrationSession) Complete(zb64resp []byte) (*Credential, error) {
	regJson, err := Zb64ToBytes(zb64resp)
	if err != nil {
		return nil, err
	}

	parsed, err := protocol.ParseCredentialCreationResponseBytes(regJson)
	if err != nil {
		return nil, fmt.Errorf("failed to parse as CredentialCreationResponse: %w", err)
	}

	cred, err := sess.wan.CreateCredential(sess.wau, *sess.was, parsed)
	if err != nil {
		return nil, fmt.Errorf("failed to webauthn.CreateCredential: %w", err)
	}

	credJson, err := json.Marshal(cred)
	if err != nil {
		return nil, fmt.Errorf("failed to json.Marshal(credential): %w", err)
	}

	return &Credential{
		UserName:        sess.wau.Name,
		ReflectorUrlStr: sess.reflectorUrl.String(),
		WacJson:         credJson,
	}, nil
}

type GetPRFSecretSession struct {
	GetPRFSecretUrlStr string

	reflectorUrl *url.URL
	wan          *gowebauthn.WebAuthn
	wau          *wauser
	was          *gowebauthn.SessionData
	ephPriv      *ecdsa.PrivateKey
}

func StartGetPRFSecretSession(cred *Credential, salt []byte, randr io.Reader) (*GetPRFSecretSession, error) {
	reflectorUrl, err := url.Parse(cred.ReflectorUrlStr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse reflector URL: %w", err)
	}

	if len(salt) != 32 {
		return nil, fmt.Errorf("salt must be 32 bytes")
	}

	wan, err := gowebauthnFromReflectorUrlStr(reflectorUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to webauthn.New: %w", err)
	}

	wau, err := userFromCredential(cred)
	if err != nil {
		return nil, err
	}

	assertion, was, err := wan.BeginLogin(
		wau,
		gowebauthn.WithAssertionExtensions(protocol.AuthenticationExtensions{
			"prf": map[string]any{
				"eval": map[string]any{
					"first": protocol.URLEncodedBase64(salt),
				},
			},
		}),
	)
	if err != nil {
		return nil, err
	}

	ephPriv, err := ecdsa.GenerateKey(elliptic.P256(), randr)
	if err != nil {
		return nil, fmt.Errorf("failed to ecdsa.GenerateKey: %w", err)
	}

	loginData := struct {
		Action    string                        `json:"action"`
		Assertion *protocol.CredentialAssertion `json:"assertion"`
		UserName  string                        `json:"userName"`
		ECDHKey   jose.JSONWebKey               `json:"ecdhKey"`
	}{
		Action:    "login",
		Assertion: assertion,
		UserName:  cred.UserName,
		ECDHKey: jose.JSONWebKey{
			Key: ephPriv.Public(),
		},
	}

	jsonbs, err := json.Marshal(loginData)
	if err != nil {
		return nil, fmt.Errorf("failed to json.Marshal(loginOptions): %w", err)
	}

	generatedUrl, err := genUrl(reflectorUrl, jsonbs)
	if err != nil {
		return nil, fmt.Errorf("failed to genUrl: %w", err)
	}

	return &GetPRFSecretSession{
		GetPRFSecretUrlStr: generatedUrl,

		reflectorUrl: reflectorUrl,
		wan:          wan,
		wau:          wau,
		was:          was,
		ephPriv:      ephPriv,
	}, nil
}

var LoginSalt = []byte("h132_login")

func (sess *GetPRFSecretSession) deriveCEK(pubECDSA *ecdsa.PublicKey) ([]byte, error) {
	pub, err := pubECDSA.ECDH()
	if err != nil {
		return nil, fmt.Errorf("failed to pubECDSA.ECDH: %w", err)
	}
	priv, err := sess.ephPriv.ECDH()
	if err != nil {
		return nil, fmt.Errorf("failed to ephPriv.ECDH: %w", err)
	}
	ss, err := priv.ECDH(pub)
	if err != nil {
		return nil, fmt.Errorf("failed to priv.ECDH(pub): %w", err)
	}
	hkdfr := hkdf.New(sha256.New, ss, LoginSalt, nil)
	cek := make([]byte, 32)
	if _, err := hkdfr.Read(cek); err != nil {
		return nil, fmt.Errorf("failed to HKDF.read: %w", err)
	}
	return cek, nil
}

func decryptAESGCM(cek, nonce, encBs []byte) ([]byte, error) {
	aesc, err := aes.NewCipher(cek)
	if err != nil {
		return nil, fmt.Errorf("failed to aes.NewCipher: %w", err)
	}
	aead, err := cipher.NewGCM(aesc)
	if err != nil {
		return nil, fmt.Errorf("failed to cipher.NewGCM: %w", err)
	}
	return aead.Open(nil, nonce, encBs, nil)
}

// Complete completes the GetPRFSecretSession and returns the PRF secret and the credential public key
func (sess *GetPRFSecretSession) Complete(b64resp []byte) ([]byte, error) {
	encJson, err := B64ToBytes(b64resp)
	if err != nil {
		return nil, err
	}

	loginEnc := struct {
		NonceB64     string          `json:"nonceB64"`
		EncryptedB64 string          `json:"encryptedB64"`
		PubJWK       jose.JSONWebKey `json:"pubJwk"`
	}{}
	if err := json.Unmarshal(encJson, &loginEnc); err != nil {
		return nil, fmt.Errorf("failed to unmarshal loginEnc: %w", err)
	}

	b64Enc := base64.RawStdEncoding.WithPadding(base64.StdPadding)

	nonceB64, err := b64Enc.DecodeString(loginEnc.NonceB64)
	if err != nil {
		return nil, fmt.Errorf("failed to base64 decode loginEnc.NonceB64: %w", err)
	}

	encBs, err := b64Enc.DecodeString(loginEnc.EncryptedB64)
	if err != nil {
		return nil, fmt.Errorf("failed to base64 decode loginEnc.EncryptedB64: %w", err)
	}

	pubECDSA, ok := loginEnc.PubJWK.Key.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("loginEnc.PubJWK.Key is not *ecdsa.PublicKey")
	}
	cek, err := sess.deriveCEK(pubECDSA)
	if err != nil {
		return nil, err
	}
	gzbs, err := decryptAESGCM(cek, nonceB64, encBs)
	if err != nil {
		return nil, fmt.Errorf("failed to aead.Open: %w", err)
	}
	zreader, err := gzip.NewReader(bytes.NewBuffer(gzbs))
	if err != nil {
		return nil, fmt.Errorf("failed to gzip.NewReader: %w", err)
	}
	loginJson, err := io.ReadAll(zreader)
	if err != nil {
		return nil, fmt.Errorf("failed to gzip.Read: %w", err)
	}

	parsed, err := protocol.ParseCredentialRequestResponseBytes(loginJson)
	if err != nil {
		return nil, fmt.Errorf("failed to parse as CredentialAssertionResponse: %w", err)
	}

	if _, err = sess.wan.ValidateLogin(sess.wau, *sess.was, parsed); err != nil {
		return nil, fmt.Errorf("failed to webauthn.ValidateLogin: %w", err)
	}

	prfany, ok := parsed.ClientExtensionResults["prf"]
	if !ok {
		return nil, fmt.Errorf("no prf in ClientExtensionResults")
	}
	prfmap, ok := prfany.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("prf is not a map")
	}
	prfResultsAny, ok := prfmap["results"]
	if !ok {
		return nil, fmt.Errorf("no results in prf")
	}
	prfResultsMap, ok := prfResultsAny.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("prf.results is not a map")
	}
	first, ok := prfResultsMap["first"]
	if !ok {
		return nil, fmt.Errorf("no first in prf.results")
	}
	firststr, ok := first.(string)
	if !ok {
		return nil, fmt.Errorf("prf.results.first is not a string")
	}

	prf, err := base64.RawStdEncoding.WithPadding(base64.StdPadding).DecodeString(firststr)
	if err != nil {
		return nil, fmt.Errorf("failed to base64 decode prf: %w", err)
	}
	return prf, nil
}
