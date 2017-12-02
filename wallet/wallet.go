package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
)

func CreateKey(name string) Identity {
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	publicKey := &privateKey.PublicKey

	x509Encoded, _ := x509.MarshalECPrivateKey(privateKey)
	x509EncodedPub, _ := x509.MarshalPKIXPublicKey(publicKey)

	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: x509Encoded})
	pemEncodedPub := pem.EncodeToMemory(&pem.Block{Type: "EC PUBLIC KEY", Bytes: x509EncodedPub})

	walletAddress, _ := EncodeAddress(publicKey)

	return Identity{Name: name, Private: string(pemEncoded), Address: walletAddress, Public: string(pemEncodedPub)}
}

func EncodeAddress(pubKey *ecdsa.PublicKey) (string, error) {
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
