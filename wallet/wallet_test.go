package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"testing"
)

func TestCreateKey(t *testing.T) {
	var pKey *ecdsa.PrivateKey
	var pubKey *ecdsa.PublicKey
	identity := CreateKey("test")
	if identity.Name != "test" {
		t.Logf("expected the name to be %s found %s", "test", identity.Name)
		t.FailNow()
	}
	if identity.Private == "" {
		t.Log("private key is not set")
		t.FailNow()
	}
	pKey = parsePEM(identity.Private, t)
	if identity.Public == "" {
		t.Log("public key is not found")
		t.FailNow()
	}
	pubKey = parsePubKey(identity.Public, t)
	expectedPub := pKey.PublicKey
	if expectedPub.X.Cmp(pubKey.X) != 0 || expectedPub.Y.Cmp(pubKey.Y) != 0 {
		t.Log("expected public key and the returned public key are not the same")
		t.FailNow()
	}
	if identity.Address != encodeAddress(&pKey.PublicKey, t) {
		t.Log("expected address not found")
		t.FailNow()
	}
}

func parsePEM(privKey string, t *testing.T) *ecdsa.PrivateKey {
	block, _ := pem.Decode([]byte(privKey))
	pkey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		t.Log("unable to parse private key bytes")
		t.Fail()
	}
	return pkey
}

func parsePubKey(pubKey string, t *testing.T) *ecdsa.PublicKey {
	pubBytes, err := hex.DecodeString(pubKey)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}
	x, y := elliptic.Unmarshal(elliptic.P256(), pubBytes)
	if x == nil {
		t.Log("invalid ECDSA public key detected")
		t.Fail()
	}
	return &ecdsa.PublicKey{elliptic.P256(), x, y}
}

func encodeAddress(pubKey *ecdsa.PublicKey, t *testing.T) string {
	pubBytes := elliptic.Marshal(elliptic.P256(), pubKey.X, pubKey.Y)
	hash1 := sha256.New()
	hash2 := sha256.New()
	if _, err := hash1.Write(pubBytes); err != nil {
		t.Log(err)
		t.FailNow()
	}
	if _, err := hash2.Write(hash1.Sum(nil)); err != nil {
		t.Log(err)
		t.FailNow()
	}
	return base64.StdEncoding.EncodeToString(hash2.Sum(nil))
}
