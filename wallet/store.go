package wallet

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	Public  string
	Address string
}

var wallet *Wallet

func logError(prefix string, err error) {
	fmt.Printf("%s - %s", prefix, err)
}

func init() {
	loadWallet()
}

func (w *Wallet) Save() error {
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(*w)
	if err != nil {
		return err
	}
	if wPath, err := walletPath(); err != nil {
		return err
	} else {
		return ioutil.WriteFile(wPath, buf.Bytes(), 0600)
	}
}
func walletPath() (string, error) {
	dir, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	prefPath := filepath.Join(dir, PrefPath)
	err = os.MkdirAll(prefPath, 0755)
	if err != nil {
		return "", err
	}
	return filepath.Join(prefPath, "wallet.json"), nil
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
func loadWallet() {
	wPath, err := walletPath()
	if err != nil {
		logError("unable to get walletPath", err)
		wallet = &Wallet{}
		return
	}
	_, err = os.Stat(wPath)
	if os.IsNotExist(err) {
		logError(fmt.Sprintf("walletPath does not exist %s", wPath), err)
		wallet = &Wallet{}
		return
	}
	walletBytes, err := ioutil.ReadFile(wPath)
	if err != nil {
		logError(fmt.Sprintf("Unable to read file %s", wPath), err)
		wallet = &Wallet{}
		return
	}
	wlt := Wallet{}
	if err = json.NewDecoder(bytes.NewBuffer(walletBytes)).Decode(&wlt); err != nil {
		logError(fmt.Sprintf("Unable to decode json %s", string(walletBytes)), err)
		wallet = &Wallet{}
		return
	}
	wallet = &wlt
}
