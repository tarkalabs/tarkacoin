package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
)

func CreateKey(name string) Identity {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	derKey, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		panic(err)
	}
	privateBlock := &pem.Block{Type: "PRIVATE KEY", Bytes: derKey}
	privateKey := pem.EncodeToMemory(privateBlock)

	pubBytes := elliptic.Marshal(key.PublicKey.Curve, key.PublicKey.X, key.PublicKey.Y)
	pubAddress := hex.EncodeToString(pubBytes)
	address, err := EncodeAddress(key.PublicKey)
	if err != nil {
		panic(err)
	}

	return Identity{Name: name, Private: string(privateKey), Address: address, Public: pubAddress}
}

// ParsePublicKey parses a hex encoded public key into a ECDSA public key
func ParsePublicKey(publicKey string) (*ecdsa.PublicKey, error) {
	pubBytes, err := hex.DecodeString(publicKey)
	if err != nil {
		return nil, err
	}
	x, y := elliptic.Unmarshal(elliptic.P256(), pubBytes)
	if x == nil {
		return nil, errors.New("invalid public key")
	}
	return &ecdsa.PublicKey{elliptic.P256(), x, y}, nil
}

func ParsePrivateKey(privKey string) (*ecdsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privKey))
	pkey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return pkey, nil
}

func EncodeAddress(pubKey ecdsa.PublicKey) (string, error) {
	pubBytes := elliptic.Marshal(elliptic.P256(), pubKey.X, pubKey.Y)
	hash1 := sha256.New()
	hash2 := sha256.New()
	if _, err := hash1.Write(pubBytes); err != nil {
		return "", err
	}
	if _, err := hash2.Write(hash1.Sum(nil)); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(hash2.Sum(nil)), nil
}

func SaveKey(keyName string) error {
	wallet := CurrentWallet()
	wallet.Identities = append(wallet.Identities, CreateKey(keyName))
	return wallet.Save()
}
