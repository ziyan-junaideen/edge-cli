package secrets

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/99designs/keyring"
	"github.com/adrg/xdg"
)

const serviceName = "edge-cli"

type Store struct {
	keyring keyring.Keyring
}

func Open() (Store, error) {
	fileDir, err := xdg.ConfigFile(filepath.Join("edge", "keyring"))
	if err != nil {
		return Store{}, fmt.Errorf("resolve keyring path: %w", err)
	}

	keychain, err := keyring.Open(keyring.Config{
		ServiceName:      serviceName,
		FileDir:          fileDir,
		FilePasswordFunc: filePassword,
	})
	if err != nil {
		return Store{}, fmt.Errorf("open keyring: %w", err)
	}
	return Store{keyring: keychain}, nil
}

func (store Store) SetToken(profileName string, token string) error {
	return store.keyring.Set(keyring.Item{
		Key:  tokenKey(profileName),
		Data: []byte(token),
	})
}

func (store Store) Token(profileName string) (string, error) {
	item, err := store.keyring.Get(tokenKey(profileName))
	if errors.Is(err, keyring.ErrKeyNotFound) {
		return "", ErrNotFound
	}
	if err != nil {
		return "", err
	}
	return string(item.Data), nil
}

func (store Store) DeleteToken(profileName string) error {
	err := store.keyring.Remove(tokenKey(profileName))
	if errors.Is(err, keyring.ErrKeyNotFound) {
		return nil
	}
	return err
}

var ErrNotFound = errors.New("token is not configured")

func tokenKey(profileName string) string {
	return "token:" + profileName
}

func filePassword(prompt string) (string, error) {
	password := os.Getenv("EDGE_KEYRING_PASSWORD")
	if password == "" {
		return "", errors.New("file keyring backend requires EDGE_KEYRING_PASSWORD")
	}
	return password, nil
}
