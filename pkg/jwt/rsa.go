package jwt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"os"
)

var (
	ErrUnableGeneratePrivateKey = errors.New("unable to generate private key file")
	ErrUnableCreatePrivateKey   = errors.New("unable to create private key file")
	ErrUnableEncodePrivateKey   = errors.New("unable to encode private key")
	ErrNoPrivateKeyFile         = errors.New("private key file unreadable")
	ErrInvalidPrivateKey        = errors.New("private key invalid")
	ErrInvalidPublicKey         = errors.New("public key invalid")
)

type EncryptionRsa struct {
	PrivateKeyPassphrase string

	PrivateKeyFile string
	PublicKeyFile  string

	PrivateKeyBytes []byte
	PublicKeyBytes  []byte

	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

func New(PrivateKeyFile, PrivateKeyPassphrase string) (*EncryptionRsa, error) {
	r := &EncryptionRsa{
		PrivateKeyFile:       PrivateKeyFile,
		PrivateKeyPassphrase: PrivateKeyPassphrase,
	}

	err := r.initPrivateKey()
	if err != nil {
		return nil, err
	}

	err = r.initPublicKey()
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (r *EncryptionRsa) initPrivateKey() error {
	_, err := os.Lstat(r.PrivateKeyFile)
	if err != nil {
		if os.IsNotExist(err) {
			if err = r.generatePrivateKey(); err != nil {
				return err
			}
			return nil
		}
		return ErrNoPrivateKeyFile
	}

	r.PrivateKeyBytes, err = os.ReadFile(r.PrivateKeyFile)
	if err != nil {
		return ErrNoPrivateKeyFile
	}

	r.PrivateKey, err = jwt.ParseRSAPrivateKeyFromPEM(r.PrivateKeyBytes)
	if err != nil {
		return ErrInvalidPrivateKey
	}

	return nil
}

func (r *EncryptionRsa) initPublicKey() error {
	r.PublicKey = &r.PrivateKey.PublicKey
	keyBytes, err := x509.MarshalPKIXPublicKey(r.PublicKey)
	if err != nil {
		return ErrInvalidPublicKey
	}

	keyBlock := pem.Block{
		Type:    "PUBLIC KEY",
		Headers: nil,
		Bytes:   keyBytes,
	}
	r.PublicKeyBytes = pem.EncodeToMemory(&keyBlock)

	return nil
}

func (r *EncryptionRsa) generatePrivateKey() error {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return ErrUnableGeneratePrivateKey
	}

	r.PrivateKey = key
	keyBlock := &pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   x509.MarshalPKCS1PrivateKey(r.PrivateKey),
	}

	privatePem, err := os.Create(r.PrivateKeyFile)
	if err != nil {
		return ErrUnableCreatePrivateKey
	}

	err = pem.Encode(privatePem, keyBlock)
	if err != nil {
		return ErrUnableEncodePrivateKey
	}

	r.PrivateKeyBytes = pem.EncodeToMemory(keyBlock)

	return nil
}
