package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
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

	address := hex.EncodeToString(derKey)

	return Identity{Name: name, Private: string(privateKey), Address: address}
}

func ParsePublicAddress(address string) (*ecdsa.PublicKey, error) {
	der, err := hex.DecodeString(address)
	if err != nil {
		return nil, err
	}
	pub, err := x509.ParsePKIXPublicKey(der)
	if err != nil {
		return nil, err
	}
	publicKey := pub.(ecdsa.PublicKey)
	return &publicKey, nil
}

func ParsePrivateKey(privKey string) (*ecdsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privKey))
	pkey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return pkey, nil
}

func SaveKey(keyName string) error {
	wallet := CurrentWallet()
	wallet.Identities = append(wallet.Identities, CreateKey(keyName))
	return wallet.Save()
}
