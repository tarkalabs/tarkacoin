package wallet

import (
	"bytes"
	"encoding/json"
	homedir "github.com/mitchellh/go-homedir"
	"io/ioutil"
	"os"
	"path/filepath"
)

const PrefPath = ".tarkacoin"

type Wallet struct {
	Identities []Identity
}

type Identity struct {
	Name    string
	Private string
	Address string
}

var wallet *Wallet

func init() {
	wallet = &Wallet{}
}

func (w *Wallet) Save() error {
	dir, err := homedir.Dir()
	if err != nil {
		return err
	}
	prefPath := filepath.Join(dir, PrefPath)
	err = os.MkdirAll(prefPath, 0755)
	if err != nil {
		return err
	}
	walletPath := filepath.Join(prefPath, "wallet.json")
	buf := new(bytes.Buffer)
	err = json.NewEncoder(buf).Encode(*w)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(walletPath, buf.Bytes(), 0600)
}

func CurrentWallet() *Wallet {
	return wallet
}
func ensureFile(walletPath string) (*os.File, error) {
	_, err := os.Stat(walletPath)
	if os.IsNotExist(err) {
		return os.Create(walletPath)
	}
	return os.Open(walletPath)
}
